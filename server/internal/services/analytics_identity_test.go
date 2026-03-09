package services

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/lovely-eye/server/internal/models"
	"github.com/lovely-eye/server/internal/repository"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
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

func TestAnalyticsService_CollectPageView_ReusesClientWithinSameUTCDay(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := setupAnalyticsServiceTestDB(t)
	site := createAnalyticsIdentitySite(t, db)
	service := newAnalyticsIdentityTestService(db, nil)

	currentTime := time.Date(2026, 3, 9, 10, 0, 0, 0, time.UTC)
	service.now = func() time.Time { return currentTime }

	input := analyticsIdentityCollectInput(site.PublicKey)
	require.NoError(t, service.CollectPageView(ctx, input))

	currentTime = currentTime.Add(10 * time.Minute)
	input.Path = "/pricing"
	require.NoError(t, service.CollectPageView(ctx, input))

	require.Equal(t, 1, countClientsBySite(t, db, site.ID))
	require.Equal(t, 1, countSessionsBySite(t, db, site.ID))

	client := latestClientBySite(t, db, site.ID)
	expectedHash := service.generateVisitorID(site.ID, input.IP, models.ClientBrowserChrome, models.ClientDeviceDesktop, currentTime)
	require.Equal(t, expectedHash, client.Hash)

	session := latestSessionBySite(t, db, site.ID)
	require.Equal(t, 2, session.PageViewCount)
}

func TestAnalyticsService_CollectPageView_ReusesYesterdayClientAcrossMidnight(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := setupAnalyticsServiceTestDB(t)
	site := createAnalyticsIdentitySite(t, db)
	service := newAnalyticsIdentityTestService(db, nil)

	currentTime := time.Date(2026, 3, 9, 23, 59, 50, 0, time.UTC)
	service.now = func() time.Time { return currentTime }

	input := analyticsIdentityCollectInput(site.PublicKey)
	require.NoError(t, service.CollectPageView(ctx, input))

	currentTime = time.Date(2026, 3, 10, 0, 5, 0, 0, time.UTC)
	input.Path = "/pricing"
	require.NoError(t, service.CollectPageView(ctx, input))

	require.Equal(t, 1, countClientsBySite(t, db, site.ID))
	require.Equal(t, 1, countSessionsBySite(t, db, site.ID))

	client := latestClientBySite(t, db, site.ID)
	expectedHash := service.generateVisitorID(site.ID, input.IP, models.ClientBrowserChrome, models.ClientDeviceDesktop, currentTime)
	require.Equal(t, expectedHash, client.Hash)

	session := latestSessionBySite(t, db, site.ID)
	require.Equal(t, 2, session.PageViewCount)
	require.Equal(t, time.Date(2026, 3, 9, 23, 59, 50, 0, time.UTC).Unix(), session.EnterTime)
	require.Equal(t, currentTime.Unix(), session.ExitTime)
}

func TestAnalyticsService_CollectPageView_CreatesNewClientAfterSkippingUTCDay(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := setupAnalyticsServiceTestDB(t)
	site := createAnalyticsIdentitySite(t, db)
	service := newAnalyticsIdentityTestService(db, nil)

	currentTime := time.Date(2026, 3, 9, 10, 0, 0, 0, time.UTC)
	service.now = func() time.Time { return currentTime }

	input := analyticsIdentityCollectInput(site.PublicKey)
	require.NoError(t, service.CollectPageView(ctx, input))

	currentTime = time.Date(2026, 3, 11, 9, 0, 0, 0, time.UTC)
	input.Path = "/pricing"
	require.NoError(t, service.CollectPageView(ctx, input))

	require.Equal(t, 2, countClientsBySite(t, db, site.ID))
	require.Equal(t, 2, countSessionsBySite(t, db, site.ID))
}

func TestAnalyticsService_CollectPageView_PrefersTodayHashOverYesterdayHash(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := setupAnalyticsServiceTestDB(t)
	site := createAnalyticsIdentitySite(t, db)
	service := newAnalyticsIdentityTestService(db, nil)

	now := time.Date(2026, 3, 10, 8, 0, 0, 0, time.UTC)
	service.now = func() time.Time { return now }

	hashes := service.buildClientRotationHashes(site.ID, "203.0.113.42", models.ClientBrowserChrome, models.ClientDeviceDesktop, now)
	yesterdayClient := &models.Client{
		SiteID:     site.ID,
		Hash:       hashes.Yesterday,
		Device:     models.ClientDeviceDesktop,
		Browser:    models.ClientBrowserChrome,
		OS:         models.ClientOSWindows,
		ScreenSize: models.ClientScreenSizeXL,
	}
	todayClient := &models.Client{
		SiteID:     site.ID,
		Hash:       hashes.Today,
		Device:     models.ClientDeviceDesktop,
		Browser:    models.ClientBrowserChrome,
		OS:         models.ClientOSWindows,
		ScreenSize: models.ClientScreenSizeXL,
	}
	_, err := db.NewInsert().Model(yesterdayClient).Exec(ctx)
	require.NoError(t, err)
	_, err = db.NewInsert().Model(todayClient).Exec(ctx)
	require.NoError(t, err)

	input := analyticsIdentityCollectInput(site.PublicKey)
	require.NoError(t, service.CollectPageView(ctx, input))

	session := latestSessionBySite(t, db, site.ID)
	require.Equal(t, todayClient.ID, session.ClientID)

	persistedYesterday := clientByID(t, db, yesterdayClient.ID)
	require.Equal(t, hashes.Yesterday, persistedYesterday.Hash)
}

func TestAnalyticsService_CollectPageView_DuplicatePageViewDoesNotMutateSession(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := setupAnalyticsServiceTestDB(t)
	site := createAnalyticsIdentitySite(t, db)
	service := newAnalyticsIdentityTestService(db, nil)

	firstTime := time.Date(2026, 3, 9, 10, 0, 0, 0, time.UTC)
	currentTime := firstTime
	service.now = func() time.Time { return currentTime }

	input := analyticsIdentityCollectInput(site.PublicKey)
	require.NoError(t, service.CollectPageView(ctx, input))

	currentTime = firstTime.Add(5 * time.Second)
	require.NoError(t, service.CollectPageView(ctx, input))

	require.Equal(t, 1, countClientsBySite(t, db, site.ID))
	require.Equal(t, 1, countSessionsBySite(t, db, site.ID))
	require.Equal(t, 1, countPageViewEventsBySite(t, db, site.ID))

	session := latestSessionBySite(t, db, site.ID)
	require.Equal(t, 1, session.PageViewCount)
	require.Equal(t, firstTime.Unix(), session.ExitTime)
	require.Equal(t, 0, session.Duration)
}

func TestAnalyticsService_CollectPageView_CreatesNewSessionAfterThirtyMinutes(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := setupAnalyticsServiceTestDB(t)
	site := createAnalyticsIdentitySite(t, db)
	service := newAnalyticsIdentityTestService(db, nil)

	currentTime := time.Date(2026, 3, 9, 10, 0, 0, 0, time.UTC)
	service.now = func() time.Time { return currentTime }

	input := analyticsIdentityCollectInput(site.PublicKey)
	require.NoError(t, service.CollectPageView(ctx, input))

	currentTime = currentTime.Add(31 * time.Minute)
	input.Path = "/pricing"
	require.NoError(t, service.CollectPageView(ctx, input))

	require.Equal(t, 1, countClientsBySite(t, db, site.ID))
	require.Equal(t, 2, countSessionsBySite(t, db, site.ID))
}

func TestAnalyticsService_CollectPageView_CountryTrackingDoesNotChangeIdentity(t *testing.T) {
	t.Parallel()

	db := setupAnalyticsServiceTestDB(t)
	ctx := context.Background()

	site := createAnalyticsIdentitySite(t, db)
	geoIP := &fakeGeoIPProvider{
		resolvedCountry: Country{
			ISOCode: "US",
			Name:    "United States",
		},
	}
	service := newAnalyticsIdentityTestService(db, geoIP)

	input := analyticsIdentityCollectInput(site.PublicKey)
	require.NoError(t, service.CollectPageView(ctx, input))

	_, err := db.NewUpdate().
		Model((*models.Site)(nil)).
		Set("track_country = ?", true).
		Where("id = ?", site.ID).
		Exec(ctx)
	require.NoError(t, err)

	input.Path = "/pricing"
	require.NoError(t, service.CollectPageView(ctx, input))

	require.Equal(t, 1, countClientsBySite(t, db, site.ID))

	client := latestClientBySite(t, db, site.ID)
	require.Equal(t, "US", client.Country)
}

func analyticsIdentityCollectInput(siteKey string) CollectInput {
	return CollectInput{
		SiteKey:     siteKey,
		Path:        "/home",
		ScreenWidth: 1440,
		UserAgent:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0",
		IP:          "203.0.113.42",
		Origin:      "https://identity.test",
	}
}

func newAnalyticsIdentityTestService(db *bun.DB, geoIP geoIPProvider) *AnalyticsService {
	return NewAnalyticsService(
		repository.NewAnalyticsRepository(db),
		repository.NewSiteRepository(db),
		nil,
		geoIP,
		nil,
		testAnalyticsIdentitySecret,
	)
}

func createAnalyticsIdentitySite(t *testing.T, db *bun.DB) *models.Site {
	t.Helper()

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

	return site
}

func countClientsBySite(t *testing.T, db *bun.DB, siteID int64) int {
	t.Helper()

	count, err := db.NewSelect().
		Model((*models.Client)(nil)).
		Where("site_id = ?", siteID).
		Count(context.Background())
	require.NoError(t, err)
	return count
}

func countSessionsBySite(t *testing.T, db *bun.DB, siteID int64) int {
	t.Helper()

	count, err := db.NewSelect().
		Model((*models.Session)(nil)).
		Where("site_id = ?", siteID).
		Count(context.Background())
	require.NoError(t, err)
	return count
}

func countPageViewEventsBySite(t *testing.T, db *bun.DB, siteID int64) int {
	t.Helper()

	count, err := db.NewSelect().
		Model((*models.Event)(nil)).
		Join("INNER JOIN sessions s ON e.session_id = s.id").
		Where("s.site_id = ?", siteID).
		Where("e.definition_id IS NULL").
		Count(context.Background())
	require.NoError(t, err)
	return count
}

func latestClientBySite(t *testing.T, db *bun.DB, siteID int64) *models.Client {
	t.Helper()

	client := new(models.Client)
	err := db.NewSelect().
		Model(client).
		Where("site_id = ?", siteID).
		Order("id DESC").
		Limit(1).
		Scan(context.Background())
	require.NoError(t, err)
	return client
}

func clientByID(t *testing.T, db *bun.DB, clientID int64) *models.Client {
	t.Helper()

	client := new(models.Client)
	err := db.NewSelect().
		Model(client).
		Where("id = ?", clientID).
		Scan(context.Background())
	require.NoError(t, err)
	return client
}

func latestSessionBySite(t *testing.T, db *bun.DB, siteID int64) *models.Session {
	t.Helper()

	session := new(models.Session)
	err := db.NewSelect().
		Model(session).
		Where("site_id = ?", siteID).
		Order("id DESC").
		Limit(1).
		Scan(context.Background())
	require.NoError(t, err)
	return session
}
