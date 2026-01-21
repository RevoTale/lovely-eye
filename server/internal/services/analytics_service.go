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

	// Generate anonymous visitor ID (SHA-256 hash) with daily salt rotation for privacy
	visitorHash := s.generateVisitorID(input.IP, input.UserAgent, site.PublicKey)

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

	// Find or create Client record (stores stable attributes)
	client, err := s.findOrCreateClient(ctx, site.ID, visitorHash, device, browser, os, screenSize, country)
	if err != nil {
		return fmt.Errorf("find or create client: %w", err)
	}

	// Current time for timestamps
	now := time.Now()
	nowUnix := now.Unix()

	// Try to find existing session (within 30 minutes)
	sessionTimeout := now.Add(-30 * time.Minute)
	session, err := s.getActiveSession(ctx, site.ID, client.ID, sessionTimeout)

	if err != nil || session == nil {
		// Create new session
		session = &models.Session{
			SiteID:        site.ID,
			ClientID:      client.ID,
			EnterTime:     nowUnix,
			EnterHour:     nowUnix / 3600,
			EnterDay:      nowUnix / 86400,
			EnterPath:     input.Path,
			ExitTime:      nowUnix,
			ExitHour:      nowUnix / 3600,
			ExitDay:       nowUnix / 86400,
			ExitPath:      input.Path,
			Referrer:      input.Referrer,
			UTMSource:     input.UTMSource,
			UTMMedium:     input.UTMMedium,
			UTMCampaign:   input.UTMCampaign,
			Duration:      0,
			PageViewCount: 1,
		}
		if err := s.analyticsRepo.CreateSession(ctx, session); err != nil {
			return fmt.Errorf("create session: %w", err)
		}
	} else {
		// Update existing session
		session.ExitTime = nowUnix
		session.ExitHour = nowUnix / 3600
		session.ExitDay = nowUnix / 86400
		session.ExitPath = input.Path
		session.PageViewCount++
		session.Duration = int(nowUnix - session.EnterTime)
		if err := s.analyticsRepo.UpdateSession(ctx, session); err != nil {
			return fmt.Errorf("update session: %w", err)
		}
	}

	// Deduplicate page views: check if same session viewed same page in last 10 seconds
	// This prevents duplicate counts from double-clicks, SPA route changes, or script reloads
	recentEvent, _ := s.getRecentPageViewEvent(ctx, session.ID, input.Path, nowUnix-10)
	if recentEvent != nil {
		// Same page view within 10 seconds - ignore to prevent duplicates
		return nil
	}

	// Create Event record (Type=EventTypePageview=0)
	event := &models.Event{
		SessionID: session.ID,
		Time:      nowUnix,
		Hour:      nowUnix / 3600,
		Day:       nowUnix / 86400,
		Path:      input.Path,
		Name:      input.Title,
		Type:      models.EventTypePageview,
		DefinitionID: nil,
	}

	if err := s.analyticsRepo.CreateEvent(ctx, event); err != nil {
		return fmt.Errorf("create event: %w", err)
	}
	return nil
}

// findOrCreateClient finds an existing client by hash or creates a new one
func (s *AnalyticsService) findOrCreateClient(ctx context.Context, siteID int64, hash, device, browser, os, screenSize, country string) (*models.Client, error) {
	client, err := s.analyticsRepo.FindOrCreateClient(ctx, siteID, hash, device, browser, os, screenSize, country)
	if err != nil {
		return nil, fmt.Errorf("find or create client: %w", err)
	}
	return client, nil
}

// getActiveSession finds an active session for a client (within timeout)
func (s *AnalyticsService) getActiveSession(ctx context.Context, siteID, clientID int64, since time.Time) (*models.Session, error) {
	session, err := s.analyticsRepo.GetActiveSession(ctx, siteID, clientID, since.Unix())
	if err != nil {
		return nil, fmt.Errorf("get active session: %w", err)
	}
	return session, nil
}

// getRecentPageViewEvent checks if the same page was viewed recently in the same session
func (s *AnalyticsService) getRecentPageViewEvent(ctx context.Context, sessionID int64, path string, since int64) (*models.Event, error) {
	event, err := s.analyticsRepo.GetRecentPageViewEvent(ctx, sessionID, path, since)
	if err != nil {
		return nil, fmt.Errorf("get recent pageview event: %w", err)
	}
	return event, nil
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

	// Generate anonymous visitor ID (SHA-256 hash)
	visitorHash := s.generateVisitorID(input.IP, input.UserAgent, site.PublicKey)

	// Parse user agent
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
	screenSize := "" // Not available for custom events

	// Get country from IP (only if enabled for the site)
	country := ""
	if site.TrackCountry && s.geoIPService != nil {
		country = s.geoIPService.GetCountryName(input.IP)
	}

	// Find or create Client record
	client, err := s.findOrCreateClient(ctx, site.ID, visitorHash, device, browser, os, screenSize, country)
	if err != nil {
		return fmt.Errorf("find or create client: %w", err)
	}

	// Current time for timestamps
	now := time.Now()
	nowUnix := now.Unix()

	// Try to find existing session (within 30 minutes)
	sessionTimeout := now.Add(-30 * time.Minute)
	session, _ := s.getActiveSession(ctx, site.ID, client.ID, sessionTimeout)

	if session == nil {
		// Create new session for event-only session
		entryPath := input.Path
		if entryPath == "" {
			entryPath = "/"
		}
		session = &models.Session{
			SiteID:        site.ID,
			ClientID:      client.ID,
			EnterTime:     nowUnix,
			EnterHour:     nowUnix / 3600,
			EnterDay:      nowUnix / 86400,
			EnterPath:     entryPath,
			ExitTime:      nowUnix,
			ExitHour:      nowUnix / 3600,
			ExitDay:       nowUnix / 86400,
			ExitPath:      entryPath,
			Referrer:      "",
			UTMSource:     "",
			UTMMedium:     "",
			UTMCampaign:   "",
			Duration:      0,
			PageViewCount: 0, // Event-only session starts with 0 pageviews
		}
		if err := s.analyticsRepo.CreateSession(ctx, session); err != nil {
			return fmt.Errorf("create session: %w", err)
		}
	} else {
		// Update existing session
		session.ExitTime = nowUnix
		session.ExitHour = nowUnix / 3600
		session.ExitDay = nowUnix / 86400
		if input.Path != "" {
			session.ExitPath = input.Path
		}
		session.Duration = int(nowUnix - session.EnterTime)
		if err := s.analyticsRepo.UpdateSession(ctx, session); err != nil {
			return fmt.Errorf("update session: %w", err)
		}
	}

	// Create Event record with Type=EventTypeCustom=1
	event := &models.Event{
		SessionID:    session.ID,
		Time:         nowUnix,
		Hour:         nowUnix / 3600,
		Day:          nowUnix / 86400,
		Path:         input.Path,
		Name:         input.Name,
		Type:         models.EventTypeCustom,
		DefinitionID: &definition.ID,
	}

	if err := s.analyticsRepo.CreateEvent(ctx, event); err != nil {
		return fmt.Errorf("create event: %w", err)
	}

	// Store event properties in EventData table
	if sanitizedProps != "" {
		// Parse the JSON properties
		var propsMap map[string]string
		if err := json.Unmarshal([]byte(sanitizedProps), &propsMap); err != nil {
			return fmt.Errorf("unmarshal sanitized properties: %w", err)
		}

		// Build a map of field keys to IDs for fast lookup
		fieldMap := make(map[string]int64, len(definition.Fields))
		for _, field := range definition.Fields {
			fieldMap[field.Key] = field.ID
		}

		// Create EventData records for each property
		eventDataList := make([]*models.EventData, 0, len(propsMap))
		for key, value := range propsMap {
			fieldID, exists := fieldMap[key]
			if !exists {
				continue // Skip if field not found in definition
			}

			eventDataList = append(eventDataList, &models.EventData{
				EventID: event.ID,
				FieldID: fieldID,
				Value:   value,
			})
		}

		// Batch insert all event data
		if len(eventDataList) > 0 {
			if err := s.analyticsRepo.CreateEventDataBatch(ctx, eventDataList); err != nil {
				return fmt.Errorf("create event data batch: %w", err)
			}
		}
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

// GetEventsWithTotalAndFilter retrieves events with total count for pagination and filtering
func (s *AnalyticsService) GetEventsWithTotalAndFilter(ctx context.Context, siteID int64, from, to time.Time, referrer, device, page, country []string, limit, offset int) ([]*models.Event, int, error) {
	events, err := s.analyticsRepo.GetEventsWithFilter(ctx, siteID, from, to, referrer, device, page, country, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("get events with filter: %w", err)
	}

	total, err := s.analyticsRepo.GetEventCountWithFilter(ctx, siteID, from, to, referrer, device, page, country)
	if err != nil {
		return nil, 0, fmt.Errorf("get event count with filter: %w", err)
	}

	return events, total, nil
}

// EventCountWithEvent represents an event count with its most recent event
type EventCountWithEvent struct {
	Event *models.Event
	Count int
}

// GetEventCounts retrieves aggregated event counts by name
func (s *AnalyticsService) GetEventCounts(ctx context.Context, siteID int64, from, to time.Time, referrer, device, page, country []string, limit int) ([]EventCountWithEvent, error) {
	results, err := s.analyticsRepo.GetEventCountsGrouped(ctx, siteID, from, to, referrer, device, page, country, limit)
	if err != nil {
		return nil, fmt.Errorf("get event counts grouped: %w", err)
	}

	if len(results) == 0 {
		return []EventCountWithEvent{}, nil
	}

	// Extract event IDs
	eventIDs := make([]int64, len(results))
	for i, result := range results {
		eventIDs[i] = result.EventID
	}

	// Fetch full event details
	events, err := s.analyticsRepo.GetEventsByIDs(ctx, eventIDs)
	if err != nil {
		return nil, fmt.Errorf("get events by IDs: %w", err)
	}

	// Create map of event ID to event
	eventMap := make(map[int64]*models.Event, len(events))
	for _, event := range events {
		eventMap[event.ID] = event
	}

	// Build final result
	eventCounts := make([]EventCountWithEvent, 0, len(results))
	for _, result := range results {
		event, ok := eventMap[result.EventID]
		if !ok {
			continue
		}
		eventCounts = append(eventCounts, EventCountWithEvent{
			Event: event,
			Count: result.Count,
		})
	}

	return eventCounts, nil
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
		case models.FieldTypeString:
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
		case models.FieldTypeInt:
			// Accept both int and float for int type
			switch v := value.(type) {
			case float64:
				sanitized[key] = int64(v)
			case int:
				sanitized[key] = int64(v)
			case int64:
				sanitized[key] = v
			default:
				return "", false, nil
			}
		case models.FieldTypeFloat:
			numberValue, ok := value.(float64)
			if !ok {
				return "", false, nil
			}
			sanitized[key] = numberValue
		case models.FieldTypeBool:
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
