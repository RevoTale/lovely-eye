import React, { useMemo, useEffect } from 'react';
import { useNavigate, useParams, useSearch } from '@tanstack/react-router';
import { useQuery } from '@apollo/client/react';
import {
  DashboardDocument,
  RealtimeDocument,
  EventsDocument,
  SiteDocument,
} from '@/gql/graphql';
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

const EMPTY_COUNT = 0;
const FIRST_INDEX = 0;
const FIRST_PAGE = 1;
const PAGE_INDEX_OFFSET = 1;
const DASHBOARD_POLL_INTERVAL_MS = 60000;
const REALTIME_POLL_INTERVAL_MS = 5000;
const DEFAULT_STATS_BUCKET = 'daily';
const DAILY_STATS_LIMIT = 365;
const HOURLY_STATS_LIMIT = 168;

type PageValue = string | string[] | undefined;

function parsePage(value: PageValue): number {
  const raw = Array.isArray(value) ? value[FIRST_INDEX] : value;
  const numeric = Number(raw);
  if (!Number.isFinite(numeric) || numeric < FIRST_PAGE) {
    return FIRST_PAGE;
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
  const referrers = normalizeFilterValue(search['referrer']);
  const devices = normalizeFilterValue(search['device']);
  const pages = normalizeFilterValue(search['page']);
  const countries = normalizeFilterValue(search['country']);
  const decodedSearch = {
    ...search,
    ...(referrers.length > EMPTY_COUNT ? { referrer: referrers } : {}),
    ...(devices.length > EMPTY_COUNT ? { device: devices } : {}),
    ...(pages.length > EMPTY_COUNT ? { page: pages } : {}),
    ...(countries.length > EMPTY_COUNT ? { country: countries } : {}),
  };

  const filter = {
    ...(referrers.length > EMPTY_COUNT ? { referrer: referrers } : {}),
    ...(devices.length > EMPTY_COUNT ? { device: devices } : {}),
    ...(pages.length > EMPTY_COUNT ? { page: pages } : {}),
    ...(countries.length > EMPTY_COUNT ? { country: countries } : {}),
  };

  return { referrers, devices, pages, countries, decodedSearch, filter };
}

function normalizeStatsBucket(value: PageValue): 'daily' | 'hourly' {
  const raw = Array.isArray(value) ? value[FIRST_INDEX] : value;
  return raw === 'hourly' ? 'hourly' : 'daily';
}

// eslint-disable-next-line complexity -- DashboardPage orchestrates multiple sections and filters.
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

  const hasSiteId = siteId !== '';
  const { data: siteData, loading: siteLoading } = useQuery(SiteDocument, {
    variables: { id: siteId },
    skip: !hasSiteId,
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
        if (nextPage <= FIRST_PAGE) {
          return Object.fromEntries(
            Object.entries(prev).filter(([entryKey]) => entryKey !== key)
          );
        }
        return { ...(prev as Record<string, unknown>), [key]: String(nextPage) };
      },
    });
  };

  const setStatsBucket = (bucket: 'daily' | 'hourly'): void => {
    void navigate({
      resetScroll:false,
      to: '/sites/$siteId',
      params: { siteId },
      search: (prev) => ({
        ...(prev as Record<string, unknown>),
        statsBucket: bucket === DEFAULT_STATS_BUCKET ? undefined : bucket,
      }),
    });
  };

  const dailyStatsLimit = statsBucket === 'hourly' ? HOURLY_STATS_LIMIT : DAILY_STATS_LIMIT;

  const { data: dashboardData, loading: dashboardLoading } = useQuery(DashboardDocument, {
    variables: {
      siteId,
      dateRange: dateRange ?? null,
      filter: Object.keys(filter).length > EMPTY_COUNT ? filter : null,
      topPagesPaging: {
        limit: TOP_PAGES_PAGE_SIZE,
        offset: (topPagesPage - PAGE_INDEX_OFFSET) * TOP_PAGES_PAGE_SIZE,
      },
      referrersPaging: {
        limit: REFERRERS_PAGE_SIZE,
        offset: (referrersPage - PAGE_INDEX_OFFSET) * REFERRERS_PAGE_SIZE,
      },
      devicesPaging: {
        limit: DEVICES_PAGE_SIZE,
        offset: (devicesPage - PAGE_INDEX_OFFSET) * DEVICES_PAGE_SIZE,
      },
      countriesPaging: {
        limit: COUNTRIES_PAGE_SIZE,
        offset: (countriesPage - PAGE_INDEX_OFFSET) * COUNTRIES_PAGE_SIZE,
      },
      dailyStatsBucket: statsBucket === 'hourly' ? 'HOURLY' : 'DAILY',
      dailyStatsLimit,
    },
    skip: !hasSiteId,
    pollInterval: DASHBOARD_POLL_INTERVAL_MS,
  });

  const { data: realtimeData } = useQuery(RealtimeDocument, {
    variables: { siteId },
    skip: !hasSiteId,
    pollInterval: REALTIME_POLL_INTERVAL_MS,
  });

  const { data: eventsData, loading: eventsLoading } = useQuery(EventsDocument, {
    variables: {
      siteId,
      dateRange: dateRange ?? null,
      limit: EVENTS_PAGE_SIZE,
      offset: (eventsPage - PAGE_INDEX_OFFSET) * EVENTS_PAGE_SIZE,
    },
    skip: !hasSiteId,
    pollInterval: DASHBOARD_POLL_INTERVAL_MS,
  });

  const { data: eventsCountsData } = useQuery(EventsDocument, {
    variables: {
      siteId,
      dateRange: dateRange ?? null,
      limit: EVENTS_COUNT_LIMIT,
      offset: EMPTY_COUNT,
    },
    skip: !hasSiteId,
    pollInterval: DASHBOARD_POLL_INTERVAL_MS,
  });

  if (siteLoading || dashboardLoading) {
    return <DashboardLoading />;
  }

  const site = siteData?.site;
  const stats = dashboardData?.dashboard;
  const realtime = realtimeData?.realtime;
  const eventsResult = eventsData?.events;
  const eventsCounts = eventsCountsData?.events.events ?? [];
  const topPagesResult = stats?.topPages;
  const referrersResult = stats?.topReferrers;
  const devicesResult = stats?.devices;
  const countriesResult = stats?.countries;
  const topPages = topPagesResult?.items ?? [];
  const referrersItems = referrersResult?.items ?? [];
  const devicesItems = devicesResult?.items ?? [];
  const countriesItems = countriesResult?.items ?? [];
  const devicesTotalVisitors = devicesResult?.totalVisitors ?? EMPTY_COUNT;
  const countriesTotalVisitors = countriesResult?.totalVisitors ?? EMPTY_COUNT;

  if (site === null || site === undefined) {
    return <DashboardNotFound />;
  }

  const hasStats = stats !== undefined;

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
          chartBucket={statsBucket}
          onChartBucketChange={setStatsBucket}
          realtime={realtime}
          eventsLoading={eventsLoading}
          eventsResult={eventsResult}
          eventsCounts={eventsCounts}
          eventsPage={eventsPage}
          eventsPageSize={EVENTS_PAGE_SIZE}
          onEventsPageChange={(page) => {
            setPage('eventsPage', page);
          }}
          topPages={topPages}
          topPagesTotal={topPagesResult?.total ?? EMPTY_COUNT}
          topPagesPage={topPagesPage}
          topPagesPageSize={TOP_PAGES_PAGE_SIZE}
          onTopPagesPageChange={(page) => {
            setPage('topPagesPage', page);
          }}
          referrers={referrersItems}
          referrersTotal={referrersResult?.total ?? EMPTY_COUNT}
          referrersPage={referrersPage}
          referrersPageSize={REFERRERS_PAGE_SIZE}
          onReferrersPageChange={(page) => {
            setPage('referrersPage', page);
          }}
          countries={countriesItems}
          countriesTotal={countriesResult?.total ?? EMPTY_COUNT}
          countriesTotalVisitors={countriesTotalVisitors}
          countriesPage={countriesPage}
          countriesPageSize={COUNTRIES_PAGE_SIZE}
          onCountriesPageChange={(page) => {
            setPage('countriesPage', page);
          }}
          devices={devicesItems}
          devicesTotal={devicesResult?.total ?? EMPTY_COUNT}
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
