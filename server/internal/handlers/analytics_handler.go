package handlers

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"

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
	Name        string `json:"name"`
	Properties  string `json:"properties"`
	Referrer    string `json:"referrer"`
	ScreenWidth int    `json:"screen_width"`
	Duration    int    `json:"duration"`
	UTMSource   string `json:"utm_source"`
	UTMMedium   string `json:"utm_medium"`
	UTMCampaign string `json:"utm_campaign"`
}

func (h *AnalyticsHandler) Collect(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		h.handleAnalyticsPreflight(w, r)
		return
	}

	var req collectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.SiteKey == "" {
		respondError(w, http.StatusBadRequest, "site_key is required")
		return
	}

	if !h.applyAnalyticsCORS(w, r, req.SiteKey) {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if req.Properties != "" {
		var props map[string]interface{}
		if err := json.Unmarshal([]byte(req.Properties), &props); err != nil {
			respondError(w, http.StatusBadRequest, "properties must be a JSON object")
			return
		}
	}

	var err error
	if req.Name != "" {
		err = h.analyticsService.CollectEvent(r.Context(), services.EventInput{
			SiteKey:    req.SiteKey,
			Name:       req.Name,
			Path:       req.Path,
			Properties: req.Properties,
			UserAgent:  r.UserAgent(),
			IP:         getClientIP(r),
			Origin:     r.Header.Get("Origin"),
			Referer:    r.Header.Get("Referer"),
		})
	} else {
		if req.Path == "" {
			respondError(w, http.StatusBadRequest, "path is required")
			return
		}
		err = h.analyticsService.CollectPageView(r.Context(), services.CollectInput{
			SiteKey:     req.SiteKey,
			Path:        req.Path,
			Referrer:    req.Referrer,
			ScreenWidth: req.ScreenWidth,
			Duration:    req.Duration,
			UserAgent:   r.UserAgent(),
			IP:          getClientIP(r),
			Origin:      r.Header.Get("Origin"),
			Referer:     r.Header.Get("Referer"),
			UTMSource:   req.UTMSource,
			UTMMedium:   req.UTMMedium,
			UTMCampaign: req.UTMCampaign,
		})
	}

	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AnalyticsHandler) Event(w http.ResponseWriter, r *http.Request) {
	h.Collect(w, r)
}

func (h *AnalyticsHandler) handleAnalyticsPreflight(w http.ResponseWriter, r *http.Request) {
	siteKey := r.URL.Query().Get("site_key")
	if siteKey == "" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if !h.applyAnalyticsCORS(w, r, siteKey) {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AnalyticsHandler) applyAnalyticsCORS(w http.ResponseWriter, r *http.Request, siteKey string) bool {
	site, err := h.siteService.GetByPublicKey(r.Context(), siteKey)
	if err != nil {
		return false
	}

	origin := r.Header.Get("Origin")
	referer := r.Header.Get("Referer")
	if !services.IsAllowedDomain(origin, referer, site.Domains) {
		return false
	}

	if origin == "" {
		return true
	}

	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Vary", "Origin")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Max-Age", "3600")
	return true
}

func respondError(w http.ResponseWriter, status int, message string) {
	http.Error(w, message, status)
}

func getClientIP(r *http.Request) string {

	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {

		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			if ip != "" {
				return ip
			}
		}
	}

	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {

		return r.RemoteAddr
	}
	return ip
}
