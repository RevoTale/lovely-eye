package persistence

import (
	"database/sql"
	"os"
	"testing"

	"github.com/lovely-eye/server/internal/platform/database"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func TestPagedBreakdownsPostgres(t *testing.T) {
	dsn := os.Getenv("LOVELY_EYE_TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("LOVELY_EYE_TEST_POSTGRES_DSN is not set")
	}
	t.Setenv("DB_DRIVER", "postgres")

	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqlDB, pgdialect.New())
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("close PostgreSQL database: %v", err)
		}
	})
	if err := database.Migrate(t.Context(), db); err != nil {
		t.Fatalf("migrate PostgreSQL database: %v", err)
	}

	testPagedBreakdownsReturnExactWindowTotals(t, db)
}
