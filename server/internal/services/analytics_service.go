package services

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/lovely-eye/server/internal/models"
	"github.com/lovely-eye/server/internal/repository"
	"github.com/lovely-eye/server/pkg/utils"
	"github.com/mileusna/useragent"
)

type AnalyticsService struct {
	analyticsRepo       *repository.AnalyticsRepository
	siteRepo            *repository.SiteRepository
	eventDefinitionRepo *repository.EventDefinitionRepository
	botDetector         *BotDetector
	geoIPService        *GeoIPService
}

func NewAnalyticsService(
	analyticsRepo *repository.AnalyticsRepository,
	siteRepo *repository.SiteRepository,
	eventDefinitionRepo *repository.EventDefinitionRepository,
	geoIPService *GeoIPService,
) *AnalyticsService {
	return &AnalyticsService{
		analyticsRepo:       analyticsRepo,
		siteRepo:            siteRepo,
		eventDefinitionRepo: eventDefinitionRepo,
		botDetector:         NewBotDetector(),
		geoIPService:        geoIPService,
	}
}

// CollectInput represents data collected from the tracking script
type CollectInput struct {
	SiteKey     string `json:"site_key"`
	Path        string `json:"path"`
	Title       string `json:"title"`
	Referrer    string `json:"referrer"`
	ScreenWidth int    `json:"screen_width"`
	UserAgent   string `json:"-"` // From header
	IP          string `json:"-"` // From request
	Origin      string `json:"-"` // From header
	Referer     string `json:"-"` // From header
	UTMSource   string `json:"utm_source"`
	UTMMedium   string `json:"utm_medium"`
	UTMCampaign string `json:"utm_campaign"`
}

// EventInput represents custom event data
type EventInput struct {
	SiteKey    string `json:"site_key"`
	Name       string `json:"name"`
	Path       string `json:"path"`
	Properties string `json:"properties"` // JSON string
	UserAgent  string `json:"-"`
	IP         string `json:"-"`
	Origin     string `json:"-"`
	Referer    string `json:"-"`
}

// CollectPageView records a page view and manages sessions
func (s *AnalyticsService) CollectPageView(ctx context.Context, input CollectInput) error {
	// Filter out bots
	if s.botDetector.IsBot(input.UserAgent) {
		return nil // Silently ignore bot traffic
	}

	site, err := s.siteRepo.GetByPublicKey(ctx, input.SiteKey)
	if err != nil {
		return fmt.Errorf("get site by public key: %w", err)
	}
	if !IsAllowedDomain(input.Origin, input.Referer, site.Domains) {
		return nil
	}
	if s.isBlockedRequest(site, input.IP) {
		return nil
	}

	// Generate anonymous visitor ID with daily salt rotation for privacy
	visitorID := s.generateVisitorID(input.IP, input.UserAgent, site.PublicKey)

	// Try to find existing session (within 30 minutes)
	sessionTimeout := time.Now().Add(-30 * time.Minute)
	session, err := s.analyticsRepo.GetSessionByVisitor(ctx, site.ID, visitorID, sessionTimeout)

	// Parse user agent with proper library
	ua := useragent.Parse(input.UserAgent)
	device := categorizeDevice(ua)
	browser := ua.Name
	if browser == "" {
		browser = "Other"
	}
	os := ua.OS
	if os == "" {
		os = "Other"
	}
	screenSize := categorizeScreenSize(input.ScreenWidth)

	// Get country from IP (only if enabled for the site)
	country := ""
	if site.TrackCountry && s.geoIPService != nil {
		country = s.geoIPService.GetCountryName(input.IP)
	}

	if err != nil || session == nil {
		// Create new session
		session = &models.Session{
			SiteID:      site.ID,
			VisitorID:   visitorID,
			StartedAt:   time.Now(),
			LastSeenAt:  time.Now(),
			EntryPage:   input.Path,
			ExitPage:    input.Path,
			Referrer:    input.Referrer,
			UTMSource:   input.UTMSource,
			UTMMedium:   input.UTMMedium,
			UTMCampaign: input.UTMCampaign,
			Device:      device,
			Browser:     browser,
			OS:          os,
			ScreenSize:  screenSize,
			PageViews:   1,
			IsBounce:    true,
			Country:     country,
			EventOnly:   false,
		}
		if err := s.analyticsRepo.CreateSession(ctx, session); err != nil {
			return fmt.Errorf("create session: %w", err)
		}
	} else {
		// Update existing session
		session.LastSeenAt = time.Now()
		session.ExitPage = input.Path
		session.PageViews++
		session.IsBounce = false
		if session.EventOnly {
			session.EventOnly = false
		}
		session.Duration = int(time.Since(session.StartedAt).Seconds())
		// Update country if not set
		if session.Country == "" && country != "" {
			session.Country = country
		}
		if err := s.analyticsRepo.UpdateSession(ctx, session); err != nil {
			return fmt.Errorf("update session: %w", err)
		}
	}

	// Deduplicate page views: check if same visitor viewed same page in last 10 seconds
	// This prevents duplicate counts from double-clicks, SPA route changes, or script reloads
	recentPageView, _ := s.analyticsRepo.GetRecentPageView(ctx, site.ID, visitorID, input.Path, time.Now().Add(-10*time.Second))
	if recentPageView != nil {
		// Same page view within 10 seconds - ignore to prevent duplicates
		return nil
	}

	// Record page view
	pageView := &models.PageView{
		SiteID:    site.ID,
		SessionID: session.ID,
		VisitorID: visitorID,
		Path:      input.Path,
		Title:     input.Title,
		Referrer:  input.Referrer,
	}

	if err := s.analyticsRepo.CreatePageView(ctx, pageView); err != nil {
		return fmt.Errorf("create page view: %w", err)
	}
	return nil
}

// CollectEvent records a custom event
func (s *AnalyticsService) CollectEvent(ctx context.Context, input EventInput) error {
	// Filter out bots
	if s.botDetector.IsBot(input.UserAgent) {
		return nil // Silently ignore bot traffic
	}

	site, err := s.siteRepo.GetByPublicKey(ctx, input.SiteKey)
	if err != nil {
		return fmt.Errorf("get site by public key: %w", err)
	}
	if !IsAllowedDomain(input.Origin, input.Referer, site.Domains) {
		return nil
	}
	if s.isBlockedRequest(site, input.IP) {
		return nil
	}

	if s.eventDefinitionRepo == nil {
		return nil
	}

	definition, err := s.eventDefinitionRepo.GetByName(ctx, site.ID, input.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("get event definition by name: %w", err)
	}

	sanitizedProps, ok, err := sanitizeEventProperties(input.Properties, definition.Fields)
	if err != nil {
		return fmt.Errorf("sanitize event properties: %w", err)
	}
	if !ok {
		return nil
	}

	visitorID := s.generateVisitorID(input.IP, input.UserAgent, site.PublicKey)

	// Try to find existing session
	now := time.Now()
	sessionTimeout := now.Add(-30 * time.Minute)
	session, _ := s.analyticsRepo.GetSessionByVisitor(ctx, site.ID, visitorID, sessionTimeout)

	if session == nil {
		ua := useragent.Parse(input.UserAgent)
		device := categorizeDevice(ua)
		browser := ua.Name
		if browser == "" {
			browser = "Other"
		}
		os := ua.OS
		if os == "" {
			os = "Other"
		}

		country := ""
		if site.TrackCountry && s.geoIPService != nil {
			country = s.geoIPService.GetCountryName(input.IP)
		}

		session = &models.Session{
			SiteID:     site.ID,
			VisitorID:  visitorID,
			StartedAt:  now,
			LastSeenAt: now,
			EntryPage:  input.Path,
			ExitPage:   input.Path,
			Device:     device,
			Browser:    browser,
			OS:         os,
			Country:    country,
			IsBounce:   true,
			EventOnly:  true,
		}
		if err := s.analyticsRepo.CreateSession(ctx, session); err != nil {
			return fmt.Errorf("create session: %w", err)
		}
	} else {
		session.LastSeenAt = now
		if input.Path != "" {
			session.ExitPage = input.Path
		}
		session.Duration = int(now.Sub(session.StartedAt).Seconds())
		if err := s.analyticsRepo.UpdateSession(ctx, session); err != nil {
			return fmt.Errorf("update session: %w", err)
		}
	}

	event := &models.Event{
		SiteID:     site.ID,
		SessionID:  session.ID,
		VisitorID:  visitorID,
		Name:       input.Name,
		Path:       input.Path,
		Properties: sanitizedProps,
	}

	if err := s.analyticsRepo.CreateEvent(ctx, event); err != nil {
		return fmt.Errorf("create event: %w", err)
	}
	return nil
}

type DashboardOverview struct {
	Visitors    int
	PageViews   int
	Sessions    int
	BounceRate  float64
	AvgDuration float64
}

type TimeBucket = repository.TimeBucket

const (
	TimeBucketDaily  = repository.TimeBucketDaily
	TimeBucketHourly = repository.TimeBucketHourly
)

type DashboardFilter struct {
	Referrer []string
	Device   []string
	Page     []string
	Country  []string
}

func (s *AnalyticsService) GetDashboardOverview(ctx context.Context, siteID int64, from, to time.Time) (*DashboardOverview, error) {
	return s.GetDashboardOverviewWithFilter(ctx, siteID, from, to, DashboardFilter{})
}

func (s *AnalyticsService) SyncGeoIPRequirement(ctx context.Context) error {
	if s.geoIPService == nil {
		return nil
	}
	requires, err := s.siteRepo.AnyGeoIPRequirement(ctx)
	if err != nil {
		return fmt.Errorf("check geoip requirement: %w", err)
	}
	s.geoIPService.SetEnabled(requires)
	if !requires {
		return nil
	}
	if err := s.geoIPService.EnsureAvailable(ctx); err != nil {
		return fmt.Errorf("ensure geoip available: %w", err)
	}
	return nil
}

func (s *AnalyticsService) GeoIPStatus() GeoIPStatus {
	if s.geoIPService == nil {
		return GeoIPStatus{State: geoIPStateDisabled}
	}
	return s.geoIPService.Status()
}

func (s *AnalyticsService) GeoIPCountries(search string) ([]GeoIPCountry, error) {
	if s.geoIPService == nil {
		return []GeoIPCountry{}, errors.New("country service is nil")
	}
	return s.geoIPService.ListCountries(search)
}

func (s *AnalyticsService) RefreshGeoIPDatabase(ctx context.Context) (GeoIPStatus, error) {
	if s.geoIPService == nil {
		return GeoIPStatus{State: geoIPStateDisabled}, nil
	}
	if err := s.geoIPService.EnsureAvailable(ctx); err != nil {
		return s.geoIPService.Status(), fmt.Errorf("ensure geoip available: %w", err)
	}
	return s.geoIPService.Status(), nil
}

func (s *AnalyticsService) GetDashboardOverviewWithFilter(ctx context.Context, siteID int64, from, to time.Time, filter DashboardFilter) (*DashboardOverview, error) {
	visitors, _ := s.analyticsRepo.GetVisitorCountWithFilter(ctx, siteID, from, to, filter.Referrer, filter.Device, filter.Page, filter.Country)
	pageViews, _ := s.analyticsRepo.GetPageViewCountWithFilter(ctx, siteID, from, to, filter.Referrer, filter.Device, filter.Page, filter.Country)
	sessions, _ := s.analyticsRepo.GetSessionCountWithFilter(ctx, siteID, from, to, filter.Referrer, filter.Device, filter.Page, filter.Country)
	bounceRate, _ := s.analyticsRepo.GetBounceRateWithFilter(ctx, siteID, from, to, filter.Referrer, filter.Device, filter.Page, filter.Country)
	avgDuration, _ := s.analyticsRepo.GetAvgSessionDurationWithFilter(ctx, siteID, from, to, filter.Referrer, filter.Device, filter.Page, filter.Country)

	return &DashboardOverview{
		Visitors:    visitors,
		PageViews:   pageViews,
		Sessions:    sessions,
		BounceRate:  bounceRate,
		AvgDuration: avgDuration,
	}, nil
}

func (s *AnalyticsService) GetTopPagesWithFilterPaged(ctx context.Context, siteID int64, from, to time.Time, limit, offset int, filter DashboardFilter) ([]repository.PageStats, int, error) {
	stats, total, err := s.analyticsRepo.GetTopPagesWithFilterPaged(ctx, siteID, from, to, limit, offset, filter.Referrer, filter.Device, filter.Page, filter.Country)
	if err != nil {
		return nil, 0, fmt.Errorf("get top pages with filter paged: %w", err)
	}
	return stats, total, nil
}

func (s *AnalyticsService) GetTopReferrersWithFilterPaged(ctx context.Context, siteID int64, from, to time.Time, limit, offset int, filter DashboardFilter) ([]repository.ReferrerStats, int, error) {
	stats, total, err := s.analyticsRepo.GetTopReferrersWithFilterPaged(ctx, siteID, from, to, limit, offset, filter.Referrer, filter.Device, filter.Page, filter.Country)
	if err != nil {
		return nil, 0, fmt.Errorf("get top referrers with filter paged: %w", err)
	}
	return stats, total, nil
}

func (s *AnalyticsService) GetDeviceStatsWithFilterPaged(ctx context.Context, siteID int64, from, to time.Time, limit, offset int, filter DashboardFilter) ([]repository.DeviceStats, int, int, error) {
	stats, total, totalVisitors, err := s.analyticsRepo.GetDeviceStatsWithFilterPaged(ctx, siteID, from, to, limit, offset, filter.Referrer, filter.Device, filter.Page, filter.Country)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("get device stats with filter paged: %w", err)
	}
	return stats, total, totalVisitors, nil
}

func (s *AnalyticsService) GetCountryStatsWithFilterPaged(ctx context.Context, siteID int64, from, to time.Time, limit, offset int, filter DashboardFilter) ([]repository.CountryStats, int, int, error) {
	stats, total, totalVisitors, err := s.analyticsRepo.GetCountryStatsWithFilterPaged(ctx, siteID, from, to, limit, offset, filter.Referrer, filter.Device, filter.Page, filter.Country)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("get country stats with filter paged: %w", err)
	}
	return stats, total, totalVisitors, nil
}

func (s *AnalyticsService) GetBrowserStatsWithFilter(ctx context.Context, siteID int64, from, to time.Time, limit int, filter DashboardFilter) ([]repository.BrowserStats, error) {
	stats, err := s.analyticsRepo.GetBrowserStatsWithFilter(ctx, siteID, from, to, limit, filter.Referrer, filter.Device, filter.Page, filter.Country)
	if err != nil {
		return nil, fmt.Errorf("get browser stats with filter: %w", err)
	}
	return stats, nil
}

func (s *AnalyticsService) GetTimeSeriesStatsWithFilter(ctx context.Context, siteID int64, from, to time.Time, bucket TimeBucket, limit int, filter DashboardFilter) ([]repository.DailyVisitorStats, error) {
	stats, err := s.analyticsRepo.GetTimeSeriesStatsWithFilter(ctx, siteID, from, to, bucket, limit, filter.Referrer, filter.Device, filter.Page, filter.Country)
	if err != nil {
		return nil, fmt.Errorf("get time series stats with filter: %w", err)
	}
	return stats, nil
}

func (s *AnalyticsService) GetRealtimeVisitors(ctx context.Context, siteID int64) (int, error) {
	// Visitors active in last 5 minutes
	from := time.Now().Add(-5 * time.Minute)
	to := time.Now()
	count, err := s.analyticsRepo.GetVisitorCount(ctx, siteID, from, to)
	if err != nil {
		return 0, fmt.Errorf("get visitor count: %w", err)
	}
	return count, nil
}

func (s *AnalyticsService) GetActivePages(ctx context.Context, siteID int64) ([]repository.ActivePageStats, error) {
	// Get pages viewed in last 5 minutes
	since := time.Now().Add(-5 * time.Minute)
	stats, err := s.analyticsRepo.GetActivePages(ctx, siteID, since)
	if err != nil {
		return nil, fmt.Errorf("get active pages: %w", err)
	}
	return stats, nil
}

// GetEvents retrieves events for a site within a time range
func (s *AnalyticsService) GetEvents(ctx context.Context, siteID int64, from, to time.Time, limit, offset int) ([]*models.Event, error) {
	events, err := s.analyticsRepo.GetEvents(ctx, siteID, from, to, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get events: %w", err)
	}
	return events, nil
}

// GetEventsWithTotal retrieves events with total count for pagination
func (s *AnalyticsService) GetEventsWithTotal(ctx context.Context, siteID int64, from, to time.Time, limit, offset int) ([]*models.Event, int, error) {
	events, err := s.analyticsRepo.GetEvents(ctx, siteID, from, to, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("get events: %w", err)
	}

	total, err := s.analyticsRepo.GetEventCount(ctx, siteID, from, to)
	if err != nil {
		return nil, 0, fmt.Errorf("get event count: %w", err)
	}

	return events, total, nil
}

// Helper functions for visitor identification and user agent parsing

// generateVisitorID creates a privacy-preserving visitor identifier
// Uses daily salt rotation to prevent long-term tracking while maintaining session continuity
func (s *AnalyticsService) generateVisitorID(ip, userAgent, siteKey string) string {
	// Use daily salt for privacy - visitor IDs change daily
	// This prevents long-term tracking while allowing accurate daily counts
	dateSalt := time.Now().Format("2006-01-02")

	// Hash: IP + UserAgent + SiteKey + DateSalt
	// This approach balances privacy and accuracy:
	// - Same visitor gets same ID within a day
	// - Different ID each day prevents tracking across days
	// - IP ensures different networks get different IDs
	// - UserAgent helps distinguish different devices on same network
	data := fmt.Sprintf("%s|%s|%s|%s", ip, userAgent, siteKey, dateSalt)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])[:32]
}

func sanitizeEventProperties(propsJSON string, fields []*models.EventDefinitionField) (string, bool, error) {
	if len(fields) == 0 {
		return "", true, nil
	}

	allowed := make(map[string]*models.EventDefinitionField, len(fields))
	for _, field := range fields {
		allowed[field.Key] = field
	}

	var props map[string]interface{}
	if propsJSON != "" {
		if err := json.Unmarshal([]byte(propsJSON), &props); err != nil {
			return "", false, fmt.Errorf("unmarshal properties json: %w", err)
		}
	} else {
		props = map[string]interface{}{}
	}

	sanitized := make(map[string]interface{})
	for key, value := range props {
		field, ok := allowed[key]
		if !ok {
			continue
		}

		switch field.Type {
		case eventFieldTypeString:
			strValue, ok := value.(string)
			if !ok {
				return "", false, nil
			}
			maxLen := field.MaxLength
			if maxLen <= 0 {
				maxLen = defaultEventMaxLength
			}
			if len(strValue) > maxLen {
				strValue = strValue[:maxLen]
			}
			sanitized[key] = strValue
		case eventFieldTypeNumber:
			numberValue, ok := value.(float64)
			if !ok {
				return "", false, nil
			}
			sanitized[key] = numberValue
		case eventFieldTypeBoolean:
			boolValue, ok := value.(bool)
			if !ok {
				return "", false, nil
			}
			sanitized[key] = boolValue
		default:
			return "", false, nil
		}
	}

	for _, field := range fields {
		if field.Required {
			if _, ok := sanitized[field.Key]; !ok {
				return "", false, nil
			}
		}
	}

	if len(sanitized) == 0 {
		return "", true, nil
	}

	bytes, err := json.Marshal(sanitized)
	if err != nil {
		return "", false, fmt.Errorf("marshal sanitized properties: %w", err)
	}
	return string(bytes), true, nil
}

// categorizeDevice determines device type from parsed user agent
func categorizeDevice(ua useragent.UserAgent) string {
	if ua.Tablet {
		return "tablet"
	}
	if ua.Mobile {
		return "mobile"
	}
	if ua.Desktop {
		return "desktop"
	}

	// Fallback: manual detection when library doesn't recognize the device
	// This handles simplified/incomplete user agent strings
	uaString := strings.ToLower(ua.String)
	if strings.Contains(uaString, "iphone") || strings.Contains(uaString, "android") {
		if strings.Contains(uaString, "mobile") {
			return "mobile"
		}
		// Android tablets don't have "mobile" in UA
		if strings.Contains(uaString, "tablet") {
			return "tablet"
		}
		// iPhone is always mobile
		if strings.Contains(uaString, "iphone") {
			return "mobile"
		}
		// iPad is tablet
		if strings.Contains(uaString, "ipad") {
			return "tablet"
		}
	}

	// Default to desktop for unknown user agents
	// This is more accurate than "other" since most traffic is desktop
	return "desktop"
}

func categorizeScreenSize(width int) string {
	switch {
	case width < 576:
		return "xs"
	case width < 768:
		return "sm"
	case width < 992:
		return "md"
	case width < 1200:
		return "lg"
	default:
		return "xl"
	}
}

func IsAllowedDomain(origin, referer string, domains []*models.SiteDomain) bool {
	host := hostFromHeader(origin)
	if host == "" {
		host = hostFromHeader(referer)
	}
	if host == "" {
		return false
	}

	for _, domain := range domains {
		if domain != nil && domain.Domain == host {
			return true
		}
	}

	return false
}

func hostFromHeader(raw string) string {
	if raw == "" {
		return ""
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		return ""
	}

	host := parsed.Hostname()
	if host == "" {
		return ""
	}

	normalized, err := utils.ValidateDomain(host)
	if err != nil {
		return ""
	}

	return normalized
}

func (s *AnalyticsService) isBlockedRequest(site *models.Site, ip string) bool {
	if site == nil {
		return false
	}
	if ip == "" {
		return false
	}

	if s.isIPBlocked(site.BlockedIPs, ip) {
		return true
	}

	return s.isCountryBlocked(site.BlockedCountries, ip)
}

func (s *AnalyticsService) isIPBlocked(blocked []*models.SiteBlockedIP, ip string) bool {
	if len(blocked) == 0 {
		return false
	}
	parsed := net.ParseIP(strings.TrimSpace(ip))
	if parsed == nil {
		return false
	}
	normalized := parsed.String()
	for _, entry := range blocked {
		if entry != nil && entry.IP == normalized {
			return true
		}
	}
	return false
}

func (s *AnalyticsService) isCountryBlocked(blocked []*models.SiteBlockedCountry, ip string) bool {
	if len(blocked) == 0 || s.geoIPService == nil {
		return false
	}
	code := s.geoIPService.GetCountry(ip)
	if code == "" || code == "Unknown" || code == "Local" {
		return false
	}
	for _, entry := range blocked {
		if entry != nil && strings.EqualFold(entry.CountryCode, code) {
			return true
		}
	}
	return false
}
