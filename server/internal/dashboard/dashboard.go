package dashboard

import (
	"html/template"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Config holds the runtime configuration for the dashboard
type Config struct {
	BasePath      string
	APIUrl        string
	GraphQLUrl    string
	DashboardPath string // Path to dashboard files (defaults to "dashboard")
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		BasePath:      "/",
		APIUrl:        "/api",
		GraphQLUrl:    "/graphql",
		DashboardPath: "dashboard",
	}
}

// configJSTemplate is the template for config.js
const configJSTemplate = `// Runtime configuration - injected by server
window.__ENV__ = {
  BASE_PATH: '{{.BasePath}}',
  API_URL: '{{.APIUrl}}',
  GRAPHQL_URL: '{{.GraphQLUrl}}',
};
`

// Handler returns an http.Handler that serves the dashboard with the given config
func Handler(cfg Config) http.Handler {
	// Set default dashboard path if not provided
	if cfg.DashboardPath == "" {
		cfg.DashboardPath = "dashboard"
	}

	// Normalize base path - keep empty string for root, or ensure it starts with / and has no trailing /
	if cfg.BasePath != "" {
		if !strings.HasPrefix(cfg.BasePath, "/") {
			cfg.BasePath = "/" + cfg.BasePath
		}
		cfg.BasePath = strings.TrimSuffix(cfg.BasePath, "/")
	}

	// Parse config.js template
	tmpl := template.Must(template.New("config.js").Parse(configJSTemplate))

	// Read index.html and prepare it with base path replacements
	indexPath := filepath.Join(cfg.DashboardPath, "index.html")
	indexHTML, err := os.ReadFile(indexPath) // #nosec G304 -- indexPath is constructed from validated DashboardPath config
	if err != nil {
		// Return a minimal handler that serves a placeholder
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Dashboard not available", http.StatusServiceUnavailable)
		})
	}
	// Replace {{BASE_PATH}} placeholder with actual base path
	indexContent := string(indexHTML)
	if cfg.BasePath == "" {
		// Remove base tag entirely when serving at root (base tag can interfere with SPA routing)
		indexContent = strings.ReplaceAll(indexContent, `<base href="{{BASE_PATH}}/" />`, "")
		indexContent = strings.ReplaceAll(indexContent, "{{BASE_PATH}}", "")
	} else {
		// Keep base tag for subpaths
		indexContent = strings.ReplaceAll(indexContent, "{{BASE_PATH}}", cfg.BasePath)
	}

	// Create file server for static files
	fileServer := http.FileServer(http.Dir(cfg.DashboardPath))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serve dynamically generated config.js
		if r.URL.Path == "/config.js" || r.URL.Path == "config.js" {
			w.Header().Set("Content-Type", "application/javascript")
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			if err := tmpl.Execute(w, cfg); err != nil {
				http.Error(w, "Failed to generate config", http.StatusInternalServerError)
			}
			return
		}

		// Check if the requested file exists
		urlPath := strings.TrimPrefix(r.URL.Path, "/")
		if urlPath == "" {
			urlPath = "index.html"
		}

		// Try to check if file exists on disk
		filePath := filepath.Join(cfg.DashboardPath, filepath.Clean(urlPath))
		if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
			// Special handling for index.html - serve with BASE_PATH replacements
			if urlPath == "index.html" {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
				if _, err := w.Write([]byte(indexContent)); err != nil {
					http.Error(w, "internal error", http.StatusInternalServerError)
				}
				return
			}
			// File exists, serve it normally
			fileServer.ServeHTTP(w, r)
			return
		}

		// For SPA routing: if file doesn't exist and it's not a static asset,
		// serve index.html to let the client-side router handle it
		ext := path.Ext(urlPath)
		if ext == "" || ext == ".html" {
			// Serve processed index.html for SPA routes
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			if _, err := w.Write([]byte(indexContent)); err != nil {
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
			return
		}

		// File not found for other extensions
		http.NotFound(w, r)
	})
}
