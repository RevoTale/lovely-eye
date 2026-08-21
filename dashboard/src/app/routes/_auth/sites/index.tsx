import { createFileRoute, lazyRouteComponent } from '@tanstack/react-router';
import { loadSites } from '@/features/sites/model/use-sites';

export const Route = createFileRoute('/_auth/sites/')({
  loader: async ({ context }) => {
    if (context.apolloClient === null) throw new Error('Apollo client is unavailable.');
    await loadSites(context.apolloClient);
  },
  component: lazyRouteComponent(async () => {
    const { SitesScreen } = await import('@/features/sites/screens/sites-screen');
    return { default: SitesScreen };
  }),
});
