package services

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/lovely-eye/server/internal/models"
	"github.com/lovely-eye/server/internal/repository"
	"github.com/stretchr/testify/require"
)

var testAnalyticsIdentitySecret = strings.Repeat("a", 32)

func TestTruncateVisitorIPPrefix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "ipv4",
			input: "203.0.113.42",
			want:  "203.0.113.0/24",
		},
		{
			name:  "ipv6",
			input: "2001:db8:abcd:1234:1111:2222:3333:4444",
			want:  "2001:db8:abcd:1234::/64",
		},
		{
			name:  "ipv4 mapped ipv6",
			input: "::ffff:203.0.113.42",
			want:  "203.0.113.0/24",
		},
		{
			name:  "invalid",
			input: "not-an-ip",
			want:  unknownVisitorIPPrefix,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, truncateVisitorIPPrefix(tt.input))
		})
	}
}

func TestAnalyticsService_GenerateVisitorID(t *testing.T) {
	t.Parallel()

	service := NewAnalyticsService(nil, nil, nil, nil, nil, testAnalyticsIdentitySecret)
	now := time.Date(2026, 3, 9, 10, 30, 0, 0, time.UTC)

	base := service.generateVisitorID(42, "203.0.113.42", models.ClientBrowserChrome, models.ClientDeviceDesktop, now)

	require.Equal(t, base, service.generateVisitorID(42, "203.0.113.200", models.ClientBrowserChrome, models.ClientDeviceDesktop, now))
	require.NotEqual(t, base, service.generateVisitorID(42, "203.0.114.42", models.ClientBrowserChrome, models.ClientDeviceDesktop, now))
	require.NotEqual(t, base, service.generateVisitorID(42, "203.0.113.42", models.ClientBrowserSafari, models.ClientDeviceDesktop, now))
	require.NotEqual(t, base, service.generateVisitorID(42, "203.0.113.42", models.ClientBrowserChrome, models.ClientDeviceMobile, now))
	require.NotEqual(t, base, service.generateVisitorID(84, "203.0.113.42", models.ClientBrowserChrome, models.ClientDeviceDesktop, now))
	require.NotEqual(t, base, service.generateVisitorID(42, "203.0.113.42", models.ClientBrowserChrome, models.ClientDeviceDesktop, now.Add(24*time.Hour)))
}

func TestAnalyticsService_CollectPageView_CountryTrackingDoesNotChangeIdentity(t *testing.T) {
	t.Parallel()

	db := setupAnalyticsServiceTestDB(t)
	ctx := context.Background()

	user := &models.User{
		Username:     "identity-user",
		PasswordHash: "hash",
		Role:         "admin",
	}
	_, err := db.NewInsert().Model(user).Exec(ctx)
	require.NoError(t, err)

	site := &models.Site{
		UserID:       user.ID,
		Name:         "Identity Site",
		PublicKey:    "identity-site-key",
		TrackCountry: false,
	}
	_, err = db.NewInsert().Model(site).Exec(ctx)
	require.NoError(t, err)

	domain := &models.SiteDomain{
		SiteID:   site.ID,
		Domain:   "identity.test",
		Position: 0,
	}
	_, err = db.NewInsert().Model(domain).Exec(ctx)
	require.NoError(t, err)

	geoIP := &fakeGeoIPProvider{
		resolvedCountry: Country{
			ISOCode: "US",
			Name:    "United States",
		},
	}
	service := NewAnalyticsService(
		repository.NewAnalyticsRepository(db),
		repository.NewSiteRepository(db),
		nil,
		geoIP,
		nil,
		testAnalyticsIdentitySecret,
	)

	input := CollectInput{
		SiteKey:     site.PublicKey,
		Path:        "/home",
		ScreenWidth: 1440,
		UserAgent:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0",
		IP:          "203.0.113.42",
		Origin:      "https://identity.test",
	}

	err = service.CollectPageView(ctx, input)
	require.NoError(t, err)

	_, err = db.NewUpdate().
		Model((*models.Site)(nil)).
		Set("track_country = ?", true).
		Where("id = ?", site.ID).
		Exec(ctx)
	require.NoError(t, err)

	input.Path = "/pricing"
	err = service.CollectPageView(ctx, input)
	require.NoError(t, err)

	clientCount, err := db.NewSelect().
		Model((*models.Client)(nil)).
		Where("site_id = ?", site.ID).
		Count(ctx)
	require.NoError(t, err)
	require.Equal(t, 1, clientCount)
}
