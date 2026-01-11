package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

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
	defer database.Close(db)

	migs, err := migrations.NewMigrations()
	if err != nil {
		log.Fatalf("Failed to load migrations: %v", err)
	}

	migrator := migrate.NewMigrator(db, migs)

	app := &cli.Command{
		Name:  "bun",
		Usage: "database migrations",
		Commands: []*cli.Command{
			{
				Name:  "db",
				Usage: "database migrations",
				Commands: []*cli.Command{
					{
						Name:  "init",
						Usage: "create migration tables",
						Action: func(ctx context.Context, c *cli.Command) error {
							return migrator.Init(ctx)
						},
					},
					{
						Name:  "migrate",
						Usage: "migrate database",
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
						Name:  "rollback",
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
						Name:  "lock",
						Usage: "lock migrations",
						Action: func(ctx context.Context, c *cli.Command) error {
							return migrator.Lock(ctx)
						},
					},
					{
						Name:  "unlock",
						Usage: "unlock migrations",
						Action: func(ctx context.Context, c *cli.Command) error {
							return migrator.Unlock(ctx)
						},
					},
					{
						Name:  "create_go",
						Usage: "create Go migration",
						Action: func(ctx context.Context, c *cli.Command) error {
							name := strings.Join(c.Args().Slice(), "_")
							mf, err := migrator.CreateGoMigration(ctx, name)
							if err != nil {
								return err
							}
							fmt.Printf("created migration %s (%s)\n", mf.Name, mf.Path)
							return nil
						},
					},
					{
						Name:  "create_sql",
						Usage: "create up and down SQL migrations",
						Action: func(ctx context.Context, c *cli.Command) error {
							name := strings.Join(c.Args().Slice(), "_")
							files, err := migrator.CreateSQLMigrations(ctx, name)
							if err != nil {
								return err
							}

							for _, mf := range files {
								fmt.Printf("created migration %s (%s)\n", mf.Name, mf.Path)
							}

							return nil
						},
					},
					{
						Name:  "create_tx_sql",
						Usage: "create up and down transactional SQL migrations",
						Action: func(ctx context.Context, c *cli.Command) error {
							name := strings.Join(c.Args().Slice(), "_")
							files, err := migrator.CreateTxSQLMigrations(ctx, name)
							if err != nil {
								return err
							}

							for _, mf := range files {
								fmt.Printf("created transaction migration %s (%s)\n", mf.Name, mf.Path)
							}

							return nil
						},
					},
					{
						Name:  "status",
						Usage: "print migrations status",
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
					{
						Name:  "mark_applied",
						Usage: "mark migrations as applied without actually running them",
						Action: func(ctx context.Context, c *cli.Command) error {
							group, err := migrator.Migrate(ctx, migrate.WithNopMigration())
							if err != nil {
								return err
							}
							if group.IsZero() {
								fmt.Printf("there are no new migrations to mark as applied\n")
								return nil
							}
							fmt.Printf("marked as applied %s\n", group)
							return nil
						},
					},
				},
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
