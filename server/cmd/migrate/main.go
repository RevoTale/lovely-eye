package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/lovely-eye/server/internal/platform/config"
	"github.com/lovely-eye/server/internal/platform/database"
	"github.com/lovely-eye/server/migrations"
	"github.com/uptrace/bun/migrate"
)

type command string

const (
	commandHelp   command = "help"
	commandInit   command = "init"
	commandUp     command = "up"
	commandDown   command = "down"
	commandStatus command = "status"

	usage = `Usage: migrate <command>

Commands:
  init    Create migration tables in the database
  up      Apply all pending migrations
  down    Roll back the last migration group
  status  Show migration status and history
`
)

func main() {
	if err := run(context.Background(), os.Args[1:], os.Stdout); err != nil {
		slog.Error("migration command failed", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, output io.Writer) (err error) {
	selected, err := parseCommand(args)
	if err != nil {
		return err
	}
	if selected == commandHelp {
		return writeOutput(output, "%s", usage)
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}
	db, err := database.New(ctx, cfg.Database)
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer func() {
		if closeErr := database.Close(db); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}()

	migs, err := migrations.NewMigrations()
	if err != nil {
		return fmt.Errorf("load migrations: %w", err)
	}

	return executeCommand(ctx, selected, output, migrate.NewMigrator(db, migs))
}

func parseCommand(args []string) (command, error) {
	if len(args) == 0 {
		return commandHelp, nil
	}
	if len(args) != 1 {
		return "", fmt.Errorf("expected exactly one command\n%s", usage)
	}

	selected := command(args[0])
	switch selected {
	case commandHelp, "-h", "--help":
		return commandHelp, nil
	case commandInit, commandUp, commandDown, commandStatus:
		return selected, nil
	default:
		return "", fmt.Errorf("unknown command %q\n%s", args[0], usage)
	}
}

func executeCommand(
	ctx context.Context,
	selected command,
	output io.Writer,
	migrator *migrate.Migrator,
) error {
	switch selected {
	case commandInit:
		if err := migrator.Init(ctx); err != nil {
			return fmt.Errorf("initialize migrator: %w", err)
		}
		return nil
	case commandUp:
		return withMigrationLock(ctx, migrator, func() error {
			group, err := migrator.Migrate(ctx)
			if err != nil {
				return fmt.Errorf("migrate: %w", err)
			}
			if group.IsZero() {
				return writeOutput(output, "there are no new migrations to run (database is up to date)\n")
			}
			return writeOutput(output, "migrated to %s\n", group)
		})
	case commandDown:
		return withMigrationLock(ctx, migrator, func() error {
			group, err := migrator.Rollback(ctx)
			if err != nil {
				return fmt.Errorf("rollback: %w", err)
			}
			if group.IsZero() {
				return writeOutput(output, "there are no groups to roll back\n")
			}
			return writeOutput(output, "rolled back %s\n", group)
		})
	case commandStatus:
		migrationStatus, err := migrator.MigrationsWithStatus(ctx)
		if err != nil {
			return fmt.Errorf("get migration status: %w", err)
		}
		if err = writeOutput(output, "migrations: %s\n", migrationStatus); err != nil {
			return err
		}
		if err = writeOutput(output, "unapplied migrations: %s\n", migrationStatus.Unapplied()); err != nil {
			return err
		}
		return writeOutput(output, "last migration group: %s\n", migrationStatus.LastGroup())
	default:
		return fmt.Errorf("unsupported migration command %q", selected)
	}
}

func writeOutput(output io.Writer, format string, args ...any) error {
	if _, err := fmt.Fprintf(output, format, args...); err != nil {
		return fmt.Errorf("write command output: %w", err)
	}
	return nil
}

func withMigrationLock(
	ctx context.Context,
	migrator *migrate.Migrator,
	action func() error,
) (err error) {
	if err := migrator.Lock(ctx); err != nil {
		return fmt.Errorf("acquire migration lock: %w", err)
	}
	defer func() {
		if unlockErr := migrator.Unlock(ctx); unlockErr != nil {
			err = errors.Join(err, fmt.Errorf("release migration lock: %w", unlockErr))
		}
	}()

	return action()
}
