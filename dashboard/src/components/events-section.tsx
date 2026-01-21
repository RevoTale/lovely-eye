import React from 'react';
import { Card, CardContent, CardHeader, Skeleton } from '@/components/ui';
import type { EventsResult, EventCount } from '@/gql/graphql';
import { EventsCard } from '@/components/events-card';
import { EventCountsCard } from '@/components/event-counts-card';

interface EventsSectionProps {
  loading: boolean;
  eventsResult: EventsResult | undefined;
  eventsCounts: EventCount[];
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
  if (loading) {
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

  if (eventsResult === undefined) {
    return null;
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
      <EventCountsCard eventCounts={eventsCounts} />
    </div>
  );
}
