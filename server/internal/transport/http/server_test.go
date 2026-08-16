package http

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestHTTPServerHasBoundedHeadersAndTimeouts(t *testing.T) {
	server := newHTTPServer(":0", http.NotFoundHandler())

	require.Equal(t, 5*time.Second, server.ReadHeaderTimeout)
	require.Equal(t, 15*time.Second, server.ReadTimeout)
	require.Equal(t, 15*time.Second, server.WriteTimeout)
	require.Equal(t, 60*time.Second, server.IdleTimeout)
	require.Equal(t, 1<<20, server.MaxHeaderBytes)
}
