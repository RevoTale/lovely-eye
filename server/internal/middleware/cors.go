package middleware

import (
	"net/http"
	"strings"
)

func isAnalyticsPath(path string) bool {
	return strings.HasSuffix(path, "/api/collect") || strings.HasSuffix(path, "/api/event")
}

func isTrackerPath(path string) bool {
	return strings.HasSuffix(path, "/tracker.js")
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Analytics endpoints handle their own CORS in the handler
		if isAnalyticsPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// tracker.js sets its own CORS headers
		if isTrackerPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		origin := r.Header.Get("Origin")
		if origin != "" {
			// Only allow same-origin requests for GraphQL/dashboard API
			if isSameOrigin(origin, r.Host) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type")
				w.Header().Set("Access-Control-Max-Age", "3600")
				w.Header().Set("Vary", "Origin")
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

func isSameOrigin(origin, host string) bool {
	originHost := origin
	if idx := strings.Index(origin, "://"); idx != -1 {
		originHost = origin[idx+3:]
	}
	originHost = strings.TrimSuffix(originHost, "/")
	return originHost == host
}
