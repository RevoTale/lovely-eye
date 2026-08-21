import type { ApolloClient } from '@apollo/client';
import { useQuery } from '@apollo/client/react';
import {
  type SiteSummaryFieldsFragment,
  SiteSummaryFieldsFragmentDoc,
  SitesDocument,
  type SitesQuery,
} from '@/shared/api/generated/graphql';
import { readFragment } from '@/shared/api/read-fragment';

export const SITES_PAGING = { limit: 100, offset: 0 } as const;

export interface SitesResult {
  sites: SiteSummaryFieldsFragment[];
  isInitialLoading: boolean;
  error: Error | undefined;
}

export async function loadSites(client: ApolloClient): Promise<SiteSummaryFieldsFragment[]> {
  const { data } = await client.query({
    query: SitesDocument,
    variables: { paging: SITES_PAGING },
    fetchPolicy: 'cache-first',
  });
  if (data === undefined) throw new Error('Sites query completed without data.');
  return readSites(data.sites);
}

export function useSites(): SitesResult {
  const { data, error, loading } = useQuery(SitesDocument, {
    variables: { paging: SITES_PAGING },
  });
  const sites = readSites(data?.sites ?? []);
  return { sites, isInitialLoading: loading && data === undefined, error };
}

const readSites = (sites: SitesQuery['sites']): SiteSummaryFieldsFragment[] =>
  sites.map((site) => readFragment(SiteSummaryFieldsFragmentDoc, site));
