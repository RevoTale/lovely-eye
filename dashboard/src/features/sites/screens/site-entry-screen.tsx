import { useNavigate } from '@tanstack/react-router';
import { useEffect } from 'react';
import { getRememberedSiteId } from '@/features/sites/model/recent-site';
import { selectInitialSiteId } from '@/features/sites/model/site-selection';
import { useSites } from '@/features/sites/model/use-sites';
import { Card, CardContent, CardHeader } from '@/shared/ui/card';
import { Skeleton } from '@/shared/ui/skeleton';

export function SiteEntryScreen(): React.ReactNode {
  const navigate = useNavigate();
  const { error, loading, sites } = useSites();

  useEffect(() => {
    if (loading || error !== undefined) return;
    if (sites.length === 0) {
      void navigate({ to: '/sites/new', replace: true });
      return;
    }
    const siteId = selectInitialSiteId(sites, getRememberedSiteId());
    if (siteId === null) {
      void navigate({ to: '/sites', replace: true });
      return;
    }
    void navigate({ to: '/sites/$siteId/analytics', params: { siteId }, replace: true });
  }, [error, loading, navigate, sites]);

  if (error !== undefined) {
    return <p className='text-destructive'>Error loading sites: {error.message}</p>;
  }

  return (
    <div className='space-y-6'>
      <Skeleton className='h-8 w-48' />
      <Card>
        <CardHeader>
          <Skeleton className='h-6 w-40' />
        </CardHeader>
        <CardContent>
          <Skeleton className='h-20 w-full' />
        </CardContent>
      </Card>
    </div>
  );
}
