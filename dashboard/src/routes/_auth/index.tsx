import { createFileRoute, lazyRouteComponent } from '@tanstack/react-router';

export const Route = createFileRoute('/_auth/')({
  component: lazyRouteComponent(async () => {
    const { SitesPage } = await import('../../pages/sites');
    return { default: SitesPage };
  }),
});
