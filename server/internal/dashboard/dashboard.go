package dashboard

import (
	"html/template"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Config struct {
	BasePath      string
	APIUrl        string
	GraphQLUrl    string
	DashboardPath string
}

func DefaultConfig() Config {
	return Config{
		BasePath:      "/",
		APIUrl:        "/api",
		GraphQLUrl:    "/graphql",
		DashboardPath: "dashboard",
	}
}

const configJSTemplate = `// Runtime configuration - injected by server
window.__ENV__ = {
  BASE_PATH: '{{.BasePath}}',
  API_URL: '{{.APIUrl}}',
  GRAPHQL_URL: '{{.GraphQLUrl}}',
};
`

func Handler(cfg Config) http.Handler {

	if cfg.DashboardPath == "" {
		cfg.DashboardPath = "dashboard"
	}

	if cfg.BasePath != "" {
		if !strings.HasPrefix(cfg.BasePath, "/") {
			cfg.BasePath = "/" + cfg.BasePath
		}
		cfg.BasePath = strings.TrimSuffix(cfg.BasePath, "/")
	}

	tmpl := template.Must(template.New("config.js").Parse(configJSTemplate))

	indexPath := filepath.Join(cfg.DashboardPath, "index.html")
	indexHTML, err := os.ReadFile(indexPath) // #nosec G304 -- indexPath is constructed from validated DashboardPath config
	if err != nil {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Dashboard not available", http.StatusServiceUnavailable)
		})
	}

	indexContent := string(indexHTML)
	if cfg.BasePath == "" {

		indexContent = strings.ReplaceAll(indexContent, `<base href="{{BASE_PATH}}/" />`, "")
		indexContent = strings.ReplaceAll(indexContent, "{{BASE_PATH}}", "")
	} else {

		indexContent = strings.ReplaceAll(indexContent, "{{BASE_PATH}}", cfg.BasePath)
	}

	fileServer := http.FileServer(http.Dir(cfg.DashboardPath))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/config.js" || r.URL.Path == "config.js" {
			w.Header().Set("Content-Type", "application/javascript")
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			if err := tmpl.Execute(w, cfg); err != nil {
				http.Error(w, "Failed to generate config", http.StatusInternalServerError)
			}
			return
		}

		urlPath := strings.TrimPrefix(r.URL.Path, "/")
		if urlPath == "" {
			urlPath = "index.html"
		}

		filePath := filepath.Join(cfg.DashboardPath, filepath.Clean(urlPath))
		if info, err := os.Stat(filePath); err == nil && !info.IsDir() {

			if urlPath == "index.html" {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
				if _, err := w.Write([]byte(indexContent)); err != nil {
					http.Error(w, "internal error", http.StatusInternalServerError)
				}
				return
			}

			fileServer.ServeHTTP(w, r)
			return
		}

		ext := path.Ext(urlPath)
		if ext == "" || ext == ".html" {

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			if _, err := w.Write([]byte(indexContent)); err != nil {
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
			return
		}

		http.NotFound(w, r)
	})
}
