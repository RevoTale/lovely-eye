package middleware

import (
	"net/http"
	"net/url"
	"strings"
)

func isAnalyticsPath(path string) bool {
	return strings.HasSuffix(path, "/api/collect")
}

func isTrackerPath(path string) bool {
	return strings.HasSuffix(path, "/tracker.js")
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if isAnalyticsPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		if isTrackerPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		origin := r.Header.Get("Origin")
		if origin != "" {

			if isSameOrigin(origin, r) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type")
				w.Header().Set("Access-Control-Max-Age", "3600")
				w.Header().Set("Vary", "Origin")
			} else if r.Method != http.MethodGet && r.Method != http.MethodHead && r.Method != http.MethodOptions {
				http.Error(w, "origin not allowed", http.StatusForbidden)
				return
			}
		}

		if r.Method == "OPTIONS" {
			if w.Header().Get("Access-Control-Allow-Origin") != "" {
				w.WriteHeader(http.StatusNoContent)
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isSameOrigin(origin string, r *http.Request) bool {
	parsed, err := url.Parse(origin)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return false
	}

	return strings.EqualFold(parsed.Scheme, requestScheme(r)) && strings.EqualFold(parsed.Host, r.Host)
}

func requestScheme(r *http.Request) string {
	if scheme := forwardedRequestScheme(r.Header); scheme != "" {
		return scheme
	}
	if r.TLS != nil {
		return "https"
	}
	if r.URL != nil {
		if scheme := normalizedScheme(r.URL.Scheme); scheme != "" {
			return scheme
		}
	}
	return "http"
}

func forwardedRequestScheme(header http.Header) string {
	if forwarded := header.Get("Forwarded"); forwarded != "" {
		firstForwarded := strings.Split(forwarded, ",")[0]
		for _, part := range strings.Split(firstForwarded, ";") {
			key, value, ok := strings.Cut(strings.TrimSpace(part), "=")
			if !ok || !strings.EqualFold(key, "proto") {
				continue
			}
			if scheme := normalizedScheme(strings.Trim(value, `"`)); scheme != "" {
				return scheme
			}
		}
	}

	if forwardedProto := header.Get("X-Forwarded-Proto"); forwardedProto != "" {
		firstProto := strings.Split(forwardedProto, ",")[0]
		return normalizedScheme(firstProto)
	}

	return ""
}

func normalizedScheme(raw string) string {
	scheme := strings.ToLower(strings.TrimSpace(raw))
	switch scheme {
	case "http", "https":
		return scheme
	default:
		return ""
	}
}
