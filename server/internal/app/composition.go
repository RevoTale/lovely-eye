package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/lovely-eye/server/internal/analytics"
	analyticspersistence "github.com/lovely-eye/server/internal/analytics/persistence"
	"github.com/lovely-eye/server/internal/auth"
	authpersistence "github.com/lovely-eye/server/internal/auth/persistence"
	"github.com/lovely-eye/server/internal/country"
	countrypersistence "github.com/lovely-eye/server/internal/country/persistence"
	"github.com/lovely-eye/server/internal/event"
	eventpersistence "github.com/lovely-eye/server/internal/event/persistence"
	geoipcore "github.com/lovely-eye/server/internal/geoip"
	geoipservice "github.com/lovely-eye/server/internal/geoip/service"
	"github.com/lovely-eye/server/internal/platform/config"
	"github.com/lovely-eye/server/internal/platform/database"
	"github.com/lovely-eye/server/internal/site"
	sitepersistence "github.com/lovely-eye/server/internal/site/persistence"
	transporthttp "github.com/lovely-eye/server/internal/transport/http"
	"github.com/uptrace/bun"
)

func loadTrackerJS(cfg config.Config) ([]byte, error) {
	if len(cfg.TrackerJS) != 0 {
		return cfg.TrackerJS, nil
	}
	trackerPath := filepath.Join("static", "dist", "tracker.js")
	trackerJS, err := os.ReadFile(trackerPath) // #nosec G304 -- path is composed only from static constants
	if err != nil {
		return nil, fmt.Errorf("load tracker.js: %w", err)
	}
	return trackerJS, nil
}

func openMigratedDatabase(ctx context.Context, cfg config.Config) (*bun.DB, error) {
	db, err := database.New(ctx, cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("create database connection: %w", err)
	}
	if err := database.Migrate(ctx, db); err != nil {
		if closeErr := database.Close(db); closeErr != nil {
			slog.Error("failed to close database after migration error", "error", closeErr)
		}
		return nil, fmt.Errorf("run migrations: %w", err)
	}
	return db, nil
}

func buildServices(ctx context.Context, cfg config.Config, db *bun.DB) (transporthttp.Services, error) {
	userRepo := authpersistence.New(db)
	siteRepo := sitepersistence.New(db)
	analyticsRepo := analyticspersistence.New(db)
	countryRepo := countrypersistence.New(db)
	eventDefinitionRepo := eventpersistence.New(db)
	authService := auth.NewService(userRepo, authConfig(cfg))
	geoIPService := geoipservice.NewService(geoipcore.Config{
		DBPath:            cfg.GeoIP.DBPath,
		DownloadURL:       cfg.GeoIP.DownloadURL,
		MaxMindLicenseKey: cfg.GeoIP.MaxMindLicenseKey,
	})
	siteService := site.NewService(siteRepo)
	countryService := country.NewService(countryRepo, geoIPService)
	analyticsService := analytics.NewService(
		analyticsRepo,
		siteRepo,
		eventDefinitionRepo,
		geoIPService,
		countryService,
		analyticsIdentitySecret(cfg),
	)
	analyticsService.SetMaxSinglePageDuration(cfg.Analytics.MaxSinglePageDuration)

	result := transporthttp.Services{
		Auth:            authService,
		AuthCookies:     transporthttp.NewCookieManager(authCookieConfig(cfg)),
		Site:            siteService,
		Analytics:       analyticsService,
		Country:         countryService,
		EventDefinition: event.NewService(eventDefinitionRepo),
	}
	if err := analyticsService.SyncGeoIPRequirement(ctx); err != nil {
		// Country analytics is optional at startup; the retained status keeps the failure actionable in admin UI.
		slog.Warn("geoip synchronization failed; continuing without country analytics", "error", err)
	}
	if err := authService.CreateInitialAdmin(
		ctx,
		cfg.Auth.InitialAdminUsername,
		cfg.Auth.InitialAdminPassword,
	); err != nil {
		return result, fmt.Errorf("create initial admin: %w", err)
	}
	return result, nil
}

func authConfig(cfg config.Config) auth.Config {
	return auth.Config{
		JWTSecret:         cfg.Auth.JWTSecret,
		AccessTokenExpiry: cfg.Auth.AccessTokenExpiry,
		RefreshExpiry:     cfg.Auth.RefreshExpiry,
		AllowRegistration: cfg.Auth.AllowRegistration,
	}
}

func authCookieConfig(cfg config.Config) transporthttp.CookieConfig {
	return transporthttp.CookieConfig{
		AccessTokenExpiry: cfg.Auth.AccessTokenExpiry,
		RefreshExpiry:     cfg.Auth.RefreshExpiry,
		Secure:            cfg.Auth.SecureCookies,
		Domain:            cfg.Auth.CookieDomain,
		BasePath:          cfg.Server.BasePath,
	}
}

func analyticsIdentitySecret(cfg config.Config) string {
	if cfg.Analytics.IdentitySecret != "" {
		return cfg.Analytics.IdentitySecret
	}
	return cfg.Auth.JWTSecret
}
