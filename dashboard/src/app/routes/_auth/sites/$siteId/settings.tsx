import { createFileRoute, lazyRouteComponent } from '@tanstack/react-router';

export const Route = createFileRoute('/_auth/sites/$siteId/settings')({
  component: lazyRouteComponent(async () => {
    const module = await import('@/features/site-settings/screens/site-settings-screen');
    return { default: module.SiteSettingsScreen };
  }),
});
