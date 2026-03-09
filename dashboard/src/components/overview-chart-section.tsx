import { TrendingUp } from 'lucide-react';
import type { FunctionComponent } from 'react';
import { useMemo } from 'react';
import DashboardCardState from '@/components/dashboard-card-state';
import OverviewChartPlot from '@/components/overview-chart/overview-chart-plot';
import { Card, CardContent, CardHeader, CardTitle, Progress, Select, SelectContent, SelectItem, SelectTrigger, SelectValue, Skeleton } from '@/components/ui';
import type { FilterInput } from '@/gql/graphql';
import { buildOverviewChartData } from '@/components/overview-chart/overview-chart-data';
import { useChartDataLoader } from '@/hooks/use-chart-data-loader';

interface OverviewChartSectionProps {
  siteId: string;
  dateRange: { from: Date; to: Date } | null;
  filter: FilterInput | null;
  bucket: 'daily' | 'hourly';
  onBucketChange: (bucket: 'daily' | 'hourly') => void;
}

const EMPTY_COUNT = 0;
const PROGRESS_MIN = 0;

const OverviewChartSection: FunctionComponent<OverviewChartSectionProps> = ({ siteId, dateRange, filter, bucket, onBucketChange }) => {
  const { loadedData, state, loadingMore, progress, expectedCount } = useChartDataLoader({ siteId, dateRange, filter, bucket });
  const chartData = useMemo(() => buildOverviewChartData(loadedData), [loadedData]);
  const showProgress = (state === 'refreshing' && chartData.length > EMPTY_COUNT) || loadingMore;
  const progressLabel = expectedCount === null ? `Loaded ${loadedData.length.toLocaleString()} points` : `Loaded ${Math.min(loadedData.length, expectedCount).toLocaleString()} of ${expectedCount.toLocaleString()} points`;

  if (chartData.length === EMPTY_COUNT && state === 'ready') return null;

  return (
    <Card className="transition-shadow hover:shadow-md">
      <CardHeader className="flex flex-col gap-4">
        <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <CardTitle className="flex items-center gap-2"><div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10"><TrendingUp className="h-4 w-4 text-primary" /></div>Analytics Overview</CardTitle>
          <div className="flex items-center gap-2"><span className="text-xs text-muted-foreground">Granularity</span><Select value={bucket} onValueChange={(value) => { if (value === 'daily' || value === 'hourly') onBucketChange(value); }}><SelectTrigger className="h-8 w-[140px]"><SelectValue placeholder="Daily" /></SelectTrigger><SelectContent><SelectItem value="daily">Daily</SelectItem><SelectItem value="hourly">Hourly</SelectItem></SelectContent></Select></div>
        </div>
        {showProgress ? <div className="space-y-2"><div className="flex items-center justify-between text-xs text-muted-foreground"><span>{loadingMore ? 'Loading more data' : 'Refreshing chart data'}</span><span>{progressLabel}</span></div><Progress value={progress ?? PROGRESS_MIN} className="h-2" /></div> : null}
      </CardHeader>
      <CardContent>
        <DashboardCardState state={state} overlayLabel={loadingMore ? 'Loading more' : 'Refreshing chart'} skeleton={<Skeleton className="h-[clamp(240px,44vw,320px)] w-full" />} className="min-h-[clamp(240px,44vw,320px)]">
          <OverviewChartPlot bucket={bucket} data={chartData} />
        </DashboardCardState>
      </CardContent>
    </Card>
  );
};

export default OverviewChartSection;
