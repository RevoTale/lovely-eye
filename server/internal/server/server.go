package server

import (
	"context"
	"net/http"
	"time"

	"github.com/lovely-eye/server/internal/auth"
	"github.com/lovely-eye/server/internal/config"
	"github.com/lovely-eye/server/internal/database"
	"github.com/lovely-eye/server/internal/graph"
	"github.com/lovely-eye/server/internal/handlers"
	"github.com/lovely-eye/server/internal/middleware"
	"github.com/lovely-eye/server/internal/repository"
	"github.com/lovely-eye/server/internal/services"
	"github.com/uptrace/bun"
)

// Server holds all server dependencies.
type Server struct {
	DB               *bun.DB
	AuthService      auth.Service
	SiteService      *services.SiteService
	AnalyticsService *services.AnalyticsService
	Handler          http.Handler
	HTTPServer       *http.Server
}

// New creates a new Server from config.
func New(cfg *config.Config) (*Server, error) {
	// Initialize database
	db, err := database.New(&cfg.Database)
	if err != nil {
		return nil, err
	}

	// Run migrations
	if err := database.Migrate(context.Background(), db); err != nil {
		database.Close(db)
		return nil, err
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	siteRepo := repository.NewSiteRepository(db)
	analyticsRepo := repository.NewAnalyticsRepository(db)

	// Initialize auth service
	authService := auth.NewService(userRepo, auth.Config{
		JWTSecret:         cfg.Auth.JWTSecret,
		AccessTokenExpiry: cfg.Auth.AccessTokenExpiry,
		RefreshExpiry:     cfg.Auth.RefreshExpiry,
		AllowRegistration: cfg.Auth.AllowRegistration,
		SecureCookies:     cfg.Auth.SecureCookies,
		CookieDomain:      cfg.Auth.CookieDomain,
	})

	// Initialize other services
	siteService := services.NewSiteService(siteRepo)
	analyticsService := services.NewAnalyticsService(analyticsRepo, siteRepo)

	// Create initial admin from env vars if configured
	if err := authService.CreateInitialAdmin(context.Background(), cfg.Auth.InitialAdminUsername, cfg.Auth.InitialAdminPassword); err != nil {
		database.Close(db)
		return nil, err
	}

	// Initialize handlers for tracking (REST API)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService, siteService)

	// Initialize auth middleware
	authMiddleware := auth.NewMiddleware(authService)

	// Setup GraphQL resolver
	resolver := graph.NewResolver(authService, siteService, analyticsService)

	// Setup HTTP router
	mux := http.NewServeMux()

	// REST API: Only tracking endpoints (public, no auth required)
	mux.HandleFunc("POST /api/collect", analyticsHandler.Collect)
	mux.HandleFunc("POST /api/event", analyticsHandler.Event)

	// GraphQL endpoint (CSRF protected for cookie-based auth)
	mux.Handle("POST /graphql", authMiddleware.RequireCSRF(http.HandlerFunc(graph.Handler(resolver))))
	mux.HandleFunc("GET /graphql", graphqlPlaygroundHandler)

	// Serve tracking script
	mux.HandleFunc("GET /tracker.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		http.ServeFile(w, r, "static/tracker.js")
	})

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Apply global middleware
	handler := middleware.Logging(
		middleware.CORS(
			authMiddleware.Authenticate(mux),
		),
	)

	// Create HTTP server
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	httpServer := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		DB:               db,
		AuthService:      authService,
		SiteService:      siteService,
		AnalyticsService: analyticsService,
		Handler:          handler,
		HTTPServer:       httpServer,
	}, nil
}

// Close closes the server and database connection.
func (s *Server) Close() error {
	return database.Close(s.DB)
}

func graphqlPlaygroundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<!DOCTYPE html>
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
    const fetcher = GraphiQL.createFetcher({ url: '/graphql' });
    ReactDOM.createRoot(document.getElementById('graphiql')).render(
      React.createElement(GraphiQL, { fetcher: fetcher })
    );
  </script>
</body>
</html>`))
}
