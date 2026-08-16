package persistence

import (
	"database/sql"
	"testing"
	"time"

	analyticspersistence "github.com/lovely-eye/server/internal/analytics/persistence"
	authpersistence "github.com/lovely-eye/server/internal/auth/persistence"
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

func createTestSite(t *testing.T, db *bun.DB) *Site {
	t.Helper()

	user := &authpersistence.User{Username: "site-persistence", PasswordHash: "hash", Role: "admin"}
	_, err := db.NewInsert().Model(user).Exec(t.Context())
	require.NoError(t, err)
	site := &Site{UserID: user.ID, Name: "Site Test", PublicKey: "site-persistence"}
	_, err = db.NewInsert().Model(site).Exec(t.Context())
	require.NoError(t, err)
	_, err = db.NewInsert().Model(&Domain{SiteID: site.ID, Domain: "example.com"}).Exec(t.Context())
	require.NoError(t, err)
	return site
}

func createTestClient(
	t *testing.T,
	db *bun.DB,
	siteID int64,
	hash string,
	device string,
	browser string,
	os string,
) int64 {
	t.Helper()

	client := &analyticspersistence.Client{
		SiteID:  siteID,
		Hash:    hash,
		Device:  analyticspersistence.ClientDeviceFromLegacyLabel(device),
		Browser: analyticspersistence.ClientBrowserFromLegacyLabel(browser),
		OS:      analyticspersistence.ClientOSFromLegacyLabel(os),
	}
	_, err := db.NewInsert().Model(client).Exec(t.Context())
	require.NoError(t, err)
	return client.ID
}

func insertSessionWithPath(
	t *testing.T,
	db *bun.DB,
	siteID int64,
	clientID int64,
	path string,
	enterTime time.Time,
	durationSeconds int,
	pageViewCount int,
) int64 {
	t.Helper()

	enterUnix := enterTime.Unix()
	exitUnix := enterUnix + int64(durationSeconds)
	session := &analyticspersistence.Session{
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
	_, err := db.NewInsert().Model(session).Exec(t.Context())
	require.NoError(t, err)
	return session.ID
}
