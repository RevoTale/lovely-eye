import { createFileRoute, lazyRouteComponent } from '@tanstack/react-router';
import { GeoIpStatusDocument, SiteDocument } from '@/shared/api/generated/graphql';

export const Route = createFileRoute('/_auth/sites/$siteId/settings')({
  loader: async ({ context, params }) => {
    if (context.apolloClient === null) throw new Error('Apollo client is unavailable.');
    await Promise.all([
      context.apolloClient.query({
        query: SiteDocument,
        variables: { id: params.siteId },
        fetchPolicy: 'cache-first',
        errorPolicy: 'all',
      }),
      context.apolloClient.query({ query: GeoIpStatusDocument, fetchPolicy: 'cache-first' }),
    ]);
  },
  component: lazyRouteComponent(async () => {
    const module = await import('@/features/site-settings/screens/site-settings-screen');
    return { default: module.SiteSettingsScreen };
  }),
});
