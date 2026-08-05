import { useMemo, type ReactElement } from 'react';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
  Button,
} from '@/components/ui';
import { EventDefinitionFieldsFragmentDoc } from '@/gql/graphql';
import { useFragment as getFragmentData } from '@/gql/fragment-masking';
import EventDefinitionEditor from './event-definitions/event-definition-editor';
import EventDefinitionList from './event-definitions/event-definition-list';
import type { EventDefinitionsCardProps } from './event-definitions/types';
import { useEventDefinitionEditor } from './event-definitions/use-event-definition-editor';

export const EventDefinitionsCard = ({
  definitions,
  saving,
  deleting,
  onSave,
  onDelete,
}: EventDefinitionsCardProps): ReactElement => {
  const definitionItems = getFragmentData(EventDefinitionFieldsFragmentDoc, definitions);
  const sortedDefinitions = useMemo(
    () => [...definitionItems].sort((a, b) => a.name.localeCompare(b.name)),
    [definitionItems]
  );
  const editor = useEventDefinitionEditor({ onDelete, onSave });

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex flex-wrap items-center justify-between gap-3">
          <span>Event Definitions</span>
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={() => {
              if (editor.editorOpen) {
                editor.resetDraft();
                return;
              }
              editor.setEditorOpen(true);
              editor.setError('');
            }}
          >
            {editor.editorOpen ? 'Close editor' : 'New event name'}
          </Button>
        </CardTitle>
        <CardDescription>
          Allowlist event names and metadata fields to keep tracking clean.
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        {editor.editorOpen ? (
          <EventDefinitionEditor
            eventSnippet={editor.eventSnippet}
            onAddField={editor.handleAddField}
            onCancel={editor.resetDraft}
            onSave={() => void editor.handleSave()}
            saving={saving}
            setDraftFields={editor.setDraftFields}
            setDraftName={editor.setDraftName}
            setShowSnippet={editor.setShowSnippet}
            state={editor.state}
          />
        ) : (
          <p className="text-sm text-muted-foreground">
            Create an event definition to allowlist event names and metadata.
          </p>
        )}
        <EventDefinitionList
          definitions={sortedDefinitions}
          deleting={deleting}
          editingName={editor.state.originalName}
          onDelete={onDelete}
          onEdit={editor.startEdit}
          onResetEditor={editor.resetDraft}
          setError={editor.setError}
        />
      </CardContent>
    </Card>
  );
};
