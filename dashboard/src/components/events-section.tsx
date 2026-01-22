
import { Card, CardContent, CardHeader, Skeleton } from '@/components/ui';
import { EventCountFieldsFragmentDoc } from '@/gql/graphql';
import type { EventCountsQuery, EventsQuery } from '@/gql/graphql';
import { useFragment as getFragmentData } from '@/gql/fragment-masking';
import { EventsCard } from '@/components/events-card';
import { EventCountsCard } from '@/components/event-counts-card';

interface EventsSectionProps {
  loading: boolean;
  eventsResult: EventsQuery['events'] | undefined;
  eventsCounts: EventCountsQuery['eventCounts'];
  page: number;
  pageSize: number;
  onPageChange: (page: number) => void;
}

export function EventsSection({
  loading,
  eventsResult,
  eventsCounts,
  page,
  pageSize,
  onPageChange,
}: EventsSectionProps): React.JSX.Element | null {
  const eventCountsData = getFragmentData(EventCountFieldsFragmentDoc, eventsCounts);
  if (eventsResult === undefined || loading) {
    return (
      <Card>
        <CardHeader>
          <Skeleton className="h-5 w-32" />
        </CardHeader>
        <CardContent>
          <Skeleton className="h-20 w-full" />
        </CardContent>
      </Card>
    );
  }



  return (
    <div className="grid gap-6 md:grid-cols-2">
      <EventsCard
        events={eventsResult.events}
        total={eventsResult.total}
        page={page}
        pageSize={pageSize}
        onPageChange={onPageChange}
      />
      <EventCountsCard eventCounts={eventCountsData} />
    </div>
  );
}
