package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:  "migrator",
		Usage: "development tool for preparing database migrations using Atlas",
		Commands: []*cli.Command{
			{
				Name:  "diff",
				Usage: "auto-generate migration files for SQLite and PostgreSQL by comparing current schema with target",
				Action: func(ctx context.Context, c *cli.Command) error {
					reader := bufio.NewReader(os.Stdin)
					fmt.Print("Migration name: ")
					name, err := reader.ReadString('\n')
					if err != nil {
						return fmt.Errorf("failed to read input: %w", err)
					}
					name = strings.TrimSpace(name)
					if name == "" {
						return fmt.Errorf("migration name is required")
					}

					// Generate SQLite migration
					fmt.Println("\nGenerating SQLite migration...")
					cmd := exec.Command("atlas", "migrate", "diff", name, "--env", "sqlite")
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					if err := cmd.Run(); err != nil {
						return fmt.Errorf("atlas diff failed for sqlite: %w", err)
					}

					// Generate PostgreSQL migration
					fmt.Println("\nGenerating PostgreSQL migration...")
					cmd = exec.Command("atlas", "migrate", "diff", name, "--env", "postgres", "--dev-url", "sqlite://file?mode=memory")
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					if err := cmd.Run(); err != nil {
						return fmt.Errorf("atlas diff failed for postgres: %w", err)
					}

					fmt.Println("\nMigrations generated successfully!")
					return nil
				},
			},
			{
				Name:  "hash",
				Usage: "recalculate and update atlas.sum hash files after manually editing migrations",
				Action: func(ctx context.Context, c *cli.Command) error {
					// Hash SQLite migrations
					fmt.Println("Hashing SQLite migrations...")
					cmd := exec.Command("atlas", "migrate", "hash", "--env", "sqlite")
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					if err := cmd.Run(); err != nil {
						return fmt.Errorf("atlas hash failed for sqlite: %w", err)
					}

					// Hash PostgreSQL migrations
					fmt.Println("Hashing PostgreSQL migrations...")
					cmd = exec.Command("atlas", "migrate", "hash", "--env", "postgres")
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					if err := cmd.Run(); err != nil {
						return fmt.Errorf("atlas hash failed for postgres: %w", err)
					}

					fmt.Println("\nHash files updated successfully!")
					return nil
				},
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
