import type { Dispatch, FunctionComponent, SetStateAction } from 'react';
import { Button } from '@/shared/ui/button';
import { Label } from '@/shared/ui/label';
import EventDefinitionFieldRow from './event-definition-field-row';
import { EMPTY_COUNT } from './event-definition-utils';
import type { DraftEventDefinitionFieldInput } from './types';

interface EventDefinitionFieldsProps {
  fields: DraftEventDefinitionFieldInput[];
  onAddField: () => void;
  setDraftFields: Dispatch<SetStateAction<DraftEventDefinitionFieldInput[]>>;
}

const EventDefinitionFields: FunctionComponent<EventDefinitionFieldsProps> = ({
  fields,
  onAddField,
  setDraftFields,
}) => (
  <div className='space-y-3'>
    <div className='flex items-center justify-between'>
      <Label>Fields</Label>
      <Button type='button' variant='outline' size='sm' onClick={onAddField}>
        Add field
      </Button>
    </div>
    {fields.length === EMPTY_COUNT ? (
      <p className='text-xs text-muted-foreground'>
        No fields defined. Events will be stored without metadata.
      </p>
    ) : (
      <div className='space-y-3'>
        {fields.map((field, index) => (
          <EventDefinitionFieldRow
            key={field.draftId}
            field={field}
            index={index}
            setDraftFields={setDraftFields}
          />
        ))}
      </div>
    )}
  </div>
);

export default EventDefinitionFields;
