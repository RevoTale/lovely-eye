
import { Card, CardContent, CardHeader, CardTitle, Badge } from '@/components/ui';
import type { Event } from '@/gql/graphql';
import { PaginationControls } from '@/components/pagination-controls';

interface EventsCardProps {
  events: Event[];
  total: number;
  page: number;
  pageSize: number;
  onPageChange: (page: number) => void;
}

const EMPTY_COUNT = 0;
const FALLBACK_PATH = '/';
const EMPTY_STRING = '';

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
        {events.length === EMPTY_COUNT ? (
          <p className="text-sm text-muted-foreground text-center py-6">No events recorded yet.</p>
        ) : (
          <div className="space-y-4">
            {events.map((event) => (
              <div key={event.id} className="border rounded-md p-3">
                <div className="flex items-center justify-between gap-2">
                  <div className="min-w-0">
                    <p className="text-sm font-medium truncate">{event.name}</p>
                    <p className="text-xs text-muted-foreground truncate">
                      {event.path === EMPTY_STRING ? FALLBACK_PATH : event.path}
                    </p>
                  </div>
                  <span className="text-xs text-muted-foreground shrink-0">
                    {new Date(event.createdAt).toLocaleString()}
                  </span>
                </div>
                {event.properties.length > EMPTY_COUNT ? (
                  <div className="flex flex-wrap gap-2 mt-2">
                    {event.properties.map((property) => (
                      <Badge key={`${event.id}-${property.key}`} variant="outline" className="max-w-full">
                        <span className="break-all">{property.key}: {property.value}</span>
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
