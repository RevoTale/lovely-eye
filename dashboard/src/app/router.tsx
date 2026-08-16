import { createRouter, Link, useNavigate } from '@tanstack/react-router';
import { routeTree } from '@/app/route-tree.gen';
import { createInitialContext } from '@/app/router-context';
import { getBasePath } from '@/shared/config/runtime';

export const router = createRouter({
  routeTree,
  context: createInitialContext(),
  defaultPreload: 'intent',
  basepath: getBasePath(),
});

declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router;
  }
}

export { Link, useNavigate };
