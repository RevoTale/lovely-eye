import { createFileRoute, redirect, lazyRouteComponent } from '@tanstack/react-router';

export const Route = createFileRoute('/register')({
  beforeLoad: async ({ context }) => {
    await Promise.resolve();
    if (context.auth.isLoading) {
      return;
    }
    if (context.auth.isAuthenticated) {
      throw redirect({ to: '/' });
    }
  },
  component: lazyRouteComponent(async () => {
    const module = await import('../pages/register');
    return { default: module.RegisterPage };
  }),
});
