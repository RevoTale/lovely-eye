
import { Card, CardContent, CardHeader, CardTitle, ChartContainer, ChartLegend, ChartLegendContent, ChartTooltip, ChartTooltipContent, Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui';
import { Area, AreaChart, CartesianGrid, XAxis, YAxis } from 'recharts';
import { TrendingUp } from 'lucide-react';
import { CHART_CONFIG, CHART_MARGIN, TICK_MARGIN } from '@/lib/chart-config';
import { ChartSkeleton } from '@/components/chart-skeleton';
import { useChartDataLoader } from '@/hooks/use-chart-data-loader';

interface OverviewChartSectionProps {
  siteId: string;
  dateRange: { from: Date; to: Date } | null;
  filter: Record<string, string[]> | null;
  bucket: 'daily' | 'hourly';
  onBucketChange: (bucket: 'daily' | 'hourly') => void;
}

const EMPTY_COUNT = 0;

export function OverviewChartSection({ siteId, dateRange, filter, bucket, onBucketChange }: OverviewChartSectionProps): React.JSX.Element | null {
  const { loadedData, loading } = useChartDataLoader({ siteId, dateRange, filter, bucket });

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
    return <ChartSkeleton />;
  }

  if (chartData.length === EMPTY_COUNT && !loading) {
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
        <ChartContainer config={CHART_CONFIG} className="h-[300px] w-full">
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
