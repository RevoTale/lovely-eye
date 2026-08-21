import type { FunctionComponent } from 'react';
import { PAGE_SIZES } from '@/features/analytics/model/dashboard-query';
import {
  createEmptyDashboardStats,
  extractStatsData,
} from '@/features/analytics/model/dashboard-utils';
import type { DashboardData } from '@/features/analytics/model/use-dashboard-data';
import type { AnalyticsPageState } from '@/features/analytics/model/use-dashboard-page-state';
import AnalyticsOverviewSection from '@/features/analytics/ui/analytics-overview-section';
import AnalyticsPlatformBreakdownSection from '@/features/analytics/ui/analytics-platform-breakdown-section';
import AnalyticsTrafficBreakdownSection from '@/features/analytics/ui/analytics-traffic-breakdown-section';
import EventsSection from '@/features/analytics/ui/events-section';
import { DashboardStatsFieldsFragmentDoc } from '@/shared/api/generated/graphql';
import { readFragment } from '@/shared/api/read-fragment';

interface AnalyticsContentProps {
  dashboard: DashboardData;
  state: AnalyticsPageState;
}

const AnalyticsContent: FunctionComponent<AnalyticsContentProps> = ({ dashboard, state }) => {
  const statsData =
    dashboard.stats === undefined
      ? createEmptyDashboardStats()
      : readFragment(DashboardStatsFieldsFragmentDoc, dashboard.stats);
  const statsCollections = extractStatsData(dashboard.stats);

  return (
    <>
      <AnalyticsOverviewSection
        siteId={state.siteId}
        stats={statsData}
        dashboardState={dashboard.dashboardState}
        realtime={dashboard.realtime}
        dateRange={state.dateRangeForChart}
        filter={state.filterInput}
        chartBucket={state.statsBucket}
        onChartBucketChange={state.setStatsBucket}
      />
      <EventsSection
        siteId={state.siteId}
        eventsState={dashboard.eventsState}
        eventCountsState={dashboard.eventCountsState}
        eventsResult={dashboard.eventsResult}
        eventsCounts={dashboard.eventCounts}
        eventsCountsTotal={dashboard.eventCountsTotal}
        page={state.eventsPage}
        pageSize={PAGE_SIZES.EVENTS}
        onPageChange={(page) => state.setPage('eventsPage', page)}
        countsPage={state.eventsCountsPage}
        countsPageSize={PAGE_SIZES.EVENT_COUNTS}
        onCountsPageChange={(page) => state.setPage('eventsCountsPage', page)}
      />
      <AnalyticsTrafficBreakdownSection
        siteId={state.siteId}
        dashboardState={dashboard.dashboardState}
        topPages={statsCollections.topPages}
        topPagesTotal={statsCollections.topPagesTotal}
        topPagesPage={state.topPagesPage}
        topPagesPageSize={PAGE_SIZES.TOP_PAGES}
        onTopPagesPageChange={(page) => state.setPage('topPagesPage', page)}
        pagePathContains={state.pagePathContains}
        onPagePathSearch={state.setPagePathContains}
        referrers={statsCollections.referrersItems}
        referrersTotal={statsCollections.referrersTotal}
        referrersPage={state.referrersPage}
        referrersPageSize={PAGE_SIZES.REFERRERS}
        totalVisitors={statsData.visitors}
        onReferrersPageChange={(page) => state.setPage('referrersPage', page)}
        countries={statsCollections.countriesItems}
        countriesTotal={statsCollections.countriesTotal}
        countriesTotalVisitors={statsCollections.countriesTotalVisitors}
        countriesPage={state.countriesPage}
        countriesPageSize={PAGE_SIZES.COUNTRIES}
        onCountriesPageChange={(page) => state.setPage('countriesPage', page)}
      />
      <AnalyticsPlatformBreakdownSection
        siteId={state.siteId}
        dashboardState={dashboard.dashboardState}
        totalVisitors={statsData.visitors}
        browsers={statsCollections.browsersItems}
        devices={statsCollections.devicesItems}
        devicesTotal={statsCollections.devicesTotal}
        devicesTotalVisitors={statsCollections.devicesTotalVisitors}
        devicesPage={state.devicesPage}
        devicesPageSize={PAGE_SIZES.DEVICES}
        onDevicesPageChange={(page) => state.setPage('devicesPage', page)}
        operatingSystems={statsCollections.operatingSystemsItems}
        operatingSystemsTotal={statsCollections.operatingSystemsTotal}
        operatingSystemsTotalVisitors={statsCollections.operatingSystemsTotalVisitors}
        osPage={state.osPage}
        osPageSize={PAGE_SIZES.OS}
        onOSPageChange={(page) => state.setPage('osPage', page)}
      />
    </>
  );
};

export default AnalyticsContent;
