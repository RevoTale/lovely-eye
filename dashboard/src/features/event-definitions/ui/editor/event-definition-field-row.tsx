import type { Dispatch, FunctionComponent, SetStateAction } from 'react';
import type { EventDefinitionFieldInput } from '@/shared/api/generated/graphql';
import { Button } from '@/shared/ui/button';
import { Checkbox } from '@/shared/ui/checkbox';
import { Input } from '@/shared/ui/input';
import { Label } from '@/shared/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui/select';
import {
  DEFAULT_MAX_LENGTH,
  FIELD_TYPES,
  isEventFieldType,
  MIN_FIELD_LENGTH,
} from './event-definition-utils';
import type { DraftEventDefinitionFieldInput } from './types';

interface EventDefinitionFieldRowProps {
  field: DraftEventDefinitionFieldInput;
  index: number;
  setDraftFields: Dispatch<SetStateAction<DraftEventDefinitionFieldInput[]>>;
}

const EventDefinitionFieldRow: FunctionComponent<EventDefinitionFieldRowProps> = ({
  field,
  index,
  setDraftFields,
}) => {
  const keyId = `event-field-key-${field.draftId}`;
  const typeId = `event-field-type-${field.draftId}`;
  const maxLengthId = `event-field-max-length-${field.draftId}`;
  const requiredId = `event-field-required-${field.draftId}`;

  return (
    <div className='grid grid-cols-1 items-center gap-2 md:grid-cols-[2fr_1fr_1fr_1fr_auto]'>
      <div>
        <Label className='sr-only' htmlFor={keyId}>
          Field key {index + 1}
        </Label>
        <Input
          id={keyId}
          placeholder='error_code'
          value={field.key}
          onChange={(event) =>
            updateField(setDraftFields, index, { key: event.currentTarget.value })
          }
        />
      </div>
      <div>
        <Label className='sr-only' htmlFor={typeId}>
          Field type {index + 1}
        </Label>
        <Select
          value={field.type}
          onValueChange={(value) => {
            if (isEventFieldType(value)) updateField(setDraftFields, index, { type: value });
          }}
        >
          <SelectTrigger id={typeId} className='w-full'>
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            {FIELD_TYPES.map((option) => (
              <SelectItem key={option.value} value={option.value}>
                {option.label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>
      <div>
        <Label className='sr-only' htmlFor={maxLengthId}>
          Maximum length {index + 1}
        </Label>
        <Input
          id={maxLengthId}
          type='number'
          min={MIN_FIELD_LENGTH}
          placeholder='Max length'
          value={field.maxLength ?? DEFAULT_MAX_LENGTH}
          onChange={(event) => {
            const parsed = Number(event.currentTarget.value);
            updateField(setDraftFields, index, {
              maxLength: Number.isNaN(parsed) ? DEFAULT_MAX_LENGTH : parsed,
            });
          }}
        />
      </div>
      <div className='flex items-center gap-2'>
        <Checkbox
          id={requiredId}
          checked={field.required}
          onCheckedChange={(value) =>
            updateField(setDraftFields, index, { required: value === true })
          }
        />
        <Label htmlFor={requiredId} className='text-xs'>
          Required
        </Label>
      </div>
      <Button
        type='button'
        variant='outline'
        size='sm'
        onClick={() => setDraftFields((prev) => prev.filter((_, idx) => idx !== index))}
      >
        Remove
      </Button>
    </div>
  );
};

const updateField = (
  setDraftFields: Dispatch<SetStateAction<DraftEventDefinitionFieldInput[]>>,
  index: number,
  patch: Partial<EventDefinitionFieldInput>
): void => {
  setDraftFields((prev) => prev.map((item, idx) => (idx === index ? { ...item, ...patch } : item)));
};

export default EventDefinitionFieldRow;
