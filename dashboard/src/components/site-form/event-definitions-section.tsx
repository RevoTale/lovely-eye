import React from 'react';
import { useMutation, useQuery } from '@apollo/client';
import { EVENT_DEFINITIONS_QUERY, UPSERT_EVENT_DEFINITION_MUTATION, DELETE_EVENT_DEFINITION_MUTATION } from '@/graphql';
import { EventDefinitionsCard } from '@/components/event-definitions-card';
import type { EventDefinitionInput } from '@/generated/graphql';

interface EventDefinitionsSectionProps {
  siteId: string;
}

export function EventDefinitionsSection({
  siteId,
}: EventDefinitionsSectionProps): React.JSX.Element {
  const [actionError, setActionError] = React.useState('');
  const { data: eventDefinitionsData } = useQuery(EVENT_DEFINITIONS_QUERY, {
    variables: { siteId },
  });

  const [upsertEventDefinition, { loading: savingDefinition }] = useMutation(UPSERT_EVENT_DEFINITION_MUTATION);
  const [deleteEventDefinition, { loading: deletingDefinition }] = useMutation(DELETE_EVENT_DEFINITION_MUTATION);

  const eventDefinitions = eventDefinitionsData?.eventDefinitions ?? [];

  const handleSaveEventDefinition = async (input: EventDefinitionInput): Promise<void> => {
    setActionError('');
    try {
      await upsertEventDefinition({
        variables: {
          siteId,
          input,
        },
        refetchQueries: [{ query: EVENT_DEFINITIONS_QUERY, variables: { siteId } }],
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
        refetchQueries: [{ query: EVENT_DEFINITIONS_QUERY, variables: { siteId } }],
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
