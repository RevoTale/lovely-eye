import React from 'react';
import { AnalyticsContent } from '@/components/analytics-content';
import type { RealtimeStats, DashboardQuery } from '@/gql/graphql';

const EMPTY_COUNT = 0;
const FIRST_PAGE = 1;
const EVENTS_PAGE_SIZE = 5;
const TOP_PAGES_PAGE_SIZE = 5;
const REFERRERS_PAGE_SIZE = 5;
const DEVICES_PAGE_SIZE = 6;
const COUNTRIES_PAGE_SIZE = 6;

interface AnalyticsSkeletonProps {
  siteId: string;
  dateRangeForChart: { from: Date; to: Date } | null;
  filter: Record<string, string[]> | null;
  statsBucket: 'daily' | 'hourly';
  realtime: RealtimeStats | undefined;
  onStatsBucketChange: (bucket: 'daily' | 'hourly') => void;
  onPageChange: (key: string, page: number) => void;
}

export function AnalyticsSkeleton({
  siteId,
  dateRangeForChart,
  filter,
  statsBucket,
  realtime,
  onStatsBucketChange,
  onPageChange,
}: AnalyticsSkeletonProps): React.JSX.Element {
  // eslint-disable-next-line @typescript-eslint/consistent-type-assertions, @typescript-eslint/no-unsafe-type-assertion -- Skeleton requires partial stats object
  const emptyStats: DashboardQuery['dashboard'] = { visitors: EMPTY_COUNT, pageViews: EMPTY_COUNT, sessions: EMPTY_COUNT, bounceRate: EMPTY_COUNT, avgDuration: EMPTY_COUNT } as DashboardQuery['dashboard'];

  return (
    <AnalyticsContent
      siteId={siteId}
      stats={emptyStats}
      dateRange={dateRangeForChart}
      filter={filter}
      chartBucket={statsBucket}
      onChartBucketChange={onStatsBucketChange}
      realtime={realtime}
      eventsLoading={true}
      eventsResult={undefined}
      eventsCounts={[]}
      eventsPage={FIRST_PAGE}
      eventsPageSize={EVENTS_PAGE_SIZE}
      onEventsPageChange={(page) => {
        onPageChange('eventsPage', page);
      }}
      topPages={[]}
      topPagesTotal={EMPTY_COUNT}
      topPagesPage={FIRST_PAGE}
      topPagesPageSize={TOP_PAGES_PAGE_SIZE}
      topPagesLoading={true}
      onTopPagesPageChange={(page) => {
        onPageChange('topPagesPage', page);
      }}
      referrers={[]}
      referrersTotal={EMPTY_COUNT}
      referrersPage={FIRST_PAGE}
      referrersPageSize={REFERRERS_PAGE_SIZE}
      referrersLoading={true}
      onReferrersPageChange={(page) => {
        onPageChange('referrersPage', page);
      }}
      countries={[]}
      countriesTotal={EMPTY_COUNT}
      countriesTotalVisitors={EMPTY_COUNT}
      countriesPage={FIRST_PAGE}
      countriesPageSize={COUNTRIES_PAGE_SIZE}
      countriesLoading={true}
      onCountriesPageChange={(page) => {
        onPageChange('countriesPage', page);
      }}
      devices={[]}
      devicesTotal={EMPTY_COUNT}
      devicesTotalVisitors={EMPTY_COUNT}
      devicesPage={FIRST_PAGE}
      devicesPageSize={DEVICES_PAGE_SIZE}
      devicesLoading={true}
      onDevicesPageChange={(page) => {
        onPageChange('devicesPage', page);
      }}
    />
  );
}
