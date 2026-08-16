import type { TypedDocumentNode } from '@apollo/client';
import {
  type FragmentType,
  useFragment as getFragmentData,
} from '@/shared/api/generated/fragment-masking';

/**
 * Unmasks data known to contain a complete fragment.
 *
 * GraphQL Code Generator's partial-fragment overload is selected too broadly by
 * TypeScript 7. Keeping the compatibility cast here preserves fragment masking
 * at call sites without weakening the generated schema types.
 */
export function readFragment<TData>(
  document: TypedDocumentNode<TData, unknown>,
  fragment: FragmentType<TypedDocumentNode<TData, unknown>>
): TData;
export function readFragment<TData>(
  document: TypedDocumentNode<TData, unknown>,
  fragments: FragmentType<TypedDocumentNode<TData, unknown>>[]
): TData[];
export function readFragment<TData>(
  document: TypedDocumentNode<TData, unknown>,
  fragment:
    | FragmentType<TypedDocumentNode<TData, unknown>>
    | FragmentType<TypedDocumentNode<TData, unknown>>[]
): TData | TData[] {
  return getFragmentData(document, fragment) as TData | TData[];
}
