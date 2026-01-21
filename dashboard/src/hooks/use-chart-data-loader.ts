import { useState, useCallback, useEffect, useMemo } from 'react';
import { useQuery } from '@apollo/client/react';
import { ChartDataDocument, type DailyStats } from '@/gql/graphql';

const BATCH_SIZE = 10;
const INITIAL_OFFSET = 0;
const EMPTY_COUNT = 0;

interface UseChartDataLoaderParams {
  siteId: string;
  dateRange: { from: Date; to: Date } | null;
  filter: Record<string, string[]> | null;
  bucket: 'daily' | 'hourly';
}

interface ChartDataLoaderResult {
  loadedData: DailyStats[];
  loading: boolean;
}

export function useChartDataLoader({ siteId, dateRange, filter, bucket }: UseChartDataLoaderParams): ChartDataLoaderResult {
  const [loadedData, setLoadedData] = useState<DailyStats[]>([]);
  const [isLoadingMore, setIsLoadingMore] = useState(false);

  const filterKey = useMemo(
    () => JSON.stringify({ siteId, dateRange, filter, bucket }),
    [siteId, dateRange, filter, bucket]
  );

  const bucketValue = bucket === 'daily' ? 'DAILY' : 'HOURLY';

  const { data, loading, fetchMore } = useQuery(ChartDataDocument, {
    variables: {
      siteId,
      dateRange: dateRange === null ? undefined : { from: dateRange.from.toISOString(), to: dateRange.to.toISOString() },
      filter: filter === null ? undefined : {
        referrer: filter['referrer'],
        device: filter['device'],
        page: filter['page'],
        country: filter['country'],
      },
      bucket: bucketValue,
      limit: BATCH_SIZE,
      offset: INITIAL_OFFSET,
    },
  });

  useEffect(() => {
    if (data?.dashboard.dailyStats !== undefined) {
      setLoadedData(data.dashboard.dailyStats);
    }
  }, [data]);

  useEffect(() => {
    setLoadedData([]);
    setIsLoadingMore(false);
  }, [filterKey]);

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
        setLoadedData(prev => [...prev, ...newData]);
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
