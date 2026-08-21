import { Link } from '@tanstack/react-router';
import { ExternalLink, Globe, Plus, TrendingUp } from 'lucide-react';
import { rememberSite } from '@/features/sites/model/recent-site';
import { useSites } from '@/features/sites/model/use-sites';
import { SKELETON_KEYS } from '@/shared/lib/skeleton';
import { Badge } from '@/shared/ui/badge';
import { Button } from '@/shared/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card';
import { Skeleton } from '@/shared/ui/skeleton';

const EMPTY_COUNT = 0;
const FIRST_INDEX = 0;
const EXTRA_DOMAIN_OFFSET = 1;
const SKELETON_CARD_COUNT = 3;

export const SitesScreen = (): React.ReactNode => {
  const { sites, isInitialLoading, error } = useSites();

  if (isInitialLoading) {
    return (
      <div className='space-y-6'>
        <div className='flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between'>
          <div className='space-y-2'>
            <Skeleton className='h-8 w-32' />
            <Skeleton className='h-4 w-48' />
          </div>
          <Skeleton className='h-10 w-32' />
        </div>
        <div className='grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3'>
          {SKELETON_KEYS.slice(0, SKELETON_CARD_COUNT).map((key) => (
            <Card key={key}>
              <CardHeader>
                <Skeleton className='h-6 w-32' />
                <Skeleton className='h-4 w-48' />
              </CardHeader>
              <CardContent>
                <Skeleton className='h-4 w-24' />
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    );
  }

  if (error !== undefined) {
    return (
      <div className='flex items-center justify-center min-h-[400px]'>
        <div className='text-destructive'>Error loading sites: {error.message}</div>
      </div>
    );
  }

  return (
    <div className='space-y-8'>
      <div className='flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between'>
        <div>
          <h1 className='text-3xl font-bold tracking-tight'>Sites</h1>
          <p className='text-muted-foreground mt-1'>Manage your tracked websites</p>
        </div>
        <Link to='/sites/new'>
          <Button className='shadow-sm'>
            <Plus className='mr-2 h-4 w-4' />
            Add Site
          </Button>
        </Link>
      </div>

      {sites.length === EMPTY_COUNT ? (
        <Card className='border-dashed'>
          <CardContent className='flex flex-col items-center justify-center py-16'>
            <div className='h-16 w-16 rounded-full bg-primary/10 flex items-center justify-center mb-4'>
              <Globe className='h-8 w-8 text-primary' />
            </div>
            <h3 className='text-xl font-semibold mb-2'>No sites yet</h3>
            <p className='text-muted-foreground text-center mb-6 max-w-sm'>
              Add your first website to start tracking analytics and monitor visitor behavior
            </p>
            <Link to='/sites/new'>
              <Button size='lg'>
                <Plus className='mr-2 h-4 w-4' />
                Add your first site
              </Button>
            </Link>
          </CardContent>
        </Card>
      ) : (
        <div className='grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3'>
          {sites.map((site) => (
            <Link
              key={site.id}
              to='/sites/$siteId/analytics'
              params={{ siteId: site.id }}
              onClick={() => rememberSite(site.id)}
            >
              <Card className='group hover:shadow-lg hover:border-primary/50 transition-all cursor-pointer h-full'>
                <CardHeader>
                  <div className='flex items-start justify-between'>
                    <div className='h-10 w-10 rounded-lg bg-primary/10 flex items-center justify-center group-hover:bg-primary/20 transition-colors'>
                      <Globe className='h-5 w-5 text-primary' />
                    </div>
                    <Badge variant='secondary' className='flex items-center gap-1'>
                      <TrendingUp className='h-3 w-3' />
                      Active
                    </Badge>
                  </div>
                  <CardTitle className='mt-4 group-hover:text-primary transition-colors'>
                    {site.name}
                  </CardTitle>
                  <CardDescription className='flex min-w-0 items-center gap-2'>
                    <span className='truncate'>{site.domains[FIRST_INDEX] ?? ''}</span>
                    {site.domains.length > EXTRA_DOMAIN_OFFSET ? (
                      <span className='shrink-0 text-xs text-muted-foreground'>
                        +{site.domains.length - EXTRA_DOMAIN_OFFSET} more
                      </span>
                    ) : null}
                    <ExternalLink className='size-3 shrink-0 opacity-0 transition-opacity group-hover:opacity-100' />
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <p className='text-sm text-muted-foreground'>
                    Added{' '}
                    {new Date(site.createdAt).toLocaleDateString('en-US', {
                      month: 'short',
                      day: 'numeric',
                      year: 'numeric',
                    })}
                  </p>
                </CardContent>
              </Card>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
};
