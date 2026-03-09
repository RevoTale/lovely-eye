import type { FunctionComponent } from 'react';
import type { TooltipContentProps } from 'recharts';
import { compareOverviewSeriesOrder, getOverviewSeries, isOverviewSeriesKey } from '@/components/overview-chart/overview-chart-series';
import { formatOverviewValue } from '@/components/overview-chart/overview-chart-formatters';
import { buildTimeFormatter } from '@/lib/chart-date';

interface OverviewChartTooltipProps extends TooltipContentProps {
  bucket: 'daily' | 'hourly';
}

const OverviewChartTooltip: FunctionComponent<OverviewChartTooltipProps> = ({ active, bucket, label, payload }) => {
  if (!active || payload === undefined || payload.length === 0 || typeof label !== 'number') {
    return null;
  }

  const formatter = buildTimeFormatter(bucket);
  const values = payload
    .map((entry) => {
      const dataKey = entry.dataKey;
      if (!isOverviewSeriesKey(dataKey)) {
        return null;
      }

      const series = getOverviewSeries(dataKey);
      const value = typeof entry.value === 'number' ? entry.value : Number(entry.value ?? 0);
      return series === undefined ? null : { series, value };
    })
    .filter((entry): entry is NonNullable<typeof entry> => entry !== null)
    .sort((left, right) => compareOverviewSeriesOrder(left.series.key, right.series.key));

  return (
    <div className="min-w-[180px] rounded-xl border border-border/70 bg-background/96 px-3 py-2 text-xs shadow-xl backdrop-blur-sm">
      <div className="mb-2 border-b border-border/70 pb-2 font-medium text-foreground">{formatter.format(label)}</div>
      <div className="space-y-1.5">
        {values.map(({ series, value }) => (
          <div key={series.key} className="flex items-center justify-between gap-4">
            <div className="flex items-center gap-2 text-muted-foreground">
              <span className="h-2.5 w-2.5 rounded-full" style={{ backgroundColor: series.stroke }} />
              <span>{series.label}</span>
            </div>
            <span className="font-mono font-medium tabular-nums text-foreground">{formatOverviewValue(value)}</span>
          </div>
        ))}
      </div>
    </div>
  );
};

export default OverviewChartTooltip;
