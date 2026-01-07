package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/lovely-eye/server/internal/middleware"
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

// Dashboard returns analytics stats for a site (protected endpoint)
func (h *AnalyticsHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	siteID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid site ID")
		return
	}

	// Verify user owns this site
	_, err = h.siteService.GetByID(r.Context(), siteID, claims.UserID)
	if err != nil {
		switch err {
		case services.ErrSiteNotFound:
			respondError(w, http.StatusNotFound, "Site not found")
		case services.ErrNotAuthorized:
			respondError(w, http.StatusForbidden, "Not authorized")
		default:
			respondError(w, http.StatusInternalServerError, "Failed to get site")
		}
		return
	}

	// Parse date range from query params
	from, to := parseDateRange(r)

	stats, err := h.analyticsService.GetDashboardStats(r.Context(), siteID, from, to)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get stats")
		return
	}

	respondJSON(w, http.StatusOK, stats)
}

// Realtime returns current active visitors
func (h *AnalyticsHandler) Realtime(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	siteID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid site ID")
		return
	}

	// Verify user owns this site
	_, err = h.siteService.GetByID(r.Context(), siteID, claims.UserID)
	if err != nil {
		switch err {
		case services.ErrSiteNotFound:
			respondError(w, http.StatusNotFound, "Site not found")
		case services.ErrNotAuthorized:
			respondError(w, http.StatusForbidden, "Not authorized")
		default:
			respondError(w, http.StatusInternalServerError, "Failed to get site")
		}
		return
	}

	visitors, err := h.analyticsService.GetRealtimeVisitors(r.Context(), siteID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get realtime stats")
		return
	}

	respondJSON(w, http.StatusOK, map[string]int{"visitors": visitors})
}

func parseDateRange(r *http.Request) (time.Time, time.Time) {
	now := time.Now()
	defaultFrom := now.AddDate(0, 0, -30)
	defaultTo := now

	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")
	period := r.URL.Query().Get("period")

	// Handle preset periods
	switch period {
	case "today":
		return startOfDay(now), now
	case "yesterday":
		yesterday := now.AddDate(0, 0, -1)
		return startOfDay(yesterday), endOfDay(yesterday)
	case "7d":
		return now.AddDate(0, 0, -7), now
	case "30d":
		return now.AddDate(0, 0, -30), now
	case "90d":
		return now.AddDate(0, 0, -90), now
	case "12m":
		return now.AddDate(-1, 0, 0), now
	}

	// Parse custom date range
	from := defaultFrom
	to := defaultTo

	if fromStr != "" {
		if parsed, err := time.Parse("2006-01-02", fromStr); err == nil {
			from = parsed
		}
	}

	if toStr != "" {
		if parsed, err := time.Parse("2006-01-02", toStr); err == nil {
			to = endOfDay(parsed)
		}
	}

	return from, to
}

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func endOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}
