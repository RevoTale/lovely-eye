import type { FunctionComponent } from 'react';
import type { DashboardLoadState } from '@/features/analytics/model/dashboard-load-state';
import EventCountsCard from '@/features/analytics/ui/event-counts-card';
import EventsCard from '@/features/analytics/ui/events-card';
import { useFragment as getFragmentData } from '@/shared/api/generated/fragment-masking';
import type { EventCountsQuery, EventsQuery } from '@/shared/api/generated/graphql';
import { EventCountFieldsFragmentDoc } from '@/shared/api/generated/graphql';

interface EventsSectionProps {
  siteId: string;
  eventsState: DashboardLoadState;
  eventCountsState: DashboardLoadState;
  eventsResult: EventsQuery['events'] | undefined;
  eventsCounts: EventCountsQuery['eventCounts']['items'];
  eventsCountsTotal: number;
  page: number;
  pageSize: number;
  onPageChange: (page: number) => void;
  countsPage: number;
  countsPageSize: number;
  onCountsPageChange: (page: number) => void;
}

const EventsSection: FunctionComponent<EventsSectionProps> = ({
  siteId,
  eventsState,
  eventCountsState,
  eventsResult,
  eventsCounts,
  eventsCountsTotal,
  page,
  pageSize,
  onPageChange,
  countsPage,
  countsPageSize,
  onCountsPageChange,
}) => {
  const eventCountsData = getFragmentData(EventCountFieldsFragmentDoc, eventsCounts);

  return (
    <div className='grid grid-cols-1 gap-6 md:grid-cols-2'>
      <EventsCard
        siteId={siteId}
        events={eventsResult?.events ?? []}
        total={eventsResult?.total ?? 0}
        page={page}
        pageSize={pageSize}
        onPageChange={onPageChange}
        state={eventsState}
      />
      <EventCountsCard
        siteId={siteId}
        eventCounts={eventCountsData}
        total={eventsCountsTotal}
        page={countsPage}
        pageSize={countsPageSize}
        onPageChange={onCountsPageChange}
        state={eventCountsState}
      />
    </div>
  );
};

export default EventsSection;
