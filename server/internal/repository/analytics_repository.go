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
	session := new(models.Session)
	err := r.db.NewSelect().
		Model(session).
		Where("site_id = ?", siteID).
		Where("visitor_id = ?", visitorID).
		Where("last_seen_at > ?", since).
		Order("last_seen_at DESC").
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

func (r *AnalyticsRepository) CreatePageView(ctx context.Context, pageView *models.PageView) error {
	_, err := r.db.NewInsert().Model(pageView).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create page view: %w", err)
	}
	return nil
}

func (r *AnalyticsRepository) GetPageViews(ctx context.Context, siteID int64, from, to time.Time, limit, offset int) ([]*models.PageView, error) {
	var pageViews []*models.PageView
	err := r.db.NewSelect().
		Model(&pageViews).
		Where("site_id = ?", siteID).
		Where("created_at >= ?", from).
		Where("created_at <= ?", to).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get page views: %w", err)
	}
	return pageViews, nil
}

func (r *AnalyticsRepository) GetRecentPageView(ctx context.Context, siteID int64, visitorID, path string, since time.Time) (*models.PageView, error) {
	pageView := new(models.PageView)
	err := r.db.NewSelect().
		Model(pageView).
		Where("site_id = ?", siteID).
		Where("visitor_id = ?", visitorID).
		Where("path = ?", path).
		Where("created_at > ?", since).
		Order("created_at DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent page view: %w", err)
	}
	return pageView, nil
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
	err := r.db.NewSelect().
		Model(&events).
		Where("site_id = ?", siteID).
		Where("created_at >= ?", from).
		Where("created_at <= ?", to).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}
	return events, nil
}

func (r *AnalyticsRepository) GetEventCount(ctx context.Context, siteID int64, from, to time.Time) (int, error) {
	count, err := r.db.NewSelect().
		Model((*models.Event)(nil)).
		Where("site_id = ?", siteID).
		Where("created_at >= ?", from).
		Where("created_at <= ?", to).
		Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get event count: %w", err)
	}
	return count, nil
}

func (r *AnalyticsRepository) GetVisitorCount(ctx context.Context, siteID int64, from, to time.Time) (int, error) {
	var count int
	err := r.db.NewSelect().
		Model((*models.Session)(nil)).
		ColumnExpr("COUNT(DISTINCT visitor_id)").
		Where("site_id = ?", siteID).
		Where("started_at >= ?", from).
		Where("started_at <= ?", to).
		Scan(ctx, &count)
	if err != nil {
		return 0, fmt.Errorf("failed to get visitor count: %w", err)
	}
	return count, nil
}

func (r *AnalyticsRepository) GetPageViewCount(ctx context.Context, siteID int64, from, to time.Time) (int, error) {
	count, err := r.db.NewSelect().
		Model((*models.PageView)(nil)).
		Where("site_id = ?", siteID).
		Where("created_at >= ?", from).
		Where("created_at <= ?", to).
		Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get page view count: %w", err)
	}
	return count, nil
}

func (r *AnalyticsRepository) GetSessionCount(ctx context.Context, siteID int64, from, to time.Time) (int, error) {
	count, err := r.db.NewSelect().
		Model((*models.Session)(nil)).
		Where("site_id = ?", siteID).
		Where("started_at >= ?", from).
		Where("started_at <= ?", to).
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
	err := r.db.NewSelect().
		Model((*models.Session)(nil)).
		ColumnExpr("COUNT(*) as total").
		ColumnExpr("SUM(CASE WHEN is_bounce THEN 1 ELSE 0 END) as bounced").
		Where("site_id = ?", siteID).
		Where("started_at >= ?", from).
		Where("started_at <= ?", to).
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
	err := r.db.NewSelect().
		Model((*models.Session)(nil)).
		ColumnExpr("CAST(COALESCE(AVG(duration), 0) AS REAL)").
		Where("site_id = ?", siteID).
		Where("started_at >= ?", from).
		Where("started_at <= ?", to).
		Where("is_bounce = false").
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
	err := r.db.NewSelect().
		Model((*models.PageView)(nil)).
		ColumnExpr("path").
		ColumnExpr("COUNT(*) as views").
		ColumnExpr("COUNT(DISTINCT visitor_id) as visitors").
		Where("site_id = ?", siteID).
		Where("created_at >= ?", from).
		Where("created_at <= ?", to).
		Group("path").
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
	err := r.db.NewSelect().
		Model((*models.Session)(nil)).
		ColumnExpr("COALESCE(NULLIF(referrer, ''), '(direct)') as referrer").
		ColumnExpr("COUNT(DISTINCT visitor_id) as visitors").
		Where("site_id = ?", siteID).
		Where("started_at >= ?", from).
		Where("started_at <= ?", to).
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
	err := r.db.NewSelect().
		Model((*models.Session)(nil)).
		ColumnExpr("browser").
		ColumnExpr("COUNT(DISTINCT visitor_id) as visitors").
		Where("site_id = ?", siteID).
		Where("started_at >= ?", from).
		Where("started_at <= ?", to).
		Where("browser != ''").
		Group("browser").
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
	err := r.db.NewSelect().
		Model((*models.Session)(nil)).
		ColumnExpr("device").
		ColumnExpr("COUNT(DISTINCT visitor_id) as visitors").
		Where("site_id = ?", siteID).
		Where("started_at >= ?", from).
		Where("started_at <= ?", to).
		Where("device != ''").
		Group("device").
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
	err := r.db.NewSelect().
		Model((*models.Session)(nil)).
		ColumnExpr("COALESCE(NULLIF(country, ''), 'Unknown') as country").
		ColumnExpr("COUNT(DISTINCT visitor_id) as visitors").
		Where("site_id = ?", siteID).
		Where("started_at >= ?", from).
		Where("started_at <= ?", to).
		Group("country").
		Order("visitors DESC").
		Limit(limit).
		Scan(ctx, &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to get country stats: %w", err)
	}
	return stats, nil
}

type DailyVisitorStats struct {
	Date      time.Time
	Visitors  int
	PageViews int
	Sessions  int
}

func (r *AnalyticsRepository) GetDailyStats(ctx context.Context, siteID int64, from, to time.Time) ([]DailyVisitorStats, error) {
	var stats []DailyVisitorStats
	err := r.db.NewSelect().
		Model((*models.Session)(nil)).
		ColumnExpr("DATE(started_at) as date").
		ColumnExpr("COUNT(DISTINCT visitor_id) as visitors").
		ColumnExpr("SUM(page_views) as page_views").
		ColumnExpr("COUNT(*) as sessions").
		Where("site_id = ?", siteID).
		Where("started_at >= ?", from).
		Where("started_at <= ?", to).
		Group("DATE(started_at)").
		Order("date ASC").
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
	err := r.db.NewSelect().
		Model((*models.PageView)(nil)).
		ColumnExpr("path").
		ColumnExpr("COUNT(DISTINCT visitor_id) as visitors").
		Where("site_id = ?", siteID).
		Where("created_at >= ?", since).
		Group("path").
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
		q = q.Where("referrer IN (?)", bun.In(referrer))
	}
	if len(device) > 0 {
		q = q.Where("device IN (?)", bun.In(device))
	}
	if len(page) > 0 {
		// Need to join with page_views to filter by page
		q = q.Where("id IN (SELECT DISTINCT session_id FROM page_views WHERE path IN (?))", bun.In(page))
	}
	if len(country) > 0 {
		q = q.Where("country IN (?)", bun.In(normalizeCountryValues(country)))
	}
	return q
}

func applyPageViewFilters(q *bun.SelectQuery, referrer, device, page, country []string) *bun.SelectQuery {
	if len(page) > 0 {
		q = q.Where("path IN (?)", bun.In(page))
	}
	if len(referrer) > 0 || len(device) > 0 {
		// Join with sessions for referrer/device filters
		if len(referrer) > 0 {
			// Apply referrer filter (empty string filters for direct traffic)
			q = q.Where("session_id IN (SELECT id FROM sessions WHERE referrer IN (?))", bun.In(referrer))
		}
		if len(device) > 0 {
			q = q.Where("session_id IN (SELECT id FROM sessions WHERE device IN (?))", bun.In(device))
		}
	}
	if len(country) > 0 {
		q = q.Where("session_id IN (SELECT id FROM sessions WHERE country IN (?))", bun.In(normalizeCountryValues(country)))
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
	q := r.db.NewSelect().
		Model((*models.Session)(nil)).
		ColumnExpr("COUNT(DISTINCT visitor_id)").
		Where("site_id = ?", siteID).
		Where("started_at >= ?", from).
		Where("started_at <= ?", to)
	q = applySessionFilters(q, referrer, device, page, country)
	err := q.Scan(ctx, &count)
	if err != nil {
		return 0, fmt.Errorf("failed to get visitor count with filter: %w", err)
	}
	return count, nil
}

func (r *AnalyticsRepository) GetPageViewCountWithFilter(ctx context.Context, siteID int64, from, to time.Time, referrer, device, page, country []string) (int, error) {
	q := r.db.NewSelect().
		Model((*models.PageView)(nil)).
		Where("site_id = ?", siteID).
		Where("created_at >= ?", from).
		Where("created_at <= ?", to)
	q = applyPageViewFilters(q, referrer, device, page, country)
	count, err := q.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get page view count with filter: %w", err)
	}
	return count, nil
}

func (r *AnalyticsRepository) GetSessionCountWithFilter(ctx context.Context, siteID int64, from, to time.Time, referrer, device, page, country []string) (int, error) {
	q := r.db.NewSelect().
		Model((*models.Session)(nil)).
		Where("site_id = ?", siteID).
		Where("started_at >= ?", from).
		Where("started_at <= ?", to)
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
	q := r.db.NewSelect().
		Model((*models.Session)(nil)).
		ColumnExpr("COUNT(*) as total").
		ColumnExpr("SUM(CASE WHEN is_bounce THEN 1 ELSE 0 END) as bounced").
		Where("site_id = ?", siteID).
		Where("started_at >= ?", from).
		Where("started_at <= ?", to)
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
	q := r.db.NewSelect().
		Model((*models.Session)(nil)).
		ColumnExpr("COALESCE(AVG(duration * 1.0), 0.0) ").
		Where("site_id = ?", siteID).
		Where("started_at >= ?", from).
		Where("started_at <= ?", to).
		Where("is_bounce = false")
	q = applySessionFilters(q, referrer, device, page, country)
	err := q.Scan(ctx, &avg)
	if err != nil {
		return 0, fmt.Errorf("failed to get average session duration with filter: %w", err)
	}
	return avg, nil
}

func (r *AnalyticsRepository) GetTopPagesWithFilter(ctx context.Context, siteID int64, from, to time.Time, limit int, referrer, device, page, country []string) ([]PageStats, error) {
	var stats []PageStats
	q := r.db.NewSelect().
		Model((*models.PageView)(nil)).
		ColumnExpr("path").
		ColumnExpr("COUNT(*) as views").
		ColumnExpr("COUNT(DISTINCT visitor_id) as visitors").
		Where("site_id = ?", siteID).
		Where("created_at >= ?", from).
		Where("created_at <= ?", to)
	q = applyPageViewFilters(q, referrer, device, page, country)
	err := q.Group("path").
		Order("views DESC", "path ASC").
		Limit(limit).
		Scan(ctx, &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to get top pages with filter: %w", err)
	}
	return stats, nil
}

func (r *AnalyticsRepository) GetTopReferrersWithFilter(ctx context.Context, siteID int64, from, to time.Time, limit int, referrer, device, page, country []string) ([]ReferrerStats, error) {
	var stats []ReferrerStats
	q := r.db.NewSelect().
		Model((*models.Session)(nil)).
		ColumnExpr("COALESCE(NULLIF(referrer, ''), '(direct)') as referrer").
		ColumnExpr("COUNT(DISTINCT visitor_id) as visitors").
		Where("site_id = ?", siteID).
		Where("started_at >= ?", from).
		Where("started_at <= ?", to)
	q = applySessionFilters(q, referrer, device, page, country)
	err := q.Group("referrer").
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
	q := r.db.NewSelect().
		Model((*models.Session)(nil)).
		ColumnExpr("browser").
		ColumnExpr("COUNT(DISTINCT visitor_id) as visitors").
		Where("site_id = ?", siteID).
		Where("started_at >= ?", from).
		Where("started_at <= ?", to).
		Where("browser != ''")
	q = applySessionFilters(q, referrer, device, page, country)
	err := q.Group("browser").
		Order("visitors DESC", "browser ASC").
		Limit(limit).
		Scan(ctx, &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to get browser stats with filter: %w", err)
	}
	return stats, nil
}

func (r *AnalyticsRepository) GetDeviceStatsWithFilter(ctx context.Context, siteID int64, from, to time.Time, limit int, referrer, device, page, country []string) ([]DeviceStats, error) {
	var stats []DeviceStats
	q := r.db.NewSelect().
		Model((*models.Session)(nil)).
		ColumnExpr("device").
		ColumnExpr("COUNT(DISTINCT visitor_id) as visitors").
		Where("site_id = ?", siteID).
		Where("started_at >= ?", from).
		Where("started_at <= ?", to).
		Where("device != ''")
	q = applySessionFilters(q, referrer, device, page, country)
	err := q.Group("device").
		Order("visitors DESC", "device ASC").
		Limit(limit).
		Scan(ctx, &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to get device stats with filter: %w", err)
	}
	return stats, nil
}

func (r *AnalyticsRepository) GetCountryStatsWithFilter(ctx context.Context, siteID int64, from, to time.Time, limit int, referrer, device, page, country []string) ([]CountryStats, error) {
	var stats []CountryStats
	q := r.db.NewSelect().
		Model((*models.Session)(nil)).
		ColumnExpr("COALESCE(NULLIF(country, ''), 'Unknown') as country").
		ColumnExpr("COUNT(DISTINCT visitor_id) as visitors").
		Where("site_id = ?", siteID).
		Where("started_at >= ?", from).
		Where("started_at <= ?", to)
	q = applySessionFilters(q, referrer, device, page, country)
	err := q.Group("country").
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
	q := r.db.NewSelect().
		Model((*models.PageView)(nil)).
		ColumnExpr("path").
		ColumnExpr("COUNT(*) as views").
		ColumnExpr("COUNT(DISTINCT visitor_id) as visitors").
		Where("site_id = ?", siteID).
		Where("created_at >= ?", from).
		Where("created_at <= ?", to)
	q = applyPageViewFilters(q, referrer, device, page, country)
	err := q.Group("path").
		Order("views DESC", "path ASC").
		Limit(limit).
		Offset(offset).
		Scan(ctx, &stats)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get top pages with filter paged: %w", err)
	}

	countQuery := r.db.NewSelect().
		Model((*models.PageView)(nil)).
		ColumnExpr("COUNT(DISTINCT path)").
		Where("site_id = ?", siteID).
		Where("created_at >= ?", from).
		Where("created_at <= ?", to)
	countQuery = applyPageViewFilters(countQuery, referrer, device, page, country)
	err = countQuery.Scan(ctx, &total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count top pages with filter: %w", err)
	}
	return stats, total, nil
}

func (r *AnalyticsRepository) GetTopReferrersWithFilterPaged(ctx context.Context, siteID int64, from, to time.Time, limit, offset int, referrer, device, page, country []string) ([]ReferrerStats, int, error) {
	var stats []ReferrerStats
	var total int
	q := r.db.NewSelect().
		Model((*models.Session)(nil)).
		ColumnExpr("COALESCE(NULLIF(referrer, ''), '(direct)') as referrer").
		ColumnExpr("COUNT(DISTINCT visitor_id) as visitors").
		Where("site_id = ?", siteID).
		Where("started_at >= ?", from).
		Where("started_at <= ?", to)
	q = applySessionFilters(q, referrer, device, page, country)
	err := q.Group("referrer").
		Order("visitors DESC", "referrer ASC").
		Limit(limit).
		Offset(offset).
		Scan(ctx, &stats)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get top referrers with filter paged: %w", err)
	}

	countQuery := r.db.NewSelect().
		Model((*models.Session)(nil)).
		ColumnExpr("COUNT(DISTINCT COALESCE(NULLIF(referrer, ''), '(direct)'))").
		Where("site_id = ?", siteID).
		Where("started_at >= ?", from).
		Where("started_at <= ?", to)
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
	q := r.db.NewSelect().
		Model((*models.Session)(nil)).
		ColumnExpr("device").
		ColumnExpr("COUNT(DISTINCT visitor_id) as visitors").
		Where("site_id = ?", siteID).
		Where("started_at >= ?", from).
		Where("started_at <= ?", to).
		Where("device != ''")
	q = applySessionFilters(q, referrer, device, page, country)
	err := q.Group("device").
		Order("visitors DESC", "device ASC").
		Limit(limit).
		Offset(offset).
		Scan(ctx, &stats)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get device stats with filter paged: %w", err)
	}

	countQuery := r.db.NewSelect().
		Model((*models.Session)(nil)).
		ColumnExpr("COUNT(DISTINCT device)").
		Where("site_id = ?", siteID).
		Where("started_at >= ?", from).
		Where("started_at <= ?", to).
		Where("device != ''")
	countQuery = applySessionFilters(countQuery, referrer, device, page, country)
	err = countQuery.Scan(ctx, &total)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to count devices with filter: %w", err)
	}

	deviceCounts := r.db.NewSelect().
		Model((*models.Session)(nil)).
		ColumnExpr("COUNT(DISTINCT visitor_id) as visitors").
		Where("site_id = ?", siteID).
		Where("started_at >= ?", from).
		Where("started_at <= ?", to).
		Where("device != ''")
	deviceCounts = applySessionFilters(deviceCounts, referrer, device, page, country)
	deviceCounts = deviceCounts.Group("device")

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
	q := r.db.NewSelect().
		Model((*models.Session)(nil)).
		ColumnExpr("COALESCE(NULLIF(country, ''), 'Unknown') as country").
		ColumnExpr("COUNT(DISTINCT visitor_id) as visitors").
		Where("site_id = ?", siteID).
		Where("started_at >= ?", from).
		Where("started_at <= ?", to)
	q = applySessionFilters(q, referrer, device, page, country)
	err := q.Group("country").
		Order("visitors DESC", "country ASC").
		Limit(limit).
		Offset(offset).
		Scan(ctx, &stats)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get country stats with filter paged: %w", err)
	}

	countQuery := r.db.NewSelect().
		Model((*models.Session)(nil)).
		ColumnExpr("COUNT(DISTINCT COALESCE(NULLIF(country, ''), 'Unknown'))").
		Where("site_id = ?", siteID).
		Where("started_at >= ?", from).
		Where("started_at <= ?", to)
	countQuery = applySessionFilters(countQuery, referrer, device, page, country)
	err = countQuery.Scan(ctx, &total)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to count countries with filter: %w", err)
	}

	countryCounts := r.db.NewSelect().
		Model((*models.Session)(nil)).
		ColumnExpr("COUNT(DISTINCT visitor_id) as visitors").
		Where("site_id = ?", siteID).
		Where("started_at >= ?", from).
		Where("started_at <= ?", to)
	countryCounts = applySessionFilters(countryCounts, referrer, device, page, country)
	countryCounts = countryCounts.Group("country")

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
	bucketExpr := r.timeBucketExpression(bucket)
	q := r.db.NewSelect().
		Model((*models.Session)(nil)).
		ColumnExpr(bucketExpr+" as date").
		ColumnExpr("COUNT(DISTINCT visitor_id) as visitors").
		ColumnExpr("SUM(page_views) as page_views").
		ColumnExpr("COUNT(*) as sessions").
		Where("site_id = ?", siteID).
		Where("started_at >= ?", from).
		Where("started_at <= ?", to)
	q = applySessionFilters(q, referrer, device, page, country)
	q = q.Group(bucketExpr)
	if limit > 0 {
		q = q.Order("date DESC").Limit(limit)
	} else {
		q = q.Order("date ASC")
	}
	err := q.Scan(ctx, &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to get time series stats with filter: %w", err)
	}
	return stats, nil
}

func (r *AnalyticsRepository) timeBucketExpression(bucket TimeBucket) string {
	dialect := fmt.Sprint(r.db.Dialect().Name())
	if bucket == TimeBucketHourly {
		if dialect == "pg" || dialect == "postgres" || dialect == "postgresql" {
			return "date_trunc('hour', started_at)"
		}
		return "strftime('%Y-%m-%d %H:00:00', started_at)"
	}
	if dialect == "pg" || dialect == "postgres" || dialect == "postgresql" {
		return "date_trunc('day', started_at)"
	}
	return "DATE(started_at)"
}
