package database

import (
	"context"
	"fmt"

	"github.com/lovely-eye/server/migrations"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

func Migrate(ctx context.Context, db *bun.DB) error {
	migs, err := migrations.NewMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	migrator := migrate.NewMigrator(db, migs)

	if err := migrator.Init(ctx); err != nil {
		return fmt.Errorf("failed to init migrator: %w", err)
	}

	group, err := migrator.Migrate(ctx)
	if err != nil {
		return fmt.Errorf("failed to migrate: %w", err)
	}

	if group.IsZero() {
		fmt.Println("No new migrations to run")
	} else {
		fmt.Printf("Migrated to %s\n", group)
	}

	return nil
}
