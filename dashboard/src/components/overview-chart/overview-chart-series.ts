export const OVERVIEW_SERIES_KEYS = ['visitors', 'pageViews', 'sessions'] as const;

export type OverviewSeriesKey = (typeof OVERVIEW_SERIES_KEYS)[number];

export type OverviewPoint = { timestamp: number } & Record<OverviewSeriesKey, number>;

interface OverviewSeriesDefinition {
  key: OverviewSeriesKey;
  label: string;
  stroke: string;
  fillStart: string;
  fillEnd: string;
  strokeWidth: number;
}

export const OVERVIEW_SERIES = [
  {
    key: 'visitors',
    label: 'Visitors',
    stroke: 'hsl(var(--primary))',
    fillStart: 'hsl(var(--primary) / 0.34)',
    fillEnd: 'hsl(var(--primary) / 0.04)',
    strokeWidth: 2.75,
  },
  {
    key: 'pageViews',
    label: 'Page Views',
    stroke: 'hsl(var(--chart-2))',
    fillStart: 'hsl(var(--chart-2) / 0.18)',
    fillEnd: 'hsl(var(--chart-2) / 0.03)',
    strokeWidth: 2,
  },
  {
    key: 'sessions',
    label: 'Sessions',
    stroke: 'hsl(var(--chart-3))',
    fillStart: 'hsl(var(--chart-3) / 0.16)',
    fillEnd: 'hsl(var(--chart-3) / 0.03)',
    strokeWidth: 2,
  },
] as const satisfies readonly OverviewSeriesDefinition[];

export const OVERVIEW_RENDER_SERIES = [...OVERVIEW_SERIES].reverse();

export const isOverviewSeriesKey = (value: unknown): value is OverviewSeriesKey =>
  typeof value === 'string' && OVERVIEW_SERIES_KEYS.some((seriesKey) => seriesKey === value);

export const getOverviewSeries = (value: unknown) =>
  OVERVIEW_SERIES.find((series) => series.key === value);

export const getOverviewSeriesIndex = (value: unknown): number =>
  OVERVIEW_SERIES.findIndex((series) => series.key === value);

export const compareOverviewSeriesOrder = (left: unknown, right: unknown): number =>
  getOverviewSeriesIndex(left) - getOverviewSeriesIndex(right);
