import type {
  AnalyticsFilterKey,
  AnalyticsSearch,
} from '@/features/analytics/model/analytics-search';
import type {
  BrowserStatsFieldsFragment,
  CountryStatsFieldsFragment,
  DashboardQuery,
  DashboardStatsFieldsFragment,
  DeviceStatsFieldsFragment,
  FilterInput,
  OperatingSystemStatsFieldsFragment,
  PageStatsFieldsFragment,
  ReferrerStatsFieldsFragment,
} from '@/shared/api/generated/graphql';
import {
  BrowserStatsFieldsFragmentDoc,
  CountryStatsFieldsFragmentDoc,
  DashboardStatsFieldsFragmentDoc,
  DeviceStatsFieldsFragmentDoc,
  OperatingSystemStatsFieldsFragmentDoc,
  PageStatsFieldsFragmentDoc,
  ReferrerStatsFieldsFragmentDoc,
} from '@/shared/api/generated/graphql';
import { readFragment } from '@/shared/api/read-fragment';

const EMPTY_COUNT = 0;

interface FilterResult {
  referrers: string[];
  browsers: string[];
  devices: string[];
  operatingSystems: string[];
  pages: string[];
  countries: string[];
  eventNames: string[];
  eventPaths: string[];
  pagePathContains: string | undefined;
  decodedSearch: AnalyticsSearch;
  filter: Record<string, string[]>;
}

export function buildFilterInput(
  filter: Record<string, string[]>,
  pagePathContains?: string
): FilterInput | null {
  if (Object.keys(filter).length === EMPTY_COUNT && pagePathContains === undefined) return null;
  const getFilter = (key: string): string[] | null => filter[key] ?? null;
  return {
    referrer: getFilter('referrer'),
    browser: getFilter('browser'),
    device: getFilter('device'),
    os: getFilter('os'),
    page: getFilter('page'),
    pagePathContains: pagePathContains ?? null,
    country: getFilter('country'),
    eventType: null,
    eventDefinitionId: getFilter('eventDefinitionId'),
    eventName: getFilter('eventName'),
    eventPath: getFilter('eventPath'),
  };
}

export function buildFilters(search: AnalyticsSearch): FilterResult {
  const getFilter = (key: AnalyticsFilterKey): string[] => search[key] ?? [];
  const referrers = getFilter('referrer');
  const browsers = getFilter('browser');
  const devices = getFilter('device');
  const operatingSystems = getFilter('os');
  const pages = getFilter('page');
  const countries = getFilter('country');
  const eventNames = getFilter('eventName');
  const eventPaths = getFilter('eventPath');
  const pagePathContains = search.pagePathContains;

  const decodedSearch = {
    ...search,
    ...(referrers.length > EMPTY_COUNT ? { referrer: referrers } : {}),
    ...(browsers.length > EMPTY_COUNT ? { browser: browsers } : {}),
    ...(devices.length > EMPTY_COUNT ? { device: devices } : {}),
    ...(operatingSystems.length > EMPTY_COUNT ? { os: operatingSystems } : {}),
    ...(pages.length > EMPTY_COUNT ? { page: pages } : {}),
    ...(countries.length > EMPTY_COUNT ? { country: countries } : {}),
    ...(eventNames.length > EMPTY_COUNT ? { eventName: eventNames } : {}),
    ...(eventPaths.length > EMPTY_COUNT ? { eventPath: eventPaths } : {}),
  };

  const filter = {
    ...(referrers.length > EMPTY_COUNT ? { referrer: referrers } : {}),
    ...(browsers.length > EMPTY_COUNT ? { browser: browsers } : {}),
    ...(devices.length > EMPTY_COUNT ? { device: devices } : {}),
    ...(operatingSystems.length > EMPTY_COUNT ? { os: operatingSystems } : {}),
    ...(pages.length > EMPTY_COUNT ? { page: pages } : {}),
    ...(countries.length > EMPTY_COUNT ? { country: countries } : {}),
    ...(eventNames.length > EMPTY_COUNT ? { eventName: eventNames } : {}),
    ...(eventPaths.length > EMPTY_COUNT ? { eventPath: eventPaths } : {}),
  };

  return {
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
  };
}

interface StatsDataResult {
  browsersItems: BrowserStatsFieldsFragment[];
  topPages: PageStatsFieldsFragment[];
  topPagesTotal: number;
  referrersItems: ReferrerStatsFieldsFragment[];
  referrersTotal: number;
  devicesItems: DeviceStatsFieldsFragment[];
  devicesTotal: number;
  devicesTotalVisitors: number;
  operatingSystemsItems: OperatingSystemStatsFieldsFragment[];
  operatingSystemsTotal: number;
  operatingSystemsTotalVisitors: number;
  countriesItems: CountryStatsFieldsFragment[];
  countriesTotal: number;
  countriesTotalVisitors: number;
}

export function createEmptyDashboardStats(): DashboardStatsFieldsFragment {
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
    operatingSystems: {
      __typename: 'PagedOperatingSystemStats',
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
  const normalizedStats =
    stats === undefined
      ? createEmptyDashboardStats()
      : readFragment(DashboardStatsFieldsFragmentDoc, stats);
  const topPages = readFragment(PageStatsFieldsFragmentDoc, normalizedStats.topPages.items);
  const referrersItems = readFragment(
    ReferrerStatsFieldsFragmentDoc,
    normalizedStats.topReferrers.items
  );
  const browsersItems = readFragment(BrowserStatsFieldsFragmentDoc, normalizedStats.browsers);
  const devicesItems = readFragment(DeviceStatsFieldsFragmentDoc, normalizedStats.devices.items);
  const operatingSystemsItems = readFragment(
    OperatingSystemStatsFieldsFragmentDoc,
    normalizedStats.operatingSystems.items
  );
  const countriesItems = readFragment(
    CountryStatsFieldsFragmentDoc,
    normalizedStats.countries.items
  );

  return {
    browsersItems,
    topPages,
    topPagesTotal: normalizedStats.topPages.total,
    referrersItems,
    referrersTotal: normalizedStats.topReferrers.total,
    devicesItems,
    devicesTotal: normalizedStats.devices.total,
    devicesTotalVisitors: normalizedStats.devices.totalVisitors,
    operatingSystemsItems,
    operatingSystemsTotal: normalizedStats.operatingSystems.total,
    operatingSystemsTotalVisitors: normalizedStats.operatingSystems.totalVisitors,
    countriesItems,
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
