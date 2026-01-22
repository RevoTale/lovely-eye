import React, { useMemo, useState } from 'react';
import { Button, Card, CardContent, CardDescription, CardHeader, CardTitle, Checkbox, Input, Label } from '@/components/ui';
import type { EventDefinition, EventDefinitionInput, EventDefinitionFieldInput, EventFieldType } from '@/gql/graphql';

const DEFAULT_MAX_LENGTH = 500;
const EMPTY_COUNT = 0;
const EXAMPLE_NUMBER = 42;
const FIRST_INDEX_OFFSET = 1;
const INLINE_EXAMPLE_MAX = 20;
const MIN_FIELD_LENGTH = 1;

interface EventFieldTypeOption {
  label: string;
  value: EventFieldType;
}

const FIELD_TYPES: EventFieldTypeOption[] = [
  { label: 'String', value: 'STRING' },
  { label: 'Int', value: 'INT' },
  { label: 'Boolean', value: 'BOOLEAN' },
];

function isEventFieldType(value: string): value is EventFieldType {
  return value === 'STRING' || value === 'INT' || value === 'BOOLEAN';
}

interface EventDefinitionsCardProps {
  definitions: EventDefinition[];
  saving: boolean;
  deleting: boolean;
  onSave: (input: EventDefinitionInput) => Promise<void>;
  onDelete: (name: string) => Promise<void>;
}

export function EventDefinitionsCard({
  definitions,
  saving,
  deleting,
  onSave,
  onDelete,
}: EventDefinitionsCardProps): React.JSX.Element {
  const [draftName, setDraftName] = useState('');
  const [draftFields, setDraftFields] = useState<EventDefinitionFieldInput[]>([]);
  const [originalName, setOriginalName] = useState<string | null>(null);
  const [editorOpen, setEditorOpen] = useState(false);
  const [showSnippet, setShowSnippet] = useState(false);
  const [error, setError] = useState('');
  const [pendingDeleteName, setPendingDeleteName] = useState<string | null>(null);

  const hasOriginalName = originalName !== null && originalName !== '';

  const sortedDefinitions = useMemo(
    () => [...definitions].sort((a, b) => a.name.localeCompare(b.name)),
    [definitions]
  );

  const resetDraft = (): void => {
    setDraftName('');
    setDraftFields([]);
    setOriginalName(null);
    setEditorOpen(false);
    setShowSnippet(false);
    setError('');
  };

  const startEdit = (definition: EventDefinition): void => {
    setDraftName(definition.name);
    setDraftFields(
      definition.fields.map((field) => ({
        key: field.key,
        type: field.type,
        required: field.required,
        maxLength: field.maxLength,
      }))
    );
    setOriginalName(definition.name);
    setEditorOpen(true);
    setShowSnippet(false);
    setError('');
  };

  const eventSnippet = useMemo(() => {
    const trimmedEventName = draftName.trim();
    const eventName = trimmedEventName === '' ? 'event_name' : trimmedEventName;
    const fieldEntries = draftFields.map((field, index) => {
      const trimmedKey = field.key.trim();
      const key = trimmedKey === '' ? `field_${index + FIRST_INDEX_OFFSET}` : trimmedKey;
      switch (field.type) {
        case 'INT':
          return `${key}: ${EXAMPLE_NUMBER}`;
        case 'BOOLEAN':
          return `${key}: true`;
        default:
          if (field.maxLength !== null && field.maxLength !== undefined && field.maxLength > EMPTY_COUNT) {
            if (field.maxLength <= INLINE_EXAMPLE_MAX) {
              return `${key}: '${'a'.repeat(field.maxLength)}'`;
            }
            return `${key}: 'a'.repeat(${field.maxLength})`;
          }
          return `${key}: 'example'`;
      }
    });
    const properties = fieldEntries.length > EMPTY_COUNT
      ? `{\n  ${fieldEntries.join(',\n  ')}\n}`
      : '{}';
    return `window.lovelyEye?.track('${eventName}', ${properties});`;
  }, [draftFields, draftName]);

  const handleAddField = (): void => {
    setDraftFields((prev) => [
      ...prev,
      { key: '', type: 'STRING', required: false, maxLength: DEFAULT_MAX_LENGTH },
    ]);
  };

  const handleSave = async (): Promise<void> => {
    const trimmedName = draftName.trim();
    if (trimmedName === '') {
      setError('Event name is required.');
      return;
    }

    const normalizedFields = draftFields.map((field) => ({
      key: field.key.trim(),
      type: field.type,
      required: field.required,
      maxLength: field.maxLength ?? DEFAULT_MAX_LENGTH,
    }));

    if (normalizedFields.some((field) => field.key === '')) {
      setError('Field key cannot be empty.');
      return;
    }

    const keySet = new Set<string>();
    for (const field of normalizedFields) {
      if (keySet.has(field.key)) {
        setError('Field keys must be unique.');
        return;
      }
      keySet.add(field.key);
    }

    try {
      await onSave({
        name: trimmedName,
        fields: normalizedFields,
      });
      if (originalName !== null && originalName !== trimmedName) {
        await onDelete(originalName);
      }
      resetDraft();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save event definition.');
    }
  };

  const handleDelete = async (eventName: string): Promise<void> => {
    try {
      await onDelete(eventName);
      if (originalName === eventName) {
        resetDraft();
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete event definition.');
    } finally {
      setPendingDeleteName(null);
    }
  };

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
              if (editorOpen) {
                resetDraft();
                return;
              }
              setEditorOpen(true);
              setError('');
            }}
          >
            {editorOpen ? 'Close editor' : 'New event name'}
          </Button>
        </CardTitle>
        <CardDescription>
          Allowlist event names and metadata fields to keep tracking clean.
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        {editorOpen ? (
          <div className="space-y-6 rounded-lg border bg-muted/30 p-4">
            <div className="flex flex-wrap items-center justify-between gap-3">
              <div>
                <h4 className="text-sm font-medium">
                  {hasOriginalName ? `Editing: ${originalName}` : 'New event name'}
                </h4>
                <p className="text-xs text-muted-foreground">
                  Events not listed here will be ignored.
                </p>
              </div>
              <Button type="button" variant="ghost" size="sm" onClick={resetDraft}>
                Close
              </Button>
            </div>

            {error === '' ? null : (
              <div className="text-sm text-destructive">{error}</div>
            )}

            <div className="space-y-3">
              <Label htmlFor="event-name">Event Name</Label>
              <Input
                id="event-name"
                placeholder="signup_error"
                value={draftName}
                onChange={(e) => {
                  setDraftName(e.target.value);
                }}
              />
            </div>

            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <Label>Fields</Label>
                <Button type="button" variant="outline" size="sm" onClick={handleAddField}>
                  Add field
                </Button>
              </div>

              {draftFields.length === EMPTY_COUNT ? (
                <p className="text-xs text-muted-foreground">
                  No fields defined. Events will be stored without metadata.
                </p>
              ) : (
                <div className="space-y-3">
                  {draftFields.map((field, index) => (
                    <div key={index} className="grid gap-2 md:grid-cols-[2fr_1fr_1fr_1fr_auto] items-center">
                      <Input
                        placeholder="error_code"
                        value={field.key}
                        onChange={(e) => {
                          const { currentTarget } = e;
                          const { value } = currentTarget;
                          setDraftFields((prev) =>
                            prev.map((item, idx) => (idx === index ? { ...item, key: value } : item))
                          );
                        }}
                      />

                      <select
                        value={field.type}
                        className="h-10 w-full rounded-md border border-input bg-background px-3 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                        onChange={(e) => {
                          const { currentTarget } = e;
                          const { value } = currentTarget;
                          if (isEventFieldType(value)) {
                            setDraftFields((prev) =>
                              prev.map((item, idx) => (idx === index ? { ...item, type: value } : item))
                            );
                          }
                        }}
                      >
                        {FIELD_TYPES.map((option) => (
                          <option key={option.value} value={option.value}>
                            {option.label}
                          </option>
                        ))}
                      </select>

                      <Input
                        type="number"
                        min={MIN_FIELD_LENGTH}
                        placeholder="Max length"
                        value={field.maxLength ?? DEFAULT_MAX_LENGTH}
                        onChange={(e) => {
                          const { currentTarget } = e;
                          const { value } = currentTarget;
                          const parsed = Number(value);
                          setDraftFields((prev) =>
                            prev.map((item, idx) =>
                              idx === index
                                ? {
                                    ...item,
                                    maxLength: Number.isNaN(parsed) ? DEFAULT_MAX_LENGTH : parsed,
                                  }
                                : item
                            )
                          );
                        }}
                      />

                      <div className="flex items-center gap-2">
                        <Checkbox
                          id={`required-${index}`}
                          checked={field.required}
                          onCheckedChange={(value) => {
                            setDraftFields((prev) =>
                              prev.map((item, idx) =>
                                idx === index ? { ...item, required: value === true } : item
                              )
                            );
                          }}
                        />
                        <Label htmlFor={`required-${index}`} className="text-xs">
                          Required
                        </Label>
                      </div>

                      <Button
                        type="button"
                        variant="outline"
                        size="sm"
                        onClick={() => {
                          setDraftFields((prev) => prev.filter((_, idx) => idx !== index));
                        }}
                      >
                        Remove
                      </Button>
                    </div>
                  ))}
                </div>
              )}
            </div>

            <div className="space-y-2">
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={() => {
                  setShowSnippet((prev) => !prev);
                }}
              >
                {showSnippet ? 'Hide snippet' : 'Show snippet'}
              </Button>
              <div
                className={`overflow-hidden rounded-md border bg-background transition-[max-height,opacity] duration-300 ease-out ${showSnippet ? 'max-h-64 opacity-100' : 'max-h-0 opacity-0'}`}
              >
                <pre className="p-3 text-xs">
                  <code>{eventSnippet}</code>
                </pre>
              </div>
            </div>

            <div className="flex gap-3">
              <Button
                type="button"
                onClick={() => {
                  void handleSave();
                }}
                disabled={saving}
              >
                {saving ? 'Saving...' : hasOriginalName ? 'Update Definition' : 'Save Definition'}
              </Button>
              <Button type="button" variant="outline" onClick={resetDraft} disabled={saving}>
                Cancel
              </Button>
            </div>
          </div>
        ) : (
          <p className="text-sm text-muted-foreground">
            Create an event definition to allowlist event names and metadata.
          </p>
        )}

        <div className="border-t pt-4 space-y-3">
          <h4 className="text-sm font-medium">Existing definitions</h4>
          {sortedDefinitions.length === EMPTY_COUNT ? (
            <p className="text-xs text-muted-foreground">No event definitions yet.</p>
          ) : (
            <div className="space-y-2">
              {sortedDefinitions.map((definition) => (
                <div key={definition.id} className="flex flex-wrap items-center justify-between gap-3 border rounded-md p-3">
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
                        startEdit(definition);
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
                          onClick={() => {
                            void handleDelete(definition.name);
                          }}
                          disabled={deleting}
                        >
                          Confirm delete
                        </Button>
                        <Button
                          type="button"
                          variant="outline"
                          size="sm"
                          onClick={() => {
                            setPendingDeleteName(null);
                          }}
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
                        onClick={() => {
                          setPendingDeleteName(definition.name);
                        }}
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
      </CardContent>
    </Card>
  );
}
