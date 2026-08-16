package graph

import (
	"context"
	"fmt"
	"strconv"
	"time"

	analyticfeature "github.com/lovely-eye/server/internal/analytics"
	"github.com/lovely-eye/server/internal/auth"
	"github.com/lovely-eye/server/internal/graph/model"
)

// TopPages is the resolver for the topPages field.
func (r *dashboardStatsResolver) TopPages(ctx context.Context, obj *model.DashboardStats, paging model.PagingInput) (*model.PagedPageStats, error) {
	limit, offset := normalizePaging(paging)
	query := analyticfeature.Query{
		SiteID: obj.SiteID,
		From:   obj.From,
		To:     obj.To,
		Limit:  limit,
		Offset: offset,
		Filter: obj.Filter,
	}
	stats, total, err := r.AnalyticsService.GetTopPagesWithFilterPaged(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get top pages: %w", err)
	}

	items := make([]*model.PageStats, 0, len(stats))
	for _, stat := range stats {
		items = append(items, &model.PageStats{
			Path:     stat.Path,
			Views:    stat.Views,
			Visitors: stat.Visitors,
		})
	}

	return &model.PagedPageStats{
		Items: items,
		Total: total,
	}, nil
}

// TopReferrers is the resolver for the topReferrers field.
func (r *dashboardStatsResolver) TopReferrers(ctx context.Context, obj *model.DashboardStats, paging model.PagingInput) (*model.PagedReferrerStats, error) {
	limit, offset := normalizePaging(paging)
	query := analyticfeature.Query{
		SiteID: obj.SiteID,
		From:   obj.From,
		To:     obj.To,
		Limit:  limit,
		Offset: offset,
		Filter: obj.Filter,
	}
	stats, total, err := r.AnalyticsService.GetTopReferrersWithFilterPaged(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get top referrers: %w", err)
	}

	items := make([]*model.ReferrerStats, 0, len(stats))
	for _, stat := range stats {
		items = append(items, &model.ReferrerStats{
			Referrer: stat.Referrer,
			Visitors: stat.Visitors,
		})
	}

	return &model.PagedReferrerStats{
		Items: items,
		Total: total,
	}, nil
}

// Browsers is the resolver for the browsers field.
func (r *dashboardStatsResolver) Browsers(ctx context.Context, obj *model.DashboardStats, paging model.PagingInput) ([]*model.BrowserStats, error) {
	limit, offset := normalizePaging(paging)
	query := analyticfeature.Query{
		SiteID: obj.SiteID,
		From:   obj.From,
		To:     obj.To,
		Limit:  limit,
		Offset: offset,
		Filter: obj.Filter,
	}
	stats, err := r.AnalyticsService.GetBrowserStatsWithFilter(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get browser stats: %w", err)
	}

	items := make([]*model.BrowserStats, 0, len(stats))
	for _, stat := range stats {
		items = append(items, &model.BrowserStats{
			Browser:  stat.Browser,
			Visitors: stat.Visitors,
		})
	}
	return items, nil
}

// Devices is the resolver for the devices field.
func (r *dashboardStatsResolver) Devices(ctx context.Context, obj *model.DashboardStats, paging model.PagingInput) (*model.PagedDeviceStats, error) {
	limit, offset := normalizePaging(paging)
	query := analyticfeature.Query{
		SiteID: obj.SiteID,
		From:   obj.From,
		To:     obj.To,
		Limit:  limit,
		Offset: offset,
		Filter: obj.Filter,
	}
	stats, total, totalVisitors, err := r.AnalyticsService.GetDeviceStatsWithFilterPaged(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get device stats: %w", err)
	}

	items := make([]*model.DeviceStats, 0, len(stats))
	for _, stat := range stats {
		items = append(items, &model.DeviceStats{
			Device:   stat.Device,
			Visitors: stat.Visitors,
		})
	}

	return &model.PagedDeviceStats{
		Items:         items,
		Total:         total,
		TotalVisitors: totalVisitors,
	}, nil
}

// OperatingSystems is the resolver for the operatingSystems field.
func (r *dashboardStatsResolver) OperatingSystems(ctx context.Context, obj *model.DashboardStats, paging model.PagingInput) (*model.PagedOperatingSystemStats, error) {
	limit, offset := normalizePaging(paging)
	query := analyticfeature.Query{
		SiteID: obj.SiteID,
		From:   obj.From,
		To:     obj.To,
		Limit:  limit,
		Offset: offset,
		Filter: obj.Filter,
	}
	stats, total, totalVisitors, err := r.AnalyticsService.GetOperatingSystemStatsWithFilterPaged(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get operating system stats: %w", err)
	}

	items := make([]*model.OperatingSystemStats, 0, len(stats))
	for _, stat := range stats {
		items = append(items, &model.OperatingSystemStats{
			OS:       stat.OS,
			Visitors: stat.Visitors,
		})
	}

	return &model.PagedOperatingSystemStats{
		Items:         items,
		Total:         total,
		TotalVisitors: totalVisitors,
	}, nil
}

// Countries is the resolver for the countries field.
func (r *dashboardStatsResolver) Countries(ctx context.Context, obj *model.DashboardStats, paging model.PagingInput) (*model.PagedCountryStats, error) {
	limit, offset := normalizePaging(paging)
	query := analyticfeature.Query{
		SiteID: obj.SiteID,
		From:   obj.From,
		To:     obj.To,
		Limit:  limit,
		Offset: offset,
		Filter: obj.Filter,
	}
	stats, total, totalVisitors, err := r.AnalyticsService.GetCountryStatsWithFilterPaged(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get country stats: %w", err)
	}

	items := make([]*model.CountryStats, 0, len(stats))
	for _, stat := range stats {
		items = append(items, &model.CountryStats{
			Country:  newGraphQLCountry(stat.CountryCode, ""),
			Visitors: stat.Visitors,
		})
	}

	return &model.PagedCountryStats{
		Items:         items,
		Total:         total,
		TotalVisitors: totalVisitors,
	}, nil
}

// DailyStats is the resolver for the dailyStats field.
func (r *dashboardStatsResolver) DailyStats(ctx context.Context, obj *model.DashboardStats, bucket *model.TimeBucket, paging model.PagingInput) ([]*model.DailyStats, error) {
	var selectedBucket analyticfeature.TimeBucket
	switch bucketValue := bucketValueOrDefault(bucket); bucketValue {
	case model.TimeBucketHourly:
		selectedBucket = analyticfeature.TimeBucketHourly
	case model.TimeBucketDaily:
		selectedBucket = analyticfeature.TimeBucketDaily
	default:
		return nil, badUserInput("invalid time bucket")
	}

	pointLimit := clampLimit(paging.Limit, maxTimeSeriesPoints)
	maxRangeDays := r.DashboardLimits.MaxDailyRangeDays
	if selectedBucket == analyticfeature.TimeBucketHourly {
		maxRangeDays = r.DashboardLimits.MaxHourlyRangeDays
	}
	if err := validateDateRange(obj.From, obj.To, maxRangeDays); err != nil {
		return nil, err
	}

	offsetValue := paging.Offset
	if offsetValue < 0 {
		offsetValue = 0
	}

	query := analyticfeature.Query{
		SiteID: obj.SiteID,
		From:   obj.From,
		To:     obj.To,
		Limit:  pointLimit,
		Offset: offsetValue,
		Bucket: selectedBucket,
		Filter: obj.Filter,
	}
	stats, err := r.AnalyticsService.GetTimeSeriesStatsWithFilter(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get time series stats: %w", err)
	}

	items := make([]*model.DailyStats, 0, len(stats))
	for _, stat := range stats {
		bucketSeconds := stat.DateBucket
		switch selectedBucket {
		case analyticfeature.TimeBucketDaily:
			bucketSeconds = stat.DateBucket * 86400
		case analyticfeature.TimeBucketHourly:
			bucketSeconds = stat.DateBucket * 3600
		}
		items = append(items, &model.DailyStats{
			Date:      time.Unix(bucketSeconds, 0),
			Visitors:  stat.Visitors,
			PageViews: stat.PageViews,
			Sessions:  stat.Sessions,
		})
	}

	return items, nil
}

// Dashboard is the resolver for the dashboard field.
func (r *queryResolver) Dashboard(ctx context.Context, siteID string, dateRange *model.DateRangeInput, filter *model.FilterInput) (*model.DashboardStats, error) {
	claims := auth.GetUserFromContext(ctx)
	if claims == nil {
		return nil, unauthenticated()
	}

	id, err := strconv.ParseInt(siteID, 10, 64)
	if err != nil {
		return nil, badUserInput("invalid site ID")
	}

	if err := r.SiteService.RequireOwnership(ctx, id, claims.UserID); err != nil {
		return nil, fmt.Errorf("failed to get site: %w", err)
	}

	from, to, err := parseDateRangeInput(dateRange, r.DashboardLimits.MaxDailyRangeDays)
	if err != nil {
		return nil, err
	}
	filterOpts, err := parseFilterInput(filter, r.DashboardLimits)
	if err != nil {
		return nil, err
	}

	stats, err := r.AnalyticsService.GetDashboardOverviewWithFilter(ctx, analyticfeature.Query{
		SiteID: id,
		From:   from,
		To:     to,
		Filter: filterOpts,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard overview: %w", err)
	}

	return &model.DashboardStats{
		Visitors:    stats.Visitors,
		PageViews:   stats.PageViews,
		Sessions:    stats.Sessions,
		BounceRate:  stats.BounceRate,
		AvgDuration: stats.AvgDuration,
		SiteID:      id,
		From:        from,
		To:          to,
		Filter:      filterOpts,
	}, nil
}

// Realtime is the resolver for the realtime field.
func (r *queryResolver) Realtime(ctx context.Context, siteID string) (*model.RealtimeStats, error) {
	claims := auth.GetUserFromContext(ctx)
	if claims == nil {
		return nil, unauthenticated()
	}

	id, err := strconv.ParseInt(siteID, 10, 64)
	if err != nil {
		return nil, badUserInput("invalid site ID")
	}

	if err := r.SiteService.RequireOwnership(ctx, id, claims.UserID); err != nil {
		return nil, fmt.Errorf("failed to get site: %w", err)
	}

	visitors, err := r.AnalyticsService.GetRealtimeVisitors(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get realtime visitors: %w", err)
	}

	return &model.RealtimeStats{
		Visitors: visitors,
		SiteID:   id,
	}, nil
}

// ActivePages is the resolver for the activePages field.
func (r *realtimeStatsResolver) ActivePages(ctx context.Context, obj *model.RealtimeStats, paging model.PagingInput) ([]*model.ActivePageStats, error) {
	limit, offset := normalizePaging(paging)
	activePages, err := r.AnalyticsService.GetActivePages(ctx, obj.SiteID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get active pages: %w", err)
	}

	pages := make([]*model.ActivePageStats, len(activePages))
	for i, page := range activePages {
		pages[i] = &model.ActivePageStats{
			Path:     page.Path,
			Visitors: page.Visitors,
		}
	}

	return pages, nil
}

// DashboardStats returns DashboardStatsResolver implementation.
func (r *Resolver) DashboardStats() DashboardStatsResolver { return &dashboardStatsResolver{r} }

// RealtimeStats returns RealtimeStatsResolver implementation.
func (r *Resolver) RealtimeStats() RealtimeStatsResolver { return &realtimeStatsResolver{r} }

type (
	dashboardStatsResolver struct{ *Resolver }
	realtimeStatsResolver  struct{ *Resolver }
)
