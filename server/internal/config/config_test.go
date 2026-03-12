package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoad_ResolvesAllowRegistrationDefaultFromInitialAdmin(t *testing.T) {
	t.Setenv("JWT_SECRET", strings.Repeat("j", 32))

	tests := []struct {
		name                 string
		initialAdminUsername string
		initialAdminPassword string
		explicitAllow        *string
		expected             bool
	}{
		{
			name:                 "defaults to false when both initial admin credentials are set",
			initialAdminUsername: "admin",
			initialAdminPassword: "password123",
			expected:             false,
		},
		{
			name:                 "defaults to true when username is missing",
			initialAdminPassword: "password123",
			expected:             true,
		},
		{
			name:                 "defaults to true when password is missing",
			initialAdminUsername: "admin",
			expected:             true,
		},
		{
			name:                 "explicit true overrides derived default",
			initialAdminUsername: "admin",
			initialAdminPassword: "password123",
			explicitAllow:        ptr("true"),
			expected:             true,
		},
		{
			name:                 "explicit false overrides derived default",
			initialAdminUsername: "admin",
			initialAdminPassword: "password123",
			explicitAllow:        ptr("false"),
			expected:             false,
		},
		{
			name:                 "empty allow registration behaves as unset",
			initialAdminUsername: "admin",
			initialAdminPassword: "password123",
			explicitAllow:        ptr(""),
			expected:             false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("INITIAL_ADMIN_USERNAME", tt.initialAdminUsername)
			t.Setenv("INITIAL_ADMIN_PASSWORD", tt.initialAdminPassword)
			if tt.explicitAllow != nil {
				t.Setenv("ALLOW_REGISTRATION", *tt.explicitAllow)
			} else {
				t.Setenv("ALLOW_REGISTRATION", "")
			}

			cfg := Load()

			require.Equal(t, tt.expected, cfg.Auth.AllowRegistration)
		})
	}
}

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

func ptr[T any](value T) *T {
	return &value
}
