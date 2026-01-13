package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/uptrace/bun"
)
type HealthHandler struct {
	db            *bun.DB
	dashboardPath string
}
func NewHealthHandler(db *bun.DB, dashboardPath string) *HealthHandler {
	return &HealthHandler{
		db:            db,
		dashboardPath: dashboardPath,
	}
}
func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Check database connection
		if err := h.db.PingContext(r.Context()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status":"unhealthy","error":"database connection failed"}`))
			return
		}

		// Check dashboard files exist (skip check if dashboard path is empty for tests)
		if h.dashboardPath != "" {
			if _, err := os.Stat(filepath.Join(h.dashboardPath, "index.html")); err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte(`{"status":"unhealthy","error":"dashboard files not found"}`))
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	}