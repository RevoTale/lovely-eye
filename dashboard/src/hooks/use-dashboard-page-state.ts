import { useEffect, useMemo, useRef } from 'react';
import { useNavigate, useParams, useSearch } from '@tanstack/react-router';
import type { DatePreset } from '@/lib/date-range';
import { siteDetailRoute } from '@/router';
import { useDateRange } from '@/hooks/use-date-range';
import { clearPaginationParams, updatePageParam } from '@/lib/dashboard-navigation';
import { buildFilters, normalizeStatsBucket, parsePage } from '@/lib/dashboard-utils';
import type { FilterInput } from '@/gql/graphql';

const DEFAULT_STATS_BUCKET = 'daily';
const EMPTY_COUNT = 0;

export function useDashboardPageState(): {
  siteId: string;
  decodedSearch: Record<string, unknown>;
  filterInput: FilterInput | null;
  dateRange: { from: string; to: string } | null | undefined;
  dateRangeForChart: { from: Date; to: Date } | null;
  preset: DatePreset;
  fromDate: string;
  toDate: string;
  fromTime: string;
  toTime: string;
  setPreset: (preset: DatePreset) => void;
  applyCustomRange: (range: { fromDate: string; toDate: string; fromTime: string; toTime: string }) => boolean;
  eventsPage: number;
  eventsCountsPage: number;
  topPagesPage: number;
  referrersPage: number;
  devicesPage: number;
  osPage: number;
  countriesPage: number;
  statsBucket: 'daily' | 'hourly';
  setStatsBucket: (bucket: 'daily' | 'hourly') => void;
  setPage: (key: string, page: number) => void;
} {
  const { siteId } = useParams({ from: siteDetailRoute.id });
  const search = useSearch({ from: siteDetailRoute.id });
  const navigate = useNavigate();
  const { preset, fromDate, toDate, fromTime, toTime, dateRange, setPreset, applyCustomRange } = useDateRange();
  const eventsPage = useMemo(() => parsePage(search.eventsPage), [search.eventsPage]);
  const eventsCountsPage = useMemo(() => parsePage(search.eventsCountsPage), [search.eventsCountsPage]);
  const topPagesPage = useMemo(() => parsePage(search.topPagesPage), [search.topPagesPage]);
  const referrersPage = useMemo(() => parsePage(search.referrersPage), [search.referrersPage]);
  const devicesPage = useMemo(() => parsePage(search.devicesPage), [search.devicesPage]);
  const osPage = useMemo(() => parsePage(search.osPage), [search.osPage]);
  const countriesPage = useMemo(() => parsePage(search.countriesPage), [search.countriesPage]);
  const statsBucket = useMemo(() => normalizeStatsBucket(search.statsBucket), [search.statsBucket]);
  const { referrers, browsers, devices, operatingSystems, pages, countries, eventNames, eventPaths, decodedSearch, filter } = useMemo(() => buildFilters(search), [search]);
  const filterInput = useMemo<FilterInput | null>(() => {
    if (Object.keys(filter).length === EMPTY_COUNT) return null;
    const getFilter = (key: string): string[] | null => filter[key] ?? null;
    return {
      referrer: getFilter('referrer'),
      browser: getFilter('browser'),
      device: getFilter('device'),
      os: getFilter('os'),
      page: getFilter('page'),
      country: getFilter('country'),
      eventType: null,
      eventDefinitionId: getFilter('eventDefinitionId'),
      eventName: getFilter('eventName'),
      eventPath: getFilter('eventPath'),
    };
  }, [filter]);
  const filterKey = useMemo(() => [referrers, browsers, devices, operatingSystems, pages, countries, eventNames, eventPaths].map((value) => value.join(',')).join('|'), [browsers, countries, devices, eventNames, eventPaths, operatingSystems, pages, referrers]);
  const dateRangeForChart = useMemo(() => dateRange === undefined ? null : { from: new Date(dateRange.from), to: new Date(dateRange.to) }, [dateRange]);
  const paginationScopeKey = `${siteId}|${dateRange?.from ?? ''}|${dateRange?.to ?? ''}|${filterKey}`;
  const previousPaginationScopeKey = useRef(paginationScopeKey);

  useEffect(() => {
    if (previousPaginationScopeKey.current === paginationScopeKey) return;
    previousPaginationScopeKey.current = paginationScopeKey;
    void navigate({ to: '/sites/$siteId', params: { siteId }, search: clearPaginationParams });
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
    setStatsBucket: (bucket) => void navigate({ resetScroll: false, to: '/sites/$siteId', params: { siteId }, search: (prev) => ({ ...(prev as Record<string, unknown>), statsBucket: bucket === DEFAULT_STATS_BUCKET ? undefined : bucket }) }),
    setPage: (key, page) => void navigate({ to: '/sites/$siteId', params: { siteId }, search: (prev) => updatePageParam(prev as Record<string, unknown>, key, page) }),
  };
}
