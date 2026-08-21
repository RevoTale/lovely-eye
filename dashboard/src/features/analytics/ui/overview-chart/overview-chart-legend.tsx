import type { FunctionComponent } from 'react';
import type { DefaultLegendContentProps, LegendPayload } from 'recharts';
import {
  compareOverviewSeriesOrder,
  getOverviewSeries,
  isOverviewSeriesKey,
} from '@/features/analytics/ui/overview-chart/overview-chart-series';

const OverviewChartLegend: FunctionComponent<DefaultLegendContentProps> = ({ payload }) => {
  const items = (payload ?? [])
    .map((entry: LegendPayload) => {
      if (!isOverviewSeriesKey(entry.dataKey)) {
        return null;
      }

      return getOverviewSeries(entry.dataKey) ?? null;
    })
    .filter((entry): entry is NonNullable<typeof entry> => entry !== null)
    .sort((left, right) => compareOverviewSeriesOrder(left.key, right.key));

  return (
    <div className='mt-4 flex flex-wrap gap-2'>
      {items.map((series) => (
        <div
          key={series.key}
          className='inline-flex max-w-full items-center gap-2 rounded-full border border-border/70 bg-muted/30 px-3 py-1.5 text-xs text-muted-foreground'
          title={series.label}
        >
          <svg aria-hidden='true' className='h-2.5 w-5' viewBox='0 0 20 10'>
            <line
              x1='0'
              y1='5'
              x2='20'
              y2='5'
              stroke={series.stroke}
              strokeWidth={series.strokeWidth}
              strokeDasharray={series.strokeDasharray}
            />
          </svg>
          <span className='max-w-[10rem] truncate font-medium text-foreground'>{series.label}</span>
        </div>
      ))}
    </div>
  );
};

export default OverviewChartLegend;
