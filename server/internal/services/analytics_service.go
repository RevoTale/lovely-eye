package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/netip"
	"net/url"
	"strings"
	"time"

	"github.com/lovely-eye/server/internal/models"
	"github.com/lovely-eye/server/internal/repository"
	"github.com/lovely-eye/server/pkg/validation"
	"github.com/mileusna/useragent"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/hkdf"
)

type AnalyticsService struct {
	analyticsRepo       *repository.AnalyticsRepository
	countryService      countrySyncer
	siteRepo            *repository.SiteRepository
	eventDefinitionRepo *repository.EventDefinitionRepository
	botDetector         *BotDetector
	geoIPService        geoIPProvider
	identitySecret      []byte
	now                 func() time.Time
}

type geoIPProvider interface {
	SetEnabled(enabled bool)
	Status() GeoIPStatus
	EnsureAvailable(ctx context.Context) error
	Refresh(ctx context.Context) error
	ResolveCountry(ipStr string) (Country, error)
	ListCountries(search string) ([]GeoIPCountry, error)
	Close() error
}

func NewAnalyticsService(
	analyticsRepo *repository.AnalyticsRepository,
	siteRepo *repository.SiteRepository,
	eventDefinitionRepo *repository.EventDefinitionRepository,
	geoIPService geoIPProvider,
	countryService countrySyncer,
	identitySecret string,
) *AnalyticsService {
	return &AnalyticsService{
		analyticsRepo:       analyticsRepo,
		countryService:      countryService,
		siteRepo:            siteRepo,
		eventDefinitionRepo: eventDefinitionRepo,
		botDetector:         NewBotDetector(),
		geoIPService:        geoIPService,
		identitySecret:      []byte(identitySecret),
		now:                 time.Now,
	}
}

type CollectInput struct {
	SiteKey     string `json:"site_key"`
	Path        string `json:"path"`
	Referrer    string `json:"referrer"`
	ScreenWidth int    `json:"screen_width"`
	Duration    int    `json:"duration"`
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

	ua := useragent.Parse(input.UserAgent)
	device := categorizeDevice(ua)
	browser := normalizeBrowser(ua)
	os := normalizeOS(ua)
	screenSize := categorizeScreenSize(input.ScreenWidth)
	now := s.now()
	nowUnix := now.Unix()

	country := UnknownCountry
	if site.TrackCountry && s.geoIPService != nil {
		country = s.resolveCountryBestEffort(input.IP)
	}
	isDurationOnly := input.Duration > 0 &&
		input.ScreenWidth == 0 &&
		input.Referrer == "" &&
		input.UTMSource == "" &&
		input.UTMMedium == "" &&
		input.UTMCampaign == ""

	if err := s.analyticsRepo.RunInTx(ctx, func(ctx context.Context, tx bun.Tx) error {
		client, err := s.resolveClientWithRotation(ctx, tx, site.ID, input.IP, device, browser, os, screenSize, country.ISOCode, now)
		if err != nil {
			return fmt.Errorf("resolve client with rotation: %w", err)
		}

		session, err := s.analyticsRepo.GetActiveSessionTx(ctx, tx, site.ID, client.ID, now.Add(-30*time.Minute).Unix())
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("get active session: %w", err)
		}
		if errors.Is(err, sql.ErrNoRows) {
			session = nil
		}

		if session == nil {
			if isDurationOnly {
				return nil
			}

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
			if err := s.analyticsRepo.CreateSessionTx(ctx, tx, session); err != nil {
				return fmt.Errorf("create session: %w", err)
			}
		} else {
			if isDurationOnly {
				session.ExitTime = nowUnix
				session.ExitHour = nowUnix / 3600
				session.ExitDay = nowUnix / 86400
				if input.Path != "" {
					session.ExitPath = input.Path
				}
				session.Duration = int(nowUnix - session.EnterTime)
				if err := s.analyticsRepo.UpdateSessionTx(ctx, tx, session); err != nil {
					return fmt.Errorf("update duration-only session: %w", err)
				}
				return nil
			}

			recentEvent, err := s.analyticsRepo.GetRecentPageViewEventTx(ctx, tx, session.ID, input.Path, nowUnix-10)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("get recent pageview event: %w", err)
			}
			if err == nil && recentEvent != nil {
				return nil
			}

			session.ExitTime = nowUnix
			session.ExitHour = nowUnix / 3600
			session.ExitDay = nowUnix / 86400
			session.ExitPath = input.Path
			session.Duration = int(nowUnix - session.EnterTime)
			session.PageViewCount++
			if err := s.analyticsRepo.UpdateSessionTx(ctx, tx, session); err != nil {
				return fmt.Errorf("update session: %w", err)
			}
		}

		if isDurationOnly {
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
		if err := s.analyticsRepo.CreateEventTx(ctx, tx, event); err != nil {
			return fmt.Errorf("create event: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("collect page view transaction: %w", err)
	}

	return nil
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

	ua := useragent.Parse(input.UserAgent)
	device := categorizeDevice(ua)
	browser := normalizeBrowser(ua)
	os := normalizeOS(ua)
	screenSize := models.ClientScreenSizeUnknown
	now := s.now()
	nowUnix := now.Unix()

	country := UnknownCountry
	if site.TrackCountry && s.geoIPService != nil {
		country = s.resolveCountryBestEffort(input.IP)
	}
	if err := s.analyticsRepo.RunInTx(ctx, func(ctx context.Context, tx bun.Tx) error {
		client, err := s.resolveClientWithRotation(ctx, tx, site.ID, input.IP, device, browser, os, screenSize, country.ISOCode, now)
		if err != nil {
			return fmt.Errorf("resolve client with rotation: %w", err)
		}

		session, err := s.analyticsRepo.GetActiveSessionTx(ctx, tx, site.ID, client.ID, now.Add(-30*time.Minute).Unix())
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("get active session: %w", err)
		}
		if errors.Is(err, sql.ErrNoRows) {
			session = nil
		}

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
			if err := s.analyticsRepo.CreateSessionTx(ctx, tx, session); err != nil {
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
			if err := s.analyticsRepo.UpdateSessionTx(ctx, tx, session); err != nil {
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
		if err := s.analyticsRepo.CreateEventTx(ctx, tx, event); err != nil {
			return fmt.Errorf("create event: %w", err)
		}

		if sanitizedProps == "" {
			return nil
		}

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
		if err := s.analyticsRepo.CreateEventDataBatchTx(ctx, tx, eventDataList); err != nil {
			return fmt.Errorf("create event data batch: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("collect event transaction: %w", err)
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
	if s.countryService != nil {
		if err := s.countryService.SyncFromGeoIP(ctx); err != nil {
			return fmt.Errorf("sync persisted countries: %w", err)
		}
	}
	return nil
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
	if err := s.geoIPService.Refresh(ctx); err != nil {
		return s.geoIPService.Status(), fmt.Errorf("refresh geoip database: %w", err)
	}
	if s.countryService != nil {
		if err := s.countryService.SyncFromGeoIP(ctx); err != nil {
			return s.geoIPService.Status(), fmt.Errorf("sync persisted countries: %w", err)
		}
	}
	return s.geoIPService.Status(), nil
}

func (s *AnalyticsService) Close() error {
	if s.geoIPService == nil {
		return nil
	}
	if err := s.geoIPService.Close(); err != nil {
		return fmt.Errorf("close geoip service: %w", err)
	}
	return nil
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

func (s *AnalyticsService) GetOperatingSystemStatsWithFilterPaged(ctx context.Context, query AnalyticsQuery) ([]repository.OperatingSystemStats, int, int, error) {
	stats, total, totalVisitors, err := s.analyticsRepo.GetOperatingSystemStatsWithFilterPaged(ctx, query)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("get operating system stats with filter paged: %w", err)
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

const unknownVisitorIPPrefix = "unknown"

// generateVisitorID creates the site-scoped daily UTC hash used by the
// UTC-day-skipped client rotation helper. Client reuse compares today's and
// yesterday's hash to preserve continuity across adjacent UTC days only.
func (s *AnalyticsService) generateVisitorID(
	siteID int64,
	ip string,
	browser models.ClientBrowser,
	device models.ClientDevice,
	now time.Time,
) string {
	dateBucket := now.UTC().Format("2006-01-02")
	key := s.deriveVisitorIdentityKey(siteID, dateBucket)
	ipPrefix := truncateVisitorIPPrefix(ip)
	data := fmt.Sprintf("%d|%s|%s|%s", siteID, ipPrefix, browser.String(), device.String())

	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write([]byte(data))

	sum := mac.Sum(nil)
	return hex.EncodeToString(sum[:16])
}

func (s *AnalyticsService) deriveVisitorIdentityKey(siteID int64, dateBucket string) []byte {
	info := fmt.Appendf(nil, "analytics:%d:%s", siteID, dateBucket)
	reader := hkdf.New(sha256.New, s.identitySecret, nil, info)
	key := make([]byte, sha256.Size)
	_, _ = io.ReadFull(reader, key)
	return key
}

func truncateVisitorIPPrefix(ip string) string {
	addr, err := netip.ParseAddr(strings.TrimSpace(ip))
	if err != nil {
		return unknownVisitorIPPrefix
	}

	if addr.Is4In6() {
		addr = netip.AddrFrom4(addr.As4())
	}

	if addr.Is4() {
		return netip.PrefixFrom(addr, 24).Masked().String()
	}

	return netip.PrefixFrom(addr, 64).Masked().String()
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

func categorizeDevice(ua useragent.UserAgent) models.ClientDevice {
	uaString := strings.ToLower(ua.String)

	switch {
	case strings.Contains(uaString, "watchos"),
		strings.Contains(uaString, "watch os"),
		strings.Contains(uaString, "apple watch"),
		strings.Contains(uaString, "wear os"),
		strings.Contains(uaString, "smartwatch"),
		strings.Contains(uaString, "galaxy watch"):
		return models.ClientDeviceWatch
	case strings.Contains(uaString, "smart-tv"),
		strings.Contains(uaString, "smarttv"),
		strings.Contains(uaString, "android tv"),
		strings.Contains(uaString, "bravia"),
		strings.Contains(uaString, "hbbtv"),
		strings.Contains(uaString, "googletv"),
		strings.Contains(uaString, "appletv"),
		strings.Contains(uaString, "crkey"),
		strings.Contains(uaString, "aft"),
		strings.Contains(uaString, "roku"),
		strings.Contains(uaString, "viera"),
		strings.Contains(uaString, "netcast"),
		strings.Contains(uaString, "tv;"):
		return models.ClientDeviceSmartTV
	case strings.Contains(uaString, "playstation"),
		strings.Contains(uaString, "xbox"),
		strings.Contains(uaString, "nintendo switch"):
		return models.ClientDeviceConsole
	}

	if ua.Tablet {
		return models.ClientDeviceTablet
	}
	if ua.Mobile {
		return models.ClientDeviceMobile
	}
	if ua.Desktop {
		return models.ClientDeviceDesktop
	}

	if strings.Contains(uaString, "ipad") || strings.Contains(uaString, "tablet") {
		return models.ClientDeviceTablet
	}
	if strings.Contains(uaString, "iphone") || strings.Contains(uaString, "ipod") {
		return models.ClientDeviceMobile
	}
	if strings.Contains(uaString, "android") {
		if strings.Contains(uaString, "mobile") {
			return models.ClientDeviceMobile
		}
		return models.ClientDeviceTablet
	}
	if strings.Contains(uaString, "windows") ||
		strings.Contains(uaString, "macintosh") ||
		strings.Contains(uaString, "linux") ||
		strings.Contains(uaString, "cros") {
		return models.ClientDeviceDesktop
	}

	return models.ClientDeviceDesktop
}

func normalizeBrowser(ua useragent.UserAgent) models.ClientBrowser {
	uaString := strings.ToLower(ua.String)

	switch {
	case strings.Contains(uaString, "playstation"):
		return models.ClientBrowserPlayStation
	case strings.Contains(uaString, "xbox"):
		return models.ClientBrowserXbox
	case strings.Contains(uaString, "fb_iab"),
		strings.Contains(uaString, "fban"),
		strings.Contains(uaString, "fbav"):
		return models.ClientBrowserFacebookInApp
	case strings.Contains(uaString, "instagram"):
		return models.ClientBrowserInstagramInApp
	case strings.Contains(uaString, "edg/"),
		strings.Contains(uaString, "edgios"),
		ua.IsEdge():
		return models.ClientBrowserEdge
	case strings.Contains(uaString, "samsungbrowser"):
		return models.ClientBrowserSamsungInternet
	case strings.Contains(uaString, "opr/"),
		strings.Contains(uaString, "opera mini"),
		strings.Contains(uaString, "opera mobi"),
		ua.IsOpera(),
		ua.IsOperaMini():
		return models.ClientBrowserOpera
	case strings.Contains(uaString, "vivaldi"):
		return models.ClientBrowserVivaldi
	case strings.Contains(uaString, "yabrowser"),
		strings.Contains(uaString, "yowser"):
		return models.ClientBrowserYandex
	case strings.Contains(uaString, "duckduckgo"):
		return models.ClientBrowserDuckDuckGo
	case strings.Contains(uaString, "ucbrowser"),
		strings.Contains(uaString, "ucweb"):
		return models.ClientBrowserUCBrowser
	case strings.Contains(uaString, "miuibrowser"):
		return models.ClientBrowserMIUI
	case strings.Contains(uaString, "msie"),
		strings.Contains(uaString, "trident"),
		ua.IsInternetExplorer():
		return models.ClientBrowserInternetExplorer
	case strings.Contains(uaString, "wv"),
		strings.Contains(uaString, "webview"):
		return models.ClientBrowserAndroidWebView
	case strings.Contains(uaString, "crios"),
		strings.Contains(uaString, "chrome"),
		ua.IsChrome():
		return models.ClientBrowserChrome
	case strings.Contains(uaString, "fxios"),
		strings.Contains(uaString, "firefox"),
		ua.IsFirefox():
		return models.ClientBrowserFirefox
	case strings.Contains(uaString, "safari"),
		(strings.Contains(uaString, "applewebkit") &&
			(strings.Contains(uaString, "iphone") || strings.Contains(uaString, "ipad") || strings.Contains(uaString, "macintosh"))),
		ua.IsSafari():
		return models.ClientBrowserSafari
	}

	name := strings.TrimSpace(ua.Name)
	if name == "" {
		return models.ClientBrowserOther
	}

	if browser, ok := models.ClientBrowserFromLabel(name); ok {
		return browser
	}
	return models.ClientBrowserFromLegacyLabel(name)
}

func normalizeOS(ua useragent.UserAgent) models.ClientOS {
	uaString := strings.ToLower(ua.String)

	switch {
	case strings.Contains(uaString, "wear os"):
		return models.ClientOSWearOS
	case strings.Contains(uaString, "watchos"),
		strings.Contains(uaString, "watch os"),
		strings.Contains(uaString, "apple watch"):
		return models.ClientOSWatchOS
	case strings.Contains(uaString, "playstation"):
		return models.ClientOSPlayStation
	case strings.Contains(uaString, "xbox"):
		return models.ClientOSXbox
	case strings.Contains(uaString, "ipad"):
		return models.ClientOSIPadOS
	case strings.Contains(uaString, "iphone"),
		strings.Contains(uaString, "ipod"),
		ua.IsIOS():
		return models.ClientOSIOS
	case strings.Contains(uaString, "android"):
		return models.ClientOSAndroid
	case strings.Contains(uaString, "cros"),
		ua.IsChromeOS():
		return models.ClientOSChromeOS
	case strings.Contains(uaString, "windows"),
		ua.IsWindows():
		return models.ClientOSWindows
	case strings.Contains(uaString, "mac os"),
		strings.Contains(uaString, "macintosh"),
		ua.IsMacOS():
		return models.ClientOSMacOS
	case strings.Contains(uaString, "linux"),
		ua.IsLinux():
		return models.ClientOSLinux
	}

	if os, ok := models.ClientOSFromLabel(strings.TrimSpace(ua.OS)); ok {
		return os
	}
	return models.ClientOSFromLegacyLabel(strings.TrimSpace(ua.OS))
}

func categorizeScreenSize(width int) models.ClientScreenSize {
	return models.ClientScreenSizeFromWidth(width)
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

	normalized, err := validation.ValidateDomain(host)
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
	c := s.resolveCountryBestEffort(ip)

	if c == (Country{}) || c == UnknownCountry || c == LocalNetworkCountry {
		return false
	}
	for _, entry := range blocked {
		if entry != nil && strings.EqualFold(entry.CountryCode, c.ISOCode) {
			return true
		}
	}
	return false
}

func (s *AnalyticsService) resolveCountryBestEffort(ip string) Country {
	if s.geoIPService == nil {
		return UnknownCountry
	}

	country, err := s.geoIPService.ResolveCountry(ip)
	if err == nil {
		return country
	}

	if errors.Is(err, ErrNoDBReader) {
		return UnknownCountry
	}

	slog.Error("country resolve failed", "err", fmt.Errorf("get country for ip %s: %w", ip, err))
	return UnknownCountry
}
