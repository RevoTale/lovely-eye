const OVERVIEW_SERIES_KEYS = ['visitors', 'pageViews', 'sessions'] as const;

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
    stroke: 'var(--primary)',
    fillStart: 'color-mix(in oklch, var(--primary) 34%, transparent)',
    fillEnd: 'color-mix(in oklch, var(--primary) 4%, transparent)',
    strokeWidth: 2.75,
  },
  {
    key: 'pageViews',
    label: 'Page Views',
    stroke: 'var(--chart-2)',
    fillStart: 'color-mix(in oklch, var(--chart-2) 18%, transparent)',
    fillEnd: 'color-mix(in oklch, var(--chart-2) 3%, transparent)',
    strokeWidth: 2,
  },
  {
    key: 'sessions',
    label: 'Sessions',
    stroke: 'var(--chart-3)',
    fillStart: 'color-mix(in oklch, var(--chart-3) 16%, transparent)',
    fillEnd: 'color-mix(in oklch, var(--chart-3) 3%, transparent)',
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
