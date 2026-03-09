import type { OverviewPoint } from '@/components/overview-chart/overview-chart-series';

const MAX_AXIS_TICKS = 6;
const FIRST_INDEX = 0;
const MIN_TICK_COUNT = 2;
const dailyAxisFormatter = new Intl.DateTimeFormat('en-US', {
  month: 'short',
  day: 'numeric',
  timeZone: 'UTC',
});
const hourlyAxisDateFormatter = new Intl.DateTimeFormat('en-US', {
  month: 'short',
  day: 'numeric',
  timeZone: 'UTC',
});
const hourlyAxisTimeFormatter = new Intl.DateTimeFormat('en-US', {
  hour: 'numeric',
  timeZone: 'UTC',
});

export const buildOverviewAxisTicks = (data: OverviewPoint[]): number[] => {
  if (data.length <= MAX_AXIS_TICKS) {
    return data.map((point) => point.timestamp);
  }

  const segmentCount = MAX_AXIS_TICKS - 1;
  const lastIndex = data.length - 1;
  const indexes = new Set<number>([FIRST_INDEX, lastIndex]);

  for (let segment = 1; segment < segmentCount; segment += 1) {
    indexes.add(Math.round((segment * lastIndex) / segmentCount));
  }

  return [...indexes]
    .sort((left, right) => left - right)
    .slice(FIRST_INDEX, Math.max(MIN_TICK_COUNT, MAX_AXIS_TICKS))
    .map((index) => data[index]?.timestamp)
    .filter((tick): tick is number => typeof tick === 'number');
};

export const formatOverviewAxisTickLines = (bucket: 'daily' | 'hourly', value: number): string[] => {
  const date = new Date(value);
  if (bucket === 'hourly') {
    return [hourlyAxisDateFormatter.format(date), hourlyAxisTimeFormatter.format(date)];
  }

  return [dailyAxisFormatter.format(date)];
};
