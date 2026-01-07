package repository

import (
	"context"
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

// Session methods
func (r *AnalyticsRepository) CreateSession(ctx context.Context, session *models.Session) error {
	_, err := r.db.NewInsert().Model(session).Exec(ctx)
	return err
}

func (r *AnalyticsRepository) GetSession(ctx context.Context, id int64) (*models.Session, error) {
	session := new(models.Session)
	err := r.db.NewSelect().Model(session).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	return session, nil
}

func (r *AnalyticsRepository) UpdateSession(ctx context.Context, session *models.Session) error {
	_, err := r.db.NewUpdate().Model(session).WherePK().Exec(ctx)
	return err
}

// PageView methods
func (r *AnalyticsRepository) CreatePageView(ctx context.Context, pageView *models.PageView) error {
	_, err := r.db.NewInsert().Model(pageView).Exec(ctx)
	return err
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
	return pageViews, err
}

// Event methods
func (r *AnalyticsRepository) CreateEvent(ctx context.Context, event *models.Event) error {
	_, err := r.db.NewInsert().Model(event).Exec(ctx)
	return err
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
	return events, err
}

// Aggregation methods
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
		return 0, err
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
	return count, err
}

func (r *AnalyticsRepository) GetSessionCount(ctx context.Context, siteID int64, from, to time.Time) (int, error) {
	count, err := r.db.NewSelect().
		Model((*models.Session)(nil)).
		Where("site_id = ?", siteID).
		Where("started_at >= ?", from).
		Where("started_at <= ?", to).
		Count(ctx)
	return count, err
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
	if err != nil || result.Total == 0 {
		return 0, err
	}
	return float64(result.Bounced) / float64(result.Total) * 100, nil
}

func (r *AnalyticsRepository) GetAvgSessionDuration(ctx context.Context, siteID int64, from, to time.Time) (float64, error) {
	var avg float64
	err := r.db.NewSelect().
		Model((*models.Session)(nil)).
		ColumnExpr("COALESCE(AVG(duration), 0)").
		Where("site_id = ?", siteID).
		Where("started_at >= ?", from).
		Where("started_at <= ?", to).
		Where("is_bounce = false").
		Scan(ctx, &avg)
	return avg, err
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
	return stats, err
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
	return stats, err
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
	return stats, err
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
	return stats, err
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
	return stats, err
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
	return stats, err
}
