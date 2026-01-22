import { useQuery } from '@apollo/client/react';
import {
  DashboardDocument,
  RealtimeDocument,
  EventsDocument,
  EventCountsDocument,
  SiteDocument,
} from '@/gql/graphql';
import type {
  DashboardQuery,
  EventsQuery,
  FilterInput,
  EventCountsQuery,
  RealtimeQuery,
  SiteQuery,
} from '@/gql/graphql';

const EVENTS_PAGE_SIZE = 5;
const EVENTS_COUNT_LIMIT = 200;
const TOP_PAGES_PAGE_SIZE = 5;
const REFERRERS_PAGE_SIZE = 5;
const BROWSERS_PAGE_SIZE = 10;
const DEVICES_PAGE_SIZE = 6;
const COUNTRIES_PAGE_SIZE = 6;
const ACTIVE_PAGES_PAGE_SIZE = 10;
const EMPTY_COUNT = 0;
const PAGE_INDEX_OFFSET = 1;
const ZERO_OFFSET = 0;
const DASHBOARD_POLL_INTERVAL_MS = 60000;
const REALTIME_POLL_INTERVAL_MS = 5000;

interface UseDashboardDataParams {
  siteId: string;
  dateRange: { from: string; to: string } | null | undefined;
  filter: FilterInput | null;
  eventsPage: number;
  topPagesPage: number;
  referrersPage: number;
  devicesPage: number;
  countriesPage: number;
}

interface DashboardData {
  site: SiteQuery['site'] | undefined;
  stats: DashboardQuery['dashboard'] | undefined;
  realtime: RealtimeQuery['realtime'] | undefined;
  eventsResult: EventsQuery['events'] | undefined;
  eventsCounts: EventCountsQuery['eventCounts'];
  siteLoading: boolean;
  dashboardLoading: boolean;
  eventsLoading: boolean;
}

export function useDashboardData(params: UseDashboardDataParams): DashboardData {
  const { siteId, dateRange, filter, eventsPage, topPagesPage, referrersPage, devicesPage, countriesPage } = params;
  const hasSiteId = siteId !== '';

  const { data: siteData, loading: siteLoading } = useQuery(SiteDocument, {
    variables: { id: siteId },
    skip: !hasSiteId,
  });

  const { data: dashboardData, loading: dashboardLoading } = useQuery(DashboardDocument, {
    variables: {
      siteId,
      dateRange: dateRange ?? null,
      filter: filter === null || Object.keys(filter).length === EMPTY_COUNT ? null : filter,
      topPagesPaging: {
        limit: TOP_PAGES_PAGE_SIZE,
        offset: (topPagesPage - PAGE_INDEX_OFFSET) * TOP_PAGES_PAGE_SIZE,
      },
      referrersPaging: {
        limit: REFERRERS_PAGE_SIZE,
        offset: (referrersPage - PAGE_INDEX_OFFSET) * REFERRERS_PAGE_SIZE,
      },
      browsersPaging: {
        limit: BROWSERS_PAGE_SIZE,
        offset: ZERO_OFFSET,
      },
      devicesPaging: {
        limit: DEVICES_PAGE_SIZE,
        offset: (devicesPage - PAGE_INDEX_OFFSET) * DEVICES_PAGE_SIZE,
      },
      countriesPaging: {
        limit: COUNTRIES_PAGE_SIZE,
        offset: (countriesPage - PAGE_INDEX_OFFSET) * COUNTRIES_PAGE_SIZE,
      },
    },
    skip: !hasSiteId,
    pollInterval: DASHBOARD_POLL_INTERVAL_MS,
  });

  const { data: realtimeData } = useQuery(RealtimeDocument, {
    variables: {
      siteId,
      activePagesPaging: {
        limit: ACTIVE_PAGES_PAGE_SIZE,
        offset: ZERO_OFFSET,
      },
    },
    skip: !hasSiteId,
    pollInterval: REALTIME_POLL_INTERVAL_MS,
  });

  const { data: eventsData, loading: eventsLoading } = useQuery(EventsDocument, {
    variables: {
      siteId,
      dateRange: dateRange ?? null,
      filter: filter === null || Object.keys(filter).length === EMPTY_COUNT ? null : filter,
      limit: EVENTS_PAGE_SIZE,
      offset: (eventsPage - PAGE_INDEX_OFFSET) * EVENTS_PAGE_SIZE,
    },
    fetchPolicy:'cache-and-network',
    skip: !hasSiteId,
    pollInterval: DASHBOARD_POLL_INTERVAL_MS,
  });

  const { data: eventsCountsData } = useQuery(EventCountsDocument, {
    variables: {
      siteId,
      dateRange: dateRange ?? null,
      filter: filter === null || Object.keys(filter).length === EMPTY_COUNT ? null : filter,
      paging: {
        limit: EVENTS_COUNT_LIMIT,
        offset: ZERO_OFFSET,
      },
    },
    skip: !hasSiteId,
    pollInterval: DASHBOARD_POLL_INTERVAL_MS,
  });

  return {
    site: siteData?.site,
    stats: dashboardData?.dashboard,
    realtime: realtimeData?.realtime,
    eventsResult: eventsData?.events,
    eventsCounts: eventsCountsData?.eventCounts ?? [],
    siteLoading,
    dashboardLoading,
    eventsLoading,
  };
}

export const PAGE_SIZES = {
  EVENTS: EVENTS_PAGE_SIZE,
  TOP_PAGES: TOP_PAGES_PAGE_SIZE,
  REFERRERS: REFERRERS_PAGE_SIZE,
  BROWSERS: BROWSERS_PAGE_SIZE,
  DEVICES: DEVICES_PAGE_SIZE,
  COUNTRIES: COUNTRIES_PAGE_SIZE,
} as const;
