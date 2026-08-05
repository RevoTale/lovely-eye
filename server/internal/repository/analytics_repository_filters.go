package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/lovely-eye/server/internal/models"
	"github.com/uptrace/bun"
)

func eventTypeFlags(eventTypes []EventType) (bool, bool) {
	hasPageView := false
	hasPredefined := false
	for _, eventType := range eventTypes {
		switch eventType {
		case EventTypePageView:
			hasPageView = true
		case EventTypePredefined:
			hasPredefined = true
		}
	}
	return hasPageView, hasPredefined
}

func uint8FilterValues[T ~uint8](values []T) []uint8 {
	if len(values) == 0 {
		return nil
	}

	translated := make([]uint8, 0, len(values))
	for _, value := range values {
		translated = append(translated, uint8(value))
	}
	return translated
}

func applyEnumFilter[T ~uint8](q *bun.SelectQuery, rawValues []string, parse func([]string) []T, clause string) *bun.SelectQuery {
	if len(rawValues) == 0 {
		return q
	}

	translated := parse(rawValues)
	if len(translated) == 0 {
		return q.Where("1 = 0")
	}
	return q.Where(clause, bun.List(uint8FilterValues(translated)))
}

func applySessionFilters(q *bun.SelectQuery, filter AnalyticsFilter) *bun.SelectQuery {
	if len(filter.Referrer) > 0 {

		q = q.Where("s.referrer IN (?)", bun.List(filter.Referrer))
	}
	q = applyEnumFilter(q, filter.Browser, models.ParseClientBrowserFilters, "s.client_id IN (SELECT id FROM clients WHERE browser IN (?))")
	q = applyEnumFilter(q, filter.Device, models.ParseClientDeviceFilters, "s.client_id IN (SELECT id FROM clients WHERE device IN (?))")
	q = applyEnumFilter(q, filter.OS, models.ParseClientOSFilters, "s.client_id IN (SELECT id FROM clients WHERE os IN (?))")
	if len(filter.Page) > 0 {
		q = q.Where("s.id IN (SELECT DISTINCT session_id FROM events WHERE definition_id IS NULL AND path IN (?))", bun.List(filter.Page))
	}
	if len(filter.Country) > 0 {

		q = q.Where("s.client_id IN (SELECT id FROM clients WHERE COALESCE(NULLIF(country, ''), '-') IN (?))", bun.List(normalizeCountryCodes(filter.Country)))
	}
	if len(filter.EventTypes) > 0 {
		hasPageView, hasPredefined := eventTypeFlags(filter.EventTypes)
		if hasPageView && !hasPredefined {
			q = q.Where("s.id IN (SELECT DISTINCT session_id FROM events WHERE definition_id IS NULL)")
		} else if hasPredefined && !hasPageView {
			q = q.Where("s.id IN (SELECT DISTINCT session_id FROM events WHERE definition_id IS NOT NULL)")
		}
	}
	if len(filter.EventName) > 0 {
		q = q.Where("s.id IN (SELECT DISTINCT e.session_id FROM events e INNER JOIN event_definitions ed ON e.definition_id = ed.id WHERE ed.name IN (?))", bun.List(filter.EventName))
	}
	if len(filter.EventPath) > 0 {
		q = q.Where("s.id IN (SELECT DISTINCT session_id FROM events WHERE path IN (?))", bun.List(filter.EventPath))
	}
	if len(filter.EventDefinitionIDs) > 0 {
		q = q.Where("s.id IN (SELECT DISTINCT session_id FROM events WHERE definition_id IN (?))", bun.List(filter.EventDefinitionIDs))
	}
	return q
}

func applyEventFilters(q *bun.SelectQuery, filter AnalyticsFilter) *bun.SelectQuery {
	if len(filter.EventTypes) > 0 {
		hasPageView, hasPredefined := eventTypeFlags(filter.EventTypes)
		if hasPageView && !hasPredefined {
			q = q.Where("e.definition_id IS NULL")
		} else if hasPredefined && !hasPageView {
			q = q.Where("e.definition_id IS NOT NULL")
		}
	}
	if len(filter.Page) > 0 {
		q = q.Where("e.path IN (?)", bun.List(filter.Page))
	}
	if len(filter.Referrer) > 0 || len(filter.Browser) > 0 || len(filter.Device) > 0 || len(filter.OS) > 0 || len(filter.Country) > 0 || len(filter.EventTypes) > 0 || len(filter.EventName) > 0 || len(filter.EventPath) > 0 || len(filter.EventDefinitionIDs) > 0 {

		if len(filter.Referrer) > 0 {
			q = q.Where("e.session_id IN (SELECT id FROM sessions WHERE referrer IN (?))", bun.List(filter.Referrer))
		}
		q = applyEnumFilter(q, filter.Browser, models.ParseClientBrowserFilters, "e.session_id IN (SELECT s.id FROM sessions s INNER JOIN clients c ON s.client_id = c.id WHERE c.browser IN (?))")
		q = applyEnumFilter(q, filter.Device, models.ParseClientDeviceFilters, "e.session_id IN (SELECT s.id FROM sessions s INNER JOIN clients c ON s.client_id = c.id WHERE c.device IN (?))")
		q = applyEnumFilter(q, filter.OS, models.ParseClientOSFilters, "e.session_id IN (SELECT s.id FROM sessions s INNER JOIN clients c ON s.client_id = c.id WHERE c.os IN (?))")
		if len(filter.Country) > 0 {
			q = q.Where("e.session_id IN (SELECT s.id FROM sessions s INNER JOIN clients c ON s.client_id = c.id WHERE COALESCE(NULLIF(c.country, ''), '-') IN (?))", bun.List(normalizeCountryCodes(filter.Country)))
		}
		if len(filter.EventName) > 0 {
			q = q.Where("e.session_id IN (SELECT DISTINCT e.session_id FROM events e INNER JOIN event_definitions ed ON e.definition_id = ed.id WHERE ed.name IN (?))", bun.List(filter.EventName))
		}
		if len(filter.EventPath) > 0 {
			q = q.Where("e.session_id IN (SELECT DISTINCT session_id FROM events WHERE path IN (?))", bun.List(filter.EventPath))
		}
		if len(filter.EventDefinitionIDs) > 0 {
			q = q.Where("e.session_id IN (SELECT DISTINCT session_id FROM events WHERE definition_id IN (?))", bun.List(filter.EventDefinitionIDs))
		}
	}
	return q
}

func applyEventNamePathFilters(q *bun.SelectQuery, filter AnalyticsFilter) *bun.SelectQuery {
	if len(filter.EventName) > 0 {
		q = q.Join("INNER JOIN event_definitions ed ON e.definition_id = ed.id")
		q = q.Where("ed.name IN (?)", bun.List(filter.EventName))
	}
	if len(filter.EventPath) > 0 {
		q = q.Where("e.path IN (?)", bun.List(filter.EventPath))
	}
	if len(filter.EventDefinitionIDs) > 0 {
		q = q.Where("e.definition_id IN (?)", bun.List(filter.EventDefinitionIDs))
	}
	return q
}

func normalizeCountryCodes(values []string) []string {
	normalized := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))

	for _, value := range values {
		code := strings.ToUpper(strings.TrimSpace(value))
		switch code {
		case "", "UNKNOWN":
			code = "-"
		}

		if _, ok := seen[code]; ok {
			continue
		}
		seen[code] = struct{}{}
		normalized = append(normalized, code)
	}

	return normalized
}

func (r *AnalyticsRepository) GetVisitorCountWithFilter(ctx context.Context, query AnalyticsQuery) (int, error) {
	var count int
	fromUnix := query.From.Unix()
	toUnix := query.To.Unix()
	q := r.db.NewSelect().
		TableExpr("sessions s").
		ColumnExpr("COUNT(DISTINCT s.client_id)").
		Where("s.site_id = ?", query.SiteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix)
	q = applySessionFilters(q, query.Filter)
	err := q.Scan(ctx, &count)
	if err != nil {
		return 0, fmt.Errorf("failed to get visitor count with filter: %w", err)
	}
	return count, nil
}

func (r *AnalyticsRepository) GetPageViewCountWithFilter(ctx context.Context, query AnalyticsQuery) (int, error) {
	fromUnix := query.From.Unix()
	toUnix := query.To.Unix()
	q := r.db.NewSelect().
		TableExpr("events e").
		Join("INNER JOIN sessions s ON e.session_id = s.id").
		Where("s.site_id = ?", query.SiteID).
		Where("e.definition_id IS NULL").
		Where("e.time >= ?", fromUnix).
		Where("e.time <= ?", toUnix)
	q = applyEventFilters(q, query.Filter)
	count, err := q.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get page view count with filter: %w", err)
	}
	return count, nil
}

func (r *AnalyticsRepository) GetSessionCountWithFilter(ctx context.Context, query AnalyticsQuery) (int, error) {
	fromUnix := query.From.Unix()
	toUnix := query.To.Unix()
	q := r.db.NewSelect().
		TableExpr("sessions s").
		Where("s.site_id = ?", query.SiteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix)
	q = applySessionFilters(q, query.Filter)
	count, err := q.Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get session count with filter: %w", err)
	}
	return count, nil
}

func (r *AnalyticsRepository) GetBounceRateWithFilter(ctx context.Context, query AnalyticsQuery) (float64, error) {
	var result struct {
		Total   int
		Bounced int
	}
	fromUnix := query.From.Unix()
	toUnix := query.To.Unix()

	dialect := fmt.Sprint(r.db.Dialect().Name())
	var bouncedExpr string
	if dialect == "pg" || dialect == "postgres" || dialect == "postgresql" {

		bouncedExpr = "COUNT(*) FILTER (WHERE s.page_view_count = 1)"
	} else {

		bouncedExpr = "SUM(CASE WHEN s.page_view_count = 1 THEN 1 ELSE 0 END)"
	}

	q := r.db.NewSelect().
		TableExpr("sessions s").
		ColumnExpr("COUNT(*) as total").
		ColumnExpr(bouncedExpr+" as bounced").
		Where("s.site_id = ?", query.SiteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix).
		Where("s.page_view_count > 0")
	q = applySessionFilters(q, query.Filter)
	err := q.Scan(ctx, &result)
	if err != nil {
		return 0, fmt.Errorf("failed to get bounce rate with filter: %w", err)
	}
	if result.Total == 0 {
		return 0, nil
	}
	return float64(result.Bounced) / float64(result.Total) * 100, nil
}

func (r *AnalyticsRepository) GetAvgSessionDurationWithFilter(ctx context.Context, query AnalyticsQuery) (float64, error) {
	var avg float64
	fromUnix := query.From.Unix()
	toUnix := query.To.Unix()
	q := r.db.NewSelect().
		TableExpr("sessions s").
		ColumnExpr("COALESCE(AVG((s.exit_time - s.enter_time) * 1.0), 0.0)").
		Where("s.site_id = ?", query.SiteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix).
		Where("s.page_view_count > 0").
		Where("s.exit_time > s.enter_time")
	q = applySessionFilters(q, query.Filter)
	err := q.Scan(ctx, &avg)
	if err != nil {
		return 0, fmt.Errorf("failed to get average session duration with filter: %w", err)
	}
	return avg, nil
}

func (r *AnalyticsRepository) GetTopPagesWithFilter(ctx context.Context, query AnalyticsQuery) ([]PageStats, error) {
	var stats []PageStats
	fromUnix := query.From.Unix()
	toUnix := query.To.Unix()
	q := r.db.NewSelect().
		TableExpr("events e").
		Join("INNER JOIN sessions s ON e.session_id = s.id").
		ColumnExpr("e.path").
		ColumnExpr("COUNT(*) as views").
		ColumnExpr("COUNT(DISTINCT s.client_id) as visitors").
		Where("s.site_id = ?", query.SiteID).
		Where("e.definition_id IS NULL").
		Where("e.time >= ?", fromUnix).
		Where("e.time <= ?", toUnix)
	q = applyEventFilters(q, query.Filter)
	err := q.Group("e.path").
		Order("views DESC", "e.path ASC").
		Limit(query.Limit).
		Scan(ctx, &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to get top pages with filter: %w", err)
	}
	return stats, nil
}

func (r *AnalyticsRepository) GetTopReferrersWithFilter(ctx context.Context, query AnalyticsQuery) ([]ReferrerStats, error) {
	var stats []ReferrerStats
	fromUnix := query.From.Unix()
	toUnix := query.To.Unix()
	q := r.db.NewSelect().
		TableExpr("sessions s").
		ColumnExpr("COALESCE(NULLIF(s.referrer, ''), '(direct)') as referrer").
		ColumnExpr("COUNT(DISTINCT s.client_id) as visitors").
		Where("s.site_id = ?", query.SiteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix)
	q = applySessionFilters(q, query.Filter)
	err := q.Group("s.referrer").
		Order("visitors DESC", "referrer ASC").
		Limit(query.Limit).
		Scan(ctx, &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to get top referrers with filter: %w", err)
	}
	return stats, nil
}

func (r *AnalyticsRepository) GetBrowserStatsWithFilter(ctx context.Context, query AnalyticsQuery) ([]BrowserStats, error) {
	var stats []BrowserStats
	fromUnix := query.From.Unix()
	toUnix := query.To.Unix()
	q := r.db.NewSelect().
		TableExpr("sessions s").
		Join("INNER JOIN clients c ON s.client_id = c.id").
		ColumnExpr("c.browser").
		ColumnExpr("COUNT(DISTINCT s.client_id) as visitors").
		Where("s.site_id = ?", query.SiteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix).
		Where("c.browser != ?", models.ClientBrowserUnknown)
	q = applySessionFilters(q, query.Filter)
	q = q.Group("c.browser").
		Order("visitors DESC", "c.browser ASC")
	if query.Limit > 0 {
		q = q.Limit(query.Limit)
	}
	if query.Offset > 0 {
		q = q.Offset(query.Offset)
	}
	err := q.Scan(ctx, &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to get browser stats with filter: %w", err)
	}
	return stats, nil
}

func (r *AnalyticsRepository) GetDeviceStatsWithFilter(ctx context.Context, query AnalyticsQuery) ([]DeviceStats, error) {
	var stats []DeviceStats
	fromUnix := query.From.Unix()
	toUnix := query.To.Unix()
	q := r.db.NewSelect().
		TableExpr("sessions s").
		Join("INNER JOIN clients c ON s.client_id = c.id").
		ColumnExpr("c.device").
		ColumnExpr("COUNT(DISTINCT s.client_id) as visitors").
		Where("s.site_id = ?", query.SiteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix).
		Where("c.device != ?", models.ClientDeviceUnknown)
	q = applySessionFilters(q, query.Filter)
	err := q.Group("c.device").
		Order("visitors DESC", "c.device ASC").
		Limit(query.Limit).
		Scan(ctx, &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to get device stats with filter: %w", err)
	}
	return stats, nil
}

func (r *AnalyticsRepository) GetCountryStatsWithFilter(ctx context.Context, query AnalyticsQuery) ([]CountryStats, error) {
	var stats []CountryStats
	fromUnix := query.From.Unix()
	toUnix := query.To.Unix()
	q := r.db.NewSelect().
		TableExpr("sessions s").
		Join("INNER JOIN clients c ON s.client_id = c.id").
		ColumnExpr("COALESCE(NULLIF(c.country, ''), '-') as country_code").
		ColumnExpr("COUNT(DISTINCT s.client_id) as visitors").
		Where("s.site_id = ?", query.SiteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix)
	q = applySessionFilters(q, query.Filter)
	err := q.GroupExpr("COALESCE(NULLIF(c.country, ''), '-')").
		Order("visitors DESC", "country_code ASC").
		Limit(query.Limit).
		Scan(ctx, &stats)
	if err != nil {
		return nil, fmt.Errorf("failed to get country stats with filter: %w", err)
	}
	return stats, nil
}

func (r *AnalyticsRepository) GetDailyStatsWithFilter(ctx context.Context, query AnalyticsQuery) ([]DailyVisitorStats, error) {
	query.Bucket = TimeBucketDaily
	return r.GetTimeSeriesStatsWithFilter(ctx, query)
}
