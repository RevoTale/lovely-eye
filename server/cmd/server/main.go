package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lovely-eye/server/internal/config"
	"github.com/lovely-eye/server/internal/server"
)

func main() {
	cfg := config.Load()

	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.LogLevel,
	})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	srv, err := server.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	defer func() {
		err := srv.Close()
		if nil != err {
			slog.Error("server close failed", "error", err)
		}
	}()

	log.Println("Database migrations completed")

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
