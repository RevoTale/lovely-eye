import type { FunctionComponent } from 'react';
import type { DashboardLoadState } from '@/features/analytics/model/dashboard-load-state';
import DashboardCardState from '@/features/analytics/ui/dashboard-card-state';
import { FilterLink } from '@/features/analytics/ui/filter-link';
import { PaginationControls } from '@/features/analytics/ui/pagination-controls';
import type { EventCountFieldsFragment } from '@/shared/api/generated/graphql';
import { EventFieldsFragmentDoc } from '@/shared/api/generated/graphql';
import { readFragment } from '@/shared/api/read-fragment';
import { SKELETON_KEYS } from '@/shared/lib/skeleton';
import { Badge } from '@/shared/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/shared/ui/card';
import { Skeleton } from '@/shared/ui/skeleton';

interface EventCountsCardProps {
  siteId: string;
  eventCounts: EventCountFieldsFragment[];
  total: number;
  page: number;
  pageSize: number;
  onPageChange: (page: number) => void;
  state?: DashboardLoadState;
}

const EMPTY_COUNT = 0;

const EventCountsCard: FunctionComponent<EventCountsCardProps> = ({
  siteId,
  eventCounts,
  total,
  page,
  pageSize,
  onPageChange,
  state = 'ready',
}) => {
  return (
    <Card className='transition-shadow hover:shadow-md'>
      <CardHeader>
        <CardTitle className='flex items-center justify-between'>
          <span>Event Counts</span>
          <Badge variant='secondary'>{total}</Badge>
        </CardTitle>
      </CardHeader>
      <CardContent>
        <DashboardCardState
          state={state}
          overlayLabel='Refreshing counts'
          skeleton={
            <div className='space-y-3'>
              {SKELETON_KEYS.map((key) => (
                <div key={key} className='flex items-center justify-between'>
                  <Skeleton className='h-4 w-28' />
                  <Skeleton className='h-5 w-12' />
                </div>
              ))}
            </div>
          }
        >
          {eventCounts.length === EMPTY_COUNT ? (
            <p className='py-6 text-center text-sm text-muted-foreground'>
              No events recorded yet.
            </p>
          ) : (
            <div className='space-y-3'>
              {eventCounts.map((item) => {
                const event = readFragment(EventFieldsFragmentDoc, item.event);
                return (
                  <div key={event.id} className='flex min-w-0 items-center justify-between gap-2'>
                    <FilterLink
                      siteId={siteId}
                      filterKey='eventName'
                      value={event.name}
                      className='min-w-0 flex-1 truncate text-sm font-medium hover:underline'
                    >
                      {event.name}
                    </FilterLink>
                    <Badge variant='outline' className='shrink-0'>
                      {item.count}
                    </Badge>
                  </div>
                );
              })}
            </div>
          )}
          {eventCounts.length > EMPTY_COUNT ? (
            <div className='mt-4'>
              <PaginationControls
                page={page}
                pageSize={pageSize}
                total={total}
                onPageChange={onPageChange}
              />
            </div>
          ) : null}
        </DashboardCardState>
      </CardContent>
    </Card>
  );
};

export default EventCountsCard;
