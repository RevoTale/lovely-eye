package persistence

import (
	"database/sql"
	"testing"

	"github.com/lovely-eye/server/internal/platform/database"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *bun.DB {
	t.Helper()

	sqldb, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	db := bun.NewDB(sqldb, sqlitedialect.New())
	require.NoError(t, database.Migrate(t.Context(), db))
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	return db
}
