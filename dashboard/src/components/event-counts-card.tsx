import React from 'react';
import { Card, CardContent, CardHeader, CardTitle, Badge } from '@/components/ui';
import type { EventCount } from '@/gql/graphql';

interface EventCountsCardProps {
  eventCounts: EventCount[];
}

const EMPTY_COUNT = 0;

export function EventCountsCard({ eventCounts }: EventCountsCardProps): React.JSX.Element {
  return (
    <Card className="hover:shadow-md transition-shadow">
      <CardHeader>
        <CardTitle className="flex items-center justify-between">
          <span>Event Counts</span>
          <Badge variant="secondary">{eventCounts.length}</Badge>
        </CardTitle>
      </CardHeader>
      <CardContent>
        {eventCounts.length === EMPTY_COUNT ? (
          <p className="text-sm text-muted-foreground text-center py-6">No events recorded yet.</p>
        ) : (
          <div className="space-y-3">
            {eventCounts.map((item) => (
              <div key={item.event.id} className="flex items-center justify-between">
                <span className="text-sm font-medium">{item.event.name}</span>
                <Badge variant="outline">{item.count}</Badge>
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
