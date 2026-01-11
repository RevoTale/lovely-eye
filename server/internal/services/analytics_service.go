package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/lovely-eye/server/internal/models"
	"github.com/lovely-eye/server/internal/repository"
)

type AnalyticsService struct {
	analyticsRepo *repository.AnalyticsRepository
	siteRepo      *repository.SiteRepository
}

func NewAnalyticsService(analyticsRepo *repository.AnalyticsRepository, siteRepo *repository.SiteRepository) *AnalyticsService {
	return &AnalyticsService{
		analyticsRepo: analyticsRepo,
		siteRepo:      siteRepo,
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
	site, err := s.siteRepo.GetByPublicKey(ctx, input.SiteKey)
	if err != nil {
		return err
	}

	// Generate anonymous visitor ID
	visitorID := generateVisitorID(input.IP, input.UserAgent, site.PublicKey)

	// Try to find existing session (within 30 minutes)
	sessionTimeout := time.Now().Add(-30 * time.Minute)
	session, err := s.analyticsRepo.GetSessionByVisitor(ctx, site.ID, visitorID, sessionTimeout)

	device := parseDevice(input.UserAgent)
	browser := parseBrowser(input.UserAgent)
	os := parseOS(input.UserAgent)
	screenSize := categorizeScreenSize(input.ScreenWidth)

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
		session.Duration = int(time.Since(session.StartedAt).Seconds())
		if err := s.analyticsRepo.UpdateSession(ctx, session); err != nil {
			return err
		}
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
	site, err := s.siteRepo.GetByPublicKey(ctx, input.SiteKey)
	if err != nil {
		return err
	}

	visitorID := generateVisitorID(input.IP, input.UserAgent, site.PublicKey)

	// Try to find existing session
	sessionTimeout := time.Now().Add(-30 * time.Minute)
	session, _ := s.analyticsRepo.GetSessionByVisitor(ctx, site.ID, visitorID, sessionTimeout)

	var sessionID int64
	if session != nil {
		sessionID = session.ID
	}

	event := &models.Event{
		SiteID:     site.ID,
		SessionID:  sessionID,
		VisitorID:  visitorID,
		Name:       input.Name,
		Path:       input.Path,
		Properties: input.Properties,
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
	Referrer *string
	Device   *string
	Page     *string
}

func (s *AnalyticsService) GetDashboardStats(ctx context.Context, siteID int64, from, to time.Time) (*DashboardStats, error) {
	return s.GetDashboardStatsWithFilter(ctx, siteID, from, to, DashboardFilter{})
}

func (s *AnalyticsService) GetDashboardStatsWithFilter(ctx context.Context, siteID int64, from, to time.Time, filter DashboardFilter) (*DashboardStats, error) {
	visitors, _ := s.analyticsRepo.GetVisitorCountWithFilter(ctx, siteID, from, to, filter.Referrer, filter.Device, filter.Page)
	pageViews, _ := s.analyticsRepo.GetPageViewCountWithFilter(ctx, siteID, from, to, filter.Referrer, filter.Device, filter.Page)
	sessions, _ := s.analyticsRepo.GetSessionCountWithFilter(ctx, siteID, from, to, filter.Referrer, filter.Device, filter.Page)
	bounceRate, _ := s.analyticsRepo.GetBounceRateWithFilter(ctx, siteID, from, to, filter.Referrer, filter.Device, filter.Page)
	avgDuration, _ := s.analyticsRepo.GetAvgSessionDurationWithFilter(ctx, siteID, from, to, filter.Referrer, filter.Device, filter.Page)
	topPages, _ := s.analyticsRepo.GetTopPagesWithFilter(ctx, siteID, from, to, 10, filter.Referrer, filter.Device, filter.Page)
	topReferrers, _ := s.analyticsRepo.GetTopReferrersWithFilter(ctx, siteID, from, to, 10, filter.Referrer, filter.Device, filter.Page)
	browsers, _ := s.analyticsRepo.GetBrowserStatsWithFilter(ctx, siteID, from, to, 10, filter.Referrer, filter.Device, filter.Page)
	devices, _ := s.analyticsRepo.GetDeviceStatsWithFilter(ctx, siteID, from, to, 10, filter.Referrer, filter.Device, filter.Page)
	countries, _ := s.analyticsRepo.GetCountryStatsWithFilter(ctx, siteID, from, to, 10, filter.Referrer, filter.Device, filter.Page)
	dailyStats, _ := s.analyticsRepo.GetDailyStatsWithFilter(ctx, siteID, from, to, filter.Referrer, filter.Device, filter.Page)

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

// Helper functions for parsing user agent
func generateVisitorID(ip, userAgent, salt string) string {
	data := ip + userAgent + salt
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])[:32]
}

func parseDevice(userAgent string) string {
	// Simplified device detection
	if contains(userAgent, "Mobile", "Android", "iPhone", "iPad") {
		if contains(userAgent, "iPad", "Tablet") {
			return "tablet"
		}
		return "mobile"
	}
	return "desktop"
}

func parseBrowser(userAgent string) string {
	switch {
	case contains(userAgent, "Firefox"):
		return "Firefox"
	case contains(userAgent, "Edg"):
		return "Edge"
	case contains(userAgent, "Chrome"):
		return "Chrome"
	case contains(userAgent, "Safari"):
		return "Safari"
	case contains(userAgent, "Opera"):
		return "Opera"
	default:
		return "Other"
	}
}

func parseOS(userAgent string) string {
	switch {
	case contains(userAgent, "Windows"):
		return "Windows"
	case contains(userAgent, "Mac OS"):
		return "macOS"
	case contains(userAgent, "Linux"):
		return "Linux"
	case contains(userAgent, "Android"):
		return "Android"
	case contains(userAgent, "iOS", "iPhone", "iPad"):
		return "iOS"
	default:
		return "Other"
	}
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

func contains(s string, substrs ...string) bool {
	for _, substr := range substrs {
		if len(s) >= len(substr) {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
		}
	}
	return false
}
