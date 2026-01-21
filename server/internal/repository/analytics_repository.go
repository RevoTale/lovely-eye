package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/lovely-eye/server/internal/models"
	"github.com/uptrace/bun"
)

type AnalyticsRepository struct {
	db *bun.DB
}

func NewAnalyticsRepository(db *bun.DB) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

func (r *AnalyticsRepository) FindOrCreateClient(ctx context.Context, siteID int64, hash, device, browser, os, screenSize, country string) (*models.Client, error) {
	// Try to find existing client by hash
	client := new(models.Client)
	err := r.db.NewSelect().
		Model(client).
		Where("site_id = ?", siteID).
		Where("hash = ?", hash).
		Limit(1).
		Scan(ctx)

	if err == nil {
		// Client exists, return it
		return client, nil
	}

	// Client doesn't exist, create new one
	client = &models.Client{
		SiteID:     siteID,
		Hash:       hash,
		Device:     device,
		Browser:    browser,
		OS:         os,
		ScreenSize: screenSize,
		Country:    country,
	}

	_, err = r.db.NewInsert().Model(client).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return client, nil
}

func (r *AnalyticsRepository) GetActiveSession(ctx context.Context, siteID, clientID int64, sinceUnix int64) (*models.Session, error) {
	session := new(models.Session)
	err := r.db.NewSelect().
		Model(session).
		Where("site_id = ?", siteID).
		Where("client_id = ?", clientID).
		Where("exit_time > ?", sinceUnix).
		Order("exit_time DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active session: %w", err)
	}
	return session, nil
}

func (r *AnalyticsRepository) GetRecentPageViewEvent(ctx context.Context, sessionID int64, path string, since int64) (*models.Event, error) {
	event := new(models.Event)
	err := r.db.NewSelect().
		Model(event).
		Where("session_id = ?", sessionID).
		Where("path = ?", path).
		Where("type = ?", models.EventTypePageview).
		Where("time > ?", since).
		Order("time DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent page view event: %w", err)
	}
	return event, nil
}

func (r *AnalyticsRepository) CreateSession(ctx context.Context, session *models.Session) error {
	_, err := r.db.NewInsert().Model(session).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	return nil
}

func (r *AnalyticsRepository) GetSession(ctx context.Context, id int64) (*models.Session, error) {
	session := new(models.Session)
	err := r.db.NewSelect().Model(session).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	return session, nil
}

func (r *AnalyticsRepository) GetSessionByVisitor(ctx context.Context, siteID int64, visitorID string, since time.Time) (*models.Session, error) {
	// Deprecated: Use GetActiveSession with client_id instead
	// This method is kept for backward compatibility but should not be used
	session := new(models.Session)
	sinceUnix := since.Unix()
	err := r.db.NewSelect().
		Model(session).
		Where("site_id = ?", siteID).
		Where("exit_time > ?", sinceUnix).
		Order("exit_time DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get session by visitor: %w", err)
	}
	return session, nil
}

func (r *AnalyticsRepository) UpdateSession(ctx context.Context, session *models.Session) error {
	_, err := r.db.NewUpdate().Model(session).WherePK().Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}
	return nil
}

func (r *AnalyticsRepository) CreateEvent(ctx context.Context, event *models.Event) error {
	_, err := r.db.NewInsert().Model(event).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}
	return nil
}

func (r *AnalyticsRepository) GetEvents(ctx context.Context, siteID int64, from, to time.Time, limit, offset int) ([]*models.Event, error) {
	var events []*models.Event
	fromUnix := from.Unix()
	toUnix := to.Unix()
	err := r.db.NewSelect().
		Model(&events).
		Relation("Data.Field").
		Join("INNER JOIN sessions s ON e.session_id = s.id").
		Where("s.site_id = ?", siteID).
		Where("e.time >= ?", fromUnix).
		Where("e.time <= ?", toUnix).
		Order("e.time DESC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}
	return events, nil
}

func (r *AnalyticsRepository) CreateEventData(ctx context.Context, eventData *models.EventData) error {
	_, err := r.db.NewInsert().Model(eventData).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create event data: %w", err)
	}
	return nil
}

func (r *AnalyticsRepository) CreateEventDataBatch(ctx context.Context, eventDataList []*models.EventData) error {
	if len(eventDataList) == 0 {
		return nil
	}
	_, err := r.db.NewInsert().Model(&eventDataList).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create event data batch: %w", err)
	}
	return nil
}

func (r *AnalyticsRepository) GetEventCount(ctx context.Context, siteID int64, from, to time.Time) (int, error) {
	fromUnix := from.Unix()
	toUnix := to.Unix()
	count, err := r.db.NewSelect().
		Model((*models.Event)(nil)).
		Join("INNER JOIN sessions s ON e.session_id = s.id").
		Where("s.site_id = ?", siteID).
		Where("e.time >= ?", fromUnix).
		Where("e.time <= ?", toUnix).
		Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get event count: %w", err)
	}
	return count, nil
}

func (r *AnalyticsRepository) GetVisitorCount(ctx context.Context, siteID int64, from, to time.Time) (int, error) {
	var count int
	fromUnix := from.Unix()
	toUnix := to.Unix()
	err := r.db.NewSelect().
		Model((*models.Session)(nil)).
		ColumnExpr("COUNT(DISTINCT client_id)").
		Where("site_id = ?", siteID).
		Where("enter_time >= ?", fromUnix).
		Where("enter_time <= ?", toUnix).
		Scan(ctx, &count)
	if err != nil {
		return 0, fmt.Errorf("failed to get visitor count: %w", err)
	}
	return count, nil
}

func (r *AnalyticsRepository) GetPageViewCount(ctx context.Context, siteID int64, from, to time.Time) (int, error) {
	fromUnix := from.Unix()
	toUnix := to.Unix()
	count, err := r.db.NewSelect().
		Model((*models.Event)(nil)).
		Join("INNER JOIN sessions s ON e.session_id = s.id").
		Where("s.site_id = ?", siteID).
		Where("e.type = ?", models.EventTypePageview).
		Where("e.time >= ?", fromUnix).
		Where("e.time <= ?", toUnix).
		Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get page view count: %w", err)
	}
	return count, nil
}

func (r *AnalyticsRepository) GetSessionCount(ctx context.Context, siteID int64, from, to time.Time) (int, error) {
	fromUnix := from.Unix()
	toUnix := to.Unix()
	count, err := r.db.NewSelect().
		Model((*models.Session)(nil)).
		Where("site_id = ?", siteID).
		Where("enter_time >= ?", fromUnix).
		Where("enter_time <= ?", toUnix).
		Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get session count: %w", err)
	}
	return count, nil
}

func (r *AnalyticsRepository) GetBounceRate(ctx context.Context, siteID int64, from, to time.Time) (float64, error) {
	var result struct {
		Total   int
		Bounced int
	}
	fromUnix := from.Unix()
	toUnix := to.Unix()

	dialect := fmt.Sprint(r.db.Dialect().Name())
	var bouncedExpr string
	if dialect == "pg" || dialect == "postgres" || dialect == "postgresql" {
		// PostgreSQL: Use FILTER clause
		bouncedExpr = "COUNT(*) FILTER (WHERE page_view_count = 1)"
	} else {
		// SQLite: Use CASE
		bouncedExpr = "SUM(CASE WHEN page_view_count = 1 THEN 1 ELSE 0 END)"
	}

	err := r.db.NewSelect().
		Model((*models.Session)(nil)).
		ColumnExpr("COUNT(*) as total").
		ColumnExpr(bouncedExpr + " as bounced").
		Where("site_id = ?", siteID).
		Where("enter_time >= ?", fromUnix).
		Where("enter_time <= ?", toUnix).
		Scan(ctx, &result)
	if err != nil {
		return 0, fmt.Errorf("failed to get bounce rate: %w", err)
	}
	if result.Total == 0 {
		return 0, nil
	}
	return float64(result.Bounced) / float64(result.Total) * 100, nil
}

func (r *AnalyticsRepository) GetAvgSessionDuration(ctx context.Context, siteID int64, from, to time.Time) (float64, error) {
	var avg float64
	fromUnix := from.Unix()
	toUnix := to.Unix()
	err := r.db.NewSelect().
		Model((*models.Session)(nil)).
		ColumnExpr("COALESCE(AVG((exit_time - enter_time) * 1.0), 0.0)").
		Where("site_id = ?", siteID).
		Where("enter_time >= ?", fromUnix).
		Where("enter_time <= ?", toUnix).
		Where("page_view_count > 1").
		Scan(ctx, &avg)
	if err != nil {
		return 0, fmt.Errorf("failed to get average session duration: %w", err)
	}
	return avg, nil
}

type PageStats struct {
	Path     string
	Views    int
	Visitors int
}

func (r *AnalyticsRepository) GetTopPages(ctx context.Context, siteID int64, from, to time.Time, limit int) ([]PageStats, error) {
	var stats []PageStats
	fromUnix := from.Unix()
	toUnix := to.Unix()
	err := r.db.NewSelect().
		Model((*models.Event)(nil)).
		Join("INNER JOIN sessions s ON e.session_id = s.id").
		ColumnExpr("e.path").
		ColumnExpr("COUNT(*) as views").
		ColumnExpr("COUNT(DISTINCT e.session_id) as visitors").
		Where("s.site_id = ?", siteID).
		Where("e.type = ?", models.EventTypePageview).
		Where("e.time >= ?", fromUnix).
		Where("e.time <= ?", toUnix).
		Group("e.path").
		Order("views DESC").
		Limit(limit).
		Scan(ctx, &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to get top pages: %w", err)
	}
	return stats, nil
}

type ReferrerStats struct {
	Referrer string
	Visitors int
}

func (r *AnalyticsRepository) GetTopReferrers(ctx context.Context, siteID int64, from, to time.Time, limit int) ([]ReferrerStats, error) {
	var stats []ReferrerStats
	fromUnix := from.Unix()
	toUnix := to.Unix()
	err := r.db.NewSelect().
		Model((*models.Session)(nil)).
		ColumnExpr("COALESCE(NULLIF(referrer, ''), '(direct)') as referrer").
		ColumnExpr("COUNT(DISTINCT client_id) as visitors").
		Where("site_id = ?", siteID).
		Where("enter_time >= ?", fromUnix).
		Where("enter_time <= ?", toUnix).
		Group("referrer").
		Order("visitors DESC").
		Limit(limit).
		Scan(ctx, &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to get top referrers: %w", err)
	}
	return stats, nil
}

type BrowserStats struct {
	Browser  string
	Visitors int
}

func (r *AnalyticsRepository) GetBrowserStats(ctx context.Context, siteID int64, from, to time.Time, limit int) ([]BrowserStats, error) {
	var stats []BrowserStats
	fromUnix := from.Unix()
	toUnix := to.Unix()
	err := r.db.NewSelect().
		TableExpr("sessions s").
		Join("INNER JOIN clients c ON s.client_id = c.id").
		ColumnExpr("c.browser").
		ColumnExpr("COUNT(DISTINCT s.client_id) as visitors").
		Where("s.site_id = ?", siteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix).
		Where("c.browser != ''").
		Group("c.browser").
		Order("visitors DESC").
		Limit(limit).
		Scan(ctx, &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to get browser stats: %w", err)
	}
	return stats, nil
}

type DeviceStats struct {
	Device   string
	Visitors int
}

func (r *AnalyticsRepository) GetDeviceStats(ctx context.Context, siteID int64, from, to time.Time, limit int) ([]DeviceStats, error) {
	var stats []DeviceStats
	fromUnix := from.Unix()
	toUnix := to.Unix()
	err := r.db.NewSelect().
		TableExpr("sessions s").
		Join("INNER JOIN clients c ON s.client_id = c.id").
		ColumnExpr("c.device").
		ColumnExpr("COUNT(DISTINCT s.client_id) as visitors").
		Where("s.site_id = ?", siteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix).
		Where("c.device != ''").
		Group("c.device").
		Order("visitors DESC").
		Limit(limit).
		Scan(ctx, &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to get device stats: %w", err)
	}
	return stats, nil
}

type CountryStats struct {
	Country  string
	Visitors int
}

func (r *AnalyticsRepository) GetCountryStats(ctx context.Context, siteID int64, from, to time.Time, limit int) ([]CountryStats, error) {
	var stats []CountryStats
	fromUnix := from.Unix()
	toUnix := to.Unix()
	err := r.db.NewSelect().
		TableExpr("sessions s").
		Join("INNER JOIN clients c ON s.client_id = c.id").
		ColumnExpr("COALESCE(NULLIF(c.country, ''), 'Unknown') as country").
		ColumnExpr("COUNT(DISTINCT s.client_id) as visitors").
		Where("s.site_id = ?", siteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix).
		Group("c.country").
		Order("visitors DESC").
		Limit(limit).
		Scan(ctx, &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to get country stats: %w", err)
	}
	return stats, nil
}

type DailyVisitorStats struct {
	DateBucket int64 // Unix timestamp bucket (day or hour) - integer for performance
	Visitors   int
	PageViews  int
	Sessions   int
}

func (r *AnalyticsRepository) GetDailyStats(ctx context.Context, siteID int64, from, to time.Time) ([]DailyVisitorStats, error) {
	var stats []DailyVisitorStats
	fromUnix := from.Unix()
	toUnix := to.Unix()
	bucketExpr := r.timeBucketExpression(TimeBucketDaily) // Just "enter_day" - integer column
	err := r.db.NewSelect().
		Model((*models.Session)(nil)).
		ColumnExpr(bucketExpr + " as date_bucket"). // Return integer, no conversion!
		ColumnExpr("COUNT(DISTINCT client_id) as visitors").
		ColumnExpr("SUM(page_view_count) as page_views").
		ColumnExpr("COUNT(*) as sessions").
		Where("site_id = ?", siteID).
		Where("enter_time >= ?", fromUnix).
		Where("enter_time <= ?", toUnix).
		GroupExpr(bucketExpr).
		Order("date_bucket ASC").
		Scan(ctx, &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily stats: %w", err)
	}
	return stats, nil
}

type ActivePageStats struct {
	Path     string
	Visitors int
}

func (r *AnalyticsRepository) GetActivePages(ctx context.Context, siteID int64, since time.Time) ([]ActivePageStats, error) {
	var stats []ActivePageStats
	sinceUnix := since.Unix()
	err := r.db.NewSelect().
		Model((*models.Event)(nil)).
		Join("INNER JOIN sessions s ON e.session_id = s.id").
		ColumnExpr("e.path").
		ColumnExpr("COUNT(DISTINCT e.session_id) as visitors").
		Where("s.site_id = ?", siteID).
		Where("e.type = ?", models.EventTypePageview).
		Where("e.time >= ?", sinceUnix).
		Group("e.path").
		Order("visitors DESC").
		Limit(10).
		Scan(ctx, &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to get active pages: %w", err)
	}
	return stats, nil
}

func applySessionFilters(q *bun.SelectQuery, referrer, device, page, country []string) *bun.SelectQuery {
	if len(referrer) > 0 {
		// Apply referrer filter (empty string filters for direct traffic)
		q = q.Where("s.referrer IN (?)", bun.In(referrer))
	}
	if len(device) > 0 {
		// Join with clients if not already joined
		q = q.Where("s.client_id IN (SELECT id FROM clients WHERE device IN (?))", bun.In(device))
	}
	if len(page) > 0 {
		// Need to join with events to filter by page
		q = q.Where("s.id IN (SELECT DISTINCT session_id FROM events WHERE type = ? AND path IN (?))", models.EventTypePageview, bun.In(page))
	}
	if len(country) > 0 {
		// Join with clients if not already joined
		q = q.Where("s.client_id IN (SELECT id FROM clients WHERE country IN (?))", bun.In(normalizeCountryValues(country)))
	}
	return q
}

func applyEventFilters(q *bun.SelectQuery, referrer, device, page, country []string) *bun.SelectQuery {
	if len(page) > 0 {
		q = q.Where("e.path IN (?)", bun.In(page))
	}
	if len(referrer) > 0 || len(device) > 0 || len(country) > 0 {
		// Join with sessions for referrer/device/country filters
		if len(referrer) > 0 {
			q = q.Where("e.session_id IN (SELECT id FROM sessions WHERE referrer IN (?))", bun.In(referrer))
		}
		if len(device) > 0 {
			q = q.Where("e.session_id IN (SELECT s.id FROM sessions s INNER JOIN clients c ON s.client_id = c.id WHERE c.device IN (?))", bun.In(device))
		}
		if len(country) > 0 {
			q = q.Where("e.session_id IN (SELECT s.id FROM sessions s INNER JOIN clients c ON s.client_id = c.id WHERE c.country IN (?))", bun.In(normalizeCountryValues(country)))
		}
	}
	return q
}

func normalizeCountryValues(values []string) []string {
	normalized := make([]string, 0, len(values)+1)
	seen := make(map[string]struct{}, len(values)+1)
	hasUnknown := false

	for _, value := range values {
		if value == "" {
			continue
		}
		if value == "Unknown" {
			hasUnknown = true
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		normalized = append(normalized, value)
	}

	if hasUnknown {
		if _, ok := seen[""]; !ok {
			normalized = append(normalized, "")
		}
		if _, ok := seen["Unknown"]; !ok {
			normalized = append(normalized, "Unknown")
		}
	}

	return normalized
}

func (r *AnalyticsRepository) GetVisitorCountWithFilter(ctx context.Context, siteID int64, from, to time.Time, referrer, device, page, country []string) (int, error) {
	var count int
	fromUnix := from.Unix()
	toUnix := to.Unix()
	q := r.db.NewSelect().
		TableExpr("sessions s").
		ColumnExpr("COUNT(DISTINCT s.client_id)").
		Where("s.site_id = ?", siteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix)
	q = applySessionFilters(q, referrer, device, page, country)
	err := q.Scan(ctx, &count)
	if err != nil {
		return 0, fmt.Errorf("failed to get visitor count with filter: %w", err)
	}
	return count, nil
}

func (r *AnalyticsRepository) GetPageViewCountWithFilter(ctx context.Context, siteID int64, from, to time.Time, referrer, device, page, country []string) (int, error) {
	fromUnix := from.Unix()
	toUnix := to.Unix()
	q := r.db.NewSelect().
		TableExpr("events e").
		Join("INNER JOIN sessions s ON e.session_id = s.id").
		Where("s.site_id = ?", siteID).
		Where("e.type = ?", models.EventTypePageview).
		Where("e.time >= ?", fromUnix).
		Where("e.time <= ?", toUnix)
	q = applyEventFilters(q, referrer, device, page, country)
	count, err := q.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get page view count with filter: %w", err)
	}
	return count, nil
}

func (r *AnalyticsRepository) GetSessionCountWithFilter(ctx context.Context, siteID int64, from, to time.Time, referrer, device, page, country []string) (int, error) {
	fromUnix := from.Unix()
	toUnix := to.Unix()
	q := r.db.NewSelect().
		TableExpr("sessions s").
		Where("s.site_id = ?", siteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix)
	q = applySessionFilters(q, referrer, device, page, country)
	count, err := q.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get session count with filter: %w", err)
	}
	return count, nil
}

func (r *AnalyticsRepository) GetBounceRateWithFilter(ctx context.Context, siteID int64, from, to time.Time, referrer, device, page, country []string) (float64, error) {
	var result struct {
		Total   int
		Bounced int
	}
	fromUnix := from.Unix()
	toUnix := to.Unix()

	dialect := fmt.Sprint(r.db.Dialect().Name())
	var bouncedExpr string
	if dialect == "pg" || dialect == "postgres" || dialect == "postgresql" {
		// PostgreSQL: Use FILTER clause
		bouncedExpr = "COUNT(*) FILTER (WHERE s.page_view_count = 1)"
	} else {
		// SQLite: Use CASE
		bouncedExpr = "SUM(CASE WHEN s.page_view_count = 1 THEN 1 ELSE 0 END)"
	}

	q := r.db.NewSelect().
		TableExpr("sessions s").
		ColumnExpr("COUNT(*) as total").
		ColumnExpr(bouncedExpr + " as bounced").
		Where("s.site_id = ?", siteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix)
	q = applySessionFilters(q, referrer, device, page, country)
	err := q.Scan(ctx, &result)
	if err != nil {
		return 0, fmt.Errorf("failed to get bounce rate with filter: %w", err)
	}
	if result.Total == 0 {
		return 0, nil
	}
	return float64(result.Bounced) / float64(result.Total) * 100, nil
}

func (r *AnalyticsRepository) GetAvgSessionDurationWithFilter(ctx context.Context, siteID int64, from, to time.Time, referrer, device, page, country []string) (float64, error) {
	var avg float64
	fromUnix := from.Unix()
	toUnix := to.Unix()
	q := r.db.NewSelect().
		TableExpr("sessions s").
		ColumnExpr("COALESCE(AVG((s.exit_time - s.enter_time) * 1.0), 0.0)").
		Where("s.site_id = ?", siteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix).
		Where("s.page_view_count > 1")
	q = applySessionFilters(q, referrer, device, page, country)
	err := q.Scan(ctx, &avg)
	if err != nil {
		return 0, fmt.Errorf("failed to get average session duration with filter: %w", err)
	}
	return avg, nil
}

func (r *AnalyticsRepository) GetTopPagesWithFilter(ctx context.Context, siteID int64, from, to time.Time, limit int, referrer, device, page, country []string) ([]PageStats, error) {
	var stats []PageStats
	fromUnix := from.Unix()
	toUnix := to.Unix()
	q := r.db.NewSelect().
		TableExpr("events e").
		Join("INNER JOIN sessions s ON e.session_id = s.id").
		ColumnExpr("e.path").
		ColumnExpr("COUNT(*) as views").
		ColumnExpr("COUNT(DISTINCT e.session_id) as visitors").
		Where("s.site_id = ?", siteID).
		Where("e.type = ?", models.EventTypePageview).
		Where("e.time >= ?", fromUnix).
		Where("e.time <= ?", toUnix)
	q = applyEventFilters(q, referrer, device, page, country)
	err := q.Group("e.path").
		Order("views DESC", "e.path ASC").
		Limit(limit).
		Scan(ctx, &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to get top pages with filter: %w", err)
	}
	return stats, nil
}

func (r *AnalyticsRepository) GetTopReferrersWithFilter(ctx context.Context, siteID int64, from, to time.Time, limit int, referrer, device, page, country []string) ([]ReferrerStats, error) {
	var stats []ReferrerStats
	fromUnix := from.Unix()
	toUnix := to.Unix()
	q := r.db.NewSelect().
		TableExpr("sessions s").
		ColumnExpr("COALESCE(NULLIF(s.referrer, ''), '(direct)') as referrer").
		ColumnExpr("COUNT(DISTINCT s.client_id) as visitors").
		Where("s.site_id = ?", siteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix)
	q = applySessionFilters(q, referrer, device, page, country)
	err := q.Group("s.referrer").
		Order("visitors DESC", "referrer ASC").
		Limit(limit).
		Scan(ctx, &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to get top referrers with filter: %w", err)
	}
	return stats, nil
}

func (r *AnalyticsRepository) GetBrowserStatsWithFilter(ctx context.Context, siteID int64, from, to time.Time, limit int, referrer, device, page, country []string) ([]BrowserStats, error) {
	var stats []BrowserStats
	fromUnix := from.Unix()
	toUnix := to.Unix()
	q := r.db.NewSelect().
		TableExpr("sessions s").
		Join("INNER JOIN clients c ON s.client_id = c.id").
		ColumnExpr("c.browser").
		ColumnExpr("COUNT(DISTINCT s.client_id) as visitors").
		Where("s.site_id = ?", siteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix).
		Where("c.browser != ''")
	q = applySessionFilters(q, referrer, device, page, country)
	err := q.Group("c.browser").
		Order("visitors DESC", "c.browser ASC").
		Limit(limit).
		Scan(ctx, &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to get browser stats with filter: %w", err)
	}
	return stats, nil
}

func (r *AnalyticsRepository) GetDeviceStatsWithFilter(ctx context.Context, siteID int64, from, to time.Time, limit int, referrer, device, page, country []string) ([]DeviceStats, error) {
	var stats []DeviceStats
	fromUnix := from.Unix()
	toUnix := to.Unix()
	q := r.db.NewSelect().
		TableExpr("sessions s").
		Join("INNER JOIN clients c ON s.client_id = c.id").
		ColumnExpr("c.device").
		ColumnExpr("COUNT(DISTINCT s.client_id) as visitors").
		Where("s.site_id = ?", siteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix).
		Where("c.device != ''")
	q = applySessionFilters(q, referrer, device, page, country)
	err := q.Group("c.device").
		Order("visitors DESC", "c.device ASC").
		Limit(limit).
		Scan(ctx, &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to get device stats with filter: %w", err)
	}
	return stats, nil
}

func (r *AnalyticsRepository) GetCountryStatsWithFilter(ctx context.Context, siteID int64, from, to time.Time, limit int, referrer, device, page, country []string) ([]CountryStats, error) {
	var stats []CountryStats
	fromUnix := from.Unix()
	toUnix := to.Unix()
	q := r.db.NewSelect().
		TableExpr("sessions s").
		Join("INNER JOIN clients c ON s.client_id = c.id").
		ColumnExpr("COALESCE(NULLIF(c.country, ''), 'Unknown') as country").
		ColumnExpr("COUNT(DISTINCT s.client_id) as visitors").
		Where("s.site_id = ?", siteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix)
	q = applySessionFilters(q, referrer, device, page, country)
	err := q.Group("c.country").
		Order("visitors DESC", "country ASC").
		Limit(limit).
		Scan(ctx, &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to get country stats with filter: %w", err)
	}
	return stats, nil
}

func (r *AnalyticsRepository) GetDailyStatsWithFilter(ctx context.Context, siteID int64, from, to time.Time, referrer, device, page, country []string) ([]DailyVisitorStats, error) {
	return r.GetTimeSeriesStatsWithFilter(ctx, siteID, from, to, TimeBucketDaily, 0, referrer, device, page, country)
}

func (r *AnalyticsRepository) GetTopPagesWithFilterPaged(ctx context.Context, siteID int64, from, to time.Time, limit, offset int, referrer, device, page, country []string) ([]PageStats, int, error) {
	var stats []PageStats
	var total int
	fromUnix := from.Unix()
	toUnix := to.Unix()
	q := r.db.NewSelect().
		TableExpr("events e").
		Join("INNER JOIN sessions s ON e.session_id = s.id").
		ColumnExpr("e.path").
		ColumnExpr("COUNT(*) as views").
		ColumnExpr("COUNT(DISTINCT e.session_id) as visitors").
		Where("s.site_id = ?", siteID).
		Where("e.type = ?", models.EventTypePageview).
		Where("e.time >= ?", fromUnix).
		Where("e.time <= ?", toUnix)
	q = applyEventFilters(q, referrer, device, page, country)
	err := q.Group("e.path").
		Order("views DESC", "e.path ASC").
		Limit(limit).
		Offset(offset).
		Scan(ctx, &stats)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get top pages with filter paged: %w", err)
	}

	countQuery := r.db.NewSelect().
		TableExpr("events e").
		Join("INNER JOIN sessions s ON e.session_id = s.id").
		ColumnExpr("COUNT(DISTINCT e.path)").
		Where("s.site_id = ?", siteID).
		Where("e.type = ?", models.EventTypePageview).
		Where("e.time >= ?", fromUnix).
		Where("e.time <= ?", toUnix)
	countQuery = applyEventFilters(countQuery, referrer, device, page, country)
	err = countQuery.Scan(ctx, &total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count top pages with filter: %w", err)
	}
	return stats, total, nil
}

func (r *AnalyticsRepository) GetTopReferrersWithFilterPaged(ctx context.Context, siteID int64, from, to time.Time, limit, offset int, referrer, device, page, country []string) ([]ReferrerStats, int, error) {
	var stats []ReferrerStats
	var total int
	fromUnix := from.Unix()
	toUnix := to.Unix()
	q := r.db.NewSelect().
		TableExpr("sessions s").
		ColumnExpr("COALESCE(NULLIF(s.referrer, ''), '(direct)') as referrer").
		ColumnExpr("COUNT(DISTINCT s.client_id) as visitors").
		Where("s.site_id = ?", siteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix)
	q = applySessionFilters(q, referrer, device, page, country)
	err := q.Group("s.referrer").
		Order("visitors DESC", "referrer ASC").
		Limit(limit).
		Offset(offset).
		Scan(ctx, &stats)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get top referrers with filter paged: %w", err)
	}

	countQuery := r.db.NewSelect().
		TableExpr("sessions s").
		ColumnExpr("COUNT(DISTINCT COALESCE(NULLIF(s.referrer, ''), '(direct)'))").
		Where("s.site_id = ?", siteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix)
	countQuery = applySessionFilters(countQuery, referrer, device, page, country)
	err = countQuery.Scan(ctx, &total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count top referrers with filter: %w", err)
	}
	return stats, total, nil
}

func (r *AnalyticsRepository) GetDeviceStatsWithFilterPaged(ctx context.Context, siteID int64, from, to time.Time, limit, offset int, referrer, device, page, country []string) ([]DeviceStats, int, int, error) {
	var stats []DeviceStats
	var total int
	var totalVisitors int
	fromUnix := from.Unix()
	toUnix := to.Unix()
	q := r.db.NewSelect().
		TableExpr("sessions s").
		Join("INNER JOIN clients c ON s.client_id = c.id").
		ColumnExpr("c.device").
		ColumnExpr("COUNT(DISTINCT s.client_id) as visitors").
		Where("s.site_id = ?", siteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix).
		Where("c.device != ''")
	q = applySessionFilters(q, referrer, device, page, country)
	err := q.Group("c.device").
		Order("visitors DESC", "c.device ASC").
		Limit(limit).
		Offset(offset).
		Scan(ctx, &stats)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get device stats with filter paged: %w", err)
	}

	countQuery := r.db.NewSelect().
		TableExpr("sessions s").
		Join("INNER JOIN clients c ON s.client_id = c.id").
		ColumnExpr("COUNT(DISTINCT c.device)").
		Where("s.site_id = ?", siteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix).
		Where("c.device != ''")
	countQuery = applySessionFilters(countQuery, referrer, device, page, country)
	err = countQuery.Scan(ctx, &total)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to count devices with filter: %w", err)
	}

	deviceCounts := r.db.NewSelect().
		TableExpr("sessions s").
		Join("INNER JOIN clients c ON s.client_id = c.id").
		ColumnExpr("COUNT(DISTINCT s.client_id) as visitors").
		Where("s.site_id = ?", siteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix).
		Where("c.device != ''")
	deviceCounts = applySessionFilters(deviceCounts, referrer, device, page, country)
	deviceCounts = deviceCounts.Group("c.device")

	err = r.db.NewSelect().
		TableExpr("(?) as device_counts", deviceCounts).
		ColumnExpr("COALESCE(SUM(visitors), 0)").
		Scan(ctx, &totalVisitors)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to sum device visitors with filter: %w", err)
	}
	return stats, total, totalVisitors, nil
}

func (r *AnalyticsRepository) GetCountryStatsWithFilterPaged(ctx context.Context, siteID int64, from, to time.Time, limit, offset int, referrer, device, page, country []string) ([]CountryStats, int, int, error) {
	var stats []CountryStats
	var total int
	var totalVisitors int
	fromUnix := from.Unix()
	toUnix := to.Unix()
	q := r.db.NewSelect().
		TableExpr("sessions s").
		Join("INNER JOIN clients c ON s.client_id = c.id").
		ColumnExpr("COALESCE(NULLIF(c.country, ''), 'Unknown') as country").
		ColumnExpr("COUNT(DISTINCT s.client_id) as visitors").
		Where("s.site_id = ?", siteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix)
	q = applySessionFilters(q, referrer, device, page, country)
	err := q.Group("c.country").
		Order("visitors DESC", "country ASC").
		Limit(limit).
		Offset(offset).
		Scan(ctx, &stats)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get country stats with filter paged: %w", err)
	}

	countQuery := r.db.NewSelect().
		TableExpr("sessions s").
		Join("INNER JOIN clients c ON s.client_id = c.id").
		ColumnExpr("COUNT(DISTINCT COALESCE(NULLIF(c.country, ''), 'Unknown'))").
		Where("s.site_id = ?", siteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix)
	countQuery = applySessionFilters(countQuery, referrer, device, page, country)
	err = countQuery.Scan(ctx, &total)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to count countries with filter: %w", err)
	}

	countryCounts := r.db.NewSelect().
		TableExpr("sessions s").
		Join("INNER JOIN clients c ON s.client_id = c.id").
		ColumnExpr("COUNT(DISTINCT s.client_id) as visitors").
		Where("s.site_id = ?", siteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix)
	countryCounts = applySessionFilters(countryCounts, referrer, device, page, country)
	countryCounts = countryCounts.Group("c.country")

	err = r.db.NewSelect().
		TableExpr("(?) as country_counts", countryCounts).
		ColumnExpr("COALESCE(SUM(visitors), 0)").
		Scan(ctx, &totalVisitors)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to sum country visitors with filter: %w", err)
	}
	return stats, total, totalVisitors, nil
}

type TimeBucket string

const (
	TimeBucketDaily  TimeBucket = "daily"
	TimeBucketHourly TimeBucket = "hourly"
)

func (r *AnalyticsRepository) GetTimeSeriesStatsWithFilter(ctx context.Context, siteID int64, from, to time.Time, bucket TimeBucket, limit int, referrer, device, page, country []string) ([]DailyVisitorStats, error) {
	var stats []DailyVisitorStats
	fromUnix := from.Unix()
	toUnix := to.Unix()
	bucketExpr := r.timeBucketExpression(bucket) // Just "enter_day" or "enter_hour" - integer column
	q := r.db.NewSelect().
		TableExpr("sessions s").
		ColumnExpr(bucketExpr + " as date_bucket"). // Return integer, no conversion!
		ColumnExpr("COUNT(DISTINCT s.client_id) as visitors").
		ColumnExpr("SUM(s.page_view_count) as page_views").
		ColumnExpr("COUNT(*) as sessions").
		Where("s.site_id = ?", siteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix)
	q = applySessionFilters(q, referrer, device, page, country)
	q = q.GroupExpr(bucketExpr)
	if limit > 0 {
		q = q.Order("date_bucket DESC").Limit(limit)
	} else {
		q = q.Order("date_bucket ASC")
	}
	err := q.Scan(ctx, &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to get time series stats with filter: %w", err)
	}
	return stats, nil
}

func (r *AnalyticsRepository) timeBucketExpression(bucket TimeBucket) string {
	// Returns the column name for time bucketing
	// No need for date functions - we use pre-calculated integer buckets!
	if bucket == TimeBucketHourly {
		return "enter_hour"
	}
	// Daily bucket
	return "enter_day"
}
