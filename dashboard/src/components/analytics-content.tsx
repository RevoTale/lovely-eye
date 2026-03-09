import type { FunctionComponent } from 'react';
import AnalyticsOverviewSection from '@/components/analytics-overview-section';
import AnalyticsPlatformBreakdownSection from '@/components/analytics-platform-breakdown-section';
import AnalyticsTrafficBreakdownSection from '@/components/analytics-traffic-breakdown-section';
import EventsSection from '@/components/events-section';
import { DashboardStatsFieldsFragmentDoc, type DashboardQuery, type EventCountsQuery, type EventsQuery, type FilterInput, type RealtimeQuery } from '@/gql/graphql';
import { useFragment as getFragmentData } from '@/gql/fragment-masking';
import type { DashboardLoadState } from '@/lib/dashboard-load-state';
import { createEmptyDashboardStats, extractStatsData } from '@/lib/dashboard-utils';

interface AnalyticsContentProps {
  siteId: string;
  stats: DashboardQuery['dashboard'] | undefined;
  dashboardState: DashboardLoadState;
  realtime: RealtimeQuery['realtime'] | undefined;
  dateRange: { from: Date; to: Date } | null;
  filter: FilterInput | null;
  chartBucket: 'daily' | 'hourly';
  onChartBucketChange: (bucket: 'daily' | 'hourly') => void;
  eventsResult: EventsQuery['events'] | undefined;
  eventsState: DashboardLoadState;
  eventCounts: EventCountsQuery['eventCounts'];
  eventCountsState: DashboardLoadState;
  eventsPage: number;
  eventsCountsPage: number;
  topPagesPage: number;
  referrersPage: number;
  devicesPage: number;
  osPage: number;
  countriesPage: number;
  onPageChange: (key: string, page: number) => void;
  pageSizes: { EVENTS: number; EVENT_COUNTS: number; TOP_PAGES: number; REFERRERS: number; DEVICES: number; OS: number; COUNTRIES: number };
}

const AnalyticsContent: FunctionComponent<AnalyticsContentProps> = ({ siteId, stats, dashboardState, realtime, dateRange, filter, chartBucket, onChartBucketChange, eventsResult, eventsState, eventCounts, eventCountsState, eventsPage, eventsCountsPage, topPagesPage, referrersPage, devicesPage, osPage, countriesPage, onPageChange, pageSizes }) => {
  const statsData = stats === undefined ? createEmptyDashboardStats() : getFragmentData(DashboardStatsFieldsFragmentDoc, stats);
  const statsCollections = extractStatsData(stats);

  return (
    <>
      <AnalyticsOverviewSection siteId={siteId} stats={statsData} dashboardState={dashboardState} realtime={realtime} dateRange={dateRange} filter={filter} chartBucket={chartBucket} onChartBucketChange={onChartBucketChange} />
      <EventsSection siteId={siteId} eventsState={eventsState} eventCountsState={eventCountsState} eventsResult={eventsResult} eventsCounts={eventCounts} page={eventsPage} pageSize={pageSizes.EVENTS} onPageChange={(page) => onPageChange('eventsPage', page)} countsPage={eventsCountsPage} countsPageSize={pageSizes.EVENT_COUNTS} onCountsPageChange={(page) => onPageChange('eventsCountsPage', page)} />
      <AnalyticsTrafficBreakdownSection siteId={siteId} dashboardState={dashboardState} topPages={statsCollections.topPages} topPagesTotal={statsCollections.topPagesTotal} topPagesPage={topPagesPage} topPagesPageSize={pageSizes.TOP_PAGES} onTopPagesPageChange={(page) => onPageChange('topPagesPage', page)} referrers={statsCollections.referrersItems} referrersTotal={statsCollections.referrersTotal} referrersPage={referrersPage} referrersPageSize={pageSizes.REFERRERS} totalVisitors={statsData.visitors} onReferrersPageChange={(page) => onPageChange('referrersPage', page)} countries={statsCollections.countriesItems} countriesTotal={statsCollections.countriesTotal} countriesTotalVisitors={statsCollections.countriesTotalVisitors} countriesPage={countriesPage} countriesPageSize={pageSizes.COUNTRIES} onCountriesPageChange={(page) => onPageChange('countriesPage', page)} />
      <AnalyticsPlatformBreakdownSection siteId={siteId} dashboardState={dashboardState} totalVisitors={statsData.visitors} browsers={statsCollections.browsersItems} devices={statsCollections.devicesItems} devicesTotal={statsCollections.devicesTotal} devicesTotalVisitors={statsCollections.devicesTotalVisitors} devicesPage={devicesPage} devicesPageSize={pageSizes.DEVICES} onDevicesPageChange={(page) => onPageChange('devicesPage', page)} operatingSystems={statsCollections.operatingSystemsItems} operatingSystemsTotal={statsCollections.operatingSystemsTotal} operatingSystemsTotalVisitors={statsCollections.operatingSystemsTotalVisitors} osPage={osPage} osPageSize={pageSizes.OS} onOSPageChange={(page) => onPageChange('osPage', page)} />
    </>
  );
};

export default AnalyticsContent;
