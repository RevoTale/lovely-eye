package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/lovely-eye/server/internal/config"
	"github.com/lovely-eye/server/internal/database"
	"github.com/lovely-eye/server/migrations"
	"github.com/uptrace/bun/migrate"
)

func main() {
	cfg := config.Load()

	db, err := database.New(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)

	migs, err := migrations.NewMigrations()
	if err != nil {
		log.Fatalf("Failed to load migrations: %v", err)
	}

	migrator := migrate.NewMigrator(db, migs)
	ctx := context.Background()

	dbType := strings.ToUpper(cfg.Database.Driver)
	fmt.Printf("=== Testing %s Migrations ===\n", dbType)
	fmt.Printf("DB_DRIVER: %s\n", cfg.Database.Driver)
	fmt.Printf("DB_DSN: %s\n\n", cfg.Database.DSN)

	// Step 1: Initialize
	fmt.Println("Step 1: Initialize migration tables")
	if err := migrator.Init(ctx); err != nil {
		log.Fatalf("Init failed: %v", err)
	}

	// Step 2: Check initial status
	fmt.Println("\nStep 2: Check initial status (should show no applied migrations)")
	ms, err := migrator.MigrationsWithStatus(ctx)
	if err != nil {
		log.Fatalf("Status check failed: %v", err)
	}
	fmt.Printf("migrations: %s\n", ms)
	fmt.Printf("unapplied migrations: %s\n", ms.Unapplied())
	fmt.Printf("last migration group: %s\n", ms.LastGroup())

	// Step 3: Apply all migrations
	fmt.Println("\nStep 3: Apply all migrations (UP)")
	group, err := migrator.Migrate(ctx)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	if group.IsZero() {
		fmt.Println("there are no new migrations to run (database is up to date)")
	} else {
		fmt.Printf("migrated to %s\n", group)
	}

	// Step 4: Check status after migration
	fmt.Println("\nStep 4: Check status after migration (should show all applied)")
	ms, err = migrator.MigrationsWithStatus(ctx)
	if err != nil {
		log.Fatalf("Status check failed: %v", err)
	}
	fmt.Printf("migrations: %s\n", ms)
	fmt.Printf("unapplied migrations: %s\n", ms.Unapplied())
	fmt.Printf("last migration group: %s\n", ms.LastGroup())

	// Verify all migrations were actually applied
	if len(ms.Unapplied()) > 0 {
		log.Fatalf("ERROR: Not all migrations were applied. Unapplied: %s", ms.Unapplied())
	}

	// Step 5: Count and rollback all migrations
	fmt.Println("\nStep 5: Rollback all migrations (DOWN)")
	appliedCount := len(ms.Applied())
	fmt.Printf("Found %d applied migrations to rollback\n", appliedCount)

	for i := 0; i < appliedCount; i++ {
		fmt.Printf("Rollback iteration %d of %d\n", i+1, appliedCount)
		group, err := migrator.Rollback(ctx)
		if err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
		if group.IsZero() {
			fmt.Println("there are no groups to roll back")
			break
		}
		fmt.Printf("rolled back %s\n", group)
	}

	// Step 6: Check status after rollback
	fmt.Println("\nStep 6: Check status after rollback (should show no applied migrations)")
	ms, err = migrator.MigrationsWithStatus(ctx)
	if err != nil {
		log.Fatalf("Status check failed: %v", err)
	}
	fmt.Printf("migrations: %s\n", ms)
	fmt.Printf("unapplied migrations: %s\n", ms.Unapplied())
	fmt.Printf("last migration group: %s\n", ms.LastGroup())

	// Step 7: Apply migrations again
	fmt.Println("\nStep 7: Apply migrations again to verify idempotency")
	group, err = migrator.Migrate(ctx)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	if group.IsZero() {
		fmt.Println("there are no new migrations to run (database is up to date)")
	} else {
		fmt.Printf("migrated to %s\n", group)
	}

	// Step 8: Final status check
	fmt.Println("\nStep 8: Final status check")
	ms, err = migrator.MigrationsWithStatus(ctx)
	if err != nil {
		log.Fatalf("Status check failed: %v", err)
	}
	fmt.Printf("migrations: %s\n", ms)
	fmt.Printf("unapplied migrations: %s\n", ms.Unapplied())
	fmt.Printf("last migration group: %s\n", ms.LastGroup())

	// Verify all migrations were applied (idempotency check)
	if len(ms.Unapplied()) > 0 {
		log.Fatalf("ERROR: Not all migrations were applied on second run. Unapplied: %s", ms.Unapplied())
	}

	fmt.Printf("\n✓ %s migration test completed successfully\n", dbType)
}
