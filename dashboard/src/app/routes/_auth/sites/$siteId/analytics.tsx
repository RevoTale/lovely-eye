import { createFileRoute, lazyRouteComponent } from '@tanstack/react-router';
import { resolveAnalyticsDateRange } from '@/features/analytics/model/analytics-date-range';
import { analyticsSearchSchema } from '@/features/analytics/model/analytics-search';
import { buildDashboardVariables } from '@/features/analytics/model/dashboard-query';
import { buildFilterInput, buildFilters } from '@/features/analytics/model/dashboard-utils';
import { DashboardDocument, SiteDocument } from '@/shared/api/generated/graphql';

export const Route = createFileRoute('/_auth/sites/$siteId/analytics')({
  validateSearch: analyticsSearchSchema,
  staleTime: Number.POSITIVE_INFINITY,
  loader: async ({ context, location, params }) => {
    if (context.apolloClient === null) throw new Error('Apollo client is unavailable.');
    const search = analyticsSearchSchema.parse(location.search);
    const dateRange = resolveAnalyticsDateRange(search, new Date());
    const filter = buildFilterInput(buildFilters(search).filter);
    const [siteResult, dashboardResult] = await Promise.allSettled([
      context.apolloClient.query({
        query: SiteDocument,
        variables: { id: params.siteId },
        fetchPolicy: 'cache-first',
        errorPolicy: 'all',
      }),
      context.apolloClient.query({
        query: DashboardDocument,
        variables: buildDashboardVariables(
          params.siteId,
          dateRange,
          filter,
          search.topPagesPage ?? 1,
          search.referrersPage ?? 1,
          search.devicesPage ?? 1,
          search.osPage ?? 1,
          search.countriesPage ?? 1
        ),
        fetchPolicy: 'cache-first',
      }),
    ]);
    if (siteResult.status === 'rejected') throw siteResult.reason;
    if (siteResult.value.data?.site === null || siteResult.value.data?.site === undefined) return;
    if (dashboardResult.status === 'rejected') throw dashboardResult.reason;
  },
  component: lazyRouteComponent(async () => {
    const module = await import('@/features/analytics/screens/analytics-screen');
    return { default: module.AnalyticsScreen };
  }),
});
