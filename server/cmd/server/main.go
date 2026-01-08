package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lovely-eye/server/internal/auth"
	"github.com/lovely-eye/server/internal/config"
	"github.com/lovely-eye/server/internal/database"
	"github.com/lovely-eye/server/internal/graph"
	"github.com/lovely-eye/server/internal/handlers"
	"github.com/lovely-eye/server/internal/middleware"
	"github.com/lovely-eye/server/internal/repository"
	"github.com/lovely-eye/server/internal/services"
)

// REST API endpoints:
//   POST /api/collect  - Track page views (public)
//   POST /api/event    - Track custom events (public)
//
// All other operations (auth, sites, dashboard) are via GraphQL at /graphql

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.New(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)

	// Run migrations
	if err := database.Migrate(context.Background(), db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Database migrations completed")

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
		log.Fatalf("Failed to create initial admin: %v", err)
	}

	// Initialize handlers for tracking (REST API)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService, siteService)

	// Initialize auth middleware
	authMiddleware := auth.NewMiddleware(authService)

	// Setup GraphQL resolver (handles all dashboard/auth/site operations)
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

	// Create server
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on %s", addr)
		log.Printf("REST API available at http://%s/api", addr)
		log.Printf("GraphQL endpoint at http://%s/graphql", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
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
