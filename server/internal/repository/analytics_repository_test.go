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
		Domain:    "test.com",
		Name:      "Test Site",
		PublicKey: "test-key",
	}
	if _, err := db.NewInsert().Model(site).Exec(ctx); err != nil {
		t.Fatalf("failed to insert site: %v", err)
	}

	return site
}

func insertSession(t *testing.T, db *bun.DB, siteID int64, visitorID string, startedAt, lastSeenAt time.Time, duration int, isBounce bool) {
	t.Helper()

	// Use raw SQL to avoid Bun's default:true tag behavior with boolean false values
	_, err := db.Exec(
		"INSERT INTO sessions (site_id, visitor_id, started_at, last_seen_at, duration, is_bounce) VALUES (?, ?, ?, ?, ?, ?)",
		siteID, visitorID, startedAt, lastSeenAt, duration, isBounce,
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
				insertSession(t, db, site.ID, "visitor1", now.Add(-1*time.Hour), now.Add(-1*time.Hour), 10, true)
			},
		},
		{
			name: "sessions outside date range",
			setup: func() {
				pastTime := now.Add(-30 * 24 * time.Hour)
				insertSession(t, db, site.ID, "visitor2", pastTime, pastTime, 100, false)
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

	// Insert test sessions
	insertSession(t, db, site.ID, "visitor1", now.Add(-2*time.Hour), now.Add(-2*time.Hour), 60, false)
	insertSession(t, db, site.ID, "visitor2", now.Add(-3*time.Hour), now.Add(-3*time.Hour), 120, false)
	insertSession(t, db, site.ID, "visitor3", now.Add(-4*time.Hour), now.Add(-4*time.Hour), 180, false)

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
	insertSession(t, db, site.ID, "visitor1", now.Add(-2*time.Hour), now.Add(-2*time.Hour), 100, false)
	insertSession(t, db, site.ID, "visitor2", now.Add(-3*time.Hour), now.Add(-3*time.Hour), 1000, true)

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
				// Use raw SQL to insert with referrer field
				_, err := db.Exec(
					"INSERT INTO sessions (site_id, visitor_id, started_at, last_seen_at, duration, is_bounce, referrer) VALUES (?, ?, ?, ?, ?, ?, ?)",
					site.ID, "visitor1", now.Add(-1*time.Hour), now.Add(-1*time.Hour), 100, false, "google.com",
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
	_, err := db.Exec(
		"INSERT INTO sessions (site_id, visitor_id, started_at, last_seen_at, duration, is_bounce, referrer, device) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		site.ID, "visitor1", now.Add(-2*time.Hour), now.Add(-2*time.Hour), 60, false, "google.com", "desktop",
	)
	if err != nil {
		t.Fatalf("failed to insert session 1: %v", err)
	}

	_, err = db.Exec(
		"INSERT INTO sessions (site_id, visitor_id, started_at, last_seen_at, duration, is_bounce, referrer, device) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		site.ID, "visitor2", now.Add(-3*time.Hour), now.Add(-3*time.Hour), 120, false, "", "mobile",
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
