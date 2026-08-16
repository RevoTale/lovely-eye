import { createFileRoute, lazyRouteComponent, redirect } from '@tanstack/react-router';
import { getAuthRedirect } from '@/features/auth/model/auth-routing';

export const Route = createFileRoute('/register')({
  beforeLoad: ({ context }) => {
    const destination = getAuthRedirect('register', context.auth);
    if (destination !== null) throw redirect({ to: destination });
  },
  component: lazyRouteComponent(async () => {
    const module = await import('@/features/auth/screens/register-screen');
    return { default: module.RegisterScreen };
  }),
});
