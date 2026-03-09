import type { FunctionComponent } from 'react';
import type { XAxisTickContentProps } from 'recharts';
import { formatOverviewAxisTickLines } from '@/components/overview-chart/overview-chart-axis';

interface OverviewChartXAxisTickProps extends XAxisTickContentProps {
  bucket: 'daily' | 'hourly';
}

const OverviewChartXAxisTick: FunctionComponent<OverviewChartXAxisTickProps> = ({ bucket, payload, x = 0, y = 0 }) => {
  const value = payload?.value;
  if (typeof value !== 'number') {
    return null;
  }

  const lines = formatOverviewAxisTickLines(bucket, value);

  return (
    <g transform={`translate(${x},${y})`}>
      <text
        x={0}
        y={0}
        dy={12}
        textAnchor="middle"
        fill="hsl(var(--muted-foreground))"
        className="text-[11px]"
      >
        {lines.map((line, index) => (
          <tspan key={`${line}-${index}`} x={0} dy={index === 0 ? 0 : 12}>
            {line}
          </tspan>
        ))}
      </text>
    </g>
  );
};

export default OverviewChartXAxisTick;
