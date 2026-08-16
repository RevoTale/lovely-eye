package analytics

import (
	"context"
	"fmt"
	"time"
)

func buildAnalyticsQuery(siteID int64, from, to time.Time, filter Filter) Query {
	return Query{
		SiteID: siteID,
		From:   from,
		To:     to,
		Filter: filter,
	}
}

func (s *Service) GetDashboardOverview(
	ctx context.Context,
	siteID int64,
	from,
	to time.Time,
) (*Overview, error) {
	return s.GetDashboardOverviewWithFilter(ctx, buildAnalyticsQuery(siteID, from, to, Filter{}))
}

func (s *Service) SyncGeoIPRequirement(ctx context.Context) error {
	if s.geoIPService == nil {
		return nil
	}
	requires, err := s.siteRepo.AnyGeoIPRequirement(ctx)
	if err != nil {
		return s.recordGeoIPFailure(fmt.Errorf("check geoip requirement: %w", err))
	}
	if err := s.geoIPService.SetEnabled(requires); err != nil {
		return s.recordGeoIPFailure(fmt.Errorf("set geoip requirement: %w", err))
	}
	if !requires {
		return nil
	}
	if err := s.geoIPService.EnsureAvailable(ctx); err != nil {
		return s.recordGeoIPFailure(fmt.Errorf("ensure geoip available: %w", err))
	}
	if s.countryService != nil {
		if err := s.countryService.SyncFromGeoIP(ctx); err != nil {
			return s.recordGeoIPFailure(fmt.Errorf("sync persisted countries: %w", err))
		}
	}
	return nil
}

func (s *Service) GeoIPStatus() GeoIPStatus {
	if s.geoIPService == nil {
		return GeoIPStatus{State: geoIPStateDisabled}
	}
	return s.geoIPService.Status()
}

func (s *Service) RefreshGeoIPDatabase(ctx context.Context) (GeoIPStatus, error) {
	if s.geoIPService == nil {
		return GeoIPStatus{State: geoIPStateDisabled}, nil
	}
	if err := s.geoIPService.Refresh(ctx); err != nil {
		err = s.recordGeoIPFailure(fmt.Errorf("refresh geoip database: %w", err))
		return s.geoIPService.Status(), err
	}
	if s.countryService != nil {
		if err := s.countryService.SyncFromGeoIP(ctx); err != nil {
			err = s.recordGeoIPFailure(fmt.Errorf("sync persisted countries: %w", err))
			return s.geoIPService.Status(), err
		}
	}
	return s.geoIPService.Status(), nil
}

func (s *Service) recordGeoIPFailure(err error) error {
	s.geoIPService.RecordFailure(err)
	return err
}

func (s *Service) Close() error {
	if s.geoIPService == nil {
		return nil
	}
	if err := s.geoIPService.Close(); err != nil {
		return fmt.Errorf("close geoip service: %w", err)
	}
	return nil
}

func (s *Service) GetDashboardOverviewWithFilter(
	ctx context.Context,
	query Query,
) (*Overview, error) {
	repositoryQuery := repositoryAnalyticsQuery(query)
	overview, err := s.analyticsRepo.GetOverviewWithFilter(ctx, repositoryQuery)
	if err != nil {
		return nil, fmt.Errorf("get overview with filter: %w", err)
	}
	pageViews, err := s.analyticsRepo.GetPageViewCountWithFilter(ctx, repositoryQuery)
	if err != nil {
		return nil, fmt.Errorf("get page view count with filter: %w", err)
	}

	return &Overview{
		Visitors:    overview.Visitors,
		PageViews:   pageViews,
		Sessions:    overview.Sessions,
		BounceRate:  overview.BounceRate,
		AvgDuration: overview.AvgDuration,
	}, nil
}

func (s *Service) GetTopPagesWithFilterPaged(
	ctx context.Context,
	query Query,
) ([]PageStats, int, error) {
	stats, total, err := s.analyticsRepo.GetTopPagesWithFilterPaged(ctx, repositoryAnalyticsQuery(query))
	if err != nil {
		return nil, 0, fmt.Errorf("get top pages with filter paged: %w", err)
	}
	return pageStats(stats), total, nil
}

func (s *Service) GetTopReferrersWithFilterPaged(
	ctx context.Context,
	query Query,
) ([]ReferrerStats, int, error) {
	stats, total, err := s.analyticsRepo.GetTopReferrersWithFilterPaged(ctx, repositoryAnalyticsQuery(query))
	if err != nil {
		return nil, 0, fmt.Errorf("get top referrers with filter paged: %w", err)
	}
	return referrerStats(stats), total, nil
}

func (s *Service) GetDeviceStatsWithFilterPaged(
	ctx context.Context,
	query Query,
) ([]DeviceStats, int, int, error) {
	stats, total, totalVisitors, err := s.analyticsRepo.GetDeviceStatsWithFilterPaged(ctx, repositoryAnalyticsQuery(query))
	if err != nil {
		return nil, 0, 0, fmt.Errorf("get device stats with filter paged: %w", err)
	}
	return deviceStats(stats), total, totalVisitors, nil
}

func (s *Service) GetOperatingSystemStatsWithFilterPaged(
	ctx context.Context,
	query Query,
) ([]OperatingSystemStats, int, int, error) {
	stats, total, totalVisitors, err := s.analyticsRepo.GetOperatingSystemStatsWithFilterPaged(
		ctx,
		repositoryAnalyticsQuery(query),
	)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("get operating system stats with filter paged: %w", err)
	}
	return operatingSystemStats(stats), total, totalVisitors, nil
}

func (s *Service) GetCountryStatsWithFilterPaged(
	ctx context.Context,
	query Query,
) ([]CountryStats, int, int, error) {
	stats, total, totalVisitors, err := s.analyticsRepo.GetCountryStatsWithFilterPaged(ctx, repositoryAnalyticsQuery(query))
	if err != nil {
		return nil, 0, 0, fmt.Errorf("get country stats with filter paged: %w", err)
	}
	return countryStats(stats), total, totalVisitors, nil
}

func (s *Service) GetBrowserStatsWithFilter(
	ctx context.Context,
	query Query,
) ([]BrowserStats, error) {
	stats, err := s.analyticsRepo.GetBrowserStatsWithFilter(ctx, repositoryAnalyticsQuery(query))
	if err != nil {
		return nil, fmt.Errorf("get browser stats with filter: %w", err)
	}
	return browserStats(stats), nil
}

func (s *Service) GetTimeSeriesStatsWithFilter(
	ctx context.Context,
	query Query,
) ([]TimeSeriesStats, error) {
	stats, err := s.analyticsRepo.GetTimeSeriesStatsWithFilter(ctx, repositoryAnalyticsQuery(query))
	if err != nil {
		return nil, fmt.Errorf("get time series stats with filter: %w", err)
	}
	return timeSeriesStats(stats), nil
}

func (s *Service) GetRealtimeVisitors(ctx context.Context, siteID int64) (int, error) {

	from := time.Now().Add(-5 * time.Minute)
	to := time.Now()
	count, err := s.analyticsRepo.GetVisitorCount(ctx, siteID, from, to)
	if err != nil {
		return 0, fmt.Errorf("get visitor count: %w", err)
	}
	return count, nil
}

func (s *Service) GetActivePages(
	ctx context.Context,
	siteID int64,
	limit,
	offset int,
) ([]ActivePageStats, error) {

	since := time.Now().Add(-5 * time.Minute)
	stats, err := s.analyticsRepo.GetActivePages(ctx, siteID, since, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get active pages: %w", err)
	}
	return activePageStats(stats), nil
}

func (s *Service) GetEvents(
	ctx context.Context,
	siteID int64,
	from,
	to time.Time,
	limit,
	offset int,
) ([]*Event, error) {
	events, err := s.analyticsRepo.GetEvents(ctx, siteID, from, to, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get events: %w", err)
	}
	return analyticsEvents(events), nil
}

func (s *Service) GetEventsWithTotal(
	ctx context.Context,
	siteID int64,
	from,
	to time.Time,
	limit,
	offset int,
) ([]*Event, int, error) {
	events, err := s.analyticsRepo.GetEvents(ctx, siteID, from, to, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("get events: %w", err)
	}

	total, err := s.analyticsRepo.GetEventCount(ctx, siteID, from, to)
	if err != nil {
		return nil, 0, fmt.Errorf("get event count: %w", err)
	}

	return analyticsEvents(events), total, nil
}

func (s *Service) GetEventsWithTotalAndFilter(
	ctx context.Context,
	query Query,
) ([]*Event, int, error) {
	repositoryQuery := repositoryAnalyticsQuery(query)
	events, err := s.analyticsRepo.GetEventsWithFilter(ctx, repositoryQuery)
	if err != nil {
		return nil, 0, fmt.Errorf("get events with filter: %w", err)
	}

	total, err := s.analyticsRepo.GetEventCountWithFilter(ctx, repositoryQuery)
	if err != nil {
		return nil, 0, fmt.Errorf("get event count with filter: %w", err)
	}

	return analyticsEvents(events), total, nil
}

func (s *Service) GetEventCounts(
	ctx context.Context,
	query Query,
) ([]EventCount, int, error) {
	results, total, err := s.analyticsRepo.GetEventCountsGrouped(ctx, repositoryAnalyticsQuery(query))
	if err != nil {
		return nil, 0, fmt.Errorf("get event counts grouped: %w", err)
	}

	if len(results) == 0 {
		return []EventCount{}, total, nil
	}

	eventIDs := make([]int64, len(results))
	for i, result := range results {
		eventIDs[i] = result.EventID
	}

	events, err := s.analyticsRepo.GetEventsByIDs(ctx, eventIDs)
	if err != nil {
		return nil, 0, fmt.Errorf("get events by IDs: %w", err)
	}

	eventMap := make(map[int64]*Event, len(events))
	for _, event := range events {
		mapped := analyticsEvent(event)
		eventMap[mapped.ID] = mapped
	}

	eventCounts := make([]EventCount, 0, len(results))
	for _, result := range results {
		event, ok := eventMap[result.EventID]
		if !ok {
			continue
		}
		eventCounts = append(eventCounts, EventCount{
			Event: event,
			Count: result.Count,
		})
	}

	return eventCounts, total, nil
}
