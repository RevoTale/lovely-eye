import { useCallback, useEffect, useMemo, useState } from 'react';
import { useQuery } from '@apollo/client/react';
import {
  ChartDataDocument,
  DailyStatsFieldsFragmentDoc,
  type ChartDataQuery,
  type ChartDataQueryVariables,
  type DailyStatsFieldsFragment,
  type FilterInput,
  type TimeBucket,
} from '@/gql/graphql';
import { useFragment as getFragmentData, type FragmentType } from '@/gql/fragment-masking';
import type { DashboardLoadState } from '@/lib/dashboard-load-state';
import { BATCH_SIZE, buildChartVariables, calculateExpectedCount, calculateProgress } from '@/lib/chart-loader-utils';

const EMPTY_COUNT = 0;

interface UseChartDataLoaderParams {
  siteId: string;
  dateRange: { from: Date; to: Date } | null;
  filter: FilterInput | null;
  bucket: 'daily' | 'hourly';
}

interface ChartDataLoaderResult {
  loadedData: DailyStatsFieldsFragment[];
  state: DashboardLoadState;
  loadingMore: boolean;
  progress: number | null;
  expectedCount: number | null;
}

export function useChartDataLoader({ siteId, dateRange, filter, bucket }: UseChartDataLoaderParams): ChartDataLoaderResult {
  const [loadedData, setLoadedData] = useState<DailyStatsFieldsFragment[]>([]);
  const [isLoadingMore, setIsLoadingMore] = useState(false);
  const [hasMore, setHasMore] = useState(true);
  const [resolvedKey, setResolvedKey] = useState('');
  const requestKey = useMemo(() => JSON.stringify({ siteId, dateRange, filter, bucket }), [siteId, dateRange, filter, bucket]);
  const bucketValue: TimeBucket = bucket === 'daily' ? 'DAILY' : 'HOURLY';
  const variables = useMemo(() => buildChartVariables(siteId, dateRange, filter, bucketValue), [siteId, dateRange, filter, bucketValue]);
  const query = useQuery<ChartDataQuery, ChartDataQueryVariables>(ChartDataDocument, {
    variables,
    fetchPolicy: 'cache-and-network',
    notifyOnNetworkStatusChange: true,
  });

  useEffect(() => {
    if (query.data === undefined) return;
    const dailyStats = query.data.dashboard.dailyStats as Array<FragmentType<typeof DailyStatsFieldsFragmentDoc>>;
    const initialData = getFragmentData(DailyStatsFieldsFragmentDoc, dailyStats);
    setLoadedData(initialData);
    setResolvedKey(requestKey);
    setHasMore(initialData.length === BATCH_SIZE);
    setIsLoadingMore(false);
  }, [query.data, requestKey]);

  const loadNextBatch = useCallback(async () => {
    if (query.data === undefined || query.loading || isLoadingMore || !hasMore || resolvedKey !== requestKey) return;
    if (query.data.dashboard.dailyStats.length < BATCH_SIZE) return;
    setIsLoadingMore(true);
    try {
      const result = await query.fetchMore({ variables: { offset: loadedData.length } });
      if (result.data === undefined) return;
      const nextBatch = getFragmentData(DailyStatsFieldsFragmentDoc, result.data.dashboard.dailyStats as Array<FragmentType<typeof DailyStatsFieldsFragmentDoc>>);
      if (nextBatch.length > EMPTY_COUNT) setLoadedData((prev) => [...prev, ...nextBatch]);
      setHasMore(nextBatch.length === BATCH_SIZE);
    } catch {
      setHasMore(false);
    } finally {
      setIsLoadingMore(false);
    }
  }, [hasMore, isLoadingMore, loadedData.length, query, requestKey, resolvedKey]);

  useEffect(() => {
    if (resolvedKey !== requestKey || query.loading || isLoadingMore || !hasMore) return;
    if (loadedData.length > EMPTY_COUNT && loadedData.length % BATCH_SIZE === EMPTY_COUNT) void loadNextBatch();
  }, [hasMore, isLoadingMore, loadNextBatch, loadedData.length, query.loading, requestKey, resolvedKey]);

  const expectedCount = useMemo(() => calculateExpectedCount(dateRange, bucket), [dateRange, bucket]);
  const progress = useMemo(() => calculateProgress(expectedCount, hasMore, loadedData.length), [expectedCount, hasMore, loadedData.length]);
  const state: DashboardLoadState = query.loading && loadedData.length === EMPTY_COUNT ? 'initial' : query.loading ? 'refreshing' : 'ready';

  return { loadedData, state, loadingMore: isLoadingMore, progress, expectedCount };
}
