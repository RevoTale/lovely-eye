import { createRouter, createRootRouteWithContext, createRoute, Outlet, redirect, Link, useNavigate, lazyRouteComponent } from '@tanstack/react-router';
import { z } from 'zod';
import type { AuthContextType } from '@/hooks/use-auth';

interface RouterContext {
  auth: AuthContextType;
}

const rootRoute = createRootRouteWithContext<RouterContext>()({
  component: () => <Outlet />,
});

const loginRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/login',
  beforeLoad: ({ context }) => {
    if (context.auth.isLoading) {
      return;
    }
    if (context.auth.isAuthenticated) {
      throw redirect({ to: '/' });
    }
  },
  component: lazyRouteComponent(() => import('./pages/login').then(m => ({ default: m.LoginPage }))),
});

const registerRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/register',
  beforeLoad: ({ context }) => {
    if (context.auth.isLoading) {
      return;
    }
    if (context.auth.isAuthenticated) {
      throw redirect({ to: '/' });
    }
  },
  component: lazyRouteComponent(() => import('./pages/register').then(m => ({ default: m.RegisterPage }))),
});

const authLayoutRoute = createRoute({
  getParentRoute: () => rootRoute,
  id: 'auth',
  beforeLoad: ({ context }) => {
    if (context.auth.isLoading) {
      return;
    }
    if (!context.auth.isAuthenticated) {
      throw redirect({ to: '/login' });
    }
  },
  component: lazyRouteComponent(() => import('./layouts/dashboard-layout').then(m => ({ default: m.DashboardLayout }))),
});

const sitesRoute = createRoute({
  getParentRoute: () => authLayoutRoute,
  path: '/',
  component: lazyRouteComponent(() => import('./pages/sites').then(m => ({ default: m.SitesPage }))),
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
      referrer: filterValue,
      device: filterValue,
      page: filterValue,
      country: filterValue,
    });

    const parsed = searchSchema.safeParse(search);
    return parsed.success ? parsed.data : {};
  },
  component: lazyRouteComponent(() => import('./pages/site-view').then(m => ({ default: m.SiteViewPage }))),
});

const routeTree = rootRoute.addChildren([
  loginRoute,
  registerRoute,
  authLayoutRoute.addChildren([sitesRoute, siteDetailRoute]),
]);

export const router = createRouter({
  routeTree,
  context: {
    auth: undefined as unknown as AuthContextType,
  },
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
