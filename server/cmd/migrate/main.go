package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/lovely-eye/server/internal/config"
	"github.com/lovely-eye/server/internal/database"
)

func main() {
	flag.Parse()

	cmd := flag.Arg(0)
	if cmd == "" {
		cmd = "up"
	}

	cfg := config.Load()

	db, err := database.New(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)

	ctx := context.Background()

	switch cmd {
	case "up":
		if err := database.Migrate(ctx, db); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
	case "down":
		if err := database.Rollback(ctx, db); err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
	case "status":
		if err := database.MigrationStatus(ctx, db); err != nil {
			log.Fatalf("Failed to get status: %v", err)
		}
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
		fmt.Println("Available commands: up, down, status")
		os.Exit(1)
	}
}
