import React from 'react';
import { useQuery, useMutation } from '@apollo/client/react';
import {
  EventDefinitionsDocument,
  UpsertEventDefinitionDocument,
  DeleteEventDefinitionDocument,
  type EventDefinitionInput,
} from '@/gql/graphql';
import { EventDefinitionsCard } from '@/components/event-definitions-card';

interface EventDefinitionsSectionProps {
  siteId: string;
}

export function EventDefinitionsSection({
  siteId,
}: EventDefinitionsSectionProps): React.JSX.Element {
  const [actionError, setActionError] = React.useState('');
  const { data: eventDefinitionsData } = useQuery(EventDefinitionsDocument, {
    variables: { siteId },
  });

  const [upsertEventDefinition, { loading: savingDefinition }] = useMutation(UpsertEventDefinitionDocument);
  const [deleteEventDefinition, { loading: deletingDefinition }] = useMutation(DeleteEventDefinitionDocument);

  const eventDefinitions = eventDefinitionsData?.eventDefinitions ?? [];

  const handleSaveEventDefinition = async (input: EventDefinitionInput): Promise<void> => {
    setActionError('');
    try {
      await upsertEventDefinition({
        variables: {
          siteId,
          input,
        },
        refetchQueries: [{ query: EventDefinitionsDocument, variables: { siteId } }],
      });
    } catch (err) {
      setActionError(err instanceof Error ? err.message : 'Failed to save event definition');
    }
  };

  const handleDeleteEventDefinition = async (nameToDelete: string): Promise<void> => {
    setActionError('');
    try {
      await deleteEventDefinition({
        variables: {
          siteId,
          name: nameToDelete,
        },
        refetchQueries: [{ query: EventDefinitionsDocument, variables: { siteId } }],
      });
    } catch (err) {
      setActionError(err instanceof Error ? err.message : 'Failed to delete event definition');
    }
  };

  return (
    <div className="space-y-2">
      <EventDefinitionsCard
        definitions={eventDefinitions}
        saving={savingDefinition}
        deleting={deletingDefinition}
        onSave={handleSaveEventDefinition}
        onDelete={handleDeleteEventDefinition}
      />
      {actionError ? (
        <p className="text-xs text-destructive">
          {actionError}
        </p>
      ) : null}
    </div>
  );
}
