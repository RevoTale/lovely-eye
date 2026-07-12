package handlers

import (
	"sync"
	"time"
)

type CollectRateLimiter struct {
	enabled bool
	rate    float64
	burst   float64
	now     func() time.Time

	mu          sync.Mutex
	buckets     map[string]*collectBucket
	lastCleanup time.Time
}

type collectBucket struct {
	tokens float64
	seenAt time.Time
}

func NewCollectRateLimiter(enabled bool, perMinute, burst int) *CollectRateLimiter {
	if perMinute <= 0 || burst <= 0 {
		enabled = false
	}
	return &CollectRateLimiter{
		enabled: enabled,
		rate:    float64(perMinute) / 60.0,
		burst:   float64(burst),
		now:     time.Now,
		buckets: make(map[string]*collectBucket),
	}
}

func (l *CollectRateLimiter) Allow(key string) bool {
	if l == nil || !l.enabled {
		return true
	}

	now := l.now()
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.lastCleanup.IsZero() {
		l.lastCleanup = now
	} else if now.Sub(l.lastCleanup) > 10*time.Minute {
		for key, bucket := range l.buckets {
			if now.Sub(bucket.seenAt) > 15*time.Minute {
				delete(l.buckets, key)
			}
		}
		l.lastCleanup = now
	}

	bucket := l.buckets[key]
	if bucket == nil {
		bucket = &collectBucket{tokens: l.burst, seenAt: now}
		l.buckets[key] = bucket
	}

	elapsed := now.Sub(bucket.seenAt).Seconds()
	if elapsed > 0 {
		bucket.tokens = min(l.burst, bucket.tokens+elapsed*l.rate)
	}
	bucket.seenAt = now

	if bucket.tokens < 1 {
		return false
	}
	bucket.tokens--
	return true
}
