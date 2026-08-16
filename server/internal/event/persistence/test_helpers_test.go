package persistence

import (
	"database/sql"
	"testing"

	authpersistence "github.com/lovely-eye/server/internal/auth/persistence"
	"github.com/lovely-eye/server/internal/platform/database"
	sitepersistence "github.com/lovely-eye/server/internal/site/persistence"
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

func createTestSite(t *testing.T, db *bun.DB) *sitepersistence.Site {
	t.Helper()

	user := &authpersistence.User{Username: "event-test", PasswordHash: "hash", Role: "admin"}
	_, err := db.NewInsert().Model(user).Exec(t.Context())
	require.NoError(t, err)
	site := &sitepersistence.Site{UserID: user.ID, Name: "Event Test", PublicKey: "event-test"}
	_, err = db.NewInsert().Model(site).Exec(t.Context())
	require.NoError(t, err)
	return site
}
