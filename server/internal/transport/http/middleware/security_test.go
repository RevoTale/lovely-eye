package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSecuritySetsBrowserHardeningHeaders(t *testing.T) {
	handler := Security(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", nil))

	require.Equal(t, "nosniff", recorder.Header().Get("X-Content-Type-Options"))
	require.Equal(t, "DENY", recorder.Header().Get("X-Frame-Options"))
	require.Equal(t, "0", recorder.Header().Get("X-XSS-Protection"))
	require.Equal(t, "strict-origin-when-cross-origin", recorder.Header().Get("Referrer-Policy"))
	require.Equal(t, "camera=(), microphone=(), geolocation=()", recorder.Header().Get("Permissions-Policy"))
	require.Equal(t, "max-age=31536000; includeSubDomains", recorder.Header().Get("Strict-Transport-Security"))
	require.Contains(t, recorder.Header().Get("Content-Security-Policy"), "default-src 'self'")
	require.Contains(t, recorder.Header().Get("Content-Security-Policy"), "frame-ancestors 'none'")
}
