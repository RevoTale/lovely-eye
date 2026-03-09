import type {
  DashboardQueryVariables,
  EventCountsQueryVariables,
  EventsQueryVariables,
  FilterInput,
  RealtimeQueryVariables,
} from '@/gql/graphql';
import { EventType } from '@/gql/graphql';

const EMPTY_COUNT = 0;
const PAGE_INDEX_OFFSET = 1;
const ZERO_OFFSET = 0;

const EVENTS_PAGE_SIZE = 5;
const EVENTS_COUNT_PAGE_SIZE = 20;
const TOP_PAGES_PAGE_SIZE = 5;
const REFERRERS_PAGE_SIZE = 5;
const BROWSERS_PAGE_SIZE = 10;
const DEVICES_PAGE_SIZE = 6;
const OS_PAGE_SIZE = 6;
const COUNTRIES_PAGE_SIZE = 6;
export const ACTIVE_PAGES_PAGE_SIZE = 10;

export const PAGE_SIZES = {
  EVENTS: EVENTS_PAGE_SIZE,
  EVENT_COUNTS: EVENTS_COUNT_PAGE_SIZE,
  TOP_PAGES: TOP_PAGES_PAGE_SIZE,
  REFERRERS: REFERRERS_PAGE_SIZE,
  BROWSERS: BROWSERS_PAGE_SIZE,
  DEVICES: DEVICES_PAGE_SIZE,
  OS: OS_PAGE_SIZE,
  COUNTRIES: COUNTRIES_PAGE_SIZE,
} as const;

export const DASHBOARD_POLL_INTERVAL_MS = 60000;
export const REALTIME_POLL_INTERVAL_MS = 5000;

const normalizeFilter = (filter: FilterInput | null): FilterInput | null =>
  filter === null || Object.keys(filter).length === EMPTY_COUNT ? null : filter;

const pageOffset = (page: number, pageSize: number): number => (page - PAGE_INDEX_OFFSET) * pageSize;

export const buildDashboardVariables = (
  siteId: string,
  dateRange: { from: string; to: string } | null | undefined,
  filter: FilterInput | null,
  topPagesPage: number,
  referrersPage: number,
  devicesPage: number,
  osPage: number,
  countriesPage: number
): DashboardQueryVariables => ({
  siteId,
  dateRange: dateRange ?? null,
  filter: normalizeFilter(filter),
  topPagesPaging: { limit: TOP_PAGES_PAGE_SIZE, offset: pageOffset(topPagesPage, TOP_PAGES_PAGE_SIZE) },
  referrersPaging: { limit: REFERRERS_PAGE_SIZE, offset: pageOffset(referrersPage, REFERRERS_PAGE_SIZE) },
  browsersPaging: { limit: BROWSERS_PAGE_SIZE, offset: ZERO_OFFSET },
  devicesPaging: { limit: DEVICES_PAGE_SIZE, offset: pageOffset(devicesPage, DEVICES_PAGE_SIZE) },
  osPaging: { limit: OS_PAGE_SIZE, offset: pageOffset(osPage, OS_PAGE_SIZE) },
  countriesPaging: { limit: COUNTRIES_PAGE_SIZE, offset: pageOffset(countriesPage, COUNTRIES_PAGE_SIZE) },
});

export const buildEventsVariables = (
  siteId: string,
  dateRange: { from: string; to: string } | null | undefined,
  filter: FilterInput | null,
  eventsPage: number
): EventsQueryVariables => ({
  siteId,
  dateRange: dateRange ?? null,
  filter: normalizeFilter(filter),
  limit: EVENTS_PAGE_SIZE,
  offset: pageOffset(eventsPage, EVENTS_PAGE_SIZE),
});

export const buildEventCountsVariables = (
  siteId: string,
  dateRange: { from: string; to: string } | null | undefined,
  filter: FilterInput | null,
  eventsCountsPage: number
): EventCountsQueryVariables => ({
  siteId,
  dateRange: dateRange ?? null,
  filter: {
    referrer: filter?.referrer ?? null,
    browser: filter?.browser ?? null,
    device: filter?.device ?? null,
    page: filter?.page ?? null,
    country: filter?.country ?? null,
    os: filter?.os ?? null,
    eventType: [EventType.Predefined],
    eventDefinitionId: filter?.eventDefinitionId ?? null,
    eventName: filter?.eventName ?? null,
    eventPath: filter?.eventPath ?? null,
  },
  paging: { limit: EVENTS_COUNT_PAGE_SIZE, offset: pageOffset(eventsCountsPage, EVENTS_COUNT_PAGE_SIZE) },
});

export const buildRealtimeVariables = (siteId: string): RealtimeQueryVariables => ({
  siteId,
  activePagesPaging: {
    limit: ACTIVE_PAGES_PAGE_SIZE,
    offset: ZERO_OFFSET,
  },
});
