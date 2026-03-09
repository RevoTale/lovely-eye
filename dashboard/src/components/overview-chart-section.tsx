import { TrendingUp } from 'lucide-react';
import type { FunctionComponent } from 'react';
import { useMemo } from 'react';
import { Area, AreaChart, CartesianGrid, XAxis, YAxis } from 'recharts';
import DashboardCardState from '@/components/dashboard-card-state';
import { Card, CardContent, CardHeader, CardTitle, ChartContainer, ChartLegend, ChartLegendContent, ChartTooltip, ChartTooltipContent, Progress, Select, SelectContent, SelectItem, SelectTrigger, SelectValue, Skeleton } from '@/components/ui';
import type { FilterInput } from '@/gql/graphql';
import { useChartDataLoader } from '@/hooks/use-chart-data-loader';
import { buildTimeFormatter, parseChartDate } from '@/lib/chart-date';
import { CHART_CONFIG, CHART_MARGIN, TICK_MARGIN } from '@/lib/chart-config';

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
  const formatters = useMemo(() => ({ daily: buildTimeFormatter('daily'), hourly: buildTimeFormatter('hourly') }), []);
  const chartData = useMemo(() => loadedData.map((stat) => {
    const timestamp = parseChartDate(stat.date);
    return timestamp === null ? null : { timestamp, visitors: stat.visitors, pageViews: stat.pageViews, sessions: stat.sessions };
  }).filter((item): item is NonNullable<typeof item> => item !== null).sort((a, b) => a.timestamp - b.timestamp), [loadedData]);
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
        <DashboardCardState state={state} overlayLabel={loadingMore ? 'Loading more' : 'Refreshing chart'} skeleton={<Skeleton className="h-[300px] w-full" />} className="min-h-[300px]">
          <ChartContainer config={CHART_CONFIG} className="h-[300px] w-full">
            <AreaChart data={chartData} margin={CHART_MARGIN}>
              <CartesianGrid strokeDasharray="3 3" className="stroke-muted" />
              <XAxis dataKey="timestamp" type="number" scale="time" domain={['dataMin', 'dataMax']} tickFormatter={(value) => formatters[bucket].format(Number(value))} tickLine={false} axisLine={false} tickMargin={TICK_MARGIN} className="text-xs" />
              <YAxis tickLine={false} axisLine={false} tickMargin={TICK_MARGIN} className="text-xs" />
              <ChartTooltip content={<ChartTooltipContent labelFormatter={(value) => typeof value === 'number' ? formatters[bucket].format(value) : String(value)} />} />
              <ChartLegend content={<ChartLegendContent />} />
              <Area type="monotone" dataKey="visitors" stackId="1" stroke="var(--color-visitors)" fill="var(--color-visitors)" fillOpacity={0.6} />
              <Area type="monotone" dataKey="pageViews" stackId="2" stroke="var(--color-pageViews)" fill="var(--color-pageViews)" fillOpacity={0.6} />
              <Area type="monotone" dataKey="sessions" stackId="3" stroke="var(--color-sessions)" fill="var(--color-sessions)" fillOpacity={0.6} />
            </AreaChart>
          </ChartContainer>
        </DashboardCardState>
      </CardContent>
    </Card>
  );
};

export default OverviewChartSection;
