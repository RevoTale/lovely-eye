import type { AnalyticsSearch } from '@/features/analytics/model/analytics-search';
import {
  buildDateRange,
  type DatePreset,
  isValidDateInput,
  isValidTimeInput,
  presetToDates,
} from '@/features/analytics/model/date-range';

export const DEFAULT_DATE_PRESET: DatePreset = '30d';

export interface AnalyticsDateRangeState {
  preset: DatePreset;
  fromDate: string;
  toDate: string;
  fromTime: string;
  toTime: string;
}

export const resolveAnalyticsDateState = (
  search: AnalyticsSearch,
  reference: Date
): AnalyticsDateRangeState => {
  const preset = search.preset ?? DEFAULT_DATE_PRESET;
  if (preset === 'all') {
    return { preset, fromDate: '', toDate: '', fromTime: '', toTime: '' };
  }

  if (preset === 'custom') {
    const customRange = resolveCustomRange(search);
    if (customRange !== undefined) return { preset, ...customRange };
  }

  const presetDates = presetToDates(preset, reference);
  return {
    preset,
    fromDate: presetDates.fromDate,
    toDate: presetDates.toDate,
    fromTime: '00:00',
    toTime: '23:59',
  };
};

export const resolveAnalyticsDateRange = (
  search: AnalyticsSearch,
  reference: Date
): { from: string; to: string } | undefined => {
  const { preset, fromDate, toDate, fromTime, toTime } = resolveAnalyticsDateState(
    search,
    reference
  );
  if (preset === 'all') return undefined;
  return buildDateRange(fromDate, toDate, fromTime, toTime);
};

const resolveCustomRange = (
  search: AnalyticsSearch
): Omit<AnalyticsDateRangeState, 'preset'> | undefined => {
  const fromDate = search.from ?? '';
  const toDate = search.to ?? '';
  const fromTime = search.fromTime ?? '';
  const toTime = search.toTime ?? '';
  const hasValidInputs =
    isValidDateInput(fromDate) &&
    isValidDateInput(toDate) &&
    isValidTimeInput(fromTime) &&
    isValidTimeInput(toTime);
  if (!hasValidInputs) return undefined;

  const candidate = buildDateRange(fromDate, toDate, fromTime, toTime);
  if (candidate === undefined || new Date(candidate.from) > new Date(candidate.to))
    return undefined;
  return { fromDate, toDate, fromTime, toTime };
};
