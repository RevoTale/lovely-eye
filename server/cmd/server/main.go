package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lovely-eye/server/internal/config"
	"github.com/lovely-eye/server/internal/server"
)

// REST API endpoints:
//   POST /api/collect  - Track page views (public)
//   POST /api/event    - Track custom events (public)
//
// All other operations (auth, sites, dashboard) are via GraphQL at /graphql

func main() {
	cfg := config.Load()

	srv, err := server.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	defer srv.Close()

	log.Println("Database migrations completed")

	// Start server in goroutine
	go func() {
		addr := srv.HTTPServer.Addr
		basePath := cfg.Server.BasePath
		log.Printf("Server starting on %s", addr)
		log.Printf("Dashboard available at http://%s%s", addr, basePath)
		log.Printf("REST API available at http://%s%s/api", addr, basePath)
		log.Printf("GraphQL endpoint at http://%s%s/graphql", addr, basePath)
		if err := srv.HTTPServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.HTTPServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
