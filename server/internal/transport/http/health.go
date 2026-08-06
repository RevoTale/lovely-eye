package http

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/uptrace/bun"
)

type healthHandler struct {
	db             *bun.DB
	dashboardPath  string
	connectTimeout time.Duration
}

func newHealthHandler(db *bun.DB, dashboardPath string, connectTimeout time.Duration) *healthHandler {
	return &healthHandler{
		db:             db,
		dashboardPath:  dashboardPath,
		connectTimeout: connectTimeout,
	}
}
func (h *healthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	pingCtx := r.Context()
	if h.connectTimeout > 0 {
		var cancel context.CancelFunc
		pingCtx, cancel = context.WithTimeout(pingCtx, h.connectTimeout)
		defer cancel()
	}
	if err := h.db.PingContext(pingCtx); err != nil {
		http.Error(w, `{"status":"unhealthy","error":"database connection failed"}`, http.StatusServiceUnavailable)
		return
	}

	if h.dashboardPath != "" {
		if _, err := os.Stat(filepath.Join(h.dashboardPath, "index.html")); err != nil {
			http.Error(w, `{"status":"unhealthy","error":"dashboard files not found"}`, http.StatusServiceUnavailable)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`{"status":"healthy"}`)); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
