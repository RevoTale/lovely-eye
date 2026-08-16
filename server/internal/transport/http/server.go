package http

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/lovely-eye/server/internal/analytics"
	"github.com/lovely-eye/server/internal/auth"
	"github.com/lovely-eye/server/internal/country"
	"github.com/lovely-eye/server/internal/dashboard"
	"github.com/lovely-eye/server/internal/event"
	"github.com/lovely-eye/server/internal/graph"
	"github.com/lovely-eye/server/internal/platform/config"
	"github.com/lovely-eye/server/internal/site"
	"github.com/lovely-eye/server/internal/transport/http/clientip"
	"github.com/lovely-eye/server/internal/transport/http/collect"
	transportmiddleware "github.com/lovely-eye/server/internal/transport/http/middleware"
	"github.com/uptrace/bun"
)

type Server struct {
	Handler    http.Handler
	HTTPServer *http.Server
}

type Services struct {
	Auth            *auth.Service
	AuthCookies     *CookieManager
	Site            *site.Service
	Analytics       *analytics.Service
	Country         *country.Service
	EventDefinition *event.Service
}

type Options struct {
	Config     config.Config
	Database   *bun.DB
	TrackerJS  []byte
	Services   Services
	IPResolver *clientip.Resolver
}

func New(options Options) *Server {
	handler := buildHTTPHandler(
		options.Config,
		options.Database,
		options.TrackerJS,
		options.Services,
		options.IPResolver,
	)
	addr := options.Config.Server.Host + ":" + options.Config.Server.Port
	httpServer := newHTTPServer(addr, handler)

	return &Server{
		Handler:    handler,
		HTTPServer: httpServer,
	}
}

func buildHTTPHandler(
	cfg config.Config,
	db *bun.DB,
	trackerJS []byte,
	deps Services,
	ipResolver *clientip.Resolver,
) http.Handler {
	collectRateLimiter := collect.NewRateLimiter(
		cfg.Analytics.RateLimitEnabled,
		cfg.Analytics.RateLimitPerMinute,
		cfg.Analytics.RateLimitBurst,
	)
	analyticsHandler := collect.NewAnalyticsHandler(
		deps.Analytics,
		deps.Site,
		collect.AnalyticsHandlerConfig{
			MaxBodyBytes:       cfg.Analytics.MaxBodyBytes,
			MaxPropertiesBytes: cfg.Analytics.MaxPropertiesBytes,
		},
		ipResolver,
		collectRateLimiter,
	)

	resolver := graph.NewResolver(
		deps.Auth,
		deps.AuthCookies,
		deps.Site,
		deps.Analytics,
		deps.Country,
		deps.EventDefinition,
		graph.DashboardLimits{
			MaxDailyRangeDays:     cfg.Dashboard.MaxDailyRangeDays,
			MaxHourlyRangeDays:    cfg.Dashboard.MaxHourlyRangeDays,
			MaxFilterValues:       cfg.Dashboard.MaxFilterValues,
			MaxFilterStringLength: cfg.Dashboard.MaxFilterStringLength,
		},
	)

	authMiddleware := newAuthMiddleware(deps.Auth, deps.AuthCookies)
	mux := http.NewServeMux()
	basePath := cfg.Server.BasePath
	if basePath == "/" {
		basePath = ""
	}
	mux.HandleFunc("POST "+basePath+"/api/collect", analyticsHandler.Collect)
	mux.HandleFunc("OPTIONS "+basePath+"/api/collect", analyticsHandler.Collect)

	authRateLimiter := transportmiddleware.NewAuthRateLimiter(
		cfg.Auth.RateLimitEnabled,
		cfg.Auth.RateLimitAttempts,
		cfg.Auth.RateLimitWindow,
		cfg.GraphQL.MaxBodyBytes,
		ipResolver,
	)
	graphqlHandler := authRateLimiter.Middleware(graph.Handler(
		resolver,
		cfg.GraphQL.MaxBodyBytes,
		cfg.GraphQL.MaxComplexity,
	))
	graphqlPath := basePath + "/graphql"
	mux.Handle("POST "+graphqlPath, graphqlHandler)
	mux.HandleFunc("GET "+graphqlPath, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("GET "+basePath+"/tracker.js", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if _, err := w.Write(trackerJS); err != nil {
			slog.Error("failed to write tracker.js", "error", err)
		}
	})
	hh := newHealthHandler(db, cfg.Server.DashboardPath, cfg.Database.ConnectTimeout)

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

	return transportmiddleware.Logging(
		transportmiddleware.Security(
			transportmiddleware.CORS(
				authMiddleware.authenticate(mux),
			),
		),
	)
}

func newHTTPServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
}
