import {
  EventDefinitionFieldFieldsFragmentDoc,
  type EventDefinitionFieldInput,
  type EventDefinitionFieldsFragment,
  type EventFieldType,
} from '@/shared/api/generated/graphql';
import { readFragment } from '@/shared/api/read-fragment';

export const DEFAULT_MAX_LENGTH = 500;
export const EMPTY_COUNT = 0;
export const MIN_FIELD_LENGTH = 1;

const EXAMPLE_NUMBER = 42;
const FIRST_INDEX_OFFSET = 1;
const INLINE_EXAMPLE_MAX = 20;

interface EventFieldTypeOption {
  label: string;
  value: EventFieldType;
}

export const FIELD_TYPES: EventFieldTypeOption[] = [
  { label: 'String', value: 'STRING' },
  { label: 'Int', value: 'INT' },
  { label: 'Boolean', value: 'BOOLEAN' },
];

export function isEventFieldType(value: string): value is EventFieldType {
  return value === 'STRING' || value === 'INT' || value === 'BOOLEAN';
}

export function definitionFieldsToInput(
  definition: EventDefinitionFieldsFragment
): EventDefinitionFieldInput[] {
  return definition.fields.map((field) => {
    const fieldData = readFragment(EventDefinitionFieldFieldsFragmentDoc, field);
    return {
      key: fieldData.key,
      type: fieldData.type,
      required: fieldData.required,
      maxLength: fieldData.maxLength,
    };
  });
}

export function normalizeDraftFields(
  fields: EventDefinitionFieldInput[]
): EventDefinitionFieldInput[] {
  return fields.map((field) => ({
    key: field.key.trim(),
    type: field.type,
    required: field.required,
    maxLength: field.maxLength ?? DEFAULT_MAX_LENGTH,
  }));
}

export function buildEventSnippet(name: string, fields: EventDefinitionFieldInput[]): string {
  const trimmedEventName = name.trim();
  const eventName = trimmedEventName === '' ? 'event_name' : trimmedEventName;
  const properties = buildSnippetProperties(fields);
  return `window.lovelyEye?.track({
  name: '${eventName}',
  properties: ${properties},
});`;
}

function buildSnippetProperties(fields: EventDefinitionFieldInput[]): string {
  const fieldEntries = fields.map((field, index) => buildSnippetField(field, index));
  return fieldEntries.length > EMPTY_COUNT ? `{\n  ${fieldEntries.join(',\n  ')}\n}` : '{}';
}

function buildSnippetField(field: EventDefinitionFieldInput, index: number): string {
  const trimmedKey = field.key.trim();
  const key = trimmedKey === '' ? `field_${index + FIRST_INDEX_OFFSET}` : trimmedKey;
  switch (field.type) {
    case 'INT':
      return `${key}: ${EXAMPLE_NUMBER}`;
    case 'BOOLEAN':
      return `${key}: true`;
    default:
      return `${key}: ${stringExample(field.maxLength)}`;
  }
}

function stringExample(maxLength: number | null | undefined): string {
  if (maxLength !== null && maxLength !== undefined && maxLength > EMPTY_COUNT) {
    if (maxLength <= INLINE_EXAMPLE_MAX) {
      return `'${'a'.repeat(maxLength)}'`;
    }
    return `'a'.repeat(${maxLength})`;
  }
  return "'example'";
}
