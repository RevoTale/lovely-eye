import React, { useMemo, useState, useEffect } from 'react';
import { useParams, useSearch } from '@tanstack/react-router';
import { useQuery } from '@apollo/client';
import { DASHBOARD_QUERY, REALTIME_QUERY, SITE_QUERY, EVENTS_QUERY } from '@/graphql';
import type { DashboardStats, Site, RealtimeStats } from '@/generated/graphql';
import { siteDetailRoute } from '@/router';
import { ActiveFilters } from '@/components/active-filters';
import { TimeRangeCard } from '@/components/time-range-card';
import { AnalyticsContent } from '@/components/analytics-content';
import { normalizeFilterValue } from '@/lib/filter-utils';
import { useDateRange } from '@/hooks/use-date-range';
import { DashboardHeader } from '@/components/dashboard-header';
import { DashboardEmptyState, DashboardLoading, DashboardNotFound } from '@/components/dashboard-states';

const EVENTS_PAGE_SIZE = 5;
const EVENTS_COUNT_LIMIT = 200;
const TOP_PAGES_PAGE_SIZE = 5;
const REFERRERS_PAGE_SIZE = 5;
const DEVICES_PAGE_SIZE = 6;

function buildFilters(search: Record<string, string | string[] | undefined>): {
  referrers: string[];
  devices: string[];
  pages: string[];
  countries: string[];
  decodedSearch: Record<string, unknown>;
  filter: Record<string, string[]>;
} {
  const referrers = normalizeFilterValue(search.referrer);
  const devices = normalizeFilterValue(search.device);
  const pages = normalizeFilterValue(search.page);
  const countries = normalizeFilterValue(search.country);
  const decodedSearch = {
    ...search,
    ...(referrers.length ? { referrer: referrers } : {}),
    ...(devices.length ? { device: devices } : {}),
    ...(pages.length ? { page: pages } : {}),
    ...(countries.length ? { country: countries } : {}),
  };

  const filter = {
    ...(referrers.length ? { referrer: referrers } : {}),
    ...(devices.length ? { device: devices } : {}),
    ...(pages.length ? { page: pages } : {}),
    ...(countries.length ? { country: countries } : {}),
  };

  return { referrers, devices, pages, countries, decodedSearch, filter };
}

export function DashboardPage(): React.JSX.Element {
  const { siteId } = useParams({ from: siteDetailRoute.id });
  const search = useSearch({ from: siteDetailRoute.id });
  const { preset, fromDate, toDate, fromTime, toTime, dateRange, setPreset, applyCustomRange } = useDateRange();
  const [eventsPage, setEventsPage] = useState(1);
  const [topPagesPage, setTopPagesPage] = useState(1);
  const [referrersPage, setReferrersPage] = useState(1);
  const [devicesPage, setDevicesPage] = useState(1);

  const { data: siteData, loading: siteLoading } = useQuery(SITE_QUERY, {
    variables: { id: siteId },
    skip: !siteId,
  });

  const { referrers, devices, pages, countries, decodedSearch, filter } = useMemo(
    () => buildFilters(search),
    [search]
  );

  const filterKey = useMemo(
    () => [referrers, devices, pages, countries].map((value) => value.join(',')).join('|'),
    [referrers, devices, pages, countries]
  );

  useEffect(() => {
    setEventsPage(1);
    setTopPagesPage(1);
    setReferrersPage(1);
    setDevicesPage(1);
  }, [siteId, dateRange?.from, dateRange?.to, filterKey]);

  const { data: dashboardData, loading: dashboardLoading } = useQuery(DASHBOARD_QUERY, {
    variables: {
      siteId,
      dateRange,
      filter: Object.keys(filter).length > 0 ? filter : undefined,
    },
    skip: !siteId,
    pollInterval: 60000,
  });

  const { data: realtimeData } = useQuery(REALTIME_QUERY, {
    variables: { siteId },
    skip: !siteId,
    pollInterval: 5000,
  });

  const { data: eventsData, loading: eventsLoading } = useQuery(EVENTS_QUERY, {
    variables: {
      siteId,
      dateRange,
      limit: EVENTS_PAGE_SIZE,
      offset: (eventsPage - 1) * EVENTS_PAGE_SIZE,
    },
    skip: !siteId,
    pollInterval: 60000,
  });

  const { data: eventsCountsData } = useQuery(EVENTS_QUERY, {
    variables: {
      siteId,
      dateRange,
      limit: EVENTS_COUNT_LIMIT,
      offset: 0,
    },
    skip: !siteId,
    pollInterval: 60000,
  });

  if (siteLoading || dashboardLoading) {
    return <DashboardLoading />;
  }

  const site = siteData?.site as Site | undefined;
  const stats = dashboardData?.dashboard as DashboardStats | undefined;
  const realtime = realtimeData?.realtime as RealtimeStats | undefined;
  const eventsResult = eventsData?.events;
  const eventsCounts = eventsCountsData?.events?.events ?? [];
  const topPages = stats?.topPages ?? [];
  const referrersAll = stats?.topReferrers ?? [];
  const devicesAll = stats?.devices ?? [];
  const devicesTotalVisitors = devicesAll.reduce((sum, device) => sum + device.visitors, 0);
  const topPagesSlice = topPages.slice(
    (topPagesPage - 1) * TOP_PAGES_PAGE_SIZE,
    topPagesPage * TOP_PAGES_PAGE_SIZE
  );
  const referrersSlice = referrersAll.slice(
    (referrersPage - 1) * REFERRERS_PAGE_SIZE,
    referrersPage * REFERRERS_PAGE_SIZE
  );
  const devicesSlice = devicesAll.slice(
    (devicesPage - 1) * DEVICES_PAGE_SIZE,
    devicesPage * DEVICES_PAGE_SIZE
  );

  if (!site) {
    return <DashboardNotFound />;
  }

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

      {stats ? (
        <AnalyticsContent
          siteId={siteId}
          stats={stats}
          realtime={realtime}
          eventsLoading={eventsLoading}
          eventsResult={eventsResult}
          eventsCounts={eventsCounts}
          eventsPage={eventsPage}
          eventsPageSize={EVENTS_PAGE_SIZE}
          onEventsPageChange={setEventsPage}
          topPages={topPagesSlice}
          topPagesTotal={topPages.length}
          topPagesPage={topPagesPage}
          topPagesPageSize={TOP_PAGES_PAGE_SIZE}
          onTopPagesPageChange={setTopPagesPage}
          referrers={referrersSlice}
          referrersTotal={referrersAll.length}
          referrersPage={referrersPage}
          referrersPageSize={REFERRERS_PAGE_SIZE}
          onReferrersPageChange={setReferrersPage}
          devices={devicesSlice}
          devicesTotal={devicesAll.length}
          devicesTotalVisitors={devicesTotalVisitors}
          devicesPage={devicesPage}
          devicesPageSize={DEVICES_PAGE_SIZE}
          onDevicesPageChange={setDevicesPage}
        />
      ) : (
        <DashboardEmptyState />
      )}
    </div>
  );
}
