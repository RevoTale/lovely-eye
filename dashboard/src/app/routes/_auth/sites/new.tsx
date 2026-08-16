import { createFileRoute, lazyRouteComponent } from '@tanstack/react-router';

export const Route = createFileRoute('/_auth/sites/new')({
  component: lazyRouteComponent(async () => {
    const { CreateSiteScreen } = await import('@/features/sites/screens/create-site-screen');
    return { default: CreateSiteScreen };
  }),
});
