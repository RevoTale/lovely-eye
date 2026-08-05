import { useMemo, useRef, useState } from 'react';
import type { EventDefinitionFieldInput } from '@/gql/graphql';
import type {
  DraftEventDefinitionFieldInput,
  EventDefinitionEditorState,
  EventDefinitionItem,
} from './types';
import {
  DEFAULT_MAX_LENGTH,
  buildEventSnippet,
  definitionFieldsToInput,
  normalizeDraftFields,
} from './event-definition-utils';

interface UseEventDefinitionEditorInput {
  onDelete: (name: string) => Promise<void>;
  onSave: (input: { name: string; fields: EventDefinitionFieldInput[] }) => Promise<void>;
}

export function useEventDefinitionEditor({ onDelete, onSave }: UseEventDefinitionEditorInput) {
  const [draftName, setDraftName] = useState('');
  const [draftFields, setDraftFields] = useState<DraftEventDefinitionFieldInput[]>([]);
  const [originalName, setOriginalName] = useState<string | null>(null);
  const [editorOpen, setEditorOpen] = useState(false);
  const [showSnippet, setShowSnippet] = useState(false);
  const [error, setError] = useState('');
  const nextDraftFieldIdRef = useRef(1);
  const hasOriginalName = originalName !== null && originalName !== '';
  const eventSnippet = useMemo(
    () => buildEventSnippet(draftName, draftFields),
    [draftFields, draftName]
  );

  const resetDraft = (): void => {
    setDraftName('');
    setDraftFields([]);
    setOriginalName(null);
    setEditorOpen(false);
    setShowSnippet(false);
    setError('');
  };

  const startEdit = (definition: EventDefinitionItem): void => {
    setDraftName(definition.name);
    setDraftFields(definitionFieldsToInput(definition).map(toDraftField));
    setOriginalName(definition.name);
    setEditorOpen(true);
    setShowSnippet(false);
    setError('');
  };

  const handleAddField = (): void => {
    setDraftFields((prev) => [
      ...prev,
      toDraftField({ key: '', type: 'STRING', required: false, maxLength: DEFAULT_MAX_LENGTH }),
    ]);
  };

  const handleSave = async (): Promise<void> => {
    const trimmedName = draftName.trim();
    const normalizedFields = normalizeDraftFields(draftFields);
    const validationError = validateDefinition(trimmedName, normalizedFields);
    if (validationError !== '') {
      setError(validationError);
      return;
    }
    try {
      await onSave({ name: trimmedName, fields: normalizedFields });
      if (originalName !== null && originalName !== trimmedName) {
        await onDelete(originalName);
      }
      resetDraft();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save event definition.');
    }
  };

  const state: EventDefinitionEditorState = {
    draftFields,
    draftName,
    error,
    hasOriginalName,
    originalName,
    showSnippet,
  };
  return {
    editorOpen,
    eventSnippet,
    handleAddField,
    handleSave,
    resetDraft,
    setDraftFields,
    setDraftName,
    setEditorOpen,
    setError,
    setShowSnippet,
    startEdit,
    state,
  };

  function toDraftField(field: EventDefinitionFieldInput): DraftEventDefinitionFieldInput {
    const draftId = String(nextDraftFieldIdRef.current);
    nextDraftFieldIdRef.current += 1;
    return { ...field, draftId };
  }
}

function validateDefinition(name: string, fields: EventDefinitionFieldInput[]): string {
  if (name === '') return 'Event name is required.';
  if (fields.some((field) => field.key === '')) return 'Field key cannot be empty.';
  const keys = new Set<string>();
  for (const field of fields) {
    if (keys.has(field.key)) return 'Field keys must be unique.';
    keys.add(field.key);
  }
  return '';
}
