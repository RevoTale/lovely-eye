import { createFileRoute, lazyRouteComponent } from '@tanstack/react-router';
import { analyticsSearchSchema } from '@/features/analytics/model/analytics-search';

export const Route = createFileRoute('/_auth/sites/$siteId/analytics')({
  validateSearch: analyticsSearchSchema,
  component: lazyRouteComponent(async () => {
    const module = await import('@/features/analytics/screens/analytics-screen');
    return { default: module.AnalyticsScreen };
  }),
});
