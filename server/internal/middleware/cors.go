package middleware

import (
	"net/http"
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
