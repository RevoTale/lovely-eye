package collect

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/lovely-eye/server/internal/analytics"
	"github.com/lovely-eye/server/internal/site"
	"github.com/lovely-eye/server/internal/transport/http/clientip"
)

type AnalyticsHandler struct {
	analyticsService *analytics.Service
	siteService      *site.Service
	config           AnalyticsHandlerConfig
	ipResolver       *clientip.Resolver
	rateLimiter      *RateLimiter
}

type AnalyticsHandlerConfig struct {
	MaxBodyBytes       int64
	MaxPropertiesBytes int
}

func NewAnalyticsHandler(
	analyticsService *analytics.Service,
	siteService *site.Service,
	config AnalyticsHandlerConfig,
	ipResolver *clientip.Resolver,
	rateLimiter *RateLimiter,
) *AnalyticsHandler {
	if config.MaxBodyBytes <= 0 {
		config.MaxBodyBytes = 16 * 1024
	}
	if config.MaxPropertiesBytes <= 0 {
		config.MaxPropertiesBytes = 8 * 1024
	}
	if ipResolver == nil {
		ipResolver = clientip.MustNewResolver(nil)
	}
	return &AnalyticsHandler{
		analyticsService: analyticsService,
		siteService:      siteService,
		config:           config,
		ipResolver:       ipResolver,
		rateLimiter:      rateLimiter,
	}
}

type collectRequest struct {
	Path        string `json:"path"`
	Name        string `json:"name"`
	Properties  string `json:"properties"`
	Referrer    string `json:"referrer"`
	Exit        bool   `json:"exit"`
	UTMSource   string `json:"utm_source"`
	UTMMedium   string `json:"utm_medium"`
	UTMCampaign string `json:"utm_campaign"`
}

const (
	maxPathLength        = 2048
	maxReferrerLength    = 2048
	maxUTMSourceLength   = 128
	maxUTMMediumLength   = 128
	maxUTMCampaignLength = 256
)

func (h *AnalyticsHandler) Collect(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		h.handleAnalyticsPreflight(w, r)
		return
	}

	siteKey := strings.TrimSpace(r.URL.Query().Get("site_key"))
	if siteKey == "" {
		respondError(w, http.StatusBadRequest, "site_key query parameter is required")
		return
	}

	ip := h.clientIP(r)
	if !h.allowCollect("ip|" + ip) {
		respondError(w, http.StatusTooManyRequests, "rate limit exceeded")
		return
	}

	site, err := h.loadAnalyticsSite(r, siteKey)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if !h.applyAnalyticsCORSForSite(w, r, site) {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if !h.allowCollect("site|" + site.PublicKey + "|ip|" + ip) {
		respondError(w, http.StatusTooManyRequests, "rate limit exceeded")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, h.config.MaxBodyBytes)
	var req collectRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			respondError(w, http.StatusRequestEntityTooLarge, "request body is too large")
			return
		}
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Properties != "" {
		if len(req.Properties) > h.config.MaxPropertiesBytes {
			respondError(w, http.StatusRequestEntityTooLarge, "properties are too large")
			return
		}
		var props map[string]interface{}
		if err := json.Unmarshal([]byte(req.Properties), &props); err != nil || props == nil {
			respondError(w, http.StatusBadRequest, "properties must be a JSON object")
			return
		}
	}

	if req.Path == "" {
		respondError(w, http.StatusBadRequest, "path is required")
		return
	}
	if exceedsCollectPersistenceLimits(req) {
		respondError(w, http.StatusBadRequest, "request field is too long")
		return
	}

	if req.Name != "" {
		err = h.analyticsService.CollectEvent(r.Context(), analytics.EventInput{
			SiteKey:    siteKey,
			Name:       req.Name,
			Path:       req.Path,
			Properties: req.Properties,
			UserAgent:  r.UserAgent(),
			IP:         ip,
			Origin:     r.Header.Get("Origin"),
			Referer:    r.Header.Get("Referer"),
		})
	} else {
		err = h.analyticsService.CollectPageView(r.Context(), analytics.CollectInput{
			SiteKey:     siteKey,
			Path:        req.Path,
			Exit:        req.Exit,
			Referrer:    req.Referrer,
			UserAgent:   r.UserAgent(),
			IP:          ip,
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

func exceedsCollectPersistenceLimits(req collectRequest) bool {
	return utf8.RuneCountInString(req.Path) > maxPathLength ||
		utf8.RuneCountInString(req.Referrer) > maxReferrerLength ||
		utf8.RuneCountInString(req.UTMSource) > maxUTMSourceLength ||
		utf8.RuneCountInString(req.UTMMedium) > maxUTMMediumLength ||
		utf8.RuneCountInString(req.UTMCampaign) > maxUTMCampaignLength
}

func (h *AnalyticsHandler) handleAnalyticsPreflight(w http.ResponseWriter, r *http.Request) {
	siteKey := strings.TrimSpace(r.URL.Query().Get("site_key"))
	if siteKey == "" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	site, err := h.loadAnalyticsSite(r, siteKey)
	if err != nil || !h.applyAnalyticsCORSForSite(w, r, site) {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AnalyticsHandler) applyAnalyticsCORSForSite(w http.ResponseWriter, r *http.Request, site *site.Site) bool {
	origin := r.Header.Get("Origin")
	referer := r.Header.Get("Referer")
	if !analytics.IsAllowedDomain(origin, referer, site.Domains) {
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

func (h *AnalyticsHandler) loadAnalyticsSite(r *http.Request, siteKey string) (*site.Site, error) {
	site, err := h.siteService.GetByPublicKey(r.Context(), siteKey)
	if err != nil {
		return nil, fmt.Errorf("get site by public key: %w", err)
	}
	return site, nil
}

func (h *AnalyticsHandler) clientIP(r *http.Request) string {
	return h.ipResolver.GetClientIP(r.Header.Get("X-Forwarded-For"), r.Header.Get("X-Real-IP"), r.RemoteAddr)
}

func (h *AnalyticsHandler) allowCollect(key string) bool {
	return h.rateLimiter == nil || h.rateLimiter.Allow(key)
}

func respondError(w http.ResponseWriter, status int, message string) {
	http.Error(w, message, status)
}
