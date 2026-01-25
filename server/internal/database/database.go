package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lovely-eye/server/internal/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"

	_ "modernc.org/sqlite"
)

func New(cfg *config.DatabaseConfig) (*bun.DB, error) {
	var db *bun.DB

	switch cfg.Driver {
	case "postgres", "postgresql":
		sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.DSN)))
		sqldb.SetMaxOpenConns(cfg.MaxConns)
		sqldb.SetMaxIdleConns(cfg.MinConns)
		db = bun.NewDB(sqldb, pgdialect.New())

	case "sqlite", "sqlite3":
		sqldb, err := sql.Open("sqlite", cfg.DSN)
		if err != nil {
			return nil, fmt.Errorf("failed to open sqlite: %w", err)
		}
		sqldb.SetMaxOpenConns(1)
		db = bun.NewDB(sqldb, sqlitedialect.New())

	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	// Add query logging in debug mode
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(false),
	))

	pingCtx := context.Background()
	if cfg.ConnectTimeout > 0 {
		var cancel context.CancelFunc
		pingCtx, cancel = context.WithTimeout(pingCtx, cfg.ConnectTimeout)
		defer cancel()
	}
	if err := db.PingContext(pingCtx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func Close(db *bun.DB) error {
	if db != nil {
		if err := db.Close(); err != nil {
			return fmt.Errorf("failed to close database: %w", err)
		}
	}
	return nil
}
