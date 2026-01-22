
import { Card, CardContent, CardHeader, CardTitle, Badge } from '@/components/ui';
import { FilterLink } from '@/components/filter-link';
import { EventFieldsFragmentDoc } from '@/gql/graphql';
import type { EventCountFieldsFragment } from '@/gql/graphql';
import { useFragment as getFragmentData } from '@/gql/fragment-masking';

interface EventCountsCardProps {
  siteId: string;
  eventCounts: EventCountFieldsFragment[];
}

const EMPTY_COUNT = 0;

export function EventCountsCard({ siteId, eventCounts }: EventCountsCardProps): React.JSX.Element {
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
            {eventCounts.map((item) => {
              const event = getFragmentData(EventFieldsFragmentDoc, item.event);
              return (
                <div key={event.id} className="flex items-center justify-between">
                  <FilterLink
                    siteId={siteId}
                    filterKey="eventName"
                    value={event.name}
                    className="text-sm font-medium hover:underline underline-offset-2"
                  >
                    {event.name}
                  </FilterLink>
                  <Badge variant="outline">{item.count}</Badge>
                </div>
              );
            })}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
