import { createRouter, createRootRouteWithContext, createRoute, Outlet, redirect, Link, useNavigate, lazyRouteComponent } from '@tanstack/react-router';
import { z } from 'zod';
import type { AuthContextType } from '@/hooks/use-auth';

interface RouterContext {
  auth: AuthContextType;
}

// Helper to create initial router context (overridden at runtime via RouterProvider)
function createInitialContext(): RouterContext {
  // eslint-disable-next-line @typescript-eslint/no-unsafe-type-assertion, @typescript-eslint/consistent-type-assertions -- TanStack Router requires initial context placeholder
  return {} as RouterContext;
}

const rootRoute = createRootRouteWithContext<RouterContext>()({
  component: () => <Outlet />,
});

const loginRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/login',
  beforeLoad: async ({ context }) => {
    await Promise.resolve();
    if (context.auth.isLoading) {
      return;
    }
    if (context.auth.isAuthenticated) {
      // eslint-disable-next-line @typescript-eslint/only-throw-error -- TanStack Router expects redirect throws.
      throw redirect({ to: '/' });
    }
  },
  component: lazyRouteComponent(async () => {
    const module = await import('./pages/login');
    return { default: module.LoginPage };
  }),
});

const registerRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/register',
  beforeLoad: async ({ context }) => {
    await Promise.resolve();
    if (context.auth.isLoading) {
      return;
    }
    if (context.auth.isAuthenticated) {
      // eslint-disable-next-line @typescript-eslint/only-throw-error -- TanStack Router expects redirect throws.
      throw redirect({ to: '/' });
    }
  },
  component: lazyRouteComponent(async () => {
    const module = await import('./pages/register');
    return { default: module.RegisterPage };
  }),
});

const authLayoutRoute = createRoute({
  getParentRoute: () => rootRoute,
  id: 'auth',
  beforeLoad: async ({ context }) => {
    await Promise.resolve();
    if (context.auth.isLoading) {
      return;
    }
    if (!context.auth.isAuthenticated) {
      // eslint-disable-next-line @typescript-eslint/only-throw-error -- TanStack Router expects redirect throws.
      throw redirect({ to: '/login' });
    }
  },
  component: lazyRouteComponent(async () => {
    const module = await import('./layouts/dashboard-layout');
    return { default: module.DashboardLayout };
  }),
});

const sitesRoute = createRoute({
  getParentRoute: () => authLayoutRoute,
  path: '/',
  component: lazyRouteComponent(async () => {
    const module = await import('./pages/sites');
    return { default: module.SitesPage };
  }),
});

const siteDetailRoute = createRoute({
  getParentRoute: () => authLayoutRoute,
  path: '/sites/$siteId',
  validateSearch: (search: Record<string, unknown>) => {
    const filterValue = z.union([z.string(), z.array(z.string())]).optional();
    const searchSchema = z.object({
      view: z.string().optional(),
      preset: z.enum(['7d', '30d', '90d', 'custom', 'all']).optional(),
      from: z.string().optional(),
      to: z.string().optional(),
      fromTime: z.string().optional(),
      toTime: z.string().optional(),
      statsBucket: z.enum(['daily', 'hourly']).optional(),
      eventsPage: z.string().optional(),
      topPagesPage: z.string().optional(),
      referrersPage: z.string().optional(),
      devicesPage: z.string().optional(),
      countriesPage: z.string().optional(),
      referrer: filterValue,
      device: filterValue,
      page: filterValue,
      country: filterValue,
    });

    const parsed = searchSchema.safeParse(search);
    return parsed.success ? parsed.data : {};
  },
  component: lazyRouteComponent(async () => {
    const module = await import('./pages/site-view');
    return { default: module.SiteViewPage };
  }),
});

const routeTree = rootRoute.addChildren([
  loginRoute,
  registerRoute,
  authLayoutRoute.addChildren([sitesRoute, siteDetailRoute]),
]);

export const router = createRouter({
  routeTree,
  context: createInitialContext(),
  defaultPreload: 'intent',
  basepath: window.__ENV__?.BASE_PATH ?? '/',
});

declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router;
  }
}

export { Link, useNavigate };

export { siteDetailRoute };
