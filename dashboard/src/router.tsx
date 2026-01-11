import { createRouter, createRootRouteWithContext, createRoute, Outlet, redirect, Link, useNavigate, lazyRouteComponent } from '@tanstack/react-router';
import type { AuthContextType } from '@/hooks/use-auth';

// Router context type
interface RouterContext {
  auth: AuthContextType;
}

// Root route
const rootRoute = createRootRouteWithContext<RouterContext>()({
  component: () => <Outlet />,
});

// Public routes
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

// Auth layout route (protected)
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

// Protected child routes
const sitesRoute = createRoute({
  getParentRoute: () => authLayoutRoute,
  path: '/',
  component: lazyRouteComponent(() => import('./pages/sites').then(m => ({ default: m.SitesPage }))),
});

const siteDetailRoute = createRoute({
  getParentRoute: () => authLayoutRoute,
  path: '/sites/$siteId',
  component: lazyRouteComponent(() => import('./pages/site-view').then(m => ({ default: m.SiteViewPage }))),
});

// Route tree
const routeTree = rootRoute.addChildren([
  loginRoute,
  registerRoute,
  authLayoutRoute.addChildren([sitesRoute, siteDetailRoute]),
]);

// Create router
export const router = createRouter({
  routeTree,
  context: {
    auth: undefined!,
  },
  defaultPreload: 'intent',
  basepath: window.__ENV__?.BASE_PATH ?? '/',
});

// Type registration
declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router;
  }
}

// Re-export Link and useNavigate for convenience
export { Link, useNavigate };

// Export route for params access
export { siteDetailRoute };
