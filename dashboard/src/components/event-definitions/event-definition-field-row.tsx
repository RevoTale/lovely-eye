import type { Dispatch, FunctionComponent, SetStateAction } from 'react';
import { Button, Checkbox, Input, Label } from '@/components/ui';
import type { EventDefinitionFieldInput } from '@/gql/graphql';
import type { DraftEventDefinitionFieldInput } from './types';
import {
  DEFAULT_MAX_LENGTH,
  FIELD_TYPES,
  MIN_FIELD_LENGTH,
  isEventFieldType,
} from './event-definition-utils';

interface EventDefinitionFieldRowProps {
  field: DraftEventDefinitionFieldInput;
  index: number;
  setDraftFields: Dispatch<SetStateAction<DraftEventDefinitionFieldInput[]>>;
}

const EventDefinitionFieldRow: FunctionComponent<EventDefinitionFieldRowProps> = ({
  field,
  index,
  setDraftFields,
}) => (
  <div className="grid items-center gap-2 md:grid-cols-[2fr_1fr_1fr_1fr_auto]">
    <Input
      placeholder="error_code"
      value={field.key}
      onChange={(event) => updateField(setDraftFields, index, { key: event.currentTarget.value })}
    />
    <select
      value={field.type}
      className="h-10 w-full rounded-md border border-input bg-background px-3 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
      onChange={(event) => {
        const { value } = event.currentTarget;
        if (isEventFieldType(value)) updateField(setDraftFields, index, { type: value });
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
      onChange={(event) => {
        const parsed = Number(event.currentTarget.value);
        updateField(setDraftFields, index, {
          maxLength: Number.isNaN(parsed) ? DEFAULT_MAX_LENGTH : parsed,
        });
      }}
    />
    <div className="flex items-center gap-2">
      <Checkbox
        id={`required-${index}`}
        checked={field.required}
        onCheckedChange={(value) => updateField(setDraftFields, index, { required: value === true })}
      />
      <Label htmlFor={`required-${index}`} className="text-xs">
        Required
      </Label>
    </div>
    <Button
      type="button"
      variant="outline"
      size="sm"
      onClick={() => setDraftFields((prev) => prev.filter((_, idx) => idx !== index))}
    >
      Remove
    </Button>
  </div>
);

const updateField = (
  setDraftFields: Dispatch<SetStateAction<DraftEventDefinitionFieldInput[]>>,
  index: number,
  patch: Partial<EventDefinitionFieldInput>
): void => {
  setDraftFields((prev) => prev.map((item, idx) => (idx === index ? { ...item, ...patch } : item)));
};

export default EventDefinitionFieldRow;
