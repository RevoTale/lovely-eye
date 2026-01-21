package repository

import (
	"context"
	"database/sql"
	"fmt"
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
		Device:  device,
		Browser: browser,
		OS:      os,
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
			name: "only bounce sessions exist",
			setup: func() {
				clientID := createTestClient(t, db, site.ID, "hash1", "desktop", "Chrome", "Windows")
				sessionTime := now.Add(-1 * time.Hour)
				// Bounce session: page_view_count = 1
				insertSession(t, db, site.ID, clientID, sessionTime, 10, 1)
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

	// Insert test sessions (non-bounce sessions: page_view_count > 1)
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

func TestGetAvgSessionDuration_ExcludesBounces(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAnalyticsRepository(db)
	ctx := context.Background()

	site := createTestSite(t, db)

	now := time.Now()
	from := now.Add(-24 * time.Hour)
	to := now

	// Insert non-bounce and bounce sessions
	client1 := createTestClient(t, db, site.ID, "hash1", "desktop", "Chrome", "Windows")
	insertSession(t, db, site.ID, client1, now.Add(-2*time.Hour), 100, 2) // non-bounce: page_view_count = 2

	client2 := createTestClient(t, db, site.ID, "hash2", "mobile", "Safari", "iOS")
	insertSession(t, db, site.ID, client2, now.Add(-3*time.Hour), 1000, 1) // bounce: page_view_count = 1

	got, err := repo.GetAvgSessionDuration(ctx, site.ID, from, to)
	if err != nil {
		t.Fatalf("GetAvgSessionDuration() error = %v", err)
	}

	// Should only include the non-bounce session
	want := 100.0
	if got != want {
		t.Errorf("GetAvgSessionDuration() = %v, want %v (bounce sessions should be excluded)", got, want)
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
				// Create client and session with different referrer
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

			got, err := repo.GetAvgSessionDurationWithFilter(ctx, site.ID, from, to, referrers, devices, pages, nil)
			fmt.Println(got)
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

	// Insert sessions with different referrers and devices
	// Session 1: desktop, google.com referrer
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

	// Session 2: mobile, no referrer
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

	tests := []struct {
		name     string
		referrer *string
		device   *string
		page     *string
		want     float64
	}{
		{
			name: "no filter",
			want: 90.0,
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

			got, err := repo.GetAvgSessionDurationWithFilter(ctx, site.ID, from, to, referrers, devices, pages, nil)
			if err != nil {
				t.Errorf("GetAvgSessionDurationWithFilter() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("GetAvgSessionDurationWithFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
