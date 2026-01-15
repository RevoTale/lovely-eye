import React, { useMemo, useEffect } from 'react';
import { useNavigate, useParams, useSearch } from '@tanstack/react-router';
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
const COUNTRIES_PAGE_SIZE = 6;

type PageValue = string | string[] | undefined;

function parsePage(value: PageValue): number {
  const raw = Array.isArray(value) ? value[0] : value;
  const numeric = Number(raw);
  if (!Number.isFinite(numeric) || numeric < 1) {
    return 1;
  }
  return Math.floor(numeric);
}

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
  const navigate = useNavigate();
  const { preset, fromDate, toDate, fromTime, toTime, dateRange, setPreset, applyCustomRange } = useDateRange();
  const eventsPage = useMemo(() => parsePage(search.eventsPage), [search.eventsPage]);
  const topPagesPage = useMemo(() => parsePage(search.topPagesPage), [search.topPagesPage]);
  const referrersPage = useMemo(() => parsePage(search.referrersPage), [search.referrersPage]);
  const devicesPage = useMemo(() => parsePage(search.devicesPage), [search.devicesPage]);
  const countriesPage = useMemo(() => parsePage(search.countriesPage), [search.countriesPage]);

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
    void navigate({
      to: '/sites/$siteId',
      params: { siteId },
      search: (prev) => {
        const keys = new Set(['eventsPage', 'topPagesPage', 'referrersPage', 'devicesPage', 'countriesPage']);
        return Object.fromEntries(
          Object.entries(prev).filter(([key]) => !keys.has(key))
        );
      },
    });
  }, [siteId, dateRange?.from, dateRange?.to, filterKey, navigate]);

  const setPage = (key: string, nextPage: number): void => {
    void navigate({
      to: '/sites/$siteId',
      params: { siteId },
      search: (prev) => {
        if (nextPage <= 1) {
          return Object.fromEntries(
            Object.entries(prev).filter(([entryKey]) => entryKey !== key)
          );
        }
        return { ...(prev as Record<string, unknown>), [key]: String(nextPage) };
      },
    });
  };

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
  const countriesAll = stats?.countries ?? [];
  const devicesTotalVisitors = devicesAll.reduce((sum, device) => sum + device.visitors, 0);
  const countriesTotalVisitors = countriesAll.reduce((sum, country) => sum + country.visitors, 0);
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
  const countriesSlice = countriesAll.slice(
    (countriesPage - 1) * COUNTRIES_PAGE_SIZE,
    countriesPage * COUNTRIES_PAGE_SIZE
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
          onEventsPageChange={(page) => {
            setPage('eventsPage', page);
          }}
          topPages={topPagesSlice}
          topPagesTotal={topPages.length}
          topPagesPage={topPagesPage}
          topPagesPageSize={TOP_PAGES_PAGE_SIZE}
          onTopPagesPageChange={(page) => {
            setPage('topPagesPage', page);
          }}
          referrers={referrersSlice}
          referrersTotal={referrersAll.length}
          referrersPage={referrersPage}
          referrersPageSize={REFERRERS_PAGE_SIZE}
          onReferrersPageChange={(page) => {
            setPage('referrersPage', page);
          }}
          countries={countriesSlice}
          countriesTotal={countriesAll.length}
          countriesTotalVisitors={countriesTotalVisitors}
          countriesPage={countriesPage}
          countriesPageSize={COUNTRIES_PAGE_SIZE}
          onCountriesPageChange={(page) => {
            setPage('countriesPage', page);
          }}
          devices={devicesSlice}
          devicesTotal={devicesAll.length}
          devicesTotalVisitors={devicesTotalVisitors}
          devicesPage={devicesPage}
          devicesPageSize={DEVICES_PAGE_SIZE}
          onDevicesPageChange={(page) => {
            setPage('devicesPage', page);
          }}
        />
      ) : (
        <DashboardEmptyState />
      )}
    </div>
  );
}
