import React from 'react';
import { Card, CardContent, CardHeader, CardTitle, ChartContainer, ChartLegend, ChartLegendContent, ChartTooltip, ChartTooltipContent, type ChartConfig } from '@/components/ui';
import { Area, AreaChart, CartesianGrid, XAxis, YAxis } from 'recharts';
import { TrendingUp } from 'lucide-react';
import type { DailyStats } from '@/generated/graphql';

interface OverviewChartSectionProps {
  dailyStats: DailyStats[];
}

export function OverviewChartSection({ dailyStats }: OverviewChartSectionProps): React.JSX.Element | null {
  if (dailyStats.length === 0) {
    return null;
  }

  return (
    <Card className="hover:shadow-md transition-shadow">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <div className="h-8 w-8 rounded-lg bg-primary/10 flex items-center justify-center">
            <TrendingUp className="h-4 w-4 text-primary" />
          </div>
          Analytics Overview
        </CardTitle>
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
            data={dailyStats.map(stat => ({
              date: new Date(stat.date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' }),
              visitors: stat.visitors,
              pageViews: stat.pageViews,
              sessions: stat.sessions,
            }))}
            margin={{ top: 10, right: 10, left: 0, bottom: 0 }}
          >
            <CartesianGrid strokeDasharray="3 3" className="stroke-muted" />
            <XAxis
              dataKey="date"
              tickLine={false}
              axisLine={false}
              tickMargin={8}
              className="text-xs"
            />
            <YAxis
              tickLine={false}
              axisLine={false}
              tickMargin={8}
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
