package services

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lovely-eye/server/internal/models"
	"github.com/lovely-eye/server/internal/repository"
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
}

// CollectPageView records a page view and manages sessions
func (s *AnalyticsService) CollectPageView(ctx context.Context, input CollectInput) error {
	// Filter out bots
	if s.botDetector.IsBot(input.UserAgent) {
		return nil // Silently ignore bot traffic
	}

	site, err := s.siteRepo.GetByPublicKey(ctx, input.SiteKey)
	if err != nil {
		return err
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
			return err
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
			return err
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

	return s.analyticsRepo.CreatePageView(ctx, pageView)
}

// CollectEvent records a custom event
func (s *AnalyticsService) CollectEvent(ctx context.Context, input EventInput) error {
	// Filter out bots
	if s.botDetector.IsBot(input.UserAgent) {
		return nil // Silently ignore bot traffic
	}

	site, err := s.siteRepo.GetByPublicKey(ctx, input.SiteKey)
	if err != nil {
		return err
	}

	if s.eventDefinitionRepo == nil {
		return nil
	}

	definition, err := s.eventDefinitionRepo.GetByName(ctx, site.ID, input.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}

	sanitizedProps, ok, err := sanitizeEventProperties(input.Properties, definition.Fields)
	if err != nil {
		return err
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
			return err
		}
	} else {
		session.LastSeenAt = now
		if input.Path != "" {
			session.ExitPage = input.Path
		}
		session.Duration = int(now.Sub(session.StartedAt).Seconds())
		if err := s.analyticsRepo.UpdateSession(ctx, session); err != nil {
			return err
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

	return s.analyticsRepo.CreateEvent(ctx, event)
}

type DashboardStats struct {
	Visitors     int
	PageViews    int
	Sessions     int
	BounceRate   float64
	AvgDuration  float64
	TopPages     []repository.PageStats
	TopReferrers []repository.ReferrerStats
	Browsers     []repository.BrowserStats
	Devices      []repository.DeviceStats
	Countries    []repository.CountryStats
	DailyStats   []repository.DailyVisitorStats
}

type DashboardFilter struct {
	Referrer []string
	Device   []string
	Page     []string
	Country  []string
}

func (s *AnalyticsService) GetDashboardStats(ctx context.Context, siteID int64, from, to time.Time) (*DashboardStats, error) {
	return s.GetDashboardStatsWithFilter(ctx, siteID, from, to, DashboardFilter{})
}

func (s *AnalyticsService) SyncGeoIPRequirement(ctx context.Context) error {
	if s.geoIPService == nil {
		return nil
	}
	requires, err := s.siteRepo.AnyTrackCountry(ctx)
	if err != nil {
		return err
	}
	s.geoIPService.SetEnabled(requires)
	if !requires {
		return nil
	}
	return s.geoIPService.EnsureAvailable(ctx)
}

func (s *AnalyticsService) GeoIPStatus() GeoIPStatus {
	if s.geoIPService == nil {
		return GeoIPStatus{State: geoIPStateDisabled}
	}
	return s.geoIPService.Status()
}

func (s *AnalyticsService) RefreshGeoIPDatabase(ctx context.Context) (GeoIPStatus, error) {
	if s.geoIPService == nil {
		return GeoIPStatus{State: geoIPStateDisabled}, nil
	}
	err := s.geoIPService.EnsureAvailable(ctx)
	return s.geoIPService.Status(), err
}

func (s *AnalyticsService) GetDashboardStatsWithFilter(ctx context.Context, siteID int64, from, to time.Time, filter DashboardFilter) (*DashboardStats, error) {
	visitors, _ := s.analyticsRepo.GetVisitorCountWithFilter(ctx, siteID, from, to, filter.Referrer, filter.Device, filter.Page, filter.Country)
	pageViews, _ := s.analyticsRepo.GetPageViewCountWithFilter(ctx, siteID, from, to, filter.Referrer, filter.Device, filter.Page, filter.Country)
	sessions, _ := s.analyticsRepo.GetSessionCountWithFilter(ctx, siteID, from, to, filter.Referrer, filter.Device, filter.Page, filter.Country)
	bounceRate, _ := s.analyticsRepo.GetBounceRateWithFilter(ctx, siteID, from, to, filter.Referrer, filter.Device, filter.Page, filter.Country)
	avgDuration, _ := s.analyticsRepo.GetAvgSessionDurationWithFilter(ctx, siteID, from, to, filter.Referrer, filter.Device, filter.Page, filter.Country)
	topPages, _ := s.analyticsRepo.GetTopPagesWithFilter(ctx, siteID, from, to, 10, filter.Referrer, filter.Device, filter.Page, filter.Country)
	topReferrers, _ := s.analyticsRepo.GetTopReferrersWithFilter(ctx, siteID, from, to, 10, filter.Referrer, filter.Device, filter.Page, filter.Country)
	browsers, _ := s.analyticsRepo.GetBrowserStatsWithFilter(ctx, siteID, from, to, 10, filter.Referrer, filter.Device, filter.Page, filter.Country)
	devices, _ := s.analyticsRepo.GetDeviceStatsWithFilter(ctx, siteID, from, to, 10, filter.Referrer, filter.Device, filter.Page, filter.Country)
	countries, _ := s.analyticsRepo.GetCountryStatsWithFilter(ctx, siteID, from, to, 10, filter.Referrer, filter.Device, filter.Page, filter.Country)
	dailyStats, _ := s.analyticsRepo.GetDailyStatsWithFilter(ctx, siteID, from, to, filter.Referrer, filter.Device, filter.Page, filter.Country)

	return &DashboardStats{
		Visitors:     visitors,
		PageViews:    pageViews,
		Sessions:     sessions,
		BounceRate:   bounceRate,
		AvgDuration:  avgDuration,
		TopPages:     topPages,
		TopReferrers: topReferrers,
		Browsers:     browsers,
		Devices:      devices,
		Countries:    countries,
		DailyStats:   dailyStats,
	}, nil
}

func (s *AnalyticsService) GetRealtimeVisitors(ctx context.Context, siteID int64) (int, error) {
	// Visitors active in last 5 minutes
	from := time.Now().Add(-5 * time.Minute)
	to := time.Now()
	return s.analyticsRepo.GetVisitorCount(ctx, siteID, from, to)
}

func (s *AnalyticsService) GetActivePages(ctx context.Context, siteID int64) ([]repository.ActivePageStats, error) {
	// Get pages viewed in last 5 minutes
	since := time.Now().Add(-5 * time.Minute)
	return s.analyticsRepo.GetActivePages(ctx, siteID, since)
}

// GetEvents retrieves events for a site within a time range
func (s *AnalyticsService) GetEvents(ctx context.Context, siteID int64, from, to time.Time, limit, offset int) ([]*models.Event, error) {
	return s.analyticsRepo.GetEvents(ctx, siteID, from, to, limit, offset)
}

// GetEventsWithTotal retrieves events with total count for pagination
func (s *AnalyticsService) GetEventsWithTotal(ctx context.Context, siteID int64, from, to time.Time, limit, offset int) ([]*models.Event, int, error) {
	events, err := s.analyticsRepo.GetEvents(ctx, siteID, from, to, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.analyticsRepo.GetEventCount(ctx, siteID, from, to)
	if err != nil {
		return nil, 0, err
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
			return "", false, nil
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
		return "", false, err
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
