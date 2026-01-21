import type { ChartConfig } from '@/components/ui';

export const CHART_CONFIG = {
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
} satisfies ChartConfig;

const MARGIN_TOP = 10;
const MARGIN_RIGHT = 10;
const MARGIN_LEFT = 0;
const MARGIN_BOTTOM = 0;

export const CHART_MARGIN = { top: MARGIN_TOP, right: MARGIN_RIGHT, left: MARGIN_LEFT, bottom: MARGIN_BOTTOM };
export const TICK_MARGIN = 8;
