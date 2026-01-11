package database

import (
	"context"
	"fmt"

	"github.com/lovely-eye/server/migrations"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

func NewMigrator(db *bun.DB) *migrate.Migrator {
	return migrate.NewMigrator(db, migrations.Migrations)
}

func Migrate(ctx context.Context, db *bun.DB) error {
	migrator := NewMigrator(db)

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

func Rollback(ctx context.Context, db *bun.DB) error {
	migrator := NewMigrator(db)

	group, err := migrator.Rollback(ctx)
	if err != nil {
		return fmt.Errorf("failed to rollback: %w", err)
	}

	if group.IsZero() {
		fmt.Println("No migrations to rollback")
	} else {
		fmt.Printf("Rolled back %s\n", group)
	}

	return nil
}

func MigrationStatus(ctx context.Context, db *bun.DB) error {
	migrator := NewMigrator(db)

	ms, err := migrator.MigrationsWithStatus(ctx)
	if err != nil {
		return fmt.Errorf("failed to get migration status: %w", err)
	}

	fmt.Println("Migrations:")
	for _, m := range ms {
		status := "not applied"
		if !m.IsApplied() {
			status = "pending"
		} else {
			status = "applied"
		}
		fmt.Printf("  %s: %s\n", m.Name, status)
	}

	return nil
}
