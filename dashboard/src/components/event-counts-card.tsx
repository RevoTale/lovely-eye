import React, { useMemo } from 'react';
import { Card, CardContent, CardHeader, CardTitle, Badge } from '@/components/ui';
import type { Event } from '@/gql/graphql';

interface EventCountsCardProps {
  events: Event[];
}

const EMPTY_COUNT = 0;
const COUNT_INCREMENT = 1;

export function EventCountsCard({ events }: EventCountsCardProps): React.JSX.Element {
  const counts = useMemo(() => {
    const counter = new Map<string, number>();
    for (const event of events) {
      counter.set(event.name, (counter.get(event.name) ?? EMPTY_COUNT) + COUNT_INCREMENT);
    }
    return Array.from(counter.entries())
      .map(([name, count]) => ({ name, count }))
      .sort((a, b) => b.count - a.count);
  }, [events]);

  return (
    <Card className="hover:shadow-md transition-shadow">
      <CardHeader>
        <CardTitle className="flex items-center justify-between">
          <span>Event Counts</span>
          <Badge variant="secondary">{events.length}</Badge>
        </CardTitle>
      </CardHeader>
      <CardContent>
        {counts.length === EMPTY_COUNT ? (
          <p className="text-sm text-muted-foreground text-center py-6">No events recorded yet.</p>
        ) : (
          <div className="space-y-3">
            {counts.map((item) => (
              <div key={item.name} className="flex items-center justify-between">
                <span className="text-sm font-medium">{item.name}</span>
                <Badge variant="outline">{item.count}</Badge>
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
