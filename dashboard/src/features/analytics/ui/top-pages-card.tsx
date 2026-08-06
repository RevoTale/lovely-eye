import { Globe } from 'lucide-react';
import type { FunctionComponent } from 'react';
import type { DashboardLoadState } from '@/features/analytics/model/dashboard-load-state';
import BoardCard from '@/features/analytics/ui/board-card';
import { FilterLink } from '@/features/analytics/ui/filter-link';
import { ListEmptyState } from '@/features/analytics/ui/list-empty-state';
import type { PageStatsFieldsFragment } from '@/shared/api/generated/graphql';
import { SKELETON_KEYS } from '@/shared/lib/skeleton';
import { Badge } from '@/shared/ui/badge';
import { Progress } from '@/shared/ui/progress';
import { Skeleton } from '@/shared/ui/skeleton';

interface TopPagesCardProps {
  pages: PageStatsFieldsFragment[];
  total: number;
  page: number;
  pageSize: number;
  siteId: string;
  onPageChange: (page: number) => void;
  state?: DashboardLoadState;
}

const EMPTY_COUNT = 0;
const FIRST_INDEX = 0;
const PERCENT_MULTIPLIER = 100;

const TopPagesCard: FunctionComponent<TopPagesCardProps> = ({
  pages,
  total,
  page,
  pageSize,
  siteId,
  onPageChange,
  state = 'ready',
}) => {
  const maxViews =
    pages.length > EMPTY_COUNT ? (pages[FIRST_INDEX]?.views ?? EMPTY_COUNT) : EMPTY_COUNT;

  return (
    <BoardCard
      title='Top Pages'
      icon={Globe}
      state={state}
      pagination={{ page, pageSize, total, onPageChange }}
      overlayLabel='Refreshing pages'
      skeleton={
        <div className='space-y-3'>
          {SKELETON_KEYS.map((key) => (
            <div key={key} className='space-y-2'>
              <div className='flex items-center justify-between'>
                <Skeleton className='h-4 w-32' />
                <Skeleton className='h-5 w-12' />
              </div>
              <Skeleton className='h-2 w-full' />
            </div>
          ))}
        </div>
      }
    >
      <div className='space-y-3'>
        {pages.length > EMPTY_COUNT ? (
          pages.map((pageStat) => (
            <div key={pageStat.path}>
              <div className='mb-1 flex min-w-0 items-center justify-between gap-2'>
                <FilterLink
                  siteId={siteId}
                  filterKey='page'
                  value={pageStat.path}
                  className='min-w-0 flex-1 cursor-pointer truncate text-sm font-medium hover:text-primary hover:underline'
                >
                  {pageStat.path}
                </FilterLink>
                <Badge variant='secondary'>{pageStat.views.toLocaleString()}</Badge>
              </div>
              <Progress
                value={
                  maxViews > EMPTY_COUNT
                    ? (pageStat.views / maxViews) * PERCENT_MULTIPLIER
                    : EMPTY_COUNT
                }
                className='h-2'
              />
            </div>
          ))
        ) : (
          <ListEmptyState title='No page data yet' />
        )}
      </div>
    </BoardCard>
  );
};

export default TopPagesCard;
