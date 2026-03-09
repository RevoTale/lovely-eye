package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoad_UsesAnalyticsIdentitySecretOverride(t *testing.T) {
	t.Setenv("JWT_SECRET", strings.Repeat("j", 32))
	t.Setenv("ANALYTICS_IDENTITY_SECRET", strings.Repeat("a", 32))

	cfg := Load()

	require.Equal(t, strings.Repeat("a", 32), cfg.Analytics.IdentitySecret)
}

func TestLoad_FallsBackToJWTSecretForAnalyticsIdentity(t *testing.T) {
	t.Setenv("JWT_SECRET", strings.Repeat("j", 32))
	t.Setenv("ANALYTICS_IDENTITY_SECRET", "")

	cfg := Load()

	require.Equal(t, strings.Repeat("j", 32), cfg.Analytics.IdentitySecret)
}
