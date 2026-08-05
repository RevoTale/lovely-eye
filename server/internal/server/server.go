package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/lovely-eye/server/internal/auth"
	"github.com/lovely-eye/server/internal/config"
	"github.com/lovely-eye/server/internal/dashboard"
	"github.com/lovely-eye/server/internal/database"
	"github.com/lovely-eye/server/internal/graph"
	"github.com/lovely-eye/server/internal/handlers"
	"github.com/lovely-eye/server/internal/middleware"
	"github.com/lovely-eye/server/internal/repository"
	"github.com/lovely-eye/server/internal/services"
	"github.com/lovely-eye/server/pkg/clientip"
	"github.com/uptrace/bun"
)

type Server struct {
	DB               *bun.DB
	AuthService      auth.Service
	SiteService      *services.SiteService
	AnalyticsService *services.AnalyticsService
	Handler          http.Handler
	HTTPServer       *http.Server
	trackerJS        []byte
}

type serverDependencies struct {
	authService            auth.Service
	siteService            *services.SiteService
	analyticsService       *services.AnalyticsService
	countryService         *services.CountryService
	eventDefinitionService *services.EventDefinitionService
}

func New(cfg config.Config) (*Server, error) {
	trackerJS, err := loadTrackerJS(cfg)
	if err != nil {
		return nil, err
	}

	db, err := openMigratedDatabase(cfg)
	if err != nil {
		return nil, err
	}
	closeOnError := true
	defer func() {
		if closeOnError {
			if closeErr := database.Close(db); closeErr != nil {
				slog.Error("failed to close database after server init error", "error", closeErr)
			}
		}
	}()

	deps, err := buildServerDependencies(cfg, db)
	if err != nil {
		return nil, err
	}

	ipResolver, err := clientip.NewResolver(cfg.Analytics.TrustedProxyCIDRs)
	if err != nil {
		return nil, fmt.Errorf("configure trusted proxies: %w", err)
	}
	handler := buildHTTPHandler(cfg, db, trackerJS, deps, ipResolver)
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	httpServer := newHTTPServer(addr, handler)
	closeOnError = false

	return &Server{
		DB:               db,
		AuthService:      deps.authService,
		SiteService:      deps.siteService,
		AnalyticsService: deps.analyticsService,
		Handler:          handler,
		HTTPServer:       httpServer,
		trackerJS:        trackerJS,
	}, nil
}

func loadTrackerJS(cfg config.Config) ([]byte, error) {
	if len(cfg.TrackerJS) != 0 {
		return cfg.TrackerJS, nil
	}
	trackerPath := filepath.Join("static", "tracker.js")
	trackerJS, err := os.ReadFile(trackerPath) // #nosec G304 -- trackerPath is constructed from static directory constant
	if err != nil {
		return nil, fmt.Errorf("failed to load tracker.js: %w", err)
	}
	return trackerJS, nil
}

func openMigratedDatabase(cfg config.Config) (*bun.DB, error) {
	db, err := database.New(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("create database connection: %w", err)
	}
	if err := database.Migrate(context.Background(), db); err != nil {
		if closeErr := database.Close(db); closeErr != nil {
			slog.Error("failed to close database after migration error", "error", closeErr)
		}
		return nil, fmt.Errorf("run migrations: %w", err)
	}
	return db, nil
}

func buildServerDependencies(cfg config.Config, db *bun.DB) (serverDependencies, error) {
	userRepo := repository.NewUserRepository(db)
	siteRepo := repository.NewSiteRepository(db)
	analyticsRepo := repository.NewAnalyticsRepository(db)
	countryRepo := repository.NewCountryRepository(db)
	eventDefinitionRepo := repository.NewEventDefinitionRepository(db)
	authService := auth.NewService(userRepo, authConfig(cfg))
	geoIPService := services.NewGeoIPService(services.GeoIPConfig{
		DBPath:            cfg.GeoIP.DBPath,
		DownloadURL:       cfg.GeoIP.DownloadURL,
		MaxMindLicenseKey: cfg.GeoIP.MaxMindLicenseKey,
	})
	siteService := services.NewSiteService(siteRepo)
	countryService := services.NewCountryService(countryRepo, geoIPService)
	analyticsService := services.NewAnalyticsService(
		analyticsRepo,
		siteRepo,
		eventDefinitionRepo,
		geoIPService,
		countryService,
		analyticsIdentitySecret(cfg),
	)
	analyticsService.SetMaxSinglePageDuration(cfg.Analytics.MaxSinglePageDuration)
	if err := analyticsService.SyncGeoIPRequirement(context.Background()); err != nil {
		fmt.Printf("Warning: GeoIP database sync failed: %v\n", err)
	}
	if err := authService.CreateInitialAdmin(context.Background(), cfg.Auth.InitialAdminUsername, cfg.Auth.InitialAdminPassword); err != nil {
		return serverDependencies{}, fmt.Errorf("create initial admin: %w", err)
	}
	return serverDependencies{
		authService:            authService,
		siteService:            siteService,
		analyticsService:       analyticsService,
		countryService:         countryService,
		eventDefinitionService: services.NewEventDefinitionService(eventDefinitionRepo),
	}, nil
}

func authConfig(cfg config.Config) auth.Config {
	return auth.Config{
		JWTSecret:         cfg.Auth.JWTSecret,
		AccessTokenExpiry: cfg.Auth.AccessTokenExpiry,
		RefreshExpiry:     cfg.Auth.RefreshExpiry,
		AllowRegistration: cfg.Auth.AllowRegistration,
		SecureCookies:     cfg.Auth.SecureCookies,
		CookieDomain:      cfg.Auth.CookieDomain,
	}
}

func analyticsIdentitySecret(cfg config.Config) string {
	if cfg.Analytics.IdentitySecret != "" {
		return cfg.Analytics.IdentitySecret
	}
	return cfg.Auth.JWTSecret
}

func buildHTTPHandler(
	cfg config.Config,
	db *bun.DB,
	trackerJS []byte,
	deps serverDependencies,
	ipResolver *clientip.Resolver,
) http.Handler {
	collectRateLimiter := handlers.NewCollectRateLimiter(
		cfg.Analytics.RateLimitEnabled,
		cfg.Analytics.RateLimitPerMinute,
		cfg.Analytics.RateLimitBurst,
	)
	analyticsHandler := handlers.NewAnalyticsHandler(
		deps.analyticsService,
		deps.siteService,
		handlers.AnalyticsHandlerConfig{
			MaxBodyBytes:       cfg.Analytics.MaxBodyBytes,
			MaxPropertiesBytes: cfg.Analytics.MaxPropertiesBytes,
		},
		ipResolver,
		collectRateLimiter,
	)

	resolver := graph.NewResolver(
		deps.authService,
		deps.siteService,
		deps.analyticsService,
		deps.countryService,
		deps.eventDefinitionService,
		graph.DashboardLimits{
			MaxDailyRangeDays:     cfg.Dashboard.MaxDailyRangeDays,
			MaxHourlyRangeDays:    cfg.Dashboard.MaxHourlyRangeDays,
			MaxFilterValues:       cfg.Dashboard.MaxFilterValues,
			MaxFilterStringLength: cfg.Dashboard.MaxFilterStringLength,
		},
	)

	authMiddleware := auth.NewMiddleware(deps.authService)
	mux := http.NewServeMux()
	basePath := cfg.Server.BasePath
	if basePath == "/" {
		basePath = ""
	}
	mux.HandleFunc("POST "+basePath+"/api/collect", analyticsHandler.Collect)
	mux.HandleFunc("OPTIONS "+basePath+"/api/collect", analyticsHandler.Collect)

	authRateLimiter := middleware.NewAuthRateLimiter(
		cfg.Auth.RateLimitEnabled,
		cfg.Auth.RateLimitAttempts,
		cfg.Auth.RateLimitWindow,
		cfg.GraphQL.MaxBodyBytes,
		ipResolver,
	)
	graphqlHandler := authRateLimiter.Middleware(graph.Handler(resolver, cfg.GraphQL.MaxBodyBytes))
	graphqlPath := basePath + "/graphql"
	mux.Handle("POST "+graphqlPath, graphqlHandler)
	mux.HandleFunc("GET "+graphqlPath, graphqlPlaygroundHandler(graphqlPath))

	mux.HandleFunc("GET "+basePath+"/tracker.js", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if _, err := w.Write(trackerJS); err != nil {
			slog.Error("failed to write tracker.js", "error", err)
		}
	})
	hh := handlers.NewHealthHandler(db, cfg.Server.DashboardPath, cfg.Database.ConnectTimeout)

	mux.Handle("GET /health", hh)

	dashboardCfg := dashboard.Config{
		BasePath:      basePath,
		APIUrl:        basePath + "/api",
		GraphQLUrl:    basePath + "/graphql",
		DashboardPath: cfg.Server.DashboardPath,
	}
	dashboardHandler := dashboard.Handler(dashboardCfg)

	if basePath == "" {
		mux.Handle("GET /", dashboardHandler)
	} else {
		mux.Handle("GET "+basePath+"/", http.StripPrefix(basePath, dashboardHandler))
		mux.Handle("GET "+basePath, http.RedirectHandler(basePath+"/", http.StatusMovedPermanently))
	}

	return middleware.Logging(
		middleware.Security(
			middleware.CORS(
				authMiddleware.Authenticate(mux),
			),
		),
	)
}

func newHTTPServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func (s *Server) Close() error {
	var err error
	if s.AnalyticsService != nil {
		if closeErr := s.AnalyticsService.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}
	if closeErr := database.Close(s.DB); closeErr != nil {
		err = errors.Join(err, fmt.Errorf("close database: %w", closeErr))
	}
	return err
}

func graphqlPlaygroundHandler(graphqlPath string) http.HandlerFunc {
	graphqlPathJSON, err := json.Marshal(graphqlPath)
	if err != nil {
		slog.Error("failed to marshal graphql path", "error", err)
		graphqlPathJSON = []byte(`"/graphql"`)
	}

	html := []byte(`<!DOCTYPE html>
<html>
<head>
  <title>Lovely Eye GraphQL Playground</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/graphiql@3/graphiql.min.css" />
</head>
<body style="margin: 0;">
  <div id="graphiql" style="height: 100vh;"></div>
  <script crossorigin src="https://unpkg.com/react@18/umd/react.production.min.js"></script>
  <script crossorigin src="https://unpkg.com/react-dom@18/umd/react-dom.production.min.js"></script>
  <script crossorigin src="https://cdn.jsdelivr.net/npm/graphiql@3/graphiql.min.js"></script>
  <script>
    const fetcher = GraphiQL.createFetcher({ url: ` + string(graphqlPathJSON) + ` });
    ReactDOM.createRoot(document.getElementById('graphiql')).render(
      React.createElement(GraphiQL, { fetcher: fetcher })
    );
  </script>
</body>
</html>`)

	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if _, err := w.Write(html); err != nil {
			slog.Error("failed to write graphql playground", "error", err)
		}
	}
}
