import { useNavigate, useParams, useSearch } from '@tanstack/react-router';
import { useEffect, useMemo, useRef } from 'react';
import {
  ANALYTICS_ROUTE_ID,
  type AnalyticsPageKey,
  type AnalyticsSearch,
  clearPagination,
  setAnalyticsPage,
  setPagePathSearch,
} from '@/features/analytics/model/analytics-search';
import { buildFilterInput, buildFilters } from '@/features/analytics/model/dashboard-utils';
import type { DatePreset } from '@/features/analytics/model/date-range';
import { useDateRange } from '@/features/analytics/model/use-date-range';
import type { FilterInput } from '@/shared/api/generated/graphql';

const DEFAULT_STATS_BUCKET = 'daily';

export interface AnalyticsPageState {
  siteId: string;
  decodedSearch: AnalyticsSearch;
  filterInput: FilterInput | null;
  dateRange: { from: string; to: string } | null | undefined;
  dateRangeForChart: { from: Date; to: Date } | null;
  preset: DatePreset;
  fromDate: string;
  toDate: string;
  fromTime: string;
  toTime: string;
  setPreset: (preset: DatePreset) => void;
  applyCustomRange: (range: {
    fromDate: string;
    toDate: string;
    fromTime: string;
    toTime: string;
  }) => boolean;
  eventsPage: number;
  eventsCountsPage: number;
  topPagesPage: number;
  referrersPage: number;
  devicesPage: number;
  osPage: number;
  countriesPage: number;
  statsBucket: 'daily' | 'hourly';
  setStatsBucket: (bucket: 'daily' | 'hourly') => void;
  setPage: (key: AnalyticsPageKey, page: number) => void;
  pagePathContains: string;
  setPagePathContains: (value: string) => void;
}

export function useDashboardPageState(): AnalyticsPageState {
  const { siteId } = useParams({ from: ANALYTICS_ROUTE_ID });
  const search = useSearch({ from: ANALYTICS_ROUTE_ID });
  const navigate = useNavigate();
  const { preset, fromDate, toDate, fromTime, toTime, dateRange, setPreset, applyCustomRange } =
    useDateRange();
  const eventsPage = search.eventsPage ?? 1;
  const eventsCountsPage = search.eventsCountsPage ?? 1;
  const topPagesPage = search.topPagesPage ?? 1;
  const referrersPage = search.referrersPage ?? 1;
  const devicesPage = search.devicesPage ?? 1;
  const osPage = search.osPage ?? 1;
  const countriesPage = search.countriesPage ?? 1;
  const statsBucket = search.statsBucket ?? DEFAULT_STATS_BUCKET;
  const {
    referrers,
    browsers,
    devices,
    operatingSystems,
    pages,
    countries,
    eventNames,
    eventPaths,
    pagePathContains,
    decodedSearch,
    filter,
  } = useMemo(() => buildFilters(search), [search]);
  const filterInput = useMemo(
    () => buildFilterInput(filter, pagePathContains),
    [filter, pagePathContains]
  );
  const filterKey = useMemo(
    () =>
      JSON.stringify({
        referrers,
        browsers,
        devices,
        operatingSystems,
        pages,
        countries,
        eventNames,
        eventPaths,
        pagePathContains,
      }),
    [
      browsers,
      countries,
      devices,
      eventNames,
      eventPaths,
      operatingSystems,
      pagePathContains,
      pages,
      referrers,
    ]
  );
  const dateRangeForChart = useMemo(
    () =>
      dateRange === undefined
        ? null
        : { from: new Date(dateRange.from), to: new Date(dateRange.to) },
    [dateRange]
  );
  const paginationScopeKey = `${siteId}|${dateRange?.from ?? ''}|${dateRange?.to ?? ''}|${filterKey}`;
  const previousPaginationScopeKey = useRef(paginationScopeKey);

  useEffect(() => {
    if (previousPaginationScopeKey.current === paginationScopeKey) return;
    previousPaginationScopeKey.current = paginationScopeKey;
    void navigate({
      to: '/sites/$siteId/analytics',
      params: { siteId },
      search: clearPagination,
    });
  }, [navigate, paginationScopeKey, siteId]);

  return {
    siteId,
    decodedSearch,
    filterInput,
    dateRange,
    dateRangeForChart,
    preset,
    fromDate,
    toDate,
    fromTime,
    toTime,
    setPreset,
    applyCustomRange,
    eventsPage,
    eventsCountsPage,
    topPagesPage,
    referrersPage,
    devicesPage,
    osPage,
    countriesPage,
    statsBucket,
    pagePathContains: pagePathContains ?? '',
    setStatsBucket: (bucket) =>
      void navigate({
        resetScroll: false,
        to: '/sites/$siteId/analytics',
        params: { siteId },
        search: (prev) => {
          const { statsBucket: _currentBucket, ...rest } = prev;
          return bucket === DEFAULT_STATS_BUCKET ? rest : { ...rest, statsBucket: bucket };
        },
      }),
    setPage: (key, page) =>
      void navigate({
        to: '/sites/$siteId/analytics',
        params: { siteId },
        search: (prev) => setAnalyticsPage(prev, key, page),
      }),
    setPagePathContains: (value) =>
      void navigate({
        to: '/sites/$siteId/analytics',
        params: { siteId },
        search: (prev) => setPagePathSearch(prev, value),
      }),
  };
}
