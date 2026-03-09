import type { CSSProperties, FunctionComponent } from 'react';
import { Area, AreaChart, CartesianGrid, Legend, Tooltip, XAxis, YAxis } from 'recharts';
import { buildOverviewAxisTicks } from '@/components/overview-chart/overview-chart-axis';
import OverviewChartLegend from '@/components/overview-chart/overview-chart-legend';
import OverviewChartTooltip from '@/components/overview-chart/overview-chart-tooltip';
import OverviewChartXAxisTick from '@/components/overview-chart/overview-chart-x-axis-tick';
import { formatOverviewAxisValue } from '@/components/overview-chart/overview-chart-formatters';
import { OVERVIEW_RENDER_SERIES, OVERVIEW_SERIES, getOverviewSeriesIndex, type OverviewPoint } from '@/components/overview-chart/overview-chart-series';

interface OverviewChartPlotProps {
  bucket: 'daily' | 'hourly';
  data: OverviewPoint[];
}

const MAX_TICK_COUNT = 6;
const X_AXIS_MIN_TICK_GAP = 24;
const X_AXIS_TICK_MARGIN = 10;
const DAILY_X_AXIS_HEIGHT = 40;
const HOURLY_X_AXIS_HEIGHT = 52;
const LEGEND_HEIGHT = 52;

const chartStyle = {
  width: '100%',
  height: '100%',
  maxWidth: '100%',
} satisfies CSSProperties;

const OverviewChartPlot: FunctionComponent<OverviewChartPlotProps> = ({ bucket, data }) => {
  const ticks = buildOverviewAxisTicks(data).slice(0, MAX_TICK_COUNT);
  const xAxisHeight = bucket === 'hourly' ? HOURLY_X_AXIS_HEIGHT : DAILY_X_AXIS_HEIGHT;

  return (
    <div className="h-[clamp(240px,44vw,320px)] w-full">
      <AreaChart accessibilityLayer data={data} responsive style={chartStyle} margin={{ top: 12, right: 6, bottom: 0, left: -18 }}>
        <defs>
          {OVERVIEW_SERIES.map((series) => (
            <linearGradient key={series.key} id={`overview-${series.key}`} x1="0" y1="0" x2="0" y2="1">
              <stop offset="5%" stopColor={series.fillStart} />
              <stop offset="95%" stopColor={series.fillEnd} />
            </linearGradient>
          ))}
        </defs>
        <CartesianGrid vertical={false} strokeDasharray="3 3" className="stroke-border/60" />
        <XAxis
          dataKey="timestamp"
          type="number"
          scale="time"
          domain={['dataMin', 'dataMax']}
          ticks={ticks}
          tickLine={false}
          axisLine={false}
          minTickGap={X_AXIS_MIN_TICK_GAP}
          interval={0}
          tickMargin={X_AXIS_TICK_MARGIN}
          height={xAxisHeight}
          tick={(props) => <OverviewChartXAxisTick {...props} bucket={bucket} />}
          className="text-xs"
        />
        <YAxis tickFormatter={(value) => formatOverviewAxisValue(Number(value))} tickLine={false} axisLine={false} width={44} className="text-xs" />
        <Tooltip
          cursor={{ stroke: 'hsl(var(--border))', strokeDasharray: '4 4' }}
          isAnimationActive={false}
          itemSorter={(item) => getOverviewSeriesIndex(item.dataKey)}
          content={(props) => <OverviewChartTooltip {...props} bucket={bucket} />}
        />
        <Legend
          verticalAlign="bottom"
          align="left"
          height={LEGEND_HEIGHT}
          itemSorter={(item) => getOverviewSeriesIndex(item.dataKey)}
          content={(props) => <OverviewChartLegend {...props} />}
        />
        {OVERVIEW_RENDER_SERIES.map((series) => (
          <Area
            key={series.key}
            type="monotone"
            dataKey={series.key}
            name={series.label}
            stroke={series.stroke}
            fill={`url(#overview-${series.key})`}
            strokeWidth={series.strokeWidth}
            fillOpacity={1}
            dot={false}
            activeDot={{ r: 4, fill: series.stroke, stroke: 'hsl(var(--background))', strokeWidth: 2 }}
            isAnimationActive={false}
          />
        ))}
      </AreaChart>
    </div>
  );
};

export default OverviewChartPlot;
