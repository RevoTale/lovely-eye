import { useMutation, useQuery } from '@apollo/client/react';
import type { ReactElement } from 'react';
import { EventDefinitionsCard } from '@/features/event-definitions/ui/event-definitions-card';
import {
  DeleteEventDefinitionDocument,
  type EventDefinitionInput,
  EventDefinitionsDocument,
  UpsertEventDefinitionDocument,
} from '@/shared/api/generated/graphql';
import { Card, CardContent, CardHeader } from '@/shared/ui/card';
import { Skeleton } from '@/shared/ui/skeleton';

interface EventDefinitionsSectionProps {
  siteId: string;
}

const EVENT_DEFINITIONS_PAGING = { limit: 100, offset: 0 } as const;

export const EventDefinitionsSection = ({ siteId }: EventDefinitionsSectionProps): ReactElement => {
  const {
    data: eventDefinitionsData,
    error,
    loading,
  } = useQuery(EventDefinitionsDocument, {
    variables: { siteId, paging: EVENT_DEFINITIONS_PAGING },
  });

  const [upsertEventDefinition, { loading: savingDefinition }] = useMutation(
    UpsertEventDefinitionDocument
  );
  const [deleteEventDefinition, { loading: deletingDefinition }] = useMutation(
    DeleteEventDefinitionDocument
  );

  const eventDefinitions = eventDefinitionsData?.eventDefinitions ?? [];

  const handleSaveEventDefinition = async (input: EventDefinitionInput): Promise<void> => {
    await upsertEventDefinition({
      variables: { siteId, input },
      refetchQueries: [
        {
          query: EventDefinitionsDocument,
          variables: { siteId, paging: EVENT_DEFINITIONS_PAGING },
        },
      ],
    });
  };

  const handleDeleteEventDefinition = async (nameToDelete: string): Promise<void> => {
    await deleteEventDefinition({
      variables: { siteId, name: nameToDelete },
      refetchQueries: [
        {
          query: EventDefinitionsDocument,
          variables: { siteId, paging: EVENT_DEFINITIONS_PAGING },
        },
      ],
    });
  };

  if (loading && eventDefinitionsData === undefined) {
    return (
      <Card>
        <CardHeader>
          <Skeleton className='h-6 w-48' />
        </CardHeader>
        <CardContent>
          <Skeleton className='h-20 w-full' />
        </CardContent>
      </Card>
    );
  }
  if (error !== undefined) {
    return <p className='text-destructive'>Error loading event definitions: {error.message}</p>;
  }

  return (
    <EventDefinitionsCard
      definitions={eventDefinitions}
      saving={savingDefinition}
      deleting={deletingDefinition}
      onSave={handleSaveEventDefinition}
      onDelete={handleDeleteEventDefinition}
    />
  );
};
