import React from 'react';
import { Card, CardContent, CardHeader, CardTitle, Badge } from '@/components/ui';
import type { Event } from '@/generated/graphql';
import { PaginationControls } from '@/components/pagination-controls';

interface EventsCardProps {
  events: Event[];
  total: number;
  page: number;
  pageSize: number;
  onPageChange: (page: number) => void;
}

export function EventsCard({
  events,
  total,
  page,
  pageSize,
  onPageChange,
}: EventsCardProps): React.JSX.Element {
  return (
    <Card className="hover:shadow-md transition-shadow">
      <CardHeader>
        <CardTitle className="flex items-center justify-between">
          <span>Recent Events</span>
          <Badge variant="secondary">{total}</Badge>
        </CardTitle>
      </CardHeader>
      <CardContent>
        {events.length === 0 ? (
          <p className="text-sm text-muted-foreground text-center py-6">No events recorded yet.</p>
        ) : (
          <div className="space-y-4">
            {events.map((event) => (
              <div key={event.id} className="border rounded-md p-3">
                <div className="flex items-center justify-between gap-2">
                  <div>
                    <p className="text-sm font-medium">{event.name}</p>
                    <p className="text-xs text-muted-foreground">{event.path || '/'}</p>
                  </div>
                  <span className="text-xs text-muted-foreground">
                    {new Date(event.createdAt).toLocaleString()}
                  </span>
                </div>
                {event.properties.length > 0 ? (
                  <div className="flex flex-wrap gap-2 mt-2">
                    {event.properties.map((property) => (
                      <Badge key={`${event.id}-${property.key}`} variant="outline">
                        {property.key}: {property.value}
                      </Badge>
                    ))}
                  </div>
                ) : null}
              </div>
            ))}
          </div>
        )}
        {total > pageSize ? (
          <div className="mt-4">
            <PaginationControls
              page={page}
              pageSize={pageSize}
              total={total}
              onPageChange={onPageChange}
            />
          </div>
        ) : null}
      </CardContent>
    </Card>
  );
}
