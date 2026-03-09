import type { FunctionComponent } from 'react';
import { EventCountFieldsFragmentDoc } from '@/gql/graphql';
import type { EventCountsQuery, EventsQuery } from '@/gql/graphql';
import { useFragment as getFragmentData } from '@/gql/fragment-masking';
import type { DashboardLoadState } from '@/lib/dashboard-load-state';
import EventCountsCard from '@/components/event-counts-card';
import EventsCard from '@/components/events-card';

interface EventsSectionProps {
  siteId: string;
  eventsState: DashboardLoadState;
  eventCountsState: DashboardLoadState;
  eventsResult: EventsQuery['events'] | undefined;
  eventsCounts: EventCountsQuery['eventCounts'];
  page: number;
  pageSize: number;
  onPageChange: (page: number) => void;
  countsPage: number;
  countsPageSize: number;
  onCountsPageChange: (page: number) => void;
}

const EventsSection: FunctionComponent<EventsSectionProps> = ({
  siteId, eventsState, eventCountsState, eventsResult, eventsCounts, page, pageSize, onPageChange, countsPage, countsPageSize, onCountsPageChange,
}) => {
  const eventCountsData = getFragmentData(EventCountFieldsFragmentDoc, eventsCounts);

  return (
    <div className="grid gap-6 md:grid-cols-2">
      <EventsCard siteId={siteId} events={eventsResult?.events ?? []} total={eventsResult?.total ?? 0} page={page} pageSize={pageSize} onPageChange={onPageChange} state={eventsState} />
      <EventCountsCard siteId={siteId} eventCounts={eventCountsData} page={countsPage} pageSize={countsPageSize} onPageChange={onCountsPageChange} state={eventCountsState} />
    </div>
  );
};

export default EventsSection;
