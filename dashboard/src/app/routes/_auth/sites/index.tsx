import { createFileRoute, lazyRouteComponent } from '@tanstack/react-router';

export const Route = createFileRoute('/_auth/sites/')({
  component: lazyRouteComponent(async () => {
    const { SitesScreen } = await import('@/features/sites/screens/sites-screen');
    return { default: SitesScreen };
  }),
});
