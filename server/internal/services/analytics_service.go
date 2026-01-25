package services

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
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

type CollectInput struct {
	SiteKey     string `json:"site_key"`
	Path        string `json:"path"`
	Referrer    string `json:"referrer"`
	ScreenWidth int    `json:"screen_width"`
	UserAgent   string `json:"-"`
	IP          string `json:"-"`
	Origin      string `json:"-"`
	Referer     string `json:"-"`
	UTMSource   string `json:"utm_source"`
	UTMMedium   string `json:"utm_medium"`
	UTMCampaign string `json:"utm_campaign"`
}

type EventInput struct {
	SiteKey    string `json:"site_key"`
	Name       string `json:"name"`
	Path       string `json:"path"`
	Properties string `json:"properties"`
	UserAgent  string `json:"-"`
	IP         string `json:"-"`
	Origin     string `json:"-"`
	Referer    string `json:"-"`
}

func (s *AnalyticsService) CollectPageView(ctx context.Context, input CollectInput) error {

	if s.botDetector.IsBot(input.UserAgent) {
		return nil
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

	country := UnknownCountry
	if site.TrackCountry && s.geoIPService != nil {
		country, err = s.geoIPService.ResolveCountry(input.IP)
		if err != nil {
			err = fmt.Errorf("get country for ip %s: %w", input.IP, err)
			slog.Error("country resolve failed", "err", err)
		}
	}

	client, err := s.findOrCreateClient(ctx, site.ID, visitorHash, device, browser, os, screenSize, country.Name)
	if err != nil {
		return fmt.Errorf("find or create client: %w", err)
	}

	now := time.Now()
	nowUnix := now.Unix()

	sessionTimeout := now.Add(-30 * time.Minute)
	session, err := s.getActiveSession(ctx, site.ID, client.ID, sessionTimeout)

	if err != nil || session == nil {

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

	event := &models.Event{
		SessionID:    session.ID,
		Time:         nowUnix,
		Hour:         nowUnix / 3600,
		Day:          nowUnix / 86400,
		Path:         input.Path,
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

func (s *AnalyticsService) getActiveSession(ctx context.Context, siteID, clientID int64, since time.Time) (*models.Session, error) {
	session, err := s.analyticsRepo.GetActiveSession(ctx, siteID, clientID, since.Unix())
	if err != nil {
		return nil, fmt.Errorf("get active session: %w", err)
	}
	return session, nil
}

func (s *AnalyticsService) getRecentPageViewEvent(ctx context.Context, sessionID int64, path string, since int64) (*models.Event, error) {
	event, err := s.analyticsRepo.GetRecentPageViewEvent(ctx, sessionID, path, since)
	if err != nil {
		return nil, fmt.Errorf("get recent pageview event: %w", err)
	}
	return event, nil
}

func (s *AnalyticsService) CollectEvent(ctx context.Context, input EventInput) error {

	if s.botDetector.IsBot(input.UserAgent) {
		return nil
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
	screenSize := ""

	country := UnknownCountry
	if site.TrackCountry && s.geoIPService != nil {
		country, err = s.geoIPService.ResolveCountry(input.IP)
		if err != nil {
			err = fmt.Errorf("get country for ip %s: %w", input.IP, err)
			slog.Error("country resolve failed", "err", err)
		}
	}

	client, err := s.findOrCreateClient(ctx, site.ID, visitorHash, device, browser, os, screenSize, country.Name)
	if err != nil {
		return fmt.Errorf("find or create client: %w", err)
	}

	now := time.Now()
	nowUnix := now.Unix()

	sessionTimeout := now.Add(-30 * time.Minute)
	session, _ := s.getActiveSession(ctx, site.ID, client.ID, sessionTimeout)

	if session == nil {

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
			PageViewCount: 0,
		}
		if err := s.analyticsRepo.CreateSession(ctx, session); err != nil {
			return fmt.Errorf("create session: %w", err)
		}
	} else {

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

	event := &models.Event{
		SessionID:    session.ID,
		Time:         nowUnix,
		Hour:         nowUnix / 3600,
		Day:          nowUnix / 86400,
		Path:         input.Path,
		DefinitionID: &definition.ID,
	}

	if err := s.analyticsRepo.CreateEvent(ctx, event); err != nil {
		return fmt.Errorf("create event: %w", err)
	}

	if sanitizedProps != "" {

		var propsMap map[string]string
		if err := json.Unmarshal([]byte(sanitizedProps), &propsMap); err != nil {
			return fmt.Errorf("unmarshal sanitized properties: %w", err)
		}

		fieldMap := make(map[string]int64, len(definition.Fields))
		for _, field := range definition.Fields {
			fieldMap[field.Key] = field.ID
		}

		eventDataList := make([]*models.EventData, 0, len(propsMap))
		for key, value := range propsMap {
			fieldID, exists := fieldMap[key]
			if !exists {
				continue
			}

			eventDataList = append(eventDataList, &models.EventData{
				EventID: event.ID,
				FieldID: fieldID,
				Value:   value,
			})
		}

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

type DashboardFilter = repository.AnalyticsFilter
type AnalyticsQuery = repository.AnalyticsQuery

func buildAnalyticsQuery(siteID int64, from, to time.Time, filter DashboardFilter) AnalyticsQuery {
	return AnalyticsQuery{
		SiteID: siteID,
		From:   from,
		To:     to,
		Filter: filter,
	}
}

func (s *AnalyticsService) GetDashboardOverview(ctx context.Context, siteID int64, from, to time.Time) (*DashboardOverview, error) {
	return s.GetDashboardOverviewWithFilter(ctx, buildAnalyticsQuery(siteID, from, to, DashboardFilter{}))
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

func (s *AnalyticsService) GetDashboardOverviewWithFilter(ctx context.Context, query AnalyticsQuery) (*DashboardOverview, error) {
	visitors, _ := s.analyticsRepo.GetVisitorCountWithFilter(ctx, query)
	pageViews, _ := s.analyticsRepo.GetPageViewCountWithFilter(ctx, query)
	sessions, _ := s.analyticsRepo.GetSessionCountWithFilter(ctx, query)
	bounceRate, _ := s.analyticsRepo.GetBounceRateWithFilter(ctx, query)
	avgDuration, _ := s.analyticsRepo.GetAvgSessionDurationWithFilter(ctx, query)

	return &DashboardOverview{
		Visitors:    visitors,
		PageViews:   pageViews,
		Sessions:    sessions,
		BounceRate:  bounceRate,
		AvgDuration: avgDuration,
	}, nil
}

func (s *AnalyticsService) GetTopPagesWithFilterPaged(ctx context.Context, query AnalyticsQuery) ([]repository.PageStats, int, error) {
	stats, total, err := s.analyticsRepo.GetTopPagesWithFilterPaged(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("get top pages with filter paged: %w", err)
	}
	return stats, total, nil
}

func (s *AnalyticsService) GetTopReferrersWithFilterPaged(ctx context.Context, query AnalyticsQuery) ([]repository.ReferrerStats, int, error) {
	stats, total, err := s.analyticsRepo.GetTopReferrersWithFilterPaged(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("get top referrers with filter paged: %w", err)
	}
	return stats, total, nil
}

func (s *AnalyticsService) GetDeviceStatsWithFilterPaged(ctx context.Context, query AnalyticsQuery) ([]repository.DeviceStats, int, int, error) {
	stats, total, totalVisitors, err := s.analyticsRepo.GetDeviceStatsWithFilterPaged(ctx, query)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("get device stats with filter paged: %w", err)
	}
	return stats, total, totalVisitors, nil
}

func (s *AnalyticsService) GetCountryStatsWithFilterPaged(ctx context.Context, query AnalyticsQuery) ([]repository.CountryStats, int, int, error) {
	stats, total, totalVisitors, err := s.analyticsRepo.GetCountryStatsWithFilterPaged(ctx, query)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("get country stats with filter paged: %w", err)
	}
	return stats, total, totalVisitors, nil
}

func (s *AnalyticsService) GetBrowserStatsWithFilter(ctx context.Context, query AnalyticsQuery) ([]repository.BrowserStats, error) {
	stats, err := s.analyticsRepo.GetBrowserStatsWithFilter(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get browser stats with filter: %w", err)
	}
	return stats, nil
}

func (s *AnalyticsService) GetTimeSeriesStatsWithFilter(ctx context.Context, query AnalyticsQuery) ([]repository.DailyVisitorStats, error) {
	stats, err := s.analyticsRepo.GetTimeSeriesStatsWithFilter(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get time series stats with filter: %w", err)
	}
	return stats, nil
}

func (s *AnalyticsService) GetRealtimeVisitors(ctx context.Context, siteID int64) (int, error) {

	from := time.Now().Add(-5 * time.Minute)
	to := time.Now()
	count, err := s.analyticsRepo.GetVisitorCount(ctx, siteID, from, to)
	if err != nil {
		return 0, fmt.Errorf("get visitor count: %w", err)
	}
	return count, nil
}

func (s *AnalyticsService) GetActivePages(ctx context.Context, siteID int64, limit, offset int) ([]repository.ActivePageStats, error) {

	since := time.Now().Add(-5 * time.Minute)
	stats, err := s.analyticsRepo.GetActivePages(ctx, siteID, since, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get active pages: %w", err)
	}
	return stats, nil
}

func (s *AnalyticsService) GetEvents(ctx context.Context, siteID int64, from, to time.Time, limit, offset int) ([]*models.Event, error) {
	events, err := s.analyticsRepo.GetEvents(ctx, siteID, from, to, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get events: %w", err)
	}
	return events, nil
}

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

func (s *AnalyticsService) GetEventsWithTotalAndFilter(ctx context.Context, query AnalyticsQuery) ([]*models.Event, int, error) {
	events, err := s.analyticsRepo.GetEventsWithFilter(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("get events with filter: %w", err)
	}

	total, err := s.analyticsRepo.GetEventCountWithFilter(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("get event count with filter: %w", err)
	}

	return events, total, nil
}

type EventCountWithEvent struct {
	Event *models.Event
	Count int
}

func (s *AnalyticsService) GetEventCounts(ctx context.Context, query AnalyticsQuery) ([]EventCountWithEvent, error) {
	results, err := s.analyticsRepo.GetEventCountsGrouped(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get event counts grouped: %w", err)
	}

	if len(results) == 0 {
		return []EventCountWithEvent{}, nil
	}

	eventIDs := make([]int64, len(results))
	for i, result := range results {
		eventIDs[i] = result.EventID
	}

	events, err := s.analyticsRepo.GetEventsByIDs(ctx, eventIDs)
	if err != nil {
		return nil, fmt.Errorf("get events by IDs: %w", err)
	}

	eventMap := make(map[int64]*models.Event, len(events))
	for _, event := range events {
		eventMap[event.ID] = event
	}

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

	uaString := strings.ToLower(ua.String)
	if strings.Contains(uaString, "iphone") || strings.Contains(uaString, "android") {
		if strings.Contains(uaString, "mobile") {
			return "mobile"
		}

		if strings.Contains(uaString, "tablet") {
			return "tablet"
		}

		if strings.Contains(uaString, "iphone") {
			return "mobile"
		}

		if strings.Contains(uaString, "ipad") {
			return "tablet"
		}
	}

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
	c, err := s.geoIPService.ResolveCountry(ip)
	if nil != err {
		err = fmt.Errorf("get country for ip %s: %w", ip, err)
		slog.Error("country resolve failed", "err", err)
	}

	if c == (Country{}) || c == LocalNetworkCountry {
		return false
	}
	for _, entry := range blocked {
		if entry != nil && strings.EqualFold(entry.CountryCode, c.ISOCode) {
			return true
		}
	}
	return false
}
