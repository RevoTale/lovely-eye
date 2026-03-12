import { createFileRoute, redirect, lazyRouteComponent } from '@tanstack/react-router';

export const Route = createFileRoute('/login')({
  beforeLoad: async ({ context }) => {
    await Promise.resolve();
    if (context.auth.isLoading || context.auth.bootstrapError !== null) {
      return;
    }
    if (context.auth.isAuthenticated) {
      throw redirect({ to: '/' });
    }
    if (context.auth.authMode === 'register-only') {
      throw redirect({ to: '/register' });
    }
  },
  component: lazyRouteComponent(async () => {
    const module = await import('../pages/login');
    return { default: module.LoginPage };
  }),
});
