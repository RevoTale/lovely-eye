import { useState, type FunctionComponent } from 'react';
import { Button } from '@/components/ui';
import { EMPTY_COUNT, MIN_FIELD_LENGTH } from './event-definition-utils';
import type { EventDefinitionItem } from './types';

interface EventDefinitionListProps {
  definitions: EventDefinitionItem[];
  deleting: boolean;
  editingName: string | null;
  onDelete: (name: string) => Promise<void>;
  onEdit: (definition: EventDefinitionItem) => void;
  onResetEditor: () => void;
  setError: (message: string) => void;
}

const EventDefinitionList: FunctionComponent<EventDefinitionListProps> = ({
  definitions,
  deleting,
  editingName,
  onDelete,
  onEdit,
  onResetEditor,
  setError,
}) => {
  const [pendingDeleteName, setPendingDeleteName] = useState<string | null>(null);

  const handleDelete = async (eventName: string): Promise<void> => {
    try {
      await onDelete(eventName);
      if (editingName === eventName) onResetEditor();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete event definition.');
    } finally {
      setPendingDeleteName(null);
    }
  };

  return (
    <div className="space-y-3 border-t pt-4">
      <h4 className="text-sm font-medium">Existing definitions</h4>
      {definitions.length === EMPTY_COUNT ? (
        <p className="text-xs text-muted-foreground">No event definitions yet.</p>
      ) : (
        <div className="space-y-2">
          {definitions.map((definition) => (
            <div key={definition.id} className="flex flex-wrap items-center justify-between gap-3 rounded-md border p-3">
              <div>
                <p className="text-sm font-medium">{definition.name}</p>
                <p className="text-xs text-muted-foreground">
                  {definition.fields.length} field
                  {definition.fields.length === MIN_FIELD_LENGTH ? '' : 's'}
                </p>
              </div>
              <div className="flex gap-2">
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  onClick={() => {
                    setPendingDeleteName(null);
                    onEdit(definition);
                  }}
                >
                  Edit
                </Button>
                {pendingDeleteName === definition.name ? (
                  <>
                    <Button
                      type="button"
                      variant="destructive"
                      size="sm"
                      onClick={() => void handleDelete(definition.name)}
                      disabled={deleting}
                    >
                      Confirm delete
                    </Button>
                    <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      onClick={() => setPendingDeleteName(null)}
                      disabled={deleting}
                    >
                      Cancel
                    </Button>
                  </>
                ) : (
                  <Button
                    type="button"
                    variant="destructive"
                    size="sm"
                    onClick={() => setPendingDeleteName(definition.name)}
                    disabled={deleting}
                  >
                    Delete
                  </Button>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default EventDefinitionList;
