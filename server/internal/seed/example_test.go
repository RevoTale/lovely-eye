package seed

import (
	"testing"
	"time"

	"github.com/lovely-eye/server/internal/platform/config"
	"github.com/stretchr/testify/require"
)

func TestRunLoadsReusableExampleData(t *testing.T) {
	cfg := config.Config{
		Database: config.DatabaseConfig{
			Driver:         config.DBDriverSQLite,
			DSN:            "file:" + t.TempDir() + "/seed.db?mode=rwc",
			MaxConns:       1,
			MinConns:       1,
			ConnectTimeout: time.Second,
		},
		Auth: config.AuthConfig{
			InitialAdminUsername: "seed-admin",
			InitialAdminPassword: "seed-password",
		},
	}

	first, err := Run(t.Context(), cfg)
	require.NoError(t, err)
	require.True(t, first.CreatedSite)
	require.Equal(t, "Localhost", first.SiteName)
	require.NotEmpty(t, first.PublicKey)
	require.Equal(t, defaultUsers, first.Clients)
	require.GreaterOrEqual(t, first.Sessions, defaultUsers*minSessions)
	require.LessOrEqual(t, first.Sessions, defaultUsers*maxSessions)
	require.GreaterOrEqual(t, first.PageViews, first.Sessions)
	require.Positive(t, first.PredefinedEvents)

	second, err := Run(t.Context(), cfg)
	require.NoError(t, err)
	require.False(t, second.CreatedSite)
	require.Equal(t, first.PublicKey, second.PublicKey)
	require.Equal(t, defaultUsers, second.Clients)
}
