package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/lovely-eye/server/internal/dashboard"
	"github.com/uptrace/bun"
)
type HealthHandler struct {
	db *bun.DB
}
func NewHealthHandler(db *bun.DB) *HealthHandler {
	return &HealthHandler{
		db: db,
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

		// Check dashboard files exist
		if _, err := os.Stat(filepath.Join(dashboard.StaticDir, "index.html")); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status":"unhealthy","error":"dashboard files not found"}`))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	}