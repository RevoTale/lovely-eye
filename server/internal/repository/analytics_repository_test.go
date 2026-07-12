package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/lovely-eye/server/internal/database"
	"github.com/lovely-eye/server/internal/models"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *bun.DB {
	t.Helper()

	sqldb, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	db := bun.NewDB(sqldb, sqlitedialect.New())

	ctx := context.Background()
	if err := database.Migrate(ctx, db); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	})

	return db
}

func createTestSite(t *testing.T, db *bun.DB) *models.Site {
	t.Helper()

	ctx := context.Background()

	user := &models.User{
		Username:     "testuser",
		PasswordHash: "hash",
		Role:         "admin",
	}
	if _, err := db.NewInsert().Model(user).Exec(ctx); err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	site := &models.Site{
		UserID:    user.ID,
		Name:      "Test Site",
		PublicKey: "test-key",
	}
	if _, err := db.NewInsert().Model(site).Exec(ctx); err != nil {
		t.Fatalf("failed to insert site: %v", err)
	}

	siteDomain := &models.SiteDomain{
		SiteID:   site.ID,
		Domain:   "test.com",
		Position: 0,
	}
	if _, err := db.NewInsert().Model(siteDomain).Exec(ctx); err != nil {
		t.Fatalf("failed to insert site domain: %v", err)
	}

	return site
}

func createTestClient(t *testing.T, db *bun.DB, siteID int64, hash string, device string, browser string, os string) int64 {
	t.Helper()

	ctx := context.Background()
	client := &models.Client{
		SiteID:  siteID,
		Hash:    hash,
		Device:  models.ClientDeviceFromLegacyLabel(device),
		Browser: models.ClientBrowserFromLegacyLabel(browser),
		OS:      models.ClientOSFromLegacyLabel(os),
	}
	if _, err := db.NewInsert().Model(client).Exec(ctx); err != nil {
		t.Fatalf("failed to insert client: %v", err)
	}

	return client.ID
}

func insertSession(t *testing.T, db *bun.DB, siteID int64, clientID int64, enterTime time.Time, durationSeconds int, pageViewCount int) {
	t.Helper()

	enterUnix := enterTime.Unix()
	exitUnix := enterUnix + int64(durationSeconds)

	_, err := db.Exec(
		"INSERT INTO sessions (site_id, client_id, enter_time, enter_hour, enter_day, enter_path, exit_time, exit_hour, exit_day, exit_path, duration, page_view_count) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		siteID, clientID, enterUnix, enterUnix/3600, enterUnix/86400, "/", exitUnix, exitUnix/3600, exitUnix/86400, "/", durationSeconds, pageViewCount,
	)
	if err != nil {
		t.Fatalf("failed to insert session: %v", err)
	}
}

func insertSessionWithReferrer(t *testing.T, db *bun.DB, siteID int64, clientID int64, referrer string, enterTime time.Time, durationSeconds int, pageViewCount int) {
	t.Helper()

	enterUnix := enterTime.Unix()
	exitUnix := enterUnix + int64(durationSeconds)
	session := &models.Session{
		SiteID:        siteID,
		ClientID:      clientID,
		EnterTime:     enterUnix,
		EnterHour:     enterUnix / 3600,
		EnterDay:      enterUnix / 86400,
		EnterPath:     "/",
		ExitTime:      exitUnix,
		ExitHour:      exitUnix / 3600,
		ExitDay:       exitUnix / 86400,
		ExitPath:      "/",
		Referrer:      referrer,
		Duration:      durationSeconds,
		PageViewCount: pageViewCount,
	}

	if _, err := db.NewInsert().Model(session).Exec(context.Background()); err != nil {
		t.Fatalf("failed to insert session with referrer: %v", err)
	}
}

func insertSessionWithPath(t *testing.T, db *bun.DB, siteID int64, clientID int64, path string, enterTime time.Time, durationSeconds int, pageViewCount int) int64 {
	t.Helper()

	enterUnix := enterTime.Unix()
	exitUnix := enterUnix + int64(durationSeconds)
	session := &models.Session{
		SiteID:        siteID,
		ClientID:      clientID,
		EnterTime:     enterUnix,
		EnterHour:     enterUnix / 3600,
		EnterDay:      enterUnix / 86400,
		EnterPath:     path,
		ExitTime:      exitUnix,
		ExitHour:      exitUnix / 3600,
		ExitDay:       exitUnix / 86400,
		ExitPath:      path,
		Duration:      durationSeconds,
		PageViewCount: pageViewCount,
	}

	if _, err := db.NewInsert().Model(session).Exec(context.Background()); err != nil {
		t.Fatalf("failed to insert session: %v", err)
	}
	return session.ID
}

func insertPageViewEvent(t *testing.T, db *bun.DB, sessionID int64, path string, eventTime time.Time) int64 {
	t.Helper()

	eventUnix := eventTime.Unix()
	event := &models.Event{
		SessionID:    sessionID,
		Time:         eventUnix,
		Hour:         eventUnix / 3600,
		Day:          eventUnix / 86400,
		Path:         path,
		DefinitionID: nil,
	}

	if _, err := db.NewInsert().Model(event).Exec(context.Background()); err != nil {
		t.Fatalf("failed to insert pageview event: %v", err)
	}
	return event.ID
}

func TestGetBounceRateExcludesEventOnlySessions(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAnalyticsRepository(db)
	ctx := context.Background()

	site := createTestSite(t, db)
	now := time.Now()
	from := now.Add(-24 * time.Hour)
	to := now

	bouncedClient := createTestClient(t, db, site.ID, "hash-bounced", "desktop", "Chrome", "Windows")
	insertSession(t, db, site.ID, bouncedClient, now.Add(-2*time.Hour), 0, 1)

	eventOnlyClient := createTestClient(t, db, site.ID, "hash-event-only", "desktop", "Chrome", "Windows")
	insertSession(t, db, site.ID, eventOnlyClient, now.Add(-3*time.Hour), 60, 0)

	got, err := repo.GetBounceRate(ctx, site.ID, from, to)
	if err != nil {
		t.Fatalf("GetBounceRate() error = %v", err)
	}

	want := 100.0
	if got != want {
		t.Errorf("GetBounceRate() = %v, want %v", got, want)
	}
}

func TestGetBounceRateWithFilterExcludesEventOnlySessions(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAnalyticsRepository(db)
	ctx := context.Background()

	site := createTestSite(t, db)
	now := time.Now()
	from := now.Add(-24 * time.Hour)
	to := now

	bouncedClient := createTestClient(t, db, site.ID, "hash-bounced-filter", "desktop", "Chrome", "Windows")
	insertSessionWithReferrer(t, db, site.ID, bouncedClient, "https://google.com", now.Add(-2*time.Hour), 0, 1)

	eventOnlyClient := createTestClient(t, db, site.ID, "hash-event-only-filter", "desktop", "Chrome", "Windows")
	insertSessionWithReferrer(t, db, site.ID, eventOnlyClient, "https://google.com", now.Add(-3*time.Hour), 60, 0)

	otherClient := createTestClient(t, db, site.ID, "hash-other-filter", "desktop", "Chrome", "Windows")
	insertSessionWithReferrer(t, db, site.ID, otherClient, "https://example.com", now.Add(-4*time.Hour), 60, 2)

	got, err := repo.GetBounceRateWithFilter(ctx, AnalyticsQuery{
		SiteID: site.ID,
		From:   from,
		To:     to,
		Filter: AnalyticsFilter{Referrer: []string{"https://google.com"}},
	})
	if err != nil {
		t.Fatalf("GetBounceRateWithFilter() error = %v", err)
	}

	want := 100.0
	if got != want {
		t.Errorf("GetBounceRateWithFilter() = %v, want %v", got, want)
	}
}

// TestGetAvgSessionDuration_EmptyResult tests that the function returns 0.0
// when there are no matching sessions. This test catches a bug where
// COALESCE(AVG(duration), 0) returns int64 instead of float64 when there
// are no rows, causing a type mismatch error during scanning.
func TestGetAvgSessionDuration_EmptyResult(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAnalyticsRepository(db)
	ctx := context.Background()

	site := createTestSite(t, db)

	now := time.Now()
	from := now.Add(-24 * time.Hour)
	to := now

	tests := []struct {
		name  string
		setup func()
	}{
		{
			name:  "no sessions exist",
			setup: func() {},
		},
		{
			name: "only zero-duration sessions exist",
			setup: func() {
				clientID := createTestClient(t, db, site.ID, "hash1", "desktop", "Chrome", "Windows")
				sessionTime := now.Add(-1 * time.Hour)

				insertSession(t, db, site.ID, clientID, sessionTime, 0, 1)
			},
		},
		{
			name: "sessions outside date range",
			setup: func() {
				clientID := createTestClient(t, db, site.ID, "hash2", "desktop", "Chrome", "Windows")
				pastTime := now.Add(-30 * 24 * time.Hour)
				insertSession(t, db, site.ID, clientID, pastTime, 100, 2)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			got, err := repo.GetAvgSessionDuration(ctx, site.ID, from, to)
			if err != nil {
				t.Errorf("GetAvgSessionDuration() error = %v, want nil", err)
			}
			if got != 0.0 {
				t.Errorf("GetAvgSessionDuration() = %v, want 0.0", got)
			}
		})
	}
}

func TestGetAvgSessionDuration_WithSessions(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAnalyticsRepository(db)
	ctx := context.Background()

	site := createTestSite(t, db)

	now := time.Now()
	from := now.Add(-24 * time.Hour)
	to := now

	client1 := createTestClient(t, db, site.ID, "hash1", "desktop", "Chrome", "Windows")
	insertSession(t, db, site.ID, client1, now.Add(-2*time.Hour), 60, 2)

	client2 := createTestClient(t, db, site.ID, "hash2", "mobile", "Safari", "iOS")
	insertSession(t, db, site.ID, client2, now.Add(-3*time.Hour), 120, 3)

	client3 := createTestClient(t, db, site.ID, "hash3", "tablet", "Firefox", "Android")
	insertSession(t, db, site.ID, client3, now.Add(-4*time.Hour), 180, 2)

	got, err := repo.GetAvgSessionDuration(ctx, site.ID, from, to)
	if err != nil {
		t.Fatalf("GetAvgSessionDuration() error = %v", err)
	}

	want := 120.0
	if got != want {
		t.Errorf("GetAvgSessionDuration() = %v, want %v", got, want)
	}
}

func TestGetAvgSessionDuration_IncludesSinglePageDurationsAndExcludesZeroDuration(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAnalyticsRepository(db)
	ctx := context.Background()

	site := createTestSite(t, db)

	now := time.Now()
	from := now.Add(-24 * time.Hour)
	to := now

	client1 := createTestClient(t, db, site.ID, "hash1", "desktop", "Chrome", "Windows")
	insertSession(t, db, site.ID, client1, now.Add(-2*time.Hour), 100, 2)

	client2 := createTestClient(t, db, site.ID, "hash2", "mobile", "Safari", "iOS")
	insertSession(t, db, site.ID, client2, now.Add(-3*time.Hour), 1000, 1)

	client3 := createTestClient(t, db, site.ID, "hash3", "mobile", "Safari", "iOS")
	insertSession(t, db, site.ID, client3, now.Add(-4*time.Hour), 0, 1)

	client4 := createTestClient(t, db, site.ID, "hash4", "desktop", "Chrome", "Windows")
	insertSession(t, db, site.ID, client4, now.Add(-5*time.Hour), 400, 0)

	got, err := repo.GetAvgSessionDuration(ctx, site.ID, from, to)
	if err != nil {
		t.Fatalf("GetAvgSessionDuration() error = %v", err)
	}

	want := 550.0
	if got != want {
		t.Errorf("GetAvgSessionDuration() = %v, want %v", got, want)
	}
}

func TestGetAvgSessionDurationWithFilter_EmptyResult(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAnalyticsRepository(db)
	ctx := context.Background()

	site := createTestSite(t, db)

	now := time.Now()
	from := now.Add(-24 * time.Hour)
	to := now

	tests := []struct {
		name     string
		referrer *string
		device   *string
		page     *string
		setup    func()
	}{
		{
			name:  "no sessions exist",
			setup: func() {},
		},
		{
			name:     "filter matches no sessions",
			referrer: stringPtr("nonexistent.com"),
			setup: func() {

				clientID := createTestClient(t, db, site.ID, "hash1", "desktop", "Chrome", "Windows")
				enterTime := now.Add(-1 * time.Hour).Unix()
				exitTime := enterTime + 100
				_, err := db.Exec(
					"INSERT INTO sessions (site_id, client_id, enter_time, enter_hour, enter_day, enter_path, exit_time, exit_hour, exit_day, exit_path, duration, page_view_count, referrer) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
					site.ID, clientID, enterTime, enterTime/3600, enterTime/86400, "/", exitTime, exitTime/3600, exitTime/86400, "/", 100, 2, "google.com",
				)
				if err != nil {
					t.Fatalf("failed to insert session: %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			var referrers []string
			if tt.referrer != nil {
				referrers = []string{*tt.referrer}
			}
			var devices []string
			if tt.device != nil {
				devices = []string{*tt.device}
			}
			var pages []string
			if tt.page != nil {
				pages = []string{*tt.page}
			}

			got, err := repo.GetAvgSessionDurationWithFilter(ctx, AnalyticsQuery{
				SiteID: site.ID,
				From:   from,
				To:     to,
				Filter: AnalyticsFilter{
					Referrer: referrers,
					Device:   devices,
					Page:     pages,
				},
			})
			if err != nil {
				t.Errorf("GetAvgSessionDurationWithFilter() error = %v, want nil", err)
			}
			if got != 0.0 {
				t.Errorf("GetAvgSessionDurationWithFilter() = %v, want 0.0", got)
			}
		})
	}
}

func TestGetAvgSessionDurationWithFilter_WithData(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAnalyticsRepository(db)
	ctx := context.Background()

	site := createTestSite(t, db)

	now := time.Now()
	from := now.Add(-24 * time.Hour)
	to := now

	client1 := createTestClient(t, db, site.ID, "hash1", "desktop", "Chrome", "Windows")
	enterTime1 := now.Add(-2 * time.Hour).Unix()
	exitTime1 := enterTime1 + 60
	_, err := db.Exec(
		"INSERT INTO sessions (site_id, client_id, enter_time, enter_hour, enter_day, enter_path, exit_time, exit_hour, exit_day, exit_path, duration, page_view_count, referrer) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		site.ID, client1, enterTime1, enterTime1/3600, enterTime1/86400, "/", exitTime1, exitTime1/3600, exitTime1/86400, "/", 60, 2, "google.com",
	)
	if err != nil {
		t.Fatalf("failed to insert session 1: %v", err)
	}

	client2 := createTestClient(t, db, site.ID, "hash2", "mobile", "Safari", "iOS")
	enterTime2 := now.Add(-3 * time.Hour).Unix()
	exitTime2 := enterTime2 + 120
	_, err = db.Exec(
		"INSERT INTO sessions (site_id, client_id, enter_time, enter_hour, enter_day, enter_path, exit_time, exit_hour, exit_day, exit_path, duration, page_view_count, referrer) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		site.ID, client2, enterTime2, enterTime2/3600, enterTime2/86400, "/", exitTime2, exitTime2/3600, exitTime2/86400, "/", 120, 2, "",
	)
	if err != nil {
		t.Fatalf("failed to insert session 2: %v", err)
	}

	client3 := createTestClient(t, db, site.ID, "hash3", "tablet", "Firefox", "Android")
	enterTime3 := now.Add(-4 * time.Hour).Unix()
	exitTime3 := enterTime3 + 180
	_, err = db.Exec(
		"INSERT INTO sessions (site_id, client_id, enter_time, enter_hour, enter_day, enter_path, exit_time, exit_hour, exit_day, exit_path, duration, page_view_count, referrer) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		site.ID, client3, enterTime3, enterTime3/3600, enterTime3/86400, "/", exitTime3, exitTime3/3600, exitTime3/86400, "/", 180, 1, "",
	)
	if err != nil {
		t.Fatalf("failed to insert session 3: %v", err)
	}

	client4 := createTestClient(t, db, site.ID, "hash4", "tablet", "Firefox", "Android")
	enterTime4 := now.Add(-5 * time.Hour).Unix()
	_, err = db.Exec(
		"INSERT INTO sessions (site_id, client_id, enter_time, enter_hour, enter_day, enter_path, exit_time, exit_hour, exit_day, exit_path, duration, page_view_count, referrer) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		site.ID, client4, enterTime4, enterTime4/3600, enterTime4/86400, "/", enterTime4, enterTime4/3600, enterTime4/86400, "/", 0, 1, "",
	)
	if err != nil {
		t.Fatalf("failed to insert session 4: %v", err)
	}

	client5 := createTestClient(t, db, site.ID, "hash5", "desktop", "Chrome", "Windows")
	enterTime5 := now.Add(-6 * time.Hour).Unix()
	exitTime5 := enterTime5 + 600
	_, err = db.Exec(
		"INSERT INTO sessions (site_id, client_id, enter_time, enter_hour, enter_day, enter_path, exit_time, exit_hour, exit_day, exit_path, duration, page_view_count, referrer) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		site.ID, client5, enterTime5, enterTime5/3600, enterTime5/86400, "/", exitTime5, exitTime5/3600, exitTime5/86400, "/", 600, 0, "",
	)
	if err != nil {
		t.Fatalf("failed to insert session 5: %v", err)
	}

	tests := []struct {
		name     string
		referrer *string
		device   *string
		page     *string
		want     float64
	}{
		{
			name: "no filter",
			want: 120.0,
		},
		{
			name:     "filter by referrer",
			referrer: stringPtr("google.com"),
			want:     60.0,
		},
		{
			name:   "filter by device",
			device: stringPtr("mobile"),
			want:   120.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var referrers []string
			if tt.referrer != nil {
				referrers = []string{*tt.referrer}
			}
			var devices []string
			if tt.device != nil {
				devices = []string{*tt.device}
			}
			var pages []string
			if tt.page != nil {
				pages = []string{*tt.page}
			}

			got, err := repo.GetAvgSessionDurationWithFilter(ctx, AnalyticsQuery{
				SiteID: site.ID,
				From:   from,
				To:     to,
				Filter: AnalyticsFilter{
					Referrer: referrers,
					Device:   devices,
					Page:     pages,
				},
			})
			if err != nil {
				t.Errorf("GetAvgSessionDurationWithFilter() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("GetAvgSessionDurationWithFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTopPagesCountsDistinctClientsNotSessions(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAnalyticsRepository(db)
	ctx := context.Background()

	site := createTestSite(t, db)
	clientID := createTestClient(t, db, site.ID, "hash1", "desktop", "Chrome", "Windows")

	firstEventTime := time.Date(2026, 3, 9, 10, 0, 0, 0, time.UTC)
	session1 := insertSessionWithPath(t, db, site.ID, clientID, "/pricing", firstEventTime, 10, 1)
	insertPageViewEvent(t, db, session1, "/pricing", firstEventTime)

	secondEventTime := firstEventTime.Add(time.Hour)
	session2 := insertSessionWithPath(t, db, site.ID, clientID, "/pricing", secondEventTime, 10, 1)
	insertPageViewEvent(t, db, session2, "/pricing", secondEventTime)

	stats, err := repo.GetTopPages(ctx, site.ID, firstEventTime.Add(-time.Minute), secondEventTime.Add(time.Minute), 10)
	if err != nil {
		t.Fatalf("GetTopPages() error = %v", err)
	}

	if len(stats) != 1 {
		t.Fatalf("GetTopPages() returned %d rows, want 1", len(stats))
	}
	if stats[0].Path != "/pricing" {
		t.Fatalf("GetTopPages()[0].Path = %q, want /pricing", stats[0].Path)
	}
	if stats[0].Views != 2 {
		t.Fatalf("GetTopPages()[0].Views = %d, want 2", stats[0].Views)
	}
	if stats[0].Visitors != 1 {
		t.Fatalf("GetTopPages()[0].Visitors = %d, want 1", stats[0].Visitors)
	}
}

func TestGetActivePagesCountsDistinctClientsNotSessions(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAnalyticsRepository(db)
	ctx := context.Background()

	site := createTestSite(t, db)
	clientID := createTestClient(t, db, site.ID, "hash1", "desktop", "Chrome", "Windows")

	firstEventTime := time.Date(2026, 3, 9, 10, 0, 0, 0, time.UTC)
	session1 := insertSessionWithPath(t, db, site.ID, clientID, "/pricing", firstEventTime, 10, 1)
	insertPageViewEvent(t, db, session1, "/pricing", firstEventTime)

	secondEventTime := firstEventTime.Add(time.Minute)
	session2 := insertSessionWithPath(t, db, site.ID, clientID, "/pricing", secondEventTime, 10, 1)
	insertPageViewEvent(t, db, session2, "/pricing", secondEventTime)

	stats, err := repo.GetActivePages(ctx, site.ID, firstEventTime.Add(-time.Minute), 10, 0)
	if err != nil {
		t.Fatalf("GetActivePages() error = %v", err)
	}

	if len(stats) != 1 {
		t.Fatalf("GetActivePages() returned %d rows, want 1", len(stats))
	}
	if stats[0].Path != "/pricing" {
		t.Fatalf("GetActivePages()[0].Path = %q, want /pricing", stats[0].Path)
	}
	if stats[0].Visitors != 1 {
		t.Fatalf("GetActivePages()[0].Visitors = %d, want 1", stats[0].Visitors)
	}
}

func TestGetTimeSeriesStatsWithFilterBucketsPageViewsByEventTime(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAnalyticsRepository(db)
	ctx := context.Background()

	site := createTestSite(t, db)
	clientID := createTestClient(t, db, site.ID, "hash1", "desktop", "Chrome", "Windows")

	enterTime := time.Date(2026, 3, 9, 23, 59, 0, 0, time.UTC)
	eventTime := time.Date(2026, 3, 10, 0, 1, 0, 0, time.UTC)
	sessionID := insertSessionWithPath(t, db, site.ID, clientID, "/pricing", enterTime, 120, 1)
	insertPageViewEvent(t, db, sessionID, "/pricing", eventTime)

	stats, err := repo.GetTimeSeriesStatsWithFilter(ctx, AnalyticsQuery{
		SiteID: site.ID,
		From:   enterTime.Add(-time.Minute),
		To:     eventTime.Add(time.Minute),
		Bucket: TimeBucketDaily,
	})
	if err != nil {
		t.Fatalf("GetTimeSeriesStatsWithFilter() error = %v", err)
	}

	if len(stats) != 2 {
		t.Fatalf("GetTimeSeriesStatsWithFilter() returned %d rows, want 2", len(stats))
	}

	enterBucket := enterTime.Unix() / 86400
	eventBucket := eventTime.Unix() / 86400
	if stats[0].DateBucket != enterBucket {
		t.Fatalf("stats[0].DateBucket = %d, want %d", stats[0].DateBucket, enterBucket)
	}
	if stats[0].Visitors != 1 || stats[0].Sessions != 1 || stats[0].PageViews != 0 {
		t.Fatalf("stats[0] = %+v, want visitor/session without pageview", stats[0])
	}
	if stats[1].DateBucket != eventBucket {
		t.Fatalf("stats[1].DateBucket = %d, want %d", stats[1].DateBucket, eventBucket)
	}
	if stats[1].Visitors != 0 || stats[1].Sessions != 0 || stats[1].PageViews != 1 {
		t.Fatalf("stats[1] = %+v, want pageview without visitor/session", stats[1])
	}
}

func stringPtr(s string) *string {
	return &s
}
