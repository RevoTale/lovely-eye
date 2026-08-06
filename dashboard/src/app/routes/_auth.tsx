import { createFileRoute, lazyRouteComponent, redirect } from '@tanstack/react-router';
import { getAuthRedirect } from '@/features/auth/model/auth-routing';

export const Route = createFileRoute('/_auth')({
  beforeLoad: ({ context }) => {
    const destination = getAuthRedirect('protected', context.auth);
    if (destination !== null) throw redirect({ to: destination });
  },
  component: lazyRouteComponent(async () => {
    const module = await import('@/app/layouts/dashboard-layout');
    return { default: module.DashboardLayout };
  }),
});
