package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/lovely-eye/server/internal/services"
)

type AnalyticsHandler struct {
	analyticsService *services.AnalyticsService
	siteService      *services.SiteService
}

func NewAnalyticsHandler(analyticsService *services.AnalyticsService, siteService *services.SiteService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
		siteService:      siteService,
	}
}

type collectRequest struct {
	SiteKey     string `json:"site_key"`
	Path        string `json:"path"`
	Title       string `json:"title"`
	Referrer    string `json:"referrer"`
	ScreenWidth int    `json:"screen_width"`
	UTMSource   string `json:"utm_source"`
	UTMMedium   string `json:"utm_medium"`
	UTMCampaign string `json:"utm_campaign"`
}

type eventRequest struct {
	SiteKey    string `json:"site_key"`
	Name       string `json:"name"`
	Path       string `json:"path"`
	Properties string `json:"properties"`
}

// Collect handles page view tracking (public endpoint)
func (h *AnalyticsHandler) Collect(w http.ResponseWriter, r *http.Request) {
	var req collectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.SiteKey == "" || req.Path == "" {
		respondError(w, http.StatusBadRequest, "site_key and path are required")
		return
	}

	// Get client IP (handle proxies)
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}

	err := h.analyticsService.CollectPageView(r.Context(), services.CollectInput{
		SiteKey:     req.SiteKey,
		Path:        req.Path,
		Title:       req.Title,
		Referrer:    req.Referrer,
		ScreenWidth: req.ScreenWidth,
		UserAgent:   r.UserAgent(),
		IP:          ip,
		UTMSource:   req.UTMSource,
		UTMMedium:   req.UTMMedium,
		UTMCampaign: req.UTMCampaign,
	})

	if err != nil {
		// Don't expose errors to tracking clients
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Event handles custom event tracking (public endpoint)
func (h *AnalyticsHandler) Event(w http.ResponseWriter, r *http.Request) {
	var req eventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.SiteKey == "" || req.Name == "" {
		respondError(w, http.StatusBadRequest, "site_key and name are required")
		return
	}

	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}

	err := h.analyticsService.CollectEvent(r.Context(), services.EventInput{
		SiteKey:    req.SiteKey,
		Name:       req.Name,
		Path:       req.Path,
		Properties: req.Properties,
		UserAgent:  r.UserAgent(),
		IP:         ip,
	})

	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
