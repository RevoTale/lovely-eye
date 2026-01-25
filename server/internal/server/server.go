package server

import (
	"context"
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

func New(cfg *config.Config) (*Server, error) {

	var trackerJS []byte
	if len(cfg.TrackerJS) == 0 {
		trackerPath := filepath.Join("static", "tracker.js")
		var err error
		trackerJS, err = os.ReadFile(trackerPath) // #nosec G304 -- trackerPath is constructed from static directory constant
		if err != nil {
			return nil, fmt.Errorf("failed to load tracker.js: %w", err)
		}
	} else {
		trackerJS = cfg.TrackerJS
	}

	db, err := database.New(&cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("create database connection: %w", err)
	}

	if err := database.Migrate(context.Background(), db); err != nil {
		if closeErr := database.Close(db); closeErr != nil {
			slog.Error("failed to close database after migration error", "error", closeErr)
		}
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	userRepo := repository.NewUserRepository(db)
	siteRepo := repository.NewSiteRepository(db)
	analyticsRepo := repository.NewAnalyticsRepository(db)
	eventDefinitionRepo := repository.NewEventDefinitionRepository(db)

	// Initialize auth service
	authService := auth.NewService(userRepo, auth.Config{
		JWTSecret:         cfg.Auth.JWTSecret,
		AccessTokenExpiry: cfg.Auth.AccessTokenExpiry,
		RefreshExpiry:     cfg.Auth.RefreshExpiry,
		AllowRegistration: cfg.Auth.AllowRegistration,
		SecureCookies:     cfg.Auth.SecureCookies,
		CookieDomain:      cfg.Auth.CookieDomain,
	})

	geoIPService, err := services.NewGeoIPService(services.GeoIPConfig{
		DBPath:            cfg.GeoIPDBPath,
		DownloadURL:       cfg.GeoIPDownloadURL,
		MaxMindLicenseKey: cfg.GeoIPMaxMindLicenseKey,
	})
	if err != nil {
		fmt.Printf("Warning: Failed to initialize GeoIP service: %v. Country detection will be disabled.\n", err)
	}

	siteService := services.NewSiteService(siteRepo)
	eventDefinitionService := services.NewEventDefinitionService(eventDefinitionRepo)
	analyticsService := services.NewAnalyticsService(analyticsRepo, siteRepo, eventDefinitionRepo, geoIPService)
	if err := analyticsService.SyncGeoIPRequirement(context.Background()); err != nil {
		fmt.Printf("Warning: GeoIP database sync failed: %v\n", err)
	}

	if err := authService.CreateInitialAdmin(context.Background(), cfg.Auth.InitialAdminUsername, cfg.Auth.InitialAdminPassword); err != nil {
		if closeErr := database.Close(db); closeErr != nil {
			slog.Error("failed to close database after admin creation error", "error", closeErr)
		}
		return nil, fmt.Errorf("create initial admin: %w", err)
	}

	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService, siteService)

	// Initialize auth middleware
	authMiddleware := auth.NewMiddleware(authService)

	resolver := graph.NewResolver(authService, siteService, analyticsService, eventDefinitionService)

	mux := http.NewServeMux()

	basePath := cfg.Server.BasePath
	if basePath == "/" {
		basePath = ""
	}

	// REST API: Only tracking endpoints (public, no auth required)
	mux.HandleFunc("POST "+basePath+"/api/collect", analyticsHandler.Collect)
	mux.HandleFunc("OPTIONS "+basePath+"/api/collect", analyticsHandler.Collect)
	mux.HandleFunc("POST "+basePath+"/api/event", analyticsHandler.Collect)
	mux.HandleFunc("OPTIONS "+basePath+"/api/event", analyticsHandler.Collect)

	// GraphQL endpoint
	// Auth uses JWT in HttpOnly + Secure cookies with SameSite=Strict/Lax
	// No CSRF protection needed - see https://www.reddit.com/r/node/comments/1im7yj0/comment/mc0ylfd/
	graphqlHandler := http.HandlerFunc(graph.Handler(resolver))
	mux.Handle("POST "+basePath+"/graphql", graphqlHandler)
	mux.HandleFunc("GET "+basePath+"/graphql", graphqlPlaygroundHandler)

	mux.HandleFunc("GET "+basePath+"/tracker.js", func(w http.ResponseWriter, r *http.Request) {
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

	handler := middleware.Logging(
		middleware.Security(
			middleware.CORS(
				authMiddleware.Authenticate(mux),
			),
		),
	)

	addr := cfg.Server.Host + ":" + cfg.Server.Port
	httpServer := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	srv := &Server{
		DB:               db,
		AuthService:      authService,
		SiteService:      siteService,
		AnalyticsService: analyticsService,
		Handler:          handler,
		HTTPServer:       httpServer,
		trackerJS:        trackerJS,
	}

	return srv, nil
}

func (s *Server) Close() error {
	if err := database.Close(s.DB); err != nil {
		return fmt.Errorf("close database: %w", err)
	}
	return nil
}

func graphqlPlaygroundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	graphqlPath := r.URL.Path
	html := `<!DOCTYPE html>
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
    const fetcher = GraphiQL.createFetcher({ url: '` + graphqlPath + `' });
    ReactDOM.createRoot(document.getElementById('graphiql')).render(
      React.createElement(GraphiQL, { fetcher: fetcher })
    );
  </script>
</body>
</html>`
	if _, err := w.Write([]byte(html)); err != nil {
		slog.Error("failed to write graphql playground", "error", err)
	}
}
