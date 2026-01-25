
import { Card, CardContent, CardHeader, CardTitle, Badge } from '@/components/ui';
import { FilterLink } from '@/components/filter-link';
import { EventFieldsFragmentDoc, EventPropertyFieldsFragmentDoc } from '@/gql/graphql';
import { useFragment as getFragmentData, type FragmentType } from '@/gql/fragment-masking';
import { PaginationControls } from '@/components/pagination-controls';

interface EventsCardProps {
  siteId: string;
  events: Array<FragmentType<typeof EventFieldsFragmentDoc>>;
  total: number;
  page: number;
  pageSize: number;
  onPageChange: (page: number) => void;
}

const EMPTY_COUNT = 0;
const FALLBACK_PATH = '/';
const EMPTY_STRING = '';

export function EventsCard({
  siteId,
  events,
  total,
  page,
  pageSize,
  onPageChange,
}: EventsCardProps): React.JSX.Element {
  const eventItems = getFragmentData(EventFieldsFragmentDoc, events);

  return (
    <Card className="hover:shadow-md transition-shadow">
      <CardHeader>
        <CardTitle className="flex items-center justify-between">
          <span>Recent Events</span>
          <Badge variant="secondary">{total}</Badge>
        </CardTitle>
      </CardHeader>
      <CardContent>
        {eventItems.length === EMPTY_COUNT ? (
          <p className="text-sm text-muted-foreground text-center py-6">No events recorded yet.</p>
        ) : (
          <div className="space-y-4">
            {eventItems.map((event) => {
              const definitionName = event.definition?.name ?? EMPTY_STRING;
              const hasDefinitionName = definitionName !== EMPTY_STRING;
              return (
                <div key={event.id} className="border rounded-md p-3">
                  <div className="flex items-center justify-between gap-2">
                    <div className="min-w-0">
                      {hasDefinitionName ? (
                        <FilterLink
                          siteId={siteId}
                          filterKey="eventName"
                          value={definitionName}
                          className="text-sm font-medium truncate hover:underline underline-offset-2 block"
                        >
                          {definitionName}
                        </FilterLink>
                      ) : (
                        <span className="text-sm font-medium truncate block">{event.path}</span>
                      )}
                    <FilterLink
                      siteId={siteId}
                      filterKey="eventPath"
                      value={event.path === EMPTY_STRING ? FALLBACK_PATH : event.path}
                      className="text-xs text-muted-foreground truncate hover:underline underline-offset-2 block"
                    >
                      {event.path === EMPTY_STRING ? FALLBACK_PATH : event.path}
                    </FilterLink>
                  </div>
                  <div className="flex items-center gap-2 shrink-0">
                    <Badge variant="outline">
                      {hasDefinitionName ? 'Event' : 'Pageview'}
                    </Badge>
                    <span className="text-xs text-muted-foreground">
                      {new Date(event.createdAt).toLocaleString()}
                    </span>
                  </div>
                </div>
                {event.properties.length > EMPTY_COUNT ? (
                  <div className="flex flex-wrap gap-2 mt-2">
                    {event.properties.map((property) => {
                      const propertyData = getFragmentData(EventPropertyFieldsFragmentDoc, property);
                      return (
                        <Badge key={`${event.id}-${propertyData.key}`} variant="outline" className="max-w-full">
                          <span className="break-all">{propertyData.key}: {propertyData.value}</span>
                        </Badge>
                      );
                    })}
                  </div>
                ) : null}
              </div>
            );
            })}
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
