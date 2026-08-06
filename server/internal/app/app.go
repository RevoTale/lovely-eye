package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/lovely-eye/server/internal/analytics"
	"github.com/lovely-eye/server/internal/auth"
	"github.com/lovely-eye/server/internal/platform/config"
	"github.com/lovely-eye/server/internal/platform/database"
	"github.com/lovely-eye/server/internal/site"
	transporthttp "github.com/lovely-eye/server/internal/transport/http"
	"github.com/lovely-eye/server/internal/transport/http/clientip"
	"github.com/uptrace/bun"
)

const shutdownTimeout = 30 * time.Second

type App struct {
	DB               *bun.DB
	AuthService      *auth.Service
	SiteService      *site.Service
	AnalyticsService *analytics.Service
	Handler          http.Handler
	HTTPServer       *http.Server
	basePath         string
}

func New(ctx context.Context, cfg config.Config) (_ *App, err error) {
	trackerJS, err := loadTrackerJS(cfg)
	if err != nil {
		return nil, err
	}

	db, err := openMigratedDatabase(ctx, cfg)
	if err != nil {
		return nil, err
	}
	application := &App{DB: db}
	defer func() {
		if err == nil {
			return
		}
		if closeErr := application.Close(); closeErr != nil {
			slog.Error("failed to close application after construction error", "error", closeErr)
		}
	}()

	features, err := buildServices(ctx, cfg, db)
	application.AuthService = features.Auth
	application.SiteService = features.Site
	application.AnalyticsService = features.Analytics
	if err != nil {
		return nil, err
	}

	ipResolver, err := clientip.NewResolver(cfg.Analytics.TrustedProxyCIDRs)
	if err != nil {
		return nil, fmt.Errorf("configure trusted proxies: %w", err)
	}
	transport := transporthttp.New(transporthttp.Options{
		Config:     cfg,
		Database:   db,
		TrackerJS:  trackerJS,
		Services:   features,
		IPResolver: ipResolver,
	})
	application.Handler = transport.Handler
	application.HTTPServer = transport.HTTPServer
	application.basePath = cfg.Server.BasePath
	return application, nil
}

func (a *App) Run(ctx context.Context) error {
	if a.HTTPServer == nil {
		return errors.New("app: http server is not configured")
	}

	addr := a.HTTPServer.Addr
	slog.Info("server starting", "address", addr, "base_path", a.basePath)
	serveErrors := make(chan error, 1)
	go func() {
		err := a.HTTPServer.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		serveErrors <- err
	}()

	select {
	case err := <-serveErrors:
		if err != nil {
			return fmt.Errorf("serve http: %w", err)
		}
		return nil
	case <-ctx.Done():
		slog.Info("server shutting down")
	}

	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), shutdownTimeout)
	defer cancel()
	if err := a.HTTPServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shut down http server: %w", err)
	}
	if err := <-serveErrors; err != nil {
		return fmt.Errorf("serve http during shutdown: %w", err)
	}
	slog.Info("server stopped")
	return nil
}

func (a *App) Close() error {
	var err error
	if a.AnalyticsService != nil {
		if closeErr := a.AnalyticsService.Close(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close analytics service: %w", closeErr))
		}
	}
	if a.DB != nil {
		if closeErr := database.Close(a.DB); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close database: %w", closeErr))
		}
	}
	return err
}
