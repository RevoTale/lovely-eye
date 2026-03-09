import type { ChartDataQueryVariables, FilterInput, TimeBucket } from '@/gql/graphql';

const MINUTES_PER_HOUR = 60;
const SECONDS_PER_MINUTE = 60;
const MS_IN_SECOND = 1000;
const HOURS_PER_DAY = 24;
const MS_IN_HOUR = MINUTES_PER_HOUR * SECONDS_PER_MINUTE * MS_IN_SECOND;
const MS_IN_DAY = HOURS_PER_DAY * MS_IN_HOUR;
const COUNT_OFFSET = 1;
const PROGRESS_UNKNOWN = 50;
const PROGRESS_COMPLETE = 100;
const PROGRESS_MIN = 0;

export const BATCH_SIZE = 10;

export const buildChartVariables = (
  siteId: string,
  dateRange: { from: Date; to: Date } | null,
  filter: FilterInput | null,
  bucketValue: TimeBucket
): ChartDataQueryVariables => ({
  siteId,
  dateRange: dateRange === null ? null : { from: dateRange.from.toISOString(), to: dateRange.to.toISOString() },
  filter: filter === null ? null : {
    referrer: filter.referrer ?? null,
    browser: filter.browser ?? null,
    device: filter.device ?? null,
    os: filter.os ?? null,
    page: filter.page ?? null,
    country: filter.country ?? null,
    eventType: filter.eventType ?? null,
    eventDefinitionId: filter.eventDefinitionId ?? null,
    eventName: filter.eventName ?? null,
    eventPath: filter.eventPath ?? null,
  },
  bucket: bucketValue,
  limit: BATCH_SIZE,
  offset: 0,
});

export const calculateExpectedCount = (
  dateRange: { from: Date; to: Date } | null,
  bucket: 'daily' | 'hourly'
): number | null => {
  if (dateRange === null) return null;
  const fromMs = dateRange.from.getTime();
  const toMs = dateRange.to.getTime();
  if (Number.isNaN(fromMs) || Number.isNaN(toMs) || toMs < fromMs) return null;
  return Math.floor((toMs - fromMs) / (bucket === 'hourly' ? MS_IN_HOUR : MS_IN_DAY)) + COUNT_OFFSET;
};

export const calculateProgress = (expectedCount: number | null, hasMore: boolean, loadedCount: number): number | null => {
  if (expectedCount === null) {
    if (loadedCount === 0) return null;
    return hasMore ? PROGRESS_UNKNOWN : PROGRESS_COMPLETE;
  }
  if (expectedCount <= 0) return null;
  return Math.min(
    PROGRESS_COMPLETE,
    Math.max(PROGRESS_MIN, Math.round((loadedCount / expectedCount) * PROGRESS_COMPLETE))
  );
};
