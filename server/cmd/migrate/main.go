package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/lovely-eye/server/internal/config"
	"github.com/lovely-eye/server/internal/database"
	"github.com/lovely-eye/server/migrations"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v3"
)

func main() {
	cfg := config.Load()

	db, err := database.New(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func ()  {
		err := database.Close(db)
		if nil != err {
			slog.Error("DB close error","err",err)
		}
	}()

	migs, err := migrations.NewMigrations()
	if err != nil {
		log.Fatalf("Failed to load migrations: %v", err)
	}

	migrator := migrate.NewMigrator(db, migs)

	app := &cli.Command{
		Name:  "migrate",
		Usage: "production database migration runner",
		Commands: []*cli.Command{
			{
				Name:  "init",
				Usage: "create migration tables in the database",
				Action: func(ctx context.Context, c *cli.Command) error {
					return migrator.Init(ctx)
				},
			},
			{
				Name:  "up",
				Usage: "apply all pending migrations",
				Action: func(ctx context.Context, c *cli.Command) error {
					if err := migrator.Lock(ctx); err != nil {
						return err
					}
					defer migrator.Unlock(ctx) //nolint:errcheck

					group, err := migrator.Migrate(ctx)
					if err != nil {
						return err
					}
					if group.IsZero() {
						fmt.Printf("there are no new migrations to run (database is up to date)\n")
						return nil
					}
					fmt.Printf("migrated to %s\n", group)
					return nil
				},
			},
			{
				Name:  "down",
				Usage: "rollback the last migration group",
				Action: func(ctx context.Context, c *cli.Command) error {
					if err := migrator.Lock(ctx); err != nil {
						return err
					}
					defer migrator.Unlock(ctx) //nolint:errcheck

					group, err := migrator.Rollback(ctx)
					if err != nil {
						return err
					}
					if group.IsZero() {
						fmt.Printf("there are no groups to roll back\n")
						return nil
					}
					fmt.Printf("rolled back %s\n", group)
					return nil
				},
			},
			{
				Name:  "status",
				Usage: "show migration status and history",
				Action: func(ctx context.Context, c *cli.Command) error {
					ms, err := migrator.MigrationsWithStatus(ctx)
					if err != nil {
						return err
					}
					fmt.Printf("migrations: %s\n", ms)
					fmt.Printf("unapplied migrations: %s\n", ms.Unapplied())
					fmt.Printf("last migration group: %s\n", ms.LastGroup())
					return nil
				},
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
