import { createRouter, createRootRouteWithContext, createRoute, Outlet, redirect, Link, useNavigate } from '@tanstack/react-router';
import type { AuthContextType } from '@/hooks/use-auth';
import { LoginPage } from './pages/login';
import { RegisterPage } from './pages/register';
import { SitesPage } from './pages/sites';
import { DashboardPage } from './pages/dashboard';
import { DashboardLayout } from './layouts/dashboard-layout';

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
    // Don't redirect while still loading auth state
    if (context.auth.isLoading) {
      return;
    }
    if (context.auth.isAuthenticated) {
      throw redirect({ to: '/' });
    }
  },
  component: LoginPage,
});

const registerRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/register',
  beforeLoad: ({ context }) => {
    // Don't redirect while still loading auth state
    if (context.auth.isLoading) {
      return;
    }
    if (context.auth.isAuthenticated) {
      throw redirect({ to: '/' });
    }
  },
  component: RegisterPage,
});

// Auth layout route (protected)
const authLayoutRoute = createRoute({
  getParentRoute: () => rootRoute,
  id: 'auth',
  beforeLoad: ({ context }) => {
    // Don't redirect while still loading auth state
    if (context.auth.isLoading) {
      return;
    }
    if (!context.auth.isAuthenticated) {
      throw redirect({ to: '/login' });
    }
  },
  component: DashboardLayout,
});

// Protected child routes
const sitesRoute = createRoute({
  getParentRoute: () => authLayoutRoute,
  path: '/',
  component: SitesPage,
});

const siteDetailRoute = createRoute({
  getParentRoute: () => authLayoutRoute,
  path: '/sites/$siteId',
  component: DashboardPage,
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
    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion -- Set by RouterProvider
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
