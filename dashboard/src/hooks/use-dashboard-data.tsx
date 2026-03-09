import { useQuery } from '@apollo/client/react';
import {
  DashboardDocument,
  EventCountsDocument,
  EventsDocument,
  RealtimeDocument,
  SiteDocument,
} from '@/gql/graphql';
import type {
  DashboardQuery,
  EventCountsQuery,
  EventsQuery,
  FilterInput,
  RealtimeQuery,
  SiteQuery,
} from '@/gql/graphql';
import type { DashboardLoadState } from '@/lib/dashboard-load-state';
import { resolveDashboardLoadState } from '@/lib/dashboard-load-state';
import {
  buildDashboardVariables,
  buildEventCountsVariables,
  buildEventsVariables,
  buildRealtimeVariables,
  DASHBOARD_POLL_INTERVAL_MS,
  PAGE_SIZES,
  REALTIME_POLL_INTERVAL_MS,
} from '@/lib/dashboard-query';

interface UseDashboardDataParams {
  siteId: string;
  dateRange: { from: string; to: string } | null | undefined;
  filter: FilterInput | null;
  eventsPage: number;
  eventsCountsPage: number;
  topPagesPage: number;
  referrersPage: number;
  devicesPage: number;
  osPage: number;
  countriesPage: number;
}

interface DashboardData {
  site: SiteQuery['site'] | undefined;
  siteState: DashboardLoadState;
  stats: DashboardQuery['dashboard'] | undefined;
  dashboardState: DashboardLoadState;
  realtime: RealtimeQuery['realtime'] | undefined;
  eventsResult: EventsQuery['events'] | undefined;
  eventsState: DashboardLoadState;
  eventCounts: EventCountsQuery['eventCounts'];
  eventCountsState: DashboardLoadState;
}

export function useDashboardData(params: UseDashboardDataParams): DashboardData {
  const { siteId, dateRange, filter, eventsPage, eventsCountsPage, topPagesPage, referrersPage, devicesPage, osPage, countriesPage } = params;
  const hasSiteId = siteId !== '';
  const siteQuery = useQuery(SiteDocument, { variables: { id: siteId }, skip: !hasSiteId, fetchPolicy: 'cache-and-network', notifyOnNetworkStatusChange: true });
  const dashboardQuery = useQuery(DashboardDocument, {
    variables: buildDashboardVariables(siteId, dateRange, filter, topPagesPage, referrersPage, devicesPage, osPage, countriesPage),
    skip: !hasSiteId,
    fetchPolicy: 'cache-and-network',
    notifyOnNetworkStatusChange: true,
    pollInterval: DASHBOARD_POLL_INTERVAL_MS,
  });
  const realtimeQuery = useQuery(RealtimeDocument, {
    variables: buildRealtimeVariables(siteId),
    skip: !hasSiteId,
    fetchPolicy: 'cache-and-network',
    pollInterval: REALTIME_POLL_INTERVAL_MS,
  });
  const eventsQuery = useQuery(EventsDocument, {
    variables: buildEventsVariables(siteId, dateRange, filter, eventsPage),
    skip: !hasSiteId,
    fetchPolicy: 'cache-and-network',
    notifyOnNetworkStatusChange: true,
    pollInterval: DASHBOARD_POLL_INTERVAL_MS,
  });
  const eventCountsQuery = useQuery(EventCountsDocument, {
    variables: buildEventCountsVariables(siteId, dateRange, filter, eventsCountsPage),
    skip: !hasSiteId,
    fetchPolicy: 'cache-and-network',
    notifyOnNetworkStatusChange: true,
    pollInterval: DASHBOARD_POLL_INTERVAL_MS,
  });

  const siteState = resolveDashboardLoadState(siteQuery.data?.site, siteQuery.previousData?.site, siteQuery.loading);
  const dashboardState = resolveDashboardLoadState(dashboardQuery.data?.dashboard, dashboardQuery.previousData?.dashboard, dashboardQuery.loading);
  const eventsState = resolveDashboardLoadState(eventsQuery.data?.events, eventsQuery.previousData?.events, eventsQuery.loading);
  const eventCountsState = resolveDashboardLoadState(eventCountsQuery.data?.eventCounts, eventCountsQuery.previousData?.eventCounts, eventCountsQuery.loading);

  return {
    site: siteState.data,
    siteState: siteState.state,
    stats: dashboardState.data,
    dashboardState: dashboardState.state,
    realtime: realtimeQuery.data?.realtime ?? realtimeQuery.previousData?.realtime,
    eventsResult: eventsState.data,
    eventsState: eventsState.state,
    eventCounts: eventCountsState.data ?? [],
    eventCountsState: eventCountsState.state,
  };
}

export { PAGE_SIZES };
