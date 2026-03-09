package migrations

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

const (
	sqliteAnalyticsClientDimensionsUp   = "sqlite/20260309151021_analytics_client_dimensions_enums.up.sql"
	sqliteAnalyticsClientDimensionsDown = "sqlite/20260309151021_analytics_client_dimensions_enums.down.sql"
)

func TestSQLiteAnalyticsClientDimensionsMigrationCompatibility(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite database: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})

	schema := `
CREATE TABLE sites (
  id integer NOT NULL PRIMARY KEY AUTOINCREMENT
);
CREATE TABLE clients (
  id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  site_id integer NOT NULL,
  hash varchar NOT NULL,
  country varchar NULL,
  device varchar NULL,
  browser varchar NULL,
  os varchar NULL,
  screen_size varchar NULL,
  CONSTRAINT clients_site_id_fkey FOREIGN KEY (site_id) REFERENCES sites (id) ON UPDATE NO ACTION ON DELETE NO ACTION
);
`
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("create old schema: %v", err)
	}

	if _, err := db.Exec(`INSERT INTO sites (id) VALUES (1)`); err != nil {
		t.Fatalf("insert site: %v", err)
	}
	if _, err := db.Exec(`
INSERT INTO clients (id, site_id, hash, country, device, browser, os, screen_size) VALUES
  (1, 1, 'hash-1', 'US', 'desktop', 'Chrome', 'Windows', '390x844'),
  (2, 1, 'hash-2', NULL, NULL, 'Arc', 'Haiku', 'watch');
`); err != nil {
		t.Fatalf("insert legacy clients: %v", err)
	}

	execSQLiteMigrationFile(t, db, sqliteAnalyticsClientDimensionsUp)

	var device1, browser1, os1, screenSize1 int
	if err := db.QueryRow(`SELECT device, browser, os, screen_size FROM clients WHERE id = 1`).Scan(&device1, &browser1, &os1, &screenSize1); err != nil {
		t.Fatalf("select migrated client 1: %v", err)
	}
	if device1 != 2 || browser1 != 3 || os1 != 11 || screenSize1 != 2 {
		t.Fatalf("unexpected migrated codes for client 1: device=%d browser=%d os=%d screen_size=%d", device1, browser1, os1, screenSize1)
	}

	var device2, browser2, os2, screenSize2 int
	if err := db.QueryRow(`SELECT device, browser, os, screen_size FROM clients WHERE id = 2`).Scan(&device2, &browser2, &os2, &screenSize2); err != nil {
		t.Fatalf("select migrated client 2: %v", err)
	}
	if device2 != 0 || browser2 != 1 || os2 != 1 || screenSize2 != 1 {
		t.Fatalf("unexpected migrated codes for client 2: device=%d browser=%d os=%d screen_size=%d", device2, browser2, os2, screenSize2)
	}

	execSQLiteMigrationFile(t, db, sqliteAnalyticsClientDimensionsDown)

	var restoredDevice1, restoredBrowser1, restoredOS1, restoredScreenSize1 sql.NullString
	if err := db.QueryRow(`SELECT device, browser, os, screen_size FROM clients WHERE id = 1`).Scan(&restoredDevice1, &restoredBrowser1, &restoredOS1, &restoredScreenSize1); err != nil {
		t.Fatalf("select rolled back client 1: %v", err)
	}
	if restoredDevice1.String != "desktop" || restoredBrowser1.String != "Chrome" || restoredOS1.String != "Windows" || restoredScreenSize1.String != "xs" {
		t.Fatalf("unexpected rolled back labels for client 1: device=%q browser=%q os=%q screen_size=%q", restoredDevice1.String, restoredBrowser1.String, restoredOS1.String, restoredScreenSize1.String)
	}

	var restoredDevice2, restoredBrowser2, restoredOS2, restoredScreenSize2 sql.NullString
	if err := db.QueryRow(`SELECT device, browser, os, screen_size FROM clients WHERE id = 2`).Scan(&restoredDevice2, &restoredBrowser2, &restoredOS2, &restoredScreenSize2); err != nil {
		t.Fatalf("select rolled back client 2: %v", err)
	}
	if restoredDevice2.Valid || restoredBrowser2.String != "Other" || restoredOS2.String != "Other" || restoredScreenSize2.String != "watch" {
		t.Fatalf("unexpected rolled back labels for client 2: device=%q browser=%q os=%q screen_size=%q", restoredDevice2.String, restoredBrowser2.String, restoredOS2.String, restoredScreenSize2.String)
	}
}

func execSQLiteMigrationFile(t *testing.T, db *sql.DB, relativePath string) {
	t.Helper()

	var (
		bytes []byte
		err   error
	)
	switch relativePath {
	case sqliteAnalyticsClientDimensionsUp:
		bytes, err = os.ReadFile(sqliteAnalyticsClientDimensionsUp)
	case sqliteAnalyticsClientDimensionsDown:
		bytes, err = os.ReadFile(sqliteAnalyticsClientDimensionsDown)
	default:
		t.Fatalf("unsupported migration path: %s", relativePath)
	}
	if err != nil {
		t.Fatalf("read migration %s: %v", relativePath, fmt.Errorf("read migration file: %w", err))
	}
	if _, err := db.Exec(string(bytes)); err != nil {
		t.Fatalf("exec migration %s: %v", relativePath, err)
	}
}
