import { createFileRoute, redirect, lazyRouteComponent } from '@tanstack/react-router';

export const Route = createFileRoute('/_auth')({
  beforeLoad: async ({ context }) => {
    await Promise.resolve();
    if (context.auth.isLoading) {
      return;
    }
    if (!context.auth.isAuthenticated) {
      throw redirect({ to: '/login' });
    }
  },
  component: lazyRouteComponent(async () => {
    const module = await import('../layouts/dashboard-layout');
    return { default: module.DashboardLayout };
  }),
});
