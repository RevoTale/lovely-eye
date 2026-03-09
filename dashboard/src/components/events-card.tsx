import type { FunctionComponent } from 'react';
import DashboardCardState from '@/components/dashboard-card-state';
import { FilterLink } from '@/components/filter-link';
import { PaginationControls } from '@/components/pagination-controls';
import { Badge, Card, CardContent, CardHeader, CardTitle, Skeleton } from '@/components/ui';
import { EventFieldsFragmentDoc, EventPropertyFieldsFragmentDoc } from '@/gql/graphql';
import { useFragment as getFragmentData, type FragmentType } from '@/gql/fragment-masking';
import type { DashboardLoadState } from '@/lib/dashboard-load-state';

interface EventsCardProps {
  siteId: string;
  events: Array<FragmentType<typeof EventFieldsFragmentDoc>>;
  total: number;
  page: number;
  pageSize: number;
  onPageChange: (page: number) => void;
  state?: DashboardLoadState;
}

const EMPTY_COUNT = 0;
const FALLBACK_PATH = '/';

const EventsCard: FunctionComponent<EventsCardProps> = ({ siteId, events, total, page, pageSize, onPageChange, state = 'ready' }) => {
  const eventItems = getFragmentData(EventFieldsFragmentDoc, events);

  return (
    <Card className="transition-shadow hover:shadow-md">
      <CardHeader><CardTitle className="flex items-center justify-between"><span>Recent Events</span><Badge variant="secondary">{total}</Badge></CardTitle></CardHeader>
      <CardContent>
        <DashboardCardState state={state} overlayLabel="Refreshing events" skeleton={<div className="space-y-4">{Array.from({ length: 3 }, (_, index) => <div key={index} className="space-y-2 rounded-md border p-3"><Skeleton className="h-4 w-28" /><Skeleton className="h-3 w-20" /><div className="flex gap-2"><Skeleton className="h-5 w-16" /><Skeleton className="h-5 w-24" /></div></div>)}</div>}>
          <>
            {eventItems.length === EMPTY_COUNT ? <p className="py-6 text-center text-sm text-muted-foreground">No events recorded yet.</p> : <div className="space-y-4">{eventItems.map((event) => {
              const definitionName = event.definition?.name ?? '';
              const eventPath = event.path === '' ? FALLBACK_PATH : event.path;
              return <div key={event.id} className="rounded-md border p-3"><div className="flex items-center justify-between gap-2"><div className="min-w-0">{definitionName !== '' ? <FilterLink siteId={siteId} filterKey="eventName" value={definitionName} className="block truncate text-sm font-medium hover:underline">{definitionName}</FilterLink> : <span className="block truncate text-sm font-medium">{eventPath}</span>}<FilterLink siteId={siteId} filterKey="eventPath" value={eventPath} className="block truncate text-xs text-muted-foreground hover:underline">{eventPath}</FilterLink></div><div className="flex shrink-0 items-center gap-2"><Badge variant="outline">{definitionName !== '' ? 'Event' : 'Pageview'}</Badge><span className="text-xs text-muted-foreground">{new Date(event.createdAt).toLocaleString()}</span></div></div>{event.properties.length > EMPTY_COUNT ? <div className="mt-2 flex flex-wrap gap-2">{event.properties.map((property) => { const propertyData = getFragmentData(EventPropertyFieldsFragmentDoc, property); return <Badge key={`${event.id}-${propertyData.key}`} variant="outline" className="max-w-full"><span className="break-all">{propertyData.key}: {propertyData.value}</span></Badge>; })}</div> : null}</div>;
            })}</div>}
            {total > pageSize ? <div className="mt-4"><PaginationControls page={page} pageSize={pageSize} total={total} onPageChange={onPageChange} /></div> : null}
          </>
        </DashboardCardState>
      </CardContent>
    </Card>
  );
};

export default EventsCard;
