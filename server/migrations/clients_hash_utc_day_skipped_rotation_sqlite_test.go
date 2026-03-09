package migrations

import (
	"database/sql"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

const sqliteClientsHashUTCDaySkippedRotationUp = "sqlite/20260309183000_clients_hash_utc_day_skipped_rotation.up.sql"

func TestSQLiteClientsHashUTCDaySkippedRotationMigrationCompatibility(t *testing.T) {
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
  device integer NOT NULL DEFAULT 0,
  browser integer NOT NULL DEFAULT 0,
  os integer NOT NULL DEFAULT 0,
  screen_size integer NOT NULL DEFAULT 0
);
CREATE TABLE sessions (
  id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  site_id integer NOT NULL,
  client_id integer NOT NULL,
  enter_time integer NOT NULL,
  enter_hour integer NOT NULL,
  enter_day integer NOT NULL,
  enter_path varchar NOT NULL,
  exit_time integer NOT NULL,
  exit_hour integer NOT NULL,
  exit_day integer NOT NULL,
  exit_path varchar NOT NULL,
  referrer varchar NULL,
  utm_source varchar NULL,
  utm_medium varchar NULL,
  utm_campaign varchar NULL,
  duration integer NOT NULL DEFAULT 0,
  page_view_count integer NOT NULL DEFAULT 0
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
  (1, 1, 'same-hash', 'US', 2, 3, 11, 4),
  (2, 1, 'same-hash', 'US', 2, 3, 11, 4),
  (3, 1, 'other-hash', '', 2, 3, 11, 4);
`); err != nil {
		t.Fatalf("insert clients: %v", err)
	}
	if _, err := db.Exec(`
INSERT INTO sessions (id, site_id, client_id, enter_time, enter_hour, enter_day, enter_path, exit_time, exit_hour, exit_day, exit_path, duration, page_view_count) VALUES
  (1, 1, 2, 100, 0, 0, '/old', 200, 0, 0, '/old', 100, 1),
  (2, 1, 3, 300, 0, 0, '/other', 400, 0, 0, '/other', 100, 1);
`); err != nil {
		t.Fatalf("insert sessions: %v", err)
	}

	bytes, err := os.ReadFile(sqliteClientsHashUTCDaySkippedRotationUp)
	if err != nil {
		t.Fatalf("read migration: %v", err)
	}
	if _, err := db.Exec(string(bytes)); err != nil {
		t.Fatalf("exec migration: %v", err)
	}

	var clientCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM clients`).Scan(&clientCount); err != nil {
		t.Fatalf("count clients: %v", err)
	}
	if clientCount != 2 {
		t.Fatalf("expected 2 clients after merge, got %d", clientCount)
	}

	var mergedSessionClientID int
	if err := db.QueryRow(`SELECT client_id FROM sessions WHERE id = 1`).Scan(&mergedSessionClientID); err != nil {
		t.Fatalf("select merged session client_id: %v", err)
	}
	if mergedSessionClientID != 1 {
		t.Fatalf("expected duplicate session to be repointed to client 1, got %d", mergedSessionClientID)
	}

	if _, err := db.Exec(`INSERT INTO clients (site_id, hash, country, device, browser, os, screen_size) VALUES (1, 'same-hash', 'US', 2, 3, 11, 4)`); err == nil {
		t.Fatalf("expected unique site/hash enforcement after migration")
	}
}
