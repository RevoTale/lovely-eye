import React, { useMemo, useEffect } from 'react';
import { useNavigate, useParams, useSearch } from '@tanstack/react-router';
import { siteDetailRoute } from '@/router';
import { ActiveFilters } from '@/components/active-filters';
import { TimeRangeCard } from '@/components/time-range-card';
import { AnalyticsContent } from '@/components/analytics-content';
import { useDateRange } from '@/hooks/use-date-range';
import { DashboardHeader } from '@/components/dashboard-header';
import { DashboardEmptyState, DashboardLoading, DashboardNotFound } from '@/components/dashboard-states';
import { useDashboardData, PAGE_SIZES } from '@/hooks/use-dashboard-data';
import { parsePage, normalizeStatsBucket, buildFilters, extractStatsData } from '@/lib/dashboard-utils';
import { clearPaginationParams, updatePageParam } from '@/lib/dashboard-navigation';
import { AnalyticsSkeleton } from '@/components/analytics-skeleton';

const EMPTY_COUNT = 0;
const DEFAULT_STATS_BUCKET = 'daily';

export function DashboardPage(): React.JSX.Element {
  const { siteId } = useParams({ from: siteDetailRoute.id });
  const search = useSearch({ from: siteDetailRoute.id });
  const navigate = useNavigate();
  const { preset, fromDate, toDate, fromTime, toTime, dateRange, setPreset, applyCustomRange } = useDateRange();
  const eventsPage = useMemo(() => parsePage(search.eventsPage), [search.eventsPage]);
  const topPagesPage = useMemo(() => parsePage(search.topPagesPage), [search.topPagesPage]);
  const referrersPage = useMemo(() => parsePage(search.referrersPage), [search.referrersPage]);
  const devicesPage = useMemo(() => parsePage(search.devicesPage), [search.devicesPage]);
  const countriesPage = useMemo(() => parsePage(search.countriesPage), [search.countriesPage]);
  const statsBucket = useMemo(() => normalizeStatsBucket(search.statsBucket), [search.statsBucket]);

  const { referrers, devices, pages, countries, decodedSearch, filter } = useMemo(() => buildFilters(search), [search]);

  const filterKey = useMemo(
    () => [referrers, devices, pages, countries].map((v) => v.join(',')).join('|'),
    [referrers, devices, pages, countries]
  );

  const dateRangeForChart = useMemo(() => {
    if (dateRange === undefined) return null;
    return { from: new Date(dateRange.from), to: new Date(dateRange.to) };
  }, [dateRange]);

  const { site, stats, realtime, eventsResult, eventsCounts, siteLoading, dashboardLoading, eventsLoading } =
    useDashboardData({
      siteId,
      dateRange,
      filter: Object.keys(filter).length > EMPTY_COUNT ? filter : null,
      eventsPage,
      topPagesPage,
      referrersPage,
      devicesPage,
      countriesPage,
    });

  useEffect(() => {
    void navigate({
      to: '/sites/$siteId',
      params: { siteId },
      search: clearPaginationParams,
    });
  }, [siteId, dateRange?.from, dateRange?.to, filterKey, navigate]);

  const setPage = (key: string, nextPage: number): void => {
    void navigate({
      to: '/sites/$siteId',
      params: { siteId },
      search: (prev) => updatePageParam(prev as Record<string, unknown>, key, nextPage),
    });
  };

  const setStatsBucket = (bucket: 'daily' | 'hourly'): void => {
    void navigate({
      resetScroll: false,
      to: '/sites/$siteId',
      params: { siteId },
      search: (prev) => ({
        ...(prev as Record<string, unknown>),
        statsBucket: bucket === DEFAULT_STATS_BUCKET ? undefined : bucket,
      }),
    });
  };

  if (siteLoading) {
    return <DashboardLoading />;
  }

  if (site === null || site === undefined) {
    return <DashboardNotFound />;
  }

  const {
    topPages,
    topPagesTotal,
    referrersItems,
    referrersTotal,
    devicesItems,
    devicesTotal,
    devicesTotalVisitors,
    countriesItems,
    countriesTotal,
    countriesTotalVisitors,
  } = extractStatsData(stats);

  const hasStats = stats !== undefined;
  const showSkeletons = dashboardLoading && !hasStats;
  const isRefreshing = dashboardLoading && hasStats;

  return (
    <div className="space-y-8">
      <DashboardHeader site={site} siteId={siteId} realtime={realtime} />

      <TimeRangeCard
        preset={preset}
        fromDate={fromDate}
        toDate={toDate}
        fromTime={fromTime}
        toTime={toTime}
        onPresetChange={setPreset}
        onApplyRange={applyCustomRange}
      />

      <ActiveFilters siteId={siteId} search={decodedSearch} />

      {hasStats ? (
        <AnalyticsContent
          siteId={siteId}
          stats={stats}
          dateRange={dateRangeForChart}
          filter={Object.keys(filter).length > EMPTY_COUNT ? filter : null}
          chartBucket={statsBucket}
          onChartBucketChange={setStatsBucket}
          realtime={realtime}
          eventsLoading={eventsLoading}
          eventsResult={eventsResult}
          eventsCounts={eventsCounts}
          eventsPage={eventsPage}
          eventsPageSize={PAGE_SIZES.EVENTS}
          onEventsPageChange={(page) => {
            setPage('eventsPage', page);
          }}
          topPages={topPages}
          topPagesTotal={topPagesTotal}
          topPagesPage={topPagesPage}
          topPagesPageSize={PAGE_SIZES.TOP_PAGES}
          topPagesLoading={isRefreshing}
          onTopPagesPageChange={(page) => {
            setPage('topPagesPage', page);
          }}
          referrers={referrersItems}
          referrersTotal={referrersTotal}
          referrersPage={referrersPage}
          referrersPageSize={PAGE_SIZES.REFERRERS}
          referrersLoading={isRefreshing}
          onReferrersPageChange={(page) => {
            setPage('referrersPage', page);
          }}
          countries={countriesItems}
          countriesTotal={countriesTotal}
          countriesTotalVisitors={countriesTotalVisitors}
          countriesPage={countriesPage}
          countriesPageSize={PAGE_SIZES.COUNTRIES}
          countriesLoading={isRefreshing}
          onCountriesPageChange={(page) => {
            setPage('countriesPage', page);
          }}
          devices={devicesItems}
          devicesTotal={devicesTotal}
          devicesTotalVisitors={devicesTotalVisitors}
          devicesPage={devicesPage}
          devicesPageSize={PAGE_SIZES.DEVICES}
          devicesLoading={isRefreshing}
          onDevicesPageChange={(page) => {
            setPage('devicesPage', page);
          }}
        />
      ) : showSkeletons ? (
        <AnalyticsSkeleton
          siteId={siteId}
          dateRangeForChart={dateRangeForChart}
          filter={Object.keys(filter).length > EMPTY_COUNT ? filter : null}
          statsBucket={statsBucket}
          realtime={realtime}
          onStatsBucketChange={setStatsBucket}
          onPageChange={setPage}
        />
      ) : (
        <DashboardEmptyState />
      )}
    </div>
  );
}
