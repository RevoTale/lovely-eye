
import { DashboardStatsFieldsFragmentDoc, RealtimeStatsFieldsFragmentDoc } from '@/gql/graphql';
import type {
  CountryStatsFieldsFragment,
  DeviceStatsFieldsFragment,
  EventCountsQuery,
  EventsQuery,
  FilterInput,
  PageStatsFieldsFragment,
  ReferrerStatsFieldsFragment,
  RealtimeQuery,
} from '@/gql/graphql';
import { useFragment as getFragmentData, type FragmentType } from '@/gql/fragment-masking';
import { Users, Eye, Clock, TrendingDown } from 'lucide-react';
import { StatCard } from '@/components/stat-card';
import { OverviewChartSection } from '@/components/overview-chart-section';
import { ActivePagesCard } from '@/components/active-pages-card';
import { EventsSection } from '@/components/events-section';
import { TopPagesCard } from '@/components/top-pages-card';
import { ReferrersCard } from '@/components/referrers-card';
import { CountryCard } from '@/components/country-card';
import { DevicesCard } from '@/components/devices-card';
import { formatDuration } from '@/lib/dashboard-utils';

interface AnalyticsContentProps {
  siteId: string;
  stats: FragmentType<typeof DashboardStatsFieldsFragmentDoc>;
  dateRange: { from: Date; to: Date } | null;
  filter: FilterInput | null;
  chartBucket: 'daily' | 'hourly';
  onChartBucketChange: (bucket: 'daily' | 'hourly') => void;
  realtime: RealtimeQuery['realtime'] | undefined;
  eventsLoading: boolean;
  eventsResult: EventsQuery['events'] | undefined;
  eventsCounts: EventCountsQuery['eventCounts'];
  eventsPage: number;
  eventsPageSize: number;
  onEventsPageChange: (page: number) => void;
  topPages: PageStatsFieldsFragment[];
  topPagesTotal: number;
  topPagesPage: number;
  topPagesPageSize: number;
  topPagesLoading?: boolean;
  onTopPagesPageChange: (page: number) => void;
  referrers: ReferrerStatsFieldsFragment[];
  referrersTotal: number;
  referrersPage: number;
  referrersPageSize: number;
  referrersLoading?: boolean;
  onReferrersPageChange: (page: number) => void;
  countries: CountryStatsFieldsFragment[];
  countriesTotal: number;
  countriesTotalVisitors: number;
  countriesPage: number;
  countriesPageSize: number;
  countriesLoading?: boolean;
  onCountriesPageChange: (page: number) => void;
  devices: DeviceStatsFieldsFragment[];
  devicesTotal: number;
  devicesTotalVisitors: number;
  devicesPage: number;
  devicesPageSize: number;
  devicesLoading?: boolean;
  onDevicesPageChange: (page: number) => void;
}

export function AnalyticsContent(props: AnalyticsContentProps): React.JSX.Element {
  const {
    siteId, stats, dateRange, filter, chartBucket, onChartBucketChange, realtime,
    eventsLoading, eventsResult, eventsCounts, eventsPage, eventsPageSize, onEventsPageChange,
    topPages, topPagesTotal, topPagesPage, topPagesPageSize, topPagesLoading = false, onTopPagesPageChange,
    referrers, referrersTotal, referrersPage, referrersPageSize, referrersLoading = false, onReferrersPageChange,
    countries, countriesTotal, countriesTotalVisitors, countriesPage, countriesPageSize, countriesLoading = false, onCountriesPageChange,
    devices, devicesTotal, devicesTotalVisitors, devicesPage, devicesPageSize, devicesLoading = false, onDevicesPageChange,
  } = props;
  const statsData = getFragmentData(DashboardStatsFieldsFragmentDoc, stats);
  const realtimeData =
    realtime === undefined ? undefined : getFragmentData(RealtimeStatsFieldsFragmentDoc, realtime);
  const activePages = realtimeData?.activePages;
  const hasActivePages = activePages !== undefined;

  return (
    <>
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <StatCard
          title="Total Visitors"
          value={statsData.visitors.toLocaleString()}
          icon={Users}
        />
        <StatCard
          title="Page Views"
          value={statsData.pageViews.toLocaleString()}
          icon={Eye}
        />
        <StatCard
          title="Avg. Session"
          value={formatDuration(statsData.avgDuration)}
          icon={Clock}
        />
        <StatCard
          title="Bounce Rate"
          value={`${String(Math.round(statsData.bounceRate))}%`}
          icon={TrendingDown}
        />
      </div>

      <OverviewChartSection
        siteId={siteId}
        dateRange={dateRange}
        filter={filter}
        bucket={chartBucket}
        onBucketChange={onChartBucketChange}
      />

      {hasActivePages ? (
        <ActivePagesCard activePages={activePages} />
      ) : null}

      <EventsSection
        loading={eventsLoading}
        eventsResult={eventsResult}
        eventsCounts={eventsCounts}
        page={eventsPage}
        pageSize={eventsPageSize}
        onPageChange={onEventsPageChange}
      />

      <div className="grid gap-6 md:grid-cols-2">
        <TopPagesCard
          pages={topPages}
          total={topPagesTotal}
          page={topPagesPage}
          pageSize={topPagesPageSize}
          siteId={siteId}
          loading={topPagesLoading}
          onPageChange={onTopPagesPageChange}
        />
        <ReferrersCard
          referrers={referrers}
          totalCount={referrersTotal}
          totalVisitors={statsData.visitors}
          siteId={siteId}
          page={referrersPage}
          pageSize={referrersPageSize}
          loading={referrersLoading}
          onPageChange={onReferrersPageChange}
        />
      </div>

      <CountryCard
        countries={countries}
        total={countriesTotal}
        totalVisitors={countriesTotalVisitors}
        page={countriesPage}
        pageSize={countriesPageSize}
        siteId={siteId}
        loading={countriesLoading}
        onPageChange={onCountriesPageChange}
      />

      <DevicesCard
        devices={devices}
        total={devicesTotal}
        totalVisitors={devicesTotalVisitors}
        page={devicesPage}
        pageSize={devicesPageSize}
        siteId={siteId}
        loading={devicesLoading}
        onPageChange={onDevicesPageChange}
      />
    </>
  );
}
