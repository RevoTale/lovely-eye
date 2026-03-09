import type { FunctionComponent } from 'react';
import DashboardCardState from '@/components/dashboard-card-state';
import { FilterLink } from '@/components/filter-link';
import { PaginationControls } from '@/components/pagination-controls';
import { Badge, Card, CardContent, CardHeader, CardTitle, Skeleton } from '@/components/ui';
import { EventFieldsFragmentDoc } from '@/gql/graphql';
import type { EventCountFieldsFragment } from '@/gql/graphql';
import { useFragment as getFragmentData } from '@/gql/fragment-masking';
import type { DashboardLoadState } from '@/lib/dashboard-load-state';

interface EventCountsCardProps {
  siteId: string;
  eventCounts: EventCountFieldsFragment[];
  page: number;
  pageSize: number;
  onPageChange: (page: number) => void;
  state?: DashboardLoadState;
}

const EMPTY_COUNT = 0;

const EventCountsCard: FunctionComponent<EventCountsCardProps> = ({ siteId, eventCounts, page, pageSize, onPageChange, state = 'ready' }) => {
  const hasNextPage = eventCounts.length === pageSize;
  const totalEstimate = (page - 1) * pageSize + eventCounts.length + (hasNextPage ? 1 : 0);

  return (
    <Card className="transition-shadow hover:shadow-md">
      <CardHeader><CardTitle className="flex items-center justify-between"><span>Event Counts</span><Badge variant="secondary">{eventCounts.length}</Badge></CardTitle></CardHeader>
      <CardContent>
        <DashboardCardState state={state} overlayLabel="Refreshing counts" skeleton={<div className="space-y-3">{Array.from({ length: 5 }, (_, index) => <div key={index} className="flex items-center justify-between"><Skeleton className="h-4 w-28" /><Skeleton className="h-5 w-12" /></div>)}</div>}>
          <>
            {eventCounts.length === EMPTY_COUNT ? <p className="py-6 text-center text-sm text-muted-foreground">No events recorded yet.</p> : <div className="space-y-3">{eventCounts.map((item) => { const event = getFragmentData(EventFieldsFragmentDoc, item.event); return <div key={event.id} className="flex items-center justify-between"><FilterLink siteId={siteId} filterKey="eventName" value={event.name} className="text-sm font-medium hover:underline">{event.name}</FilterLink><Badge variant="outline">{item.count}</Badge></div>; })}</div>}
            {eventCounts.length > EMPTY_COUNT ? <div className="mt-4"><PaginationControls page={page} pageSize={pageSize} total={totalEstimate} onPageChange={onPageChange} /></div> : null}
          </>
        </DashboardCardState>
      </CardContent>
    </Card>
  );
};

export default EventCountsCard;
