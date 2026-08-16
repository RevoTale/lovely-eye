package app

import (
	"context"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/lovely-eye/server/internal/platform/config"
	"github.com/stretchr/testify/require"
)

func TestNewReportsConstructionConfigurationFailures(t *testing.T) {
	tests := []struct {
		name          string
		configure     func(*config.Config)
		expectedError string
	}{
		{
			name: "missing tracker artifact",
			configure: func(cfg *config.Config) {
				cfg.TrackerJS = nil
			},
			expectedError: "load tracker.js",
		},
		{
			name: "unsupported database driver",
			configure: func(cfg *config.Config) {
				cfg.Database.Driver = "unsupported"
			},
			expectedError: "unsupported database driver",
		},
		{
			name: "invalid trusted proxy",
			configure: func(cfg *config.Config) {
				cfg.Analytics.TrustedProxyCIDRs = []string{"invalid-cidr"}
			},
			expectedError: "configure trusted proxies",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := constructionTestConfig(t.Name())
			tt.configure(&cfg)
			application, err := New(t.Context(), cfg)
			require.ErrorContains(t, err, tt.expectedError)
			require.Nil(t, application)
		})
	}
}

func TestRunStopsWhenContextIsCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	application := &App{
		HTTPServer: &http.Server{Addr: "127.0.0.1:0", ReadHeaderTimeout: time.Second},
	}

	require.NoError(t, application.Run(ctx))
}

func TestRunReturnsListenFailure(t *testing.T) {
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, listener.Close()) })
	application := &App{
		HTTPServer: &http.Server{Addr: listener.Addr().String(), ReadHeaderTimeout: time.Second},
	}

	err = application.Run(t.Context())
	require.ErrorContains(t, err, "serve http")
}

func constructionTestConfig(testName string) config.Config {
	return config.Config{
		Server: config.ServerConfig{Host: "127.0.0.1", Port: "0", DashboardPath: "missing"},
		Database: config.DatabaseConfig{
			Driver:         config.DBDriverSQLite,
			DSN:            "file:" + strings.ReplaceAll(testName, "/", "_") + "?mode=memory&cache=shared",
			ConnectTimeout: time.Second,
		},
		Auth: config.AuthConfig{
			JWTSecret:         strings.Repeat("j", 32),
			AccessTokenExpiry: 15 * time.Minute,
			RefreshExpiry:     7 * 24 * time.Hour,
		},
		Analytics: config.AnalyticsConfig{
			IdentitySecret:    strings.Repeat("a", 32),
			TrustedProxyCIDRs: []string{"127.0.0.1/32"},
		},
		TrackerJS: []byte("console.log('tracker');"),
	}
}
