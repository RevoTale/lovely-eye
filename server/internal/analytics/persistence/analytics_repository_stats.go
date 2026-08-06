package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

func (r *Repository) GetEventCount(ctx context.Context, siteID int64, from, to time.Time) (int, error) {
	fromUnix := from.Unix()
	toUnix := to.Unix()
	count, err := r.db.NewSelect().
		Model((*Event)(nil)).
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

func (r *Repository) GetEventCountWithFilter(ctx context.Context, query AnalyticsQuery) (int, error) {
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
	Count   int
	Total   int
	EventID int64
}

// GetEventCountsGrouped returns event counts grouped by definition with the most recent event for each
// This is used for the eventCounts GraphQL query to avoid fetching 200 full events just for counting
func (r *Repository) GetEventCountsGrouped(ctx context.Context, query AnalyticsQuery) ([]EventCountResult, int, error) {
	fromUnix := query.From.Unix()
	toUnix := query.To.Unix()

	var results []EventCountResult
	q := r.db.NewSelect().
		TableExpr("events e").
		Join("INNER JOIN sessions s ON e.session_id = s.id").
		ColumnExpr("COUNT(*) as count").
		ColumnExpr("MAX(e.id) as event_id").
		ColumnExpr("COUNT(*) OVER() as total").
		Where("s.site_id = ?", query.SiteID).
		Where("e.definition_id IS NOT NULL").
		Where("e.time >= ?", fromUnix).
		Where("e.time <= ?", toUnix).
		Group("e.definition_id")

	q = applyEventFilters(q, query.Filter)
	q = applyEventNamePathFilters(q, query.Filter)

	pagedQuery := q.Clone().Order("count DESC", "event_id DESC")
	if query.Limit > 0 {
		pagedQuery = pagedQuery.Limit(query.Limit)
	}
	if query.Offset > 0 {
		pagedQuery = pagedQuery.Offset(query.Offset)
	}

	err := pagedQuery.Scan(ctx, &results)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get event counts grouped: %w", err)
	}

	if len(results) > 0 {
		return results, results[0].Total, nil
	}
	if query.Offset <= 0 {
		return results, 0, nil
	}

	total, err := r.groupedRowCount(ctx, q)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get event counts total: %w", err)
	}
	return results, total, nil
}

func (r *Repository) GetEventsByIDs(ctx context.Context, eventIDs []int64) ([]*Event, error) {
	if len(eventIDs) == 0 {
		return []*Event{}, nil
	}

	var events []*Event
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

func (r *Repository) GetVisitorCount(ctx context.Context, siteID int64, from, to time.Time) (int, error) {
	var count int
	fromUnix := from.Unix()
	toUnix := to.Unix()
	err := r.db.NewSelect().
		Model((*Session)(nil)).
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

func (r *Repository) GetBounceRate(ctx context.Context, siteID int64, from, to time.Time) (float64, error) {
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
		Model((*Session)(nil)).
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

func (r *Repository) GetAvgSessionDuration(ctx context.Context, siteID int64, from, to time.Time) (float64, error) {
	var avg float64
	fromUnix := from.Unix()
	toUnix := to.Unix()
	err := r.db.NewSelect().
		Model((*Session)(nil)).
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
	Total    int
}

func (r *Repository) GetTopPages(ctx context.Context, siteID int64, from, to time.Time, limit int) ([]PageStats, error) {
	var stats []PageStats
	fromUnix := from.Unix()
	toUnix := to.Unix()
	err := r.db.NewSelect().
		Model((*Event)(nil)).
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
	Total    int
}

type BrowserStats struct {
	Browser  ClientBrowser
	Visitors int
}

type DeviceStats struct {
	Device        ClientDevice
	Visitors      int
	Total         int
	TotalVisitors int
}

type OperatingSystemStats struct {
	OS            ClientOS
	Visitors      int
	Total         int
	TotalVisitors int
}

type CountryStats struct {
	CountryCode   string
	Visitors      int
	Total         int
	TotalVisitors int
}

type DailyVisitorStats struct {
	DateBucket int64 // Unix timestamp bucket (day or hour) - integer for performance
	Visitors   int
	PageViews  int
	Sessions   int
}

type ActivePageStats struct {
	Path     string
	Visitors int
}

func (r *Repository) GetActivePages(ctx context.Context, siteID int64, since time.Time, limit, offset int) ([]ActivePageStats, error) {
	var stats []ActivePageStats
	sinceUnix := since.Unix()
	q := r.db.NewSelect().
		Model((*Event)(nil)).
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
