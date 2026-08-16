import { useQuery } from '@apollo/client/react';
import {
  type SiteSummaryFieldsFragment,
  SiteSummaryFieldsFragmentDoc,
  SitesDocument,
} from '@/shared/api/generated/graphql';
import { readFragment } from '@/shared/api/read-fragment';

export const SITES_PAGING = { limit: 100, offset: 0 } as const;

export interface SitesResult {
  sites: SiteSummaryFieldsFragment[];
  loading: boolean;
  error: Error | undefined;
}

export function useSites(): SitesResult {
  const { data, error, loading } = useQuery(SitesDocument, {
    variables: { paging: SITES_PAGING },
  });
  const sites = (data?.sites ?? []).map((site) => readFragment(SiteSummaryFieldsFragmentDoc, site));
  return { sites, loading, error };
}
