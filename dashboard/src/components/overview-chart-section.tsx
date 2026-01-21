import React, { useState, useCallback, useEffect, useMemo } from 'react';
import { useQuery } from '@apollo/client/react';
import { Card, CardContent, CardHeader, CardTitle, ChartContainer, ChartLegend, ChartLegendContent, ChartTooltip, ChartTooltipContent, Select, SelectContent, SelectItem, SelectTrigger, SelectValue, Skeleton, type ChartConfig } from '@/components/ui';
import { Area, AreaChart, CartesianGrid, XAxis, YAxis } from 'recharts';
import { TrendingUp } from 'lucide-react';
import { ChartDataDocument, type DailyStats } from '@/gql/graphql';

interface OverviewChartSectionProps {
  siteId: string;
  dateRange: { from: Date; to: Date } | null;
  filter: Record<string, string[]> | null;
  bucket: 'daily' | 'hourly';
  onBucketChange: (bucket: 'daily' | 'hourly') => void;
}

const BATCH_SIZE = 10;
const INITIAL_OFFSET = 0;

const EMPTY_COUNT = 0;
const CHART_MARGIN_TOP = 10;
const CHART_MARGIN_RIGHT = 10;
const CHART_MARGIN_LEFT = 0;
const CHART_MARGIN_BOTTOM = 0;
const CHART_MARGIN = {
  top: CHART_MARGIN_TOP,
  right: CHART_MARGIN_RIGHT,
  left: CHART_MARGIN_LEFT,
  bottom: CHART_MARGIN_BOTTOM,
};
const TICK_MARGIN = 8;

export function OverviewChartSection({ siteId, dateRange, filter, bucket, onBucketChange }: OverviewChartSectionProps): React.JSX.Element | null {
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

  const formatLabel = (value: string): string => {
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return value;
    if (bucket === 'hourly') {
      return date.toLocaleString('en-US', { month: 'short', day: 'numeric', hour: 'numeric' });
    }
    return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
  };

  const chartData = loadedData.map(stat => ({
    date: formatLabel(stat.date),
    visitors: stat.visitors,
    pageViews: stat.pageViews,
    sessions: stat.sessions,
  }));

  if (loading && loadedData.length === EMPTY_COUNT) {
    return (
      <Card className="hover:shadow-md transition-shadow">
        <CardHeader>
          <Skeleton className="h-8 w-48" />
        </CardHeader>
        <CardContent>
          <Skeleton className="h-[300px] w-full" />
        </CardContent>
      </Card>
    );
  }

  if (chartData.length === EMPTY_COUNT) {
    return null;
  }

  return (
    <Card className="hover:shadow-md transition-shadow">
      <CardHeader className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <CardTitle className="flex items-center gap-2">
          <div className="h-8 w-8 rounded-lg bg-primary/10 flex items-center justify-center">
            <TrendingUp className="h-4 w-4 text-primary" />
          </div>
          Analytics Overview
        </CardTitle>
        <div className="flex items-center gap-2">
          <span className="text-xs text-muted-foreground">Granularity</span>
          <Select value={bucket} onValueChange={(value) => {
            if (value === 'daily' || value === 'hourly') {
              onBucketChange(value);
            }
          }}>
            <SelectTrigger className="h-8 w-[140px]">
              <SelectValue placeholder="Daily" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="daily">Daily</SelectItem>
              <SelectItem value="hourly">Hourly</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </CardHeader>
      <CardContent>
        <ChartContainer
          config={{
            visitors: {
              label: 'Visitors',
              color: 'hsl(var(--primary))',
            },
            pageViews: {
              label: 'Page Views',
              color: 'hsl(var(--chart-2))',
            },
            sessions: {
              label: 'Sessions',
              color: 'hsl(var(--chart-3))',
            },
          } satisfies ChartConfig}
          className="h-[300px] w-full"
        >
          <AreaChart
            data={chartData}
            margin={CHART_MARGIN}
          >
            <CartesianGrid strokeDasharray="3 3" className="stroke-muted" />
            <XAxis
              dataKey="date"
              tickLine={false}
              axisLine={false}
              tickMargin={TICK_MARGIN}
              className="text-xs"
            />
            <YAxis
              tickLine={false}
              axisLine={false}
              tickMargin={TICK_MARGIN}
              className="text-xs"
            />
            <ChartTooltip content={<ChartTooltipContent />} />
            <ChartLegend content={<ChartLegendContent />} />
            <Area
              type="monotone"
              dataKey="visitors"
              stackId="1"
              stroke="var(--color-visitors)"
              fill="var(--color-visitors)"
              fillOpacity={0.6}
            />
            <Area
              type="monotone"
              dataKey="pageViews"
              stackId="2"
              stroke="var(--color-pageViews)"
              fill="var(--color-pageViews)"
              fillOpacity={0.6}
            />
            <Area
              type="monotone"
              dataKey="sessions"
              stackId="3"
              stroke="var(--color-sessions)"
              fill="var(--color-sessions)"
              fillOpacity={0.6}
            />
          </AreaChart>
        </ChartContainer>
      </CardContent>
    </Card>
  );
}
