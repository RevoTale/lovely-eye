import { Link } from '@tanstack/react-router';
import { Activity, Settings } from 'lucide-react';
import type { FragmentType } from '@/shared/api/generated/fragment-masking';
import {
  RealtimeStatsFieldsFragmentDoc,
  SiteDetailsFieldsFragmentDoc,
} from '@/shared/api/generated/graphql';
import { readFragment } from '@/shared/api/read-fragment';
import { Badge } from '@/shared/ui/badge';

interface DashboardHeaderProps {
  site: FragmentType<typeof SiteDetailsFieldsFragmentDoc>;
  siteId: string;
  realtime: FragmentType<typeof RealtimeStatsFieldsFragmentDoc> | undefined;
}

export const DashboardHeader = ({
  site,
  siteId,
  realtime,
}: DashboardHeaderProps): React.ReactNode => {
  const EMPTY_COUNT = 0;
  const ZERO_VISITORS = 0;
  const siteData = readFragment(SiteDetailsFieldsFragmentDoc, site);
  const realtimeData =
    realtime === undefined ? undefined : readFragment(RealtimeStatsFieldsFragmentDoc, realtime);
  const domainList = siteData.domains.length > EMPTY_COUNT ? siteData.domains : [''];
  const domainLabel = domainList.join(' · ');
  const hasRealtime = realtimeData !== undefined;
  const realtimeVisitors = hasRealtime ? realtimeData.visitors : ZERO_VISITORS;

  return (
    <div className='flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between'>
      <div className='min-w-0'>
        <h1 className='text-3xl font-bold tracking-tight break-words'>{siteData.name}</h1>
        <p className='text-muted-foreground mt-1 break-all'>{domainLabel}</p>
      </div>
      <div className='flex flex-wrap items-center gap-2 sm:gap-3'>
        <Link to='/sites/$siteId/settings' params={{ siteId }}>
          <Badge
            variant='outline'
            className='flex items-center gap-2 px-3 py-2 cursor-pointer hover:bg-accent'
          >
            <Settings className='h-4 w-4' />
            <span>Settings</span>
          </Badge>
        </Link>
        {hasRealtime ? (
          <Badge variant='outline' className='flex items-center gap-2 px-3 py-2'>
            <span className='relative flex h-2 w-2'>
              <span className='animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75'></span>
              <span className='relative inline-flex rounded-full h-2 w-2 bg-green-500'></span>
            </span>
            <Activity className='h-4 w-4' />
            <span className='font-semibold'>{realtimeVisitors}</span>
            <span className='text-muted-foreground'>online</span>
          </Badge>
        ) : null}
      </div>
    </div>
  );
};
