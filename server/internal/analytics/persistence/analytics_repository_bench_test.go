package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	authpersistence "github.com/lovely-eye/server/internal/auth/persistence"
	"github.com/lovely-eye/server/internal/platform/database"
	sitepersistence "github.com/lovely-eye/server/internal/site/persistence"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"

	_ "modernc.org/sqlite"
)

func BenchmarkAnalyticsDashboardReads(b *testing.B) {
	db, siteID, from, to := newAnalyticsBenchmarkFixture(b)
	repository := New(db)
	ctx := context.Background()

	b.Run("visitor-count", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			if _, err := repository.GetVisitorCount(ctx, siteID, from, to); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("top-pages", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			if _, err := repository.GetTopPages(ctx, siteID, from, to, 25); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func newAnalyticsBenchmarkFixture(b *testing.B) (*bun.DB, int64, time.Time, time.Time) {
	b.Helper()
	sqlDB, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		b.Fatal(err)
	}
	db := bun.NewDB(sqlDB, sqlitedialect.New())
	b.Cleanup(func() {
		if err := db.Close(); err != nil {
			b.Error(err)
		}
	})
	ctx := context.Background()
	if err := database.Migrate(ctx, db); err != nil {
		b.Fatal(err)
	}
	user := &authpersistence.User{Username: "benchmark", PasswordHash: "hash", Role: "admin"}
	if _, err := db.NewInsert().Model(user).Exec(ctx); err != nil {
		b.Fatal(err)
	}
	site := &sitepersistence.Site{UserID: user.ID, Name: "Benchmark", PublicKey: "benchmark-key"}
	if _, err := db.NewInsert().Model(site).Exec(ctx); err != nil {
		b.Fatal(err)
	}

	from := time.Date(2026, time.July, 1, 0, 0, 0, 0, time.UTC)
	for index := range 100 {
		insertAnalyticsBenchmarkVisit(b, db, site.ID, index, from.Add(time.Duration(index)*time.Hour))
	}
	return db, site.ID, from.Add(-time.Hour), from.Add(101 * time.Hour)
}

func insertAnalyticsBenchmarkVisit(
	b *testing.B,
	db *bun.DB,
	siteID int64,
	index int,
	timestamp time.Time,
) {
	b.Helper()
	ctx := context.Background()
	client := &Client{SiteID: siteID, Hash: fmt.Sprintf("benchmark-%03d", index)}
	if _, err := db.NewInsert().Model(client).Exec(ctx); err != nil {
		b.Fatal(err)
	}
	timestampUnix := timestamp.Unix()
	session := &Session{
		SiteID:        siteID,
		ClientID:      client.ID,
		EnterTime:     timestampUnix,
		EnterHour:     timestampUnix / 3600,
		EnterDay:      timestampUnix / 86400,
		EnterPath:     fmt.Sprintf("/page/%d", index%10),
		ExitTime:      timestamp.Add(time.Minute).Unix(),
		ExitHour:      timestampUnix / 3600,
		ExitDay:       timestampUnix / 86400,
		ExitPath:      fmt.Sprintf("/page/%d", index%10),
		PageViewCount: 1,
	}
	if _, err := db.NewInsert().Model(session).Exec(ctx); err != nil {
		b.Fatal(err)
	}
	event := &Event{
		SessionID: session.ID,
		Time:      timestampUnix,
		Hour:      timestampUnix / 3600,
		Day:       timestampUnix / 86400,
		Path:      session.EnterPath,
	}
	if _, err := db.NewInsert().Model(event).Exec(ctx); err != nil {
		b.Fatal(err)
	}
}
