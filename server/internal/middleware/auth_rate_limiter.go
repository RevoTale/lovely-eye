package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/lovely-eye/server/pkg/clientip"
)

type AuthRateLimiter struct {
	enabled      bool
	attempts     int
	window       time.Duration
	maxBodyBytes int64
	ipResolver   *clientip.Resolver
	now          func() time.Time

	mu          sync.Mutex
	buckets     map[string]*authRateBucket
	lastCleanup time.Time
}

type authRateBucket struct {
	count   int
	resetAt time.Time
	seenAt  time.Time
}

func NewAuthRateLimiter(enabled bool, attempts int, window time.Duration, maxBodyBytes int64, ipResolver *clientip.Resolver) *AuthRateLimiter {
	if attempts <= 0 || window <= 0 {
		enabled = false
	}
	if maxBodyBytes <= 0 {
		maxBodyBytes = 1024 * 1024
	}
	if ipResolver == nil {
		ipResolver = clientip.MustNewResolver(nil)
	}
	return &AuthRateLimiter{
		enabled:      enabled,
		attempts:     attempts,
		window:       window,
		maxBodyBytes: maxBodyBytes,
		ipResolver:   ipResolver,
		now:          time.Now,
		buckets:      make(map[string]*authRateBucket),
	}
}

func (l *AuthRateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if l == nil || !l.enabled || r.Method != http.MethodPost {
			next.ServeHTTP(w, r)
			return
		}

		body, ok := readBoundedRequestBody(w, r, l.maxBodyBytes)
		if !ok {
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(body))

		if !containsAuthMutation(body) {
			next.ServeHTTP(w, r)
			return
		}

		ip := l.ipResolver.GetClientIP(r.Header.Get("X-Forwarded-For"), r.Header.Get("X-Real-IP"), r.RemoteAddr)
		if !l.allow("auth|" + ip) {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func readBoundedRequestBody(w http.ResponseWriter, r *http.Request, maxBodyBytes int64) ([]byte, bool) {
	body, err := io.ReadAll(io.LimitReader(r.Body, maxBodyBytes+1))
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return nil, false
	}
	if int64(len(body)) > maxBodyBytes {
		http.Error(w, "request body is too large", http.StatusRequestEntityTooLarge)
		return nil, false
	}
	return body, true
}

func containsAuthMutation(body []byte) bool {
	var payload struct {
		Query string `json:"query"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return false
	}
	return containsGraphQLName(payload.Query, "login") || containsGraphQLName(payload.Query, "register")
}

func containsGraphQLName(query, name string) bool {
	for start := 0; start < len(query); {
		index := strings.Index(query[start:], name)
		if index == -1 {
			return false
		}
		index += start
		before := index - 1
		after := index + len(name)
		if (before < 0 || !isGraphQLNameChar(query[before])) && (after >= len(query) || !isGraphQLNameChar(query[after])) {
			return true
		}
		start = after
	}
	return false
}

func isGraphQLNameChar(char byte) bool {
	return char == '_' || (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9')
}

func (l *AuthRateLimiter) allow(key string) bool {
	now := l.now()

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.lastCleanup.IsZero() {
		l.lastCleanup = now
	} else if now.Sub(l.lastCleanup) > l.window {
		for key, bucket := range l.buckets {
			if now.Sub(bucket.seenAt) > l.window {
				delete(l.buckets, key)
			}
		}
		l.lastCleanup = now
	}

	bucket := l.buckets[key]
	if bucket == nil || !now.Before(bucket.resetAt) {
		l.buckets[key] = &authRateBucket{
			count:   1,
			resetAt: now.Add(l.window),
			seenAt:  now,
		}
		return true
	}

	bucket.seenAt = now
	if bucket.count >= l.attempts {
		return false
	}
	bucket.count++
	return true
}
