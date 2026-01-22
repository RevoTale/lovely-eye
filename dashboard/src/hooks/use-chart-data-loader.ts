import { useState, useCallback, useEffect, useMemo } from 'react';
import { useQuery } from '@apollo/client/react';
import { ChartDataDocument, DailyStatsFieldsFragmentDoc, type DailyStatsFieldsFragment, type FilterInput } from '@/gql/graphql';
import { useFragment as getFragmentData } from '@/gql/fragment-masking';

const BATCH_SIZE = 10;
const INITIAL_OFFSET = 0;
const EMPTY_COUNT = 0;

interface UseChartDataLoaderParams {
  siteId: string;
  dateRange: { from: Date; to: Date } | null;
  filter: FilterInput | null;
  bucket: 'daily' | 'hourly';
}

interface ChartDataLoaderResult {
  loadedData: DailyStatsFieldsFragment[];
  loading: boolean;
}

export function useChartDataLoader({ siteId, dateRange, filter, bucket }: UseChartDataLoaderParams): ChartDataLoaderResult {
  const [loadedData, setLoadedData] = useState<DailyStatsFieldsFragment[]>([]);
  const [isLoadingMore, setIsLoadingMore] = useState(false);

  const filterKey = useMemo(
    () => JSON.stringify({ siteId, dateRange, filter, bucket }),
    [siteId, dateRange, filter, bucket]
  );

  const bucketValue = bucket === 'daily' ? 'DAILY' : 'HOURLY';

  const { data, loading, fetchMore } = useQuery(ChartDataDocument, {
    variables: {
      siteId,
      dateRange: dateRange === null ? null : { from: dateRange.from.toISOString(), to: dateRange.to.toISOString() },
      filter: filter === null ? null : {
        referrer: filter.referrer ?? null,
        device: filter.device ?? null,
        page: filter.page ?? null,
        country: filter.country ?? null,
      },
      bucket: bucketValue,
      limit: BATCH_SIZE,
      offset: INITIAL_OFFSET,
    },
  });

  useEffect(() => {
    setLoadedData([]);
    setIsLoadingMore(false);
  }, [filterKey]);

  useEffect(() => {
    if (data?.dashboard.dailyStats !== undefined) {
      setLoadedData(getFragmentData(DailyStatsFieldsFragmentDoc, data.dashboard.dailyStats));
    }
  }, [data]);

  const loadNextBatch = useCallback(async () => {
    if (isLoadingMore || loading || data?.dashboard === undefined) return;

    const { dashboard } = data;
    const { dailyStats } = dashboard;
    if (dailyStats.length < BATCH_SIZE) return;

    setIsLoadingMore(true);
    try {
      const result = await fetchMore({
        variables: {
          offset: loadedData.length,
        },
      });

      const { data: resultData } = result;
      if (resultData?.dashboard === undefined) return;
      const { dashboard: resultDashboard } = resultData;
      const { dailyStats: newData } = resultDashboard;
      if (newData.length > EMPTY_COUNT) {
        setLoadedData(prev => [
          ...prev,
          ...getFragmentData(DailyStatsFieldsFragmentDoc, newData),
        ]);
      }
    } catch {
      // Silently handle batch loading errors
    } finally {
      setIsLoadingMore(false);
    }
  }, [isLoadingMore, loading, data, fetchMore, loadedData.length]);

  useEffect(() => {
    if (!loading && !isLoadingMore && loadedData.length > EMPTY_COUNT && loadedData.length % BATCH_SIZE === EMPTY_COUNT) {
      void loadNextBatch();
    }
  }, [loading, isLoadingMore, loadedData.length, loadNextBatch]);

  return { loadedData, loading };
}
