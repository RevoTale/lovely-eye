package analytics

import (
	"context"
	"strings"
	"testing"
	"time"

	analyticspersistence "github.com/lovely-eye/server/internal/analytics/persistence"
	authpersistence "github.com/lovely-eye/server/internal/auth/persistence"
	sitepersistence "github.com/lovely-eye/server/internal/site/persistence"
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

func TestService_GenerateVisitorID(t *testing.T) {
	t.Parallel()

	service := NewService(nil, nil, nil, nil, nil, testAnalyticsIdentitySecret)
	now := time.Date(2026, 3, 9, 10, 30, 0, 0, time.UTC)

	base := service.generateVisitorID(42, "203.0.113.42", analyticspersistence.ClientBrowserChrome, analyticspersistence.ClientDeviceDesktop, now)

	require.Equal(t, base, service.generateVisitorID(42, "203.0.113.200", analyticspersistence.ClientBrowserChrome, analyticspersistence.ClientDeviceDesktop, now))
	require.NotEqual(t, base, service.generateVisitorID(42, "203.0.114.42", analyticspersistence.ClientBrowserChrome, analyticspersistence.ClientDeviceDesktop, now))
	require.NotEqual(t, base, service.generateVisitorID(42, "203.0.113.42", analyticspersistence.ClientBrowserSafari, analyticspersistence.ClientDeviceDesktop, now))
	require.NotEqual(t, base, service.generateVisitorID(42, "203.0.113.42", analyticspersistence.ClientBrowserChrome, analyticspersistence.ClientDeviceMobile, now))
	require.NotEqual(t, base, service.generateVisitorID(84, "203.0.113.42", analyticspersistence.ClientBrowserChrome, analyticspersistence.ClientDeviceDesktop, now))
	require.NotEqual(t, base, service.generateVisitorID(42, "203.0.113.42", analyticspersistence.ClientBrowserChrome, analyticspersistence.ClientDeviceDesktop, now.Add(24*time.Hour)))
}

func TestService_CollectPageView_ReusesClientWithinSameUTCDay(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := setupServiceTestDB(t)
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
	expectedHash := service.generateVisitorID(site.ID, input.IP, analyticspersistence.ClientBrowserChrome, analyticspersistence.ClientDeviceDesktop, currentTime)
	require.Equal(t, expectedHash, client.Hash)

	session := latestSessionBySite(t, db, site.ID)
	require.Equal(t, 2, session.PageViewCount)
}

func TestService_CollectPageView_ReusesYesterdayClientAcrossMidnight(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := setupServiceTestDB(t)
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
	expectedHash := service.generateVisitorID(site.ID, input.IP, analyticspersistence.ClientBrowserChrome, analyticspersistence.ClientDeviceDesktop, currentTime)
	require.Equal(t, expectedHash, client.Hash)

	session := latestSessionBySite(t, db, site.ID)
	require.Equal(t, 2, session.PageViewCount)
	require.Equal(t, time.Date(2026, 3, 9, 23, 59, 50, 0, time.UTC).Unix(), session.EnterTime)
	require.Equal(t, currentTime.Unix(), session.ExitTime)
}

func TestService_CollectPageView_CreatesNewClientAfterSkippingUTCDay(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := setupServiceTestDB(t)
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

func TestService_CollectPageView_PrefersTodayHashOverYesterdayHash(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := setupServiceTestDB(t)
	site := createAnalyticsIdentitySite(t, db)
	service := newAnalyticsIdentityTestService(db, nil)

	now := time.Date(2026, 3, 10, 8, 0, 0, 0, time.UTC)
	service.now = func() time.Time { return now }

	hashes := service.buildClientRotationHashes(site.ID, "203.0.113.42", analyticspersistence.ClientBrowserChrome, analyticspersistence.ClientDeviceDesktop, now)
	yesterdayClient := &analyticspersistence.Client{
		SiteID:     site.ID,
		Hash:       hashes.Yesterday,
		Device:     analyticspersistence.ClientDeviceDesktop,
		Browser:    analyticspersistence.ClientBrowserChrome,
		OS:         analyticspersistence.ClientOSWindows,
		ScreenSize: analyticspersistence.ClientScreenSizeXL,
	}
	todayClient := &analyticspersistence.Client{
		SiteID:     site.ID,
		Hash:       hashes.Today,
		Device:     analyticspersistence.ClientDeviceDesktop,
		Browser:    analyticspersistence.ClientBrowserChrome,
		OS:         analyticspersistence.ClientOSWindows,
		ScreenSize: analyticspersistence.ClientScreenSizeXL,
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

func TestService_CollectPageView_DuplicatePageViewDoesNotMutateSession(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := setupServiceTestDB(t)
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

func TestService_CollectPageView_ExitSamePathUpdatesSessionWithoutPageView(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := setupServiceTestDB(t)
	site := createAnalyticsIdentitySite(t, db)
	service := newAnalyticsIdentityTestService(db, nil)

	firstTime := time.Date(2026, 3, 9, 10, 0, 0, 0, time.UTC)
	currentTime := firstTime
	service.now = func() time.Time { return currentTime }

	input := analyticsIdentityCollectInput(site.PublicKey)
	require.NoError(t, service.CollectPageView(ctx, input))

	currentTime = firstTime.Add(60 * time.Second)
	input.Exit = true
	require.NoError(t, service.CollectPageView(ctx, input))

	require.Equal(t, 1, countClientsBySite(t, db, site.ID))
	require.Equal(t, 1, countSessionsBySite(t, db, site.ID))
	require.Equal(t, 1, countPageViewEventsBySite(t, db, site.ID))

	session := latestSessionBySite(t, db, site.ID)
	require.Equal(t, 1, session.PageViewCount)
	require.Equal(t, "/home", session.ExitPath)
	require.Equal(t, currentTime.Unix(), session.ExitTime)
	require.Equal(t, 60, session.Duration)
}

func TestService_CollectPageView_ExitSamePathAfterActiveWindowUpdatesSessionWithoutPageView(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := setupServiceTestDB(t)
	site := createAnalyticsIdentitySite(t, db)
	service := newAnalyticsIdentityTestService(db, nil)

	firstTime := time.Date(2026, 3, 9, 10, 0, 0, 0, time.UTC)
	currentTime := firstTime
	service.now = func() time.Time { return currentTime }

	input := analyticsIdentityCollectInput(site.PublicKey)
	require.NoError(t, service.CollectPageView(ctx, input))

	currentTime = firstTime.Add(2 * time.Hour)
	input.Exit = true
	require.NoError(t, service.CollectPageView(ctx, input))

	require.Equal(t, 1, countClientsBySite(t, db, site.ID))
	require.Equal(t, 1, countSessionsBySite(t, db, site.ID))
	require.Equal(t, 1, countPageViewEventsBySite(t, db, site.ID))

	session := latestSessionBySite(t, db, site.ID)
	require.Equal(t, 1, session.PageViewCount)
	require.Equal(t, "/home", session.ExitPath)
	require.Equal(t, currentTime.Unix(), session.ExitTime)
	require.Equal(t, int((2 * time.Hour).Seconds()), session.Duration)
}

func TestService_CollectPageView_ExitSamePathAfterMaxDurationNoOps(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := setupServiceTestDB(t)
	site := createAnalyticsIdentitySite(t, db)
	service := newAnalyticsIdentityTestService(db, nil)

	firstTime := time.Date(2026, 3, 9, 10, 0, 0, 0, time.UTC)
	currentTime := firstTime
	service.now = func() time.Time { return currentTime }

	input := analyticsIdentityCollectInput(site.PublicKey)
	require.NoError(t, service.CollectPageView(ctx, input))

	currentTime = firstTime.Add(4*time.Hour + time.Minute)
	input.Exit = true
	require.NoError(t, service.CollectPageView(ctx, input))

	require.Equal(t, 1, countClientsBySite(t, db, site.ID))
	require.Equal(t, 1, countSessionsBySite(t, db, site.ID))
	require.Equal(t, 1, countPageViewEventsBySite(t, db, site.ID))

	session := latestSessionBySite(t, db, site.ID)
	require.Equal(t, 1, session.PageViewCount)
	require.Equal(t, "/home", session.ExitPath)
	require.Equal(t, firstTime.Unix(), session.ExitTime)
	require.Equal(t, 0, session.Duration)
}

func TestService_CollectPageView_RepeatedSamePathExitCapsSinglePageDuration(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := setupServiceTestDB(t)
	site := createAnalyticsIdentitySite(t, db)
	service := newAnalyticsIdentityTestService(db, nil)

	firstTime := time.Date(2026, 3, 9, 10, 0, 0, 0, time.UTC)
	currentTime := firstTime
	service.now = func() time.Time { return currentTime }

	input := analyticsIdentityCollectInput(site.PublicKey)
	require.NoError(t, service.CollectPageView(ctx, input))

	currentTime = firstTime.Add(3*time.Hour + 59*time.Minute)
	input.Exit = true
	require.NoError(t, service.CollectPageView(ctx, input))

	currentTime = firstTime.Add(7*time.Hour + 58*time.Minute)
	require.NoError(t, service.CollectPageView(ctx, input))

	require.Equal(t, 1, countClientsBySite(t, db, site.ID))
	require.Equal(t, 1, countSessionsBySite(t, db, site.ID))
	require.Equal(t, 1, countPageViewEventsBySite(t, db, site.ID))

	session := latestSessionBySite(t, db, site.ID)
	require.Equal(t, 1, session.PageViewCount)
	require.Equal(t, "/home", session.ExitPath)
	require.Equal(t, firstTime.Add(4*time.Hour).Unix(), session.ExitTime)
	require.Equal(t, int((4 * time.Hour).Seconds()), session.Duration)
}

func TestService_CollectPageView_ExitDifferentPathCountsPageView(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := setupServiceTestDB(t)
	site := createAnalyticsIdentitySite(t, db)
	service := newAnalyticsIdentityTestService(db, nil)

	firstTime := time.Date(2026, 3, 9, 10, 0, 0, 0, time.UTC)
	currentTime := firstTime
	service.now = func() time.Time { return currentTime }

	input := analyticsIdentityCollectInput(site.PublicKey)
	require.NoError(t, service.CollectPageView(ctx, input))

	currentTime = firstTime.Add(60 * time.Second)
	input.Path = "/pricing"
	input.Exit = true
	require.NoError(t, service.CollectPageView(ctx, input))

	require.Equal(t, 1, countClientsBySite(t, db, site.ID))
	require.Equal(t, 1, countSessionsBySite(t, db, site.ID))
	require.Equal(t, 2, countPageViewEventsBySite(t, db, site.ID))

	session := latestSessionBySite(t, db, site.ID)
	require.Equal(t, 2, session.PageViewCount)
	require.Equal(t, "/pricing", session.ExitPath)
	require.Equal(t, currentTime.Unix(), session.ExitTime)
	require.Equal(t, 60, session.Duration)
}

func TestService_CollectPageView_ExitDifferentPathAfterActiveWindowNoOps(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := setupServiceTestDB(t)
	site := createAnalyticsIdentitySite(t, db)
	service := newAnalyticsIdentityTestService(db, nil)

	firstTime := time.Date(2026, 3, 9, 10, 0, 0, 0, time.UTC)
	currentTime := firstTime
	service.now = func() time.Time { return currentTime }

	input := analyticsIdentityCollectInput(site.PublicKey)
	require.NoError(t, service.CollectPageView(ctx, input))

	currentTime = firstTime.Add(2 * time.Hour)
	input.Path = "/pricing"
	input.Exit = true
	require.NoError(t, service.CollectPageView(ctx, input))

	require.Equal(t, 1, countClientsBySite(t, db, site.ID))
	require.Equal(t, 1, countSessionsBySite(t, db, site.ID))
	require.Equal(t, 1, countPageViewEventsBySite(t, db, site.ID))

	session := latestSessionBySite(t, db, site.ID)
	require.Equal(t, 1, session.PageViewCount)
	require.Equal(t, "/home", session.ExitPath)
	require.Equal(t, firstTime.Unix(), session.ExitTime)
	require.Equal(t, 0, session.Duration)
}

func TestService_CollectPageView_ExitWithoutActiveSessionNoOps(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := setupServiceTestDB(t)
	site := createAnalyticsIdentitySite(t, db)
	service := newAnalyticsIdentityTestService(db, nil)

	currentTime := time.Date(2026, 3, 9, 10, 0, 0, 0, time.UTC)
	service.now = func() time.Time { return currentTime }

	input := analyticsIdentityCollectInput(site.PublicKey)
	input.Exit = true
	require.NoError(t, service.CollectPageView(ctx, input))

	require.Equal(t, 0, countClientsBySite(t, db, site.ID))
	require.Equal(t, 0, countSessionsBySite(t, db, site.ID))
	require.Equal(t, 0, countPageViewEventsBySite(t, db, site.ID))
}

func TestService_CollectPageView_ExitWithoutActiveSessionDoesNotRotateClient(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := setupServiceTestDB(t)
	site := createAnalyticsIdentitySite(t, db)
	service := newAnalyticsIdentityTestService(db, nil)

	firstTime := time.Date(2026, 3, 9, 10, 0, 0, 0, time.UTC)
	currentTime := firstTime
	service.now = func() time.Time { return currentTime }

	input := analyticsIdentityCollectInput(site.PublicKey)
	require.NoError(t, service.CollectPageView(ctx, input))

	originalClient := latestClientBySite(t, db, site.ID)
	originalHash := originalClient.Hash

	currentTime = firstTime.Add(24*time.Hour + time.Minute)
	input.Exit = true
	require.NoError(t, service.CollectPageView(ctx, input))

	require.Equal(t, 1, countClientsBySite(t, db, site.ID))
	require.Equal(t, 1, countSessionsBySite(t, db, site.ID))
	require.Equal(t, 1, countPageViewEventsBySite(t, db, site.ID))

	persistedClient := clientByID(t, db, originalClient.ID)
	require.Equal(t, originalHash, persistedClient.Hash)

	session := latestSessionBySite(t, db, site.ID)
	require.Equal(t, firstTime.Unix(), session.ExitTime)
	require.Equal(t, 0, session.Duration)
}

func TestService_CollectPageView_CreatesNewSessionAfterThirtyMinutes(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := setupServiceTestDB(t)
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

func TestService_CollectPageView_CountryTrackingDoesNotChangeIdentity(t *testing.T) {
	t.Parallel()

	db := setupServiceTestDB(t)
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
		Model((*sitepersistence.Site)(nil)).
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
		SiteKey:   siteKey,
		Path:      "/home",
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0",
		IP:        "203.0.113.42",
		Origin:    "https://identity.test",
	}
}

func newAnalyticsIdentityTestService(db *bun.DB, geoIP geoIPProvider) *Service {
	return NewService(
		analyticspersistence.New(db),
		sitepersistence.New(db),
		nil,
		geoIP,
		nil,
		testAnalyticsIdentitySecret,
	)
}

func createAnalyticsIdentitySite(t *testing.T, db *bun.DB) *sitepersistence.Site {
	t.Helper()

	ctx := context.Background()
	user := &authpersistence.User{
		Username:     "identity-user",
		PasswordHash: "hash",
		Role:         "admin",
	}
	_, err := db.NewInsert().Model(user).Exec(ctx)
	require.NoError(t, err)

	site := &sitepersistence.Site{
		UserID:       user.ID,
		Name:         "Identity Site",
		PublicKey:    "identity-site-key",
		TrackCountry: false,
	}
	_, err = db.NewInsert().Model(site).Exec(ctx)
	require.NoError(t, err)

	domain := &sitepersistence.Domain{
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
		Model((*analyticspersistence.Client)(nil)).
		Where("site_id = ?", siteID).
		Count(context.Background())
	require.NoError(t, err)
	return count
}

func countSessionsBySite(t *testing.T, db *bun.DB, siteID int64) int {
	t.Helper()

	count, err := db.NewSelect().
		Model((*analyticspersistence.Session)(nil)).
		Where("site_id = ?", siteID).
		Count(context.Background())
	require.NoError(t, err)
	return count
}

func countPageViewEventsBySite(t *testing.T, db *bun.DB, siteID int64) int {
	t.Helper()

	count, err := db.NewSelect().
		Model((*analyticspersistence.Event)(nil)).
		Join("INNER JOIN sessions s ON e.session_id = s.id").
		Where("s.site_id = ?", siteID).
		Where("e.definition_id IS NULL").
		Count(context.Background())
	require.NoError(t, err)
	return count
}

func latestClientBySite(t *testing.T, db *bun.DB, siteID int64) *analyticspersistence.Client {
	t.Helper()

	client := new(analyticspersistence.Client)
	err := db.NewSelect().
		Model(client).
		Where("site_id = ?", siteID).
		Order("id DESC").
		Limit(1).
		Scan(context.Background())
	require.NoError(t, err)
	return client
}

func clientByID(t *testing.T, db *bun.DB, clientID int64) *analyticspersistence.Client {
	t.Helper()

	client := new(analyticspersistence.Client)
	err := db.NewSelect().
		Model(client).
		Where("id = ?", clientID).
		Scan(context.Background())
	require.NoError(t, err)
	return client
}

func latestSessionBySite(t *testing.T, db *bun.DB, siteID int64) *analyticspersistence.Session {
	t.Helper()

	session := new(analyticspersistence.Session)
	err := db.NewSelect().
		Model(session).
		Where("site_id = ?", siteID).
		Order("id DESC").
		Limit(1).
		Scan(context.Background())
	require.NoError(t, err)
	return session
}
