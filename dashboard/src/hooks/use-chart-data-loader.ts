import { useState, useCallback, useEffect, useMemo } from 'react';
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

const BATCH_SIZE = 10;
const INITIAL_OFFSET = 0;
const EMPTY_COUNT = 0;
const MINUTES_PER_HOUR = 60;
const SECONDS_PER_MINUTE = 60;
const MS_IN_SECOND = 1000;
const HOURS_PER_DAY = 24;
const MS_IN_HOUR = MINUTES_PER_HOUR * SECONDS_PER_MINUTE * MS_IN_SECOND;
const MS_IN_DAY = HOURS_PER_DAY * MS_IN_HOUR;
const COUNT_OFFSET = 1;
const PROGRESS_UNKNOWN = 50;
const PROGRESS_COMPLETE = 100;
const PROGRESS_MIN = 0;

interface UseChartDataLoaderParams {
  siteId: string;
  dateRange: { from: Date; to: Date } | null;
  filter: FilterInput | null;
  bucket: 'daily' | 'hourly';
}

interface ChartDataLoaderResult {
  loadedData: DailyStatsFieldsFragment[];
  loading: boolean;
  loadingMore: boolean;
  progress: number | null;
  expectedCount: number | null;
  hasMore: boolean;
}

export function useChartDataLoader({ siteId, dateRange, filter, bucket }: UseChartDataLoaderParams): ChartDataLoaderResult {
  const [loadedData, setLoadedData] = useState<DailyStatsFieldsFragment[]>([]);
  const [isLoadingMore, setIsLoadingMore] = useState(false);
  const [hasMore, setHasMore] = useState(true);

  const filterKey = useMemo(
    () => JSON.stringify({ siteId, dateRange, filter, bucket }),
    [siteId, dateRange, filter, bucket]
  );

  const bucketValue: TimeBucket = bucket === 'daily' ? 'DAILY' : 'HOURLY';

  const variables = useMemo(() => ({
    siteId,
    dateRange: dateRange === null ? null : { from: dateRange.from.toISOString(), to: dateRange.to.toISOString() },
    filter: filter === null ? null : {
      referrer: filter.referrer ?? null,
      device: filter.device ?? null,
      page: filter.page ?? null,
      country: filter.country ?? null,
      eventDefinitionId: filter.eventDefinitionId ?? null,
      eventName: filter.eventName ?? null,
      eventPath: filter.eventPath ?? null,
    },
    bucket: bucketValue,
    limit: BATCH_SIZE,
    offset: INITIAL_OFFSET,
  }), [siteId, dateRange, filter, bucketValue]);

  const { data, loading, fetchMore } = useQuery<ChartDataQuery, ChartDataQueryVariables>(ChartDataDocument, {
    variables,
    notifyOnNetworkStatusChange: true,
  });

  useEffect(() => {
    setLoadedData([]);
    setIsLoadingMore(false);
    setHasMore(true);
  }, [filterKey]);

  useEffect(() => {
    if (data === undefined) return;
    const { dashboard: { dailyStats } } = data;
    const fragmentStats: Array<FragmentType<typeof DailyStatsFieldsFragmentDoc>> = dailyStats;
    const initialData = getFragmentData(DailyStatsFieldsFragmentDoc, fragmentStats);
    setLoadedData(initialData);
    setHasMore(initialData.length === BATCH_SIZE);
  }, [data]);

  const loadNextBatch = useCallback(async () => {
    if (data === undefined) return;
    const { dashboard: { dailyStats } } = data;
    if (isLoadingMore || loading || !hasMore) return;
    if (dailyStats.length < BATCH_SIZE) return;

    setIsLoadingMore(true);
    try {
      const result = await fetchMore({
        variables: {
          offset: loadedData.length,
        },
      });

      const { data: resultData } = result;
      if (resultData === undefined) return;
      const { dashboard: { dailyStats: newData } } = resultData;
      const fragmentBatch: Array<FragmentType<typeof DailyStatsFieldsFragmentDoc>> = newData;
      const nextBatch = getFragmentData(DailyStatsFieldsFragmentDoc, fragmentBatch);
      if (nextBatch.length > EMPTY_COUNT) {
        setLoadedData(prev => [...prev, ...nextBatch]);
      }
      setHasMore(nextBatch.length === BATCH_SIZE);
    } catch {
      // Silently handle batch loading errors
    } finally {
      setIsLoadingMore(false);
    }
  }, [isLoadingMore, loading, data, fetchMore, loadedData.length, hasMore]);

  useEffect(() => {
    if (!loading && !isLoadingMore && hasMore && loadedData.length > EMPTY_COUNT && loadedData.length % BATCH_SIZE === EMPTY_COUNT) {
      void loadNextBatch();
    }
  }, [loading, isLoadingMore, loadedData.length, loadNextBatch, hasMore]);

  const expectedCount = useMemo(() => {
    if (dateRange === null) return null;
    const fromMs = dateRange.from.getTime();
    const toMs = dateRange.to.getTime();
    if (Number.isNaN(fromMs) || Number.isNaN(toMs) || toMs < fromMs) return null;
    const bucketMs = bucket === 'hourly' ? MS_IN_HOUR : MS_IN_DAY;
    return Math.floor((toMs - fromMs) / bucketMs) + COUNT_OFFSET;
  }, [dateRange, bucket]);

  const progress = useMemo(() => {
    if (expectedCount === null) {
      return hasMore ? PROGRESS_UNKNOWN : loadedData.length > EMPTY_COUNT ? PROGRESS_COMPLETE : null;
    }
    if (expectedCount <= EMPTY_COUNT) return null;
    const percent = Math.min(PROGRESS_COMPLETE, Math.max(PROGRESS_MIN, Math.round((loadedData.length / expectedCount) * PROGRESS_COMPLETE)));
    return percent;
  }, [expectedCount, hasMore, loadedData.length]);

  return { loadedData, loading, loadingMore: isLoadingMore, progress, expectedCount, hasMore };
}
