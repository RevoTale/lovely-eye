package config

import (
	"strings"
	"testing"
	"time"

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

func TestLoad_UsesAnalyticsHardeningDefaults(t *testing.T) {
	t.Setenv("JWT_SECRET", strings.Repeat("j", 32))

	cfg := Load()

	require.True(t, cfg.Auth.RateLimitEnabled)
	require.Equal(t, 10, cfg.Auth.RateLimitAttempts)
	require.Equal(t, 15*time.Minute, cfg.Auth.RateLimitWindow)
	require.Equal(t, int64(16*1024), cfg.Analytics.MaxBodyBytes)
	require.Equal(t, 8*1024, cfg.Analytics.MaxPropertiesBytes)
	require.Equal(t, 4*time.Hour, cfg.Analytics.MaxSinglePageDuration)
	require.True(t, cfg.Analytics.RateLimitEnabled)
	require.Equal(t, 120, cfg.Analytics.RateLimitPerMinute)
	require.Equal(t, 240, cfg.Analytics.RateLimitBurst)
	require.Contains(t, cfg.Analytics.TrustedProxyCIDRs, "127.0.0.1/32")
	require.Contains(t, cfg.Analytics.TrustedProxyCIDRs, "10.0.0.0/8")
	require.Equal(t, int64(1024*1024), cfg.GraphQL.MaxBodyBytes)
	require.Equal(t, 730, cfg.Dashboard.MaxDailyRangeDays)
	require.Equal(t, 31, cfg.Dashboard.MaxHourlyRangeDays)
	require.Equal(t, 100, cfg.Dashboard.MaxFilterValues)
	require.Equal(t, 2048, cfg.Dashboard.MaxFilterStringLength)
}

func TestLoad_UsesAnalyticsHardeningOverrides(t *testing.T) {
	t.Setenv("JWT_SECRET", strings.Repeat("j", 32))
	t.Setenv("ANALYTICS_MAX_BODY_BYTES", "4096")
	t.Setenv("ANALYTICS_MAX_PROPERTIES_BYTES", "1024")
	t.Setenv("ANALYTICS_MAX_SINGLE_PAGE_DURATION", "2h")
	t.Setenv("ANALYTICS_RATE_LIMIT_ENABLED", "false")
	t.Setenv("ANALYTICS_RATE_LIMIT_PER_MINUTE", "12")
	t.Setenv("ANALYTICS_RATE_LIMIT_BURST", "24")
	t.Setenv("AUTH_RATE_LIMIT_ENABLED", "false")
	t.Setenv("AUTH_RATE_LIMIT_ATTEMPTS", "3")
	t.Setenv("AUTH_RATE_LIMIT_WINDOW", "30m")
	t.Setenv("TRUSTED_PROXY_CIDRS", "203.0.113.0/24, 2001:db8::/32")
	t.Setenv("GRAPHQL_MAX_BODY_BYTES", "8192")
	t.Setenv("DASHBOARD_MAX_DAILY_RANGE_DAYS", "90")
	t.Setenv("DASHBOARD_MAX_HOURLY_RANGE_DAYS", "7")
	t.Setenv("DASHBOARD_MAX_FILTER_VALUES", "25")
	t.Setenv("DASHBOARD_MAX_FILTER_STRING_LENGTH", "128")

	cfg := Load()

	require.Equal(t, int64(4096), cfg.Analytics.MaxBodyBytes)
	require.Equal(t, 1024, cfg.Analytics.MaxPropertiesBytes)
	require.Equal(t, 2*time.Hour, cfg.Analytics.MaxSinglePageDuration)
	require.False(t, cfg.Analytics.RateLimitEnabled)
	require.Equal(t, 12, cfg.Analytics.RateLimitPerMinute)
	require.Equal(t, 24, cfg.Analytics.RateLimitBurst)
	require.False(t, cfg.Auth.RateLimitEnabled)
	require.Equal(t, 3, cfg.Auth.RateLimitAttempts)
	require.Equal(t, 30*time.Minute, cfg.Auth.RateLimitWindow)
	require.Equal(t, []string{"203.0.113.0/24", "2001:db8::/32"}, cfg.Analytics.TrustedProxyCIDRs)
	require.Equal(t, int64(8192), cfg.GraphQL.MaxBodyBytes)
	require.Equal(t, 90, cfg.Dashboard.MaxDailyRangeDays)
	require.Equal(t, 7, cfg.Dashboard.MaxHourlyRangeDays)
	require.Equal(t, 25, cfg.Dashboard.MaxFilterValues)
	require.Equal(t, 128, cfg.Dashboard.MaxFilterStringLength)
}

func ptr[T any](value T) *T {
	return &value
}
