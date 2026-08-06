import { createFileRoute, lazyRouteComponent } from '@tanstack/react-router';

export const Route = createFileRoute('/_auth/')({
  component: lazyRouteComponent(async () => {
    const { SiteEntryScreen } = await import('@/features/sites/screens/site-entry-screen');
    return { default: SiteEntryScreen };
  }),
});
