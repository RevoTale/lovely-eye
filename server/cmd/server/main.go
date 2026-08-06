package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/lovely-eye/server/internal/app"
	"github.com/lovely-eye/server/internal/platform/config"
)

func main() {
	if err := run(); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}

func run() (err error) {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.LogLevel,
	})
	slog.SetDefault(slog.New(logHandler))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	application, err := app.New(ctx, cfg)
	if err != nil {
		return fmt.Errorf("create application: %w", err)
	}
	defer func() {
		err = errors.Join(err, application.Close())
	}()

	if err := application.Run(ctx); err != nil {
		return fmt.Errorf("run application: %w", err)
	}
	return nil
}
