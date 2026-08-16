package main

import (
	"context"
	"log"

	"github.com/lovely-eye/server/internal/platform/config"
	"github.com/lovely-eye/server/internal/seed"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	result, err := seed.Run(context.Background(), cfg)
	if err != nil {
		log.Fatalf("Failed to load example data: %v", err)
	}
	if result.CreatedSite {
		log.Printf(
			"Created site %q for localhost with public key %s",
			result.SiteName,
			result.PublicKey,
		)
	}
	log.Printf(
		"Seeded %d clients, %d sessions, %d page views, %d predefined events",
		result.Clients,
		result.Sessions,
		result.PageViews,
		result.PredefinedEvents,
	)
}
