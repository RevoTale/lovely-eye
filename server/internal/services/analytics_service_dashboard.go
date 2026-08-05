package services

import (
	"context"
	"fmt"
	"time"

	"github.com/lovely-eye/server/internal/models"
	"github.com/lovely-eye/server/internal/repository"
)

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
	visitors, err := s.analyticsRepo.GetVisitorCountWithFilter(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get visitor count with filter: %w", err)
	}
	pageViews, err := s.analyticsRepo.GetPageViewCountWithFilter(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get page view count with filter: %w", err)
	}
	sessions, err := s.analyticsRepo.GetSessionCountWithFilter(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get session count with filter: %w", err)
	}
	bounceRate, err := s.analyticsRepo.GetBounceRateWithFilter(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get bounce rate with filter: %w", err)
	}
	avgDuration, err := s.analyticsRepo.GetAvgSessionDurationWithFilter(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get average session duration with filter: %w", err)
	}

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
