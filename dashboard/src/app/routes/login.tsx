import { createFileRoute, lazyRouteComponent, redirect } from '@tanstack/react-router';
import { getAuthRedirect } from '@/features/auth/model/auth-routing';

export const Route = createFileRoute('/login')({
  beforeLoad: ({ context }) => {
    const destination = getAuthRedirect('login', context.auth);
    if (destination !== null) throw redirect({ to: destination });
  },
  component: lazyRouteComponent(async () => {
    const module = await import('@/features/auth/screens/login-screen');
    return { default: module.LoginScreen };
  }),
});
