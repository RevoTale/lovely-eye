import type { Dispatch, FunctionComponent, SetStateAction } from 'react';
import { Button } from '@/shared/ui/button';
import { Input } from '@/shared/ui/input';
import { Label } from '@/shared/ui/label';
import EventDefinitionFields from './event-definition-fields';
import type { DraftEventDefinitionFieldInput, EventDefinitionEditorState } from './types';

interface EventDefinitionEditorProps {
  eventSnippet: string;
  saving: boolean;
  state: EventDefinitionEditorState;
  onAddField: () => void;
  onCancel: () => void;
  onSave: () => void;
  setDraftFields: Dispatch<SetStateAction<DraftEventDefinitionFieldInput[]>>;
  setDraftName: Dispatch<SetStateAction<string>>;
  setShowSnippet: Dispatch<SetStateAction<boolean>>;
}

const EventDefinitionEditor: FunctionComponent<EventDefinitionEditorProps> = ({
  eventSnippet,
  onAddField,
  onCancel,
  onSave,
  saving,
  setDraftFields,
  setDraftName,
  setShowSnippet,
  state,
}) => (
  <div className='space-y-6 rounded-lg border bg-muted/30 p-4'>
    <div className='flex flex-wrap items-center justify-between gap-3'>
      <div>
        <h4 className='text-sm font-medium'>
          {state.hasOriginalName ? `Editing: ${state.originalName}` : 'New event name'}
        </h4>
        <p className='text-xs text-muted-foreground'>Events not listed here will be ignored.</p>
      </div>
      <Button type='button' variant='ghost' size='sm' onClick={onCancel}>
        Close
      </Button>
    </div>
    {state.error === '' ? null : <div className='text-sm text-destructive'>{state.error}</div>}
    <div className='space-y-3'>
      <Label htmlFor='event-name'>Event Name</Label>
      <Input
        id='event-name'
        placeholder='signup_error'
        value={state.draftName}
        onChange={(event) => setDraftName(event.target.value)}
      />
    </div>
    <EventDefinitionFields
      fields={state.draftFields}
      onAddField={onAddField}
      setDraftFields={setDraftFields}
    />
    <div className='space-y-2'>
      <Button
        type='button'
        variant='outline'
        size='sm'
        onClick={() => setShowSnippet((prev) => !prev)}
      >
        {state.showSnippet ? 'Hide snippet' : 'Show snippet'}
      </Button>
      <div
        className={`overflow-hidden rounded-md border bg-background transition-[max-height,opacity] duration-300 ease-out ${state.showSnippet ? 'max-h-64 opacity-100' : 'max-h-0 opacity-0'}`}
      >
        <pre className='max-w-full overflow-x-auto p-3 text-xs'>
          <code>{eventSnippet}</code>
        </pre>
      </div>
    </div>
    <div className='flex flex-col gap-3 sm:flex-row'>
      <Button type='button' className='w-full sm:w-auto' onClick={onSave} disabled={saving}>
        {saving ? 'Saving...' : state.hasOriginalName ? 'Update Definition' : 'Save Definition'}
      </Button>
      <Button
        type='button'
        variant='outline'
        className='w-full sm:w-auto'
        onClick={onCancel}
        disabled={saving}
      >
        Cancel
      </Button>
    </div>
  </div>
);

export default EventDefinitionEditor;
