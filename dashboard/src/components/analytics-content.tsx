import React from 'react';
import type {
  CountryStats,
  DashboardQuery,
  DeviceStats,
  EventCount,
  EventsResult,
  PageStats,
  ReferrerStats,
  RealtimeStats,
} from '@/gql/graphql';
import { Users, Eye, Clock, TrendingDown } from 'lucide-react';
import { StatCard } from '@/components/stat-card';
import { OverviewChartSection } from '@/components/overview-chart-section';
import { ActivePagesCard } from '@/components/active-pages-card';
import { EventsSection } from '@/components/events-section';
import { TopPagesCard } from '@/components/top-pages-card';
import { ReferrersCard } from '@/components/referrers-card';
import { CountryCard } from '@/components/country-card';
import { DevicesCard } from '@/components/devices-card';

const SECONDS_PER_MINUTE = 60;

interface AnalyticsContentProps {
  siteId: string;
  stats: DashboardQuery['dashboard'];
  dateRange: { from: Date; to: Date } | null;
  filter: Record<string, string[]> | null;
  chartBucket: 'daily' | 'hourly';
  onChartBucketChange: (bucket: 'daily' | 'hourly') => void;
  realtime: RealtimeStats | undefined;
  eventsLoading: boolean;
  eventsResult: EventsResult | undefined;
  eventsCounts: EventCount[];
  eventsPage: number;
  eventsPageSize: number;
  onEventsPageChange: (page: number) => void;
  topPages: PageStats[];
  topPagesTotal: number;
  topPagesPage: number;
  topPagesPageSize: number;
  onTopPagesPageChange: (page: number) => void;
  referrers: ReferrerStats[];
  referrersTotal: number;
  referrersPage: number;
  referrersPageSize: number;
  onReferrersPageChange: (page: number) => void;
  countries: CountryStats[];
  countriesTotal: number;
  countriesTotalVisitors: number;
  countriesPage: number;
  countriesPageSize: number;
  onCountriesPageChange: (page: number) => void;
  devices: DeviceStats[];
  devicesTotal: number;
  devicesTotalVisitors: number;
  devicesPage: number;
  devicesPageSize: number;
  onDevicesPageChange: (page: number) => void;
}

function formatDuration(seconds: number): string {
  if (seconds < SECONDS_PER_MINUTE) {
    return `${String(Math.round(seconds))}s`;
  }
  const minutes = Math.floor(seconds / SECONDS_PER_MINUTE);
  const remainingSeconds = Math.round(seconds % SECONDS_PER_MINUTE);
  return `${String(minutes)}m ${String(remainingSeconds)}s`;
}

export function AnalyticsContent({
  siteId,
  stats,
  dateRange,
  filter,
  chartBucket,
  onChartBucketChange,
  realtime,
  eventsLoading,
  eventsResult,
  eventsCounts,
  eventsPage,
  eventsPageSize,
  onEventsPageChange,
  topPages,
  topPagesTotal,
  topPagesPage,
  topPagesPageSize,
  onTopPagesPageChange,
  referrers,
  referrersTotal,
  referrersPage,
  referrersPageSize,
  onReferrersPageChange,
  countries,
  countriesTotal,
  countriesTotalVisitors,
  countriesPage,
  countriesPageSize,
  onCountriesPageChange,
  devices,
  devicesTotal,
  devicesTotalVisitors,
  devicesPage,
  devicesPageSize,
  onDevicesPageChange,
}: AnalyticsContentProps): React.JSX.Element {
  const activePages = realtime?.activePages;
  const hasActivePages = activePages !== undefined;

  return (
    <>
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <StatCard
          title="Total Visitors"
          value={stats.visitors.toLocaleString()}
          icon={Users}
        />
        <StatCard
          title="Page Views"
          value={stats.pageViews.toLocaleString()}
          icon={Eye}
        />
        <StatCard
          title="Avg. Session"
          value={formatDuration(stats.avgDuration)}
          icon={Clock}
        />
        <StatCard
          title="Bounce Rate"
          value={`${String(Math.round(stats.bounceRate))}%`}
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
          onPageChange={onTopPagesPageChange}
        />
        <ReferrersCard
          referrers={referrers}
          totalCount={referrersTotal}
          totalVisitors={stats.visitors}
          siteId={siteId}
          page={referrersPage}
          pageSize={referrersPageSize}
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
        onPageChange={onCountriesPageChange}
      />

      <DevicesCard
        devices={devices}
        total={devicesTotal}
        totalVisitors={devicesTotalVisitors}
        page={devicesPage}
        pageSize={devicesPageSize}
        siteId={siteId}
        onPageChange={onDevicesPageChange}
      />
    </>
  );
}
