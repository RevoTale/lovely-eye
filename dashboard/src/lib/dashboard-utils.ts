import { normalizeFilterValue } from '@/lib/filter-utils';
import type { DashboardQuery } from '@/gql/graphql';

const EMPTY_COUNT = 0;
const FIRST_INDEX = 0;
const FIRST_PAGE = 1;

type PageValue = string | string[] | undefined;

export function parsePage(value: PageValue): number {
  const raw = Array.isArray(value) ? value[FIRST_INDEX] : value;
  const numeric = Number(raw);
  if (!Number.isFinite(numeric) || numeric < FIRST_PAGE) {
    return FIRST_PAGE;
  }
  return Math.floor(numeric);
}

export function normalizeStatsBucket(value: PageValue): 'daily' | 'hourly' {
  const raw = Array.isArray(value) ? value[FIRST_INDEX] : value;
  return raw === 'hourly' ? 'hourly' : 'daily';
}

interface FilterResult {
  referrers: string[];
  devices: string[];
  pages: string[];
  countries: string[];
  decodedSearch: Record<string, unknown>;
  filter: Record<string, string[]>;
}

export function buildFilters(search: Record<string, string | string[] | undefined>): FilterResult {
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

interface StatsDataResult {
  topPages: DashboardQuery['dashboard']['topPages']['items'];
  topPagesTotal: number;
  referrersItems: DashboardQuery['dashboard']['topReferrers']['items'];
  referrersTotal: number;
  devicesItems: DashboardQuery['dashboard']['devices']['items'];
  devicesTotal: number;
  devicesTotalVisitors: number;
  countriesItems: DashboardQuery['dashboard']['countries']['items'];
  countriesTotal: number;
  countriesTotalVisitors: number;
}

export function createEmptyDashboardStats(): DashboardQuery['dashboard'] {
  return {
    __typename: 'DashboardStats',
    visitors: EMPTY_COUNT,
    pageViews: EMPTY_COUNT,
    sessions: EMPTY_COUNT,
    bounceRate: EMPTY_COUNT,
    avgDuration: EMPTY_COUNT,
    topPages: {
      __typename: 'PagedPageStats',
      total: EMPTY_COUNT,
      items: [],
    },
    topReferrers: {
      __typename: 'PagedReferrerStats',
      total: EMPTY_COUNT,
      items: [],
    },
    browsers: [],
    devices: {
      __typename: 'PagedDeviceStats',
      total: EMPTY_COUNT,
      totalVisitors: EMPTY_COUNT,
      items: [],
    },
    countries: {
      __typename: 'PagedCountryStats',
      total: EMPTY_COUNT,
      totalVisitors: EMPTY_COUNT,
      items: [],
    },
  };
}

export function extractStatsData(stats: DashboardQuery['dashboard'] | undefined): StatsDataResult {
  const normalizedStats = stats ?? createEmptyDashboardStats();

  return {
    topPages: normalizedStats.topPages.items,
    topPagesTotal: normalizedStats.topPages.total,
    referrersItems: normalizedStats.topReferrers.items,
    referrersTotal: normalizedStats.topReferrers.total,
    devicesItems: normalizedStats.devices.items,
    devicesTotal: normalizedStats.devices.total,
    devicesTotalVisitors: normalizedStats.devices.totalVisitors,
    countriesItems: normalizedStats.countries.items,
    countriesTotal: normalizedStats.countries.total,
    countriesTotalVisitors: normalizedStats.countries.totalVisitors,
  };
}

const SECONDS_PER_MINUTE = 60;

export function formatDuration(seconds: number): string {
  if (seconds < SECONDS_PER_MINUTE) {
    return `${String(Math.round(seconds))}s`;
  }
  const minutes = Math.floor(seconds / SECONDS_PER_MINUTE);
  const remainingSeconds = Math.round(seconds % SECONDS_PER_MINUTE);
  return `${String(minutes)}m ${String(remainingSeconds)}s`;
}
