import type {
  EventDefinitionFieldInput,
  EventDefinitionFieldsFragment,
  EventDefinitionInput,
} from '@/gql/graphql';
import type { FragmentType } from '@/gql/fragment-masking';
import type { EventDefinitionFieldsFragmentDoc } from '@/gql/graphql';

export interface EventDefinitionsCardProps {
  definitions: Array<FragmentType<typeof EventDefinitionFieldsFragmentDoc>>;
  saving: boolean;
  deleting: boolean;
  onSave: (input: EventDefinitionInput) => Promise<void>;
  onDelete: (name: string) => Promise<void>;
}

export interface DraftEventDefinitionFieldInput extends EventDefinitionFieldInput {
  draftId: string;
}

export interface EventDefinitionEditorState {
  draftFields: DraftEventDefinitionFieldInput[];
  draftName: string;
  error: string;
  hasOriginalName: boolean;
  originalName: string | null;
  showSnippet: boolean;
}

export type EventDefinitionItem = EventDefinitionFieldsFragment;
