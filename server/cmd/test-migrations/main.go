package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lovely-eye/server/internal/platform/config"
	"github.com/lovely-eye/server/internal/platform/database"
	"github.com/lovely-eye/server/migrations"
	"github.com/uptrace/bun/migrate"
)

func main() {
	os.Exit(run())
}

func run() int {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		return 1
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	db, err := database.New(ctx, cfg.Database)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		return 1
	}
	defer func() {
		if err := database.Close(db); err != nil {
			slog.Error("db close error", "error", err)
		}
	}()

	migs, err := migrations.NewMigrations()
	if err != nil {
		slog.Error("failed to load migrations", "error", err)
		return 1
	}

	// The integration suite treats "applied" as "successfully applied" so a failed
	// migration remains pending and the rollback/reapply assertions stay accurate.
	migrator := migrate.NewMigrator(db, migs, migrate.WithMarkAppliedOnSuccess(true))

	dbType := (cfg.Database.Driver)
	fmt.Printf("=== Testing %s Migrations ===\n", dbType)
	fmt.Printf("DB_DRIVER: %s\n", cfg.Database.Driver)
	fmt.Printf("DB_DSN: %s\n\n", cfg.Database.DSN)

	fmt.Println("Step 1: Initialize migration tables")
	if err := migrator.Init(ctx); err != nil {
		slog.Error("init failed", "error", err)
		return 1
	}

	printStatus := func(title string) (migrate.MigrationSlice, int) {
		ms, err := migrator.MigrationsWithStatus(ctx)
		if err != nil {
			slog.Error("status check failed", "error", err)
			return nil, 1
		}
		fmt.Printf("\n%s\n", title)
		fmt.Printf("applied: %d\n", len(ms.Applied()))
		fmt.Printf("unapplied: %d\n", len(ms.Unapplied()))
		fmt.Printf("last group: %v\n", ms.LastGroup())
		return ms, 0
	}

	if _, code := printStatus("Step 2: Initial status"); code != 0 {
		return code
	}

	fmt.Println("\nStep 3: Apply all migrations (UP)")
	group, err := migrator.Migrate(ctx)
	if err != nil {
		slog.Error("migration up failed", "error", err)
		return 1
	}
	if group.IsZero() {
		fmt.Println("no new migrations to run (database is up to date)")
	} else {
		fmt.Printf("migrated to %s\n", group)
	}

	ms, code := printStatus("Step 4: Status after migration (should show all applied)")
	if code != 0 {
		return code
	}
	if len(ms.Unapplied()) > 0 {
		slog.Error("not all migrations were applied", "unapplied", fmt.Sprint(ms.Unapplied()))
		return 1
	}

	fmt.Println("\nStep 5: Rollback all migrations (DOWN)")
	for {
		group, err := migrator.Rollback(ctx)
		if err != nil {
			slog.Error("rollback failed", "error", err)
			return 1
		}
		if group.IsZero() {
			fmt.Println("no more groups to roll back")
			break
		}
		fmt.Printf("rolled back %s\n", group)
	}

	if _, code := printStatus("Step 6: Status after rollback (should show none applied)"); code != 0 {
		return code
	}

	fmt.Println("\nStep 7: Apply migrations again (idempotency)")
	group, err = migrator.Migrate(ctx)
	if err != nil {
		slog.Error("migration up (2nd run) failed", "error", err)
		return 1
	}
	if group.IsZero() {
		fmt.Println("no new migrations to run (database is up to date)")
	} else {
		fmt.Printf("migrated to %s\n", group)
	}

	ms, code = printStatus("Step 8: Final status check")
	if code != 0 {
		return code
	}
	if len(ms.Unapplied()) > 0 {
		slog.Error("not all migrations were applied on second run", "unapplied", fmt.Sprint(ms.Unapplied()))
		return 1
	}

	fmt.Printf("\n✓ %s migration test completed successfully\n", dbType)
	return 0
}
