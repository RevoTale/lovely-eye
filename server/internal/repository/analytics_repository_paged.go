package repository

import (
	"context"
	"fmt"
	"sort"

	"github.com/lovely-eye/server/internal/models"
)

func (r *AnalyticsRepository) GetTopPagesWithFilterPaged(ctx context.Context, query AnalyticsQuery) ([]PageStats, int, error) {
	var stats []PageStats
	var total int
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
		Offset(query.Offset).
		Scan(ctx, &stats)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get top pages with filter paged: %w", err)
	}

	countQuery := r.db.NewSelect().
		TableExpr("events e").
		Join("INNER JOIN sessions s ON e.session_id = s.id").
		ColumnExpr("COUNT(DISTINCT e.path)").
		Where("s.site_id = ?", query.SiteID).
		Where("e.definition_id IS NULL").
		Where("e.time >= ?", fromUnix).
		Where("e.time <= ?", toUnix)
	countQuery = applyEventFilters(countQuery, query.Filter)
	err = countQuery.Scan(ctx, &total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count top pages with filter: %w", err)
	}
	return stats, total, nil
}

func (r *AnalyticsRepository) GetTopReferrersWithFilterPaged(ctx context.Context, query AnalyticsQuery) ([]ReferrerStats, int, error) {
	var stats []ReferrerStats
	var total int
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
		Offset(query.Offset).
		Scan(ctx, &stats)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get top referrers with filter paged: %w", err)
	}

	countQuery := r.db.NewSelect().
		TableExpr("sessions s").
		ColumnExpr("COUNT(DISTINCT COALESCE(NULLIF(s.referrer, ''), '(direct)'))").
		Where("s.site_id = ?", query.SiteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix)
	countQuery = applySessionFilters(countQuery, query.Filter)
	err = countQuery.Scan(ctx, &total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count top referrers with filter: %w", err)
	}
	return stats, total, nil
}

func (r *AnalyticsRepository) GetDeviceStatsWithFilterPaged(ctx context.Context, query AnalyticsQuery) ([]DeviceStats, int, int, error) {
	var stats []DeviceStats
	var total int
	var totalVisitors int
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
		Offset(query.Offset).
		Scan(ctx, &stats)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get device stats with filter paged: %w", err)
	}

	countQuery := r.db.NewSelect().
		TableExpr("sessions s").
		Join("INNER JOIN clients c ON s.client_id = c.id").
		ColumnExpr("COUNT(DISTINCT c.device)").
		Where("s.site_id = ?", query.SiteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix).
		Where("c.device != ?", models.ClientDeviceUnknown)
	countQuery = applySessionFilters(countQuery, query.Filter)
	err = countQuery.Scan(ctx, &total)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to count devices with filter: %w", err)
	}

	deviceCounts := r.db.NewSelect().
		TableExpr("sessions s").
		Join("INNER JOIN clients c ON s.client_id = c.id").
		ColumnExpr("COUNT(DISTINCT s.client_id) as visitors").
		Where("s.site_id = ?", query.SiteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix).
		Where("c.device != ?", models.ClientDeviceUnknown)
	deviceCounts = applySessionFilters(deviceCounts, query.Filter)
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

func (r *AnalyticsRepository) GetOperatingSystemStatsWithFilterPaged(ctx context.Context, query AnalyticsQuery) ([]OperatingSystemStats, int, int, error) {
	var stats []OperatingSystemStats
	var total int
	var totalVisitors int
	fromUnix := query.From.Unix()
	toUnix := query.To.Unix()
	q := r.db.NewSelect().
		TableExpr("sessions s").
		Join("INNER JOIN clients c ON s.client_id = c.id").
		ColumnExpr("c.os").
		ColumnExpr("COUNT(DISTINCT s.client_id) as visitors").
		Where("s.site_id = ?", query.SiteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix).
		Where("c.os != ?", models.ClientOSUnknown)
	q = applySessionFilters(q, query.Filter)
	err := q.Group("c.os").
		Order("visitors DESC", "c.os ASC").
		Limit(query.Limit).
		Offset(query.Offset).
		Scan(ctx, &stats)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get operating system stats with filter paged: %w", err)
	}

	countQuery := r.db.NewSelect().
		TableExpr("sessions s").
		Join("INNER JOIN clients c ON s.client_id = c.id").
		ColumnExpr("COUNT(DISTINCT c.os)").
		Where("s.site_id = ?", query.SiteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix).
		Where("c.os != ?", models.ClientOSUnknown)
	countQuery = applySessionFilters(countQuery, query.Filter)
	err = countQuery.Scan(ctx, &total)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to count operating systems with filter: %w", err)
	}

	osCounts := r.db.NewSelect().
		TableExpr("sessions s").
		Join("INNER JOIN clients c ON s.client_id = c.id").
		ColumnExpr("COUNT(DISTINCT s.client_id) as visitors").
		Where("s.site_id = ?", query.SiteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix).
		Where("c.os != ?", models.ClientOSUnknown)
	osCounts = applySessionFilters(osCounts, query.Filter)
	osCounts = osCounts.Group("c.os")

	err = r.db.NewSelect().
		TableExpr("(?) as os_counts", osCounts).
		ColumnExpr("COALESCE(SUM(visitors), 0)").
		Scan(ctx, &totalVisitors)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to sum operating system visitors with filter: %w", err)
	}
	return stats, total, totalVisitors, nil
}

func (r *AnalyticsRepository) GetCountryStatsWithFilterPaged(ctx context.Context, query AnalyticsQuery) ([]CountryStats, int, int, error) {
	var stats []CountryStats
	var total int
	var totalVisitors int
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
		Offset(query.Offset).
		Scan(ctx, &stats)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get country stats with filter paged: %w", err)
	}

	countQuery := r.db.NewSelect().
		TableExpr("sessions s").
		Join("INNER JOIN clients c ON s.client_id = c.id").
		ColumnExpr("COUNT(DISTINCT COALESCE(NULLIF(c.country, ''), '-'))").
		Where("s.site_id = ?", query.SiteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix)
	countQuery = applySessionFilters(countQuery, query.Filter)
	err = countQuery.Scan(ctx, &total)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to count countries with filter: %w", err)
	}

	countryCounts := r.db.NewSelect().
		TableExpr("sessions s").
		Join("INNER JOIN clients c ON s.client_id = c.id").
		ColumnExpr("COUNT(DISTINCT s.client_id) as visitors").
		Where("s.site_id = ?", query.SiteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix)
	countryCounts = applySessionFilters(countryCounts, query.Filter)
	countryCounts = countryCounts.GroupExpr("COALESCE(NULLIF(c.country, ''), '-')")

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

func (r *AnalyticsRepository) GetTimeSeriesStatsWithFilter(ctx context.Context, query AnalyticsQuery) ([]DailyVisitorStats, error) {
	fromUnix := query.From.Unix()
	toUnix := query.To.Unix()
	sessionBucketExpr := r.sessionTimeBucketExpression(query.Bucket)
	eventBucketExpr := r.eventTimeBucketExpression(query.Bucket)

	var sessionStats []struct {
		DateBucket int64
		Visitors   int
		Sessions   int
	}
	sessionQuery := r.db.NewSelect().
		TableExpr("sessions s").
		ColumnExpr(sessionBucketExpr+" as date_bucket").
		ColumnExpr("COUNT(DISTINCT s.client_id) as visitors").
		ColumnExpr("COUNT(*) as sessions").
		Where("s.site_id = ?", query.SiteID).
		Where("s.enter_time >= ?", fromUnix).
		Where("s.enter_time <= ?", toUnix)
	sessionQuery = applySessionFilters(sessionQuery, query.Filter)
	if err := sessionQuery.GroupExpr(sessionBucketExpr).Scan(ctx, &sessionStats); err != nil {
		return nil, fmt.Errorf("failed to get time series session stats with filter: %w", err)
	}

	var pageViewStats []struct {
		DateBucket int64
		PageViews  int
	}
	pageViewQuery := r.db.NewSelect().
		TableExpr("events e").
		Join("INNER JOIN sessions s ON e.session_id = s.id").
		ColumnExpr(eventBucketExpr+" as date_bucket").
		ColumnExpr("COUNT(*) as page_views").
		Where("s.site_id = ?", query.SiteID).
		Where("e.definition_id IS NULL").
		Where("e.time >= ?", fromUnix).
		Where("e.time <= ?", toUnix)
	pageViewQuery = applyEventFilters(pageViewQuery, query.Filter)
	if err := pageViewQuery.GroupExpr(eventBucketExpr).Scan(ctx, &pageViewStats); err != nil {
		return nil, fmt.Errorf("failed to get time series pageview stats with filter: %w", err)
	}

	statsByBucket := make(map[int64]*DailyVisitorStats, len(sessionStats)+len(pageViewStats))
	for _, stat := range sessionStats {
		statsByBucket[stat.DateBucket] = &DailyVisitorStats{
			DateBucket: stat.DateBucket,
			Visitors:   stat.Visitors,
			Sessions:   stat.Sessions,
		}
	}
	for _, stat := range pageViewStats {
		merged := statsByBucket[stat.DateBucket]
		if merged == nil {
			merged = &DailyVisitorStats{DateBucket: stat.DateBucket}
			statsByBucket[stat.DateBucket] = merged
		}
		merged.PageViews = stat.PageViews
	}

	buckets := make([]int64, 0, len(statsByBucket))
	for bucket := range statsByBucket {
		buckets = append(buckets, bucket)
	}
	sort.Slice(buckets, func(i, j int) bool { return buckets[i] < buckets[j] })

	if query.Limit > 0 {
		descBuckets := append([]int64(nil), buckets...)
		sort.Slice(descBuckets, func(i, j int) bool { return descBuckets[i] > descBuckets[j] })
		offset := max(query.Offset, 0)
		if offset >= len(descBuckets) {
			buckets = nil
		} else {
			end := min(offset+query.Limit, len(descBuckets))
			buckets = descBuckets[offset:end]
			sort.Slice(buckets, func(i, j int) bool { return buckets[i] < buckets[j] })
		}
	} else if query.Offset > 0 {
		if query.Offset >= len(buckets) {
			buckets = nil
		} else {
			buckets = buckets[query.Offset:]
		}
	}

	stats := make([]DailyVisitorStats, 0, len(buckets))
	for _, bucket := range buckets {
		stats = append(stats, *statsByBucket[bucket])
	}
	return stats, nil
}

func (r *AnalyticsRepository) sessionTimeBucketExpression(bucket TimeBucket) string {
	if bucket == TimeBucketHourly {
		return "s.enter_hour"
	}

	return "s.enter_day"
}

func (r *AnalyticsRepository) eventTimeBucketExpression(bucket TimeBucket) string {
	if bucket == TimeBucketHourly {
		return "e.hour"
	}

	return "e.day"
}
