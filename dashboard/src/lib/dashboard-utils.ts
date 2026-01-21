import { normalizeFilterValue } from '@/lib/filter-utils';

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
  topPages: never[];
  topPagesTotal: number;
  referrersItems: never[];
  referrersTotal: number;
  devicesItems: never[];
  devicesTotal: number;
  devicesTotalVisitors: number;
  countriesItems: never[];
  countriesTotal: number;
  countriesTotalVisitors: number;
}

export function extractStatsData(stats: unknown): StatsDataResult {
  // eslint-disable-next-line @typescript-eslint/no-unsafe-type-assertion -- Stats data structure from GraphQL query, types narrowed for component usage
  const statsData: Record<string, { items?: unknown[]; total?: number; totalVisitors?: number }> | null | undefined = stats as Record<string, { items?: unknown[]; total?: number; totalVisitors?: number }> | null | undefined;
  const topPagesResult = statsData?.['topPages'];
  const referrersResult = statsData?.['topReferrers'];
  const devicesResult = statsData?.['devices'];
  const countriesResult = statsData?.['countries'];

  return {
    // eslint-disable-next-line @typescript-eslint/no-unsafe-type-assertion -- GraphQL types are narrowed for component usage
    topPages: (topPagesResult?.items ?? []) as never[],
    topPagesTotal: topPagesResult?.total ?? EMPTY_COUNT,
    // eslint-disable-next-line @typescript-eslint/no-unsafe-type-assertion -- GraphQL types are narrowed for component usage
    referrersItems: (referrersResult?.items ?? []) as never[],
    referrersTotal: referrersResult?.total ?? EMPTY_COUNT,
    // eslint-disable-next-line @typescript-eslint/no-unsafe-type-assertion -- GraphQL types are narrowed for component usage
    devicesItems: (devicesResult?.items ?? []) as never[],
    devicesTotal: devicesResult?.total ?? EMPTY_COUNT,
    devicesTotalVisitors: devicesResult?.totalVisitors ?? EMPTY_COUNT,
    // eslint-disable-next-line @typescript-eslint/no-unsafe-type-assertion -- GraphQL types are narrowed for component usage
    countriesItems: (countriesResult?.items ?? []) as never[],
    countriesTotal: countriesResult?.total ?? EMPTY_COUNT,
    countriesTotalVisitors: countriesResult?.totalVisitors ?? EMPTY_COUNT,
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
