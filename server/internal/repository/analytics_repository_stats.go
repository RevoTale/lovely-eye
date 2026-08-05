package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/lovely-eye/server/internal/models"
	"github.com/uptrace/bun"
)

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

func (r *AnalyticsRepository) GetEventCountWithFilter(ctx context.Context, query AnalyticsQuery) (int, error) {
	fromUnix := query.From.Unix()
	toUnix := query.To.Unix()
	q := r.db.NewSelect().
		TableExpr("events e").
		Join("INNER JOIN sessions s ON e.session_id = s.id").
		Where("s.site_id = ?", query.SiteID).
		Where("e.time >= ?", fromUnix).
		Where("e.time <= ?", toUnix)
	q = applyEventFilters(q, query.Filter)
	q = applyEventNamePathFilters(q, query.Filter)
	count, err := q.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get event count with filter: %w", err)
	}
	return count, nil
}

type EventCountResult struct {
	Count int

	EventID int64
}

// GetEventCountsGrouped returns event counts grouped by definition with the most recent event for each
// This is used for the eventCounts GraphQL query to avoid fetching 200 full events just for counting
func (r *AnalyticsRepository) GetEventCountsGrouped(ctx context.Context, query AnalyticsQuery) ([]EventCountResult, error) {
	fromUnix := query.From.Unix()
	toUnix := query.To.Unix()

	var results []EventCountResult
	q := r.db.NewSelect().
		TableExpr("events e").
		Join("INNER JOIN sessions s ON e.session_id = s.id").
		ColumnExpr("COUNT(*) as count").
		ColumnExpr("MAX(e.id) as event_id").
		Where("s.site_id = ?", query.SiteID).
		Where("e.definition_id IS NOT NULL").
		Where("e.time >= ?", fromUnix).
		Where("e.time <= ?", toUnix).
		Group("e.definition_id").
		Order("count DESC")

	q = applyEventFilters(q, query.Filter)
	q = applyEventNamePathFilters(q, query.Filter)

	if query.Limit > 0 {
		q = q.Limit(query.Limit)
	}
	if query.Offset > 0 {
		q = q.Offset(query.Offset)
	}

	err := q.Scan(ctx, &results)
	if err != nil {
		return nil, fmt.Errorf("failed to get event counts grouped: %w", err)
	}

	return results, nil
}

func (r *AnalyticsRepository) GetEventsByIDs(ctx context.Context, eventIDs []int64) ([]*models.Event, error) {
	if len(eventIDs) == 0 {
		return []*models.Event{}, nil
	}

	var events []*models.Event
	err := r.db.NewSelect().
		Model(&events).
		Relation("Data.Field").
		Relation("Definition.Fields").
		Where("e.id IN (?)", bun.List(eventIDs)).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get events by IDs: %w", err)
	}
	return events, nil
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
		Where("e.definition_id IS NULL").
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

		bouncedExpr = "COUNT(*) FILTER (WHERE page_view_count = 1)"
	} else {

		bouncedExpr = "SUM(CASE WHEN page_view_count = 1 THEN 1 ELSE 0 END)"
	}

	err := r.db.NewSelect().
		Model((*models.Session)(nil)).
		ColumnExpr("COUNT(*) as total").
		ColumnExpr(bouncedExpr+" as bounced").
		Where("site_id = ?", siteID).
		Where("enter_time >= ?", fromUnix).
		Where("enter_time <= ?", toUnix).
		Where("page_view_count > 0").
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
		Where("page_view_count > 0").
		Where("exit_time > enter_time").
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
		ColumnExpr("COUNT(DISTINCT s.client_id) as visitors").
		Where("s.site_id = ?", siteID).
		Where("e.definition_id IS NULL").
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
	Browser  models.ClientBrowser
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
		Where("c.browser != ?", models.ClientBrowserUnknown).
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
	Device   models.ClientDevice
	Visitors int
}

type OperatingSystemStats struct {
	OS       models.ClientOS
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
		Where("c.device != ?", models.ClientDeviceUnknown).
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
	CountryCode string
	Visitors    int
}

func (r *AnalyticsRepository) GetCountryStats(ctx context.Context, siteID int64, from, to time.Time, limit int) ([]CountryStats, error) {
	var stats []CountryStats
	fromUnix := from.Unix()
	toUnix := to.Unix()
	err := r.db.NewSelect().
		TableExpr("sessions s").
		Join("INNER JOIN clients c ON s.client_id = c.id").
		ColumnExpr("COALESCE(NULLIF(c.country, ''), '-') as country_code").
		ColumnExpr("COUNT(DISTINCT s.client_id) as visitors").
		Where("s.site_id = ?", siteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix).
		GroupExpr("COALESCE(NULLIF(c.country, ''), '-')").
		Order("visitors DESC", "country_code ASC").
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
	stats, err := r.GetTimeSeriesStatsWithFilter(ctx, AnalyticsQuery{
		SiteID: siteID,
		From:   from,
		To:     to,
		Bucket: TimeBucketDaily,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get daily stats: %w", err)
	}
	return stats, nil
}

type ActivePageStats struct {
	Path     string
	Visitors int
}

func (r *AnalyticsRepository) GetActivePages(ctx context.Context, siteID int64, since time.Time, limit, offset int) ([]ActivePageStats, error) {
	var stats []ActivePageStats
	sinceUnix := since.Unix()
	q := r.db.NewSelect().
		Model((*models.Event)(nil)).
		Join("INNER JOIN sessions s ON e.session_id = s.id").
		ColumnExpr("e.path").
		ColumnExpr("COUNT(DISTINCT s.client_id) as visitors").
		Where("s.site_id = ?", siteID).
		Where("e.definition_id IS NULL").
		Where("e.time >= ?", sinceUnix).
		Group("e.path").
		Order("visitors DESC", "e.path ASC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}
	err := q.Scan(ctx, &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to get active pages: %w", err)
	}
	return stats, nil
}
