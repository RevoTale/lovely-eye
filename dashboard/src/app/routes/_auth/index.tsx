import { createFileRoute, redirect } from '@tanstack/react-router';
import { getRememberedSiteId } from '@/features/sites/model/recent-site';
import { selectInitialSiteId } from '@/features/sites/model/site-selection';
import { loadSites } from '@/features/sites/model/use-sites';

export const Route = createFileRoute('/_auth/')({
  loader: async ({ context }) => {
    if (context.apolloClient === null) throw new Error('Apollo client is unavailable.');
    const sites = await loadSites(context.apolloClient);
    if (sites.length === 0) throw redirect({ to: '/sites/new', replace: true });
    const siteId = selectInitialSiteId(sites, getRememberedSiteId());
    if (siteId === null) throw redirect({ to: '/sites', replace: true });
    throw redirect({ to: '/sites/$siteId/analytics', params: { siteId }, replace: true });
  },
});
