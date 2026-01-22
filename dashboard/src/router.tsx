import { createRouter, Link, useNavigate } from '@tanstack/react-router';
import { routeTree } from './routeTree.gen';
import { createInitialContext } from './router-context';
import { Route as SiteDetailRoute } from './routes/_auth/sites/$siteId';

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

export const siteDetailRoute = SiteDetailRoute;
