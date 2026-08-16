import { useNavigate, useParams, useSearch } from '@tanstack/react-router';
import { useMemo } from 'react';
import {
  ANALYTICS_ROUTE_ID,
  type AnalyticsSearch,
} from '@/features/analytics/model/analytics-search';
import {
  buildDateRange,
  type DatePreset,
  isValidDateInput,
  isValidTimeInput,
  presetToDates,
} from '@/features/analytics/model/date-range';

interface DateRangeInput {
  fromDate: string;
  toDate: string;
  fromTime: string;
  toTime: string;
}

interface DateRangeState extends DateRangeInput {
  preset: DatePreset;
  dateRange: { from: string; to: string } | undefined;
  setPreset: (preset: DatePreset) => void;
  applyCustomRange: (range: DateRangeInput) => boolean;
}

const DEFAULT_PRESET: DatePreset = '30d';

const resolveCustomRange = (raw: AnalyticsSearch): DateRangeInput | undefined => {
  const fromDate = raw.from ?? '';
  const toDate = raw.to ?? '';
  const fromTime = raw.fromTime ?? '';
  const toTime = raw.toTime ?? '';
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

const resolveSearchState = (raw: AnalyticsSearch): { preset: DatePreset } & DateRangeInput => {
  const preset = raw.preset ?? DEFAULT_PRESET;
  if (preset === 'all') {
    return { preset, fromDate: '', toDate: '', fromTime: '', toTime: '' };
  }

  if (preset === 'custom') {
    const customRange = resolveCustomRange(raw);
    if (customRange !== undefined) return { preset, ...customRange };
  }

  const presetDates = presetToDates(preset, new Date());
  return {
    preset,
    fromDate: presetDates.fromDate,
    toDate: presetDates.toDate,
    fromTime: '00:00',
    toTime: '23:59',
  };
};

export function useDateRange(): DateRangeState {
  const search = useSearch({ from: ANALYTICS_ROUTE_ID });
  const { siteId } = useParams({ from: ANALYTICS_ROUTE_ID });
  const navigate = useNavigate();

  const resolvedState = useMemo(() => resolveSearchState(search), [search]);

  const { preset, fromDate, toDate, fromTime, toTime } = resolvedState;

  const dateRange = useMemo(() => {
    if (preset === 'all') {
      return undefined;
    }
    if (preset === 'custom') {
      return buildDateRange(fromDate, toDate, fromTime, toTime);
    }
    const presetDates = presetToDates(preset, new Date());
    return buildDateRange(presetDates.fromDate, presetDates.toDate, '00:00', '23:59');
  }, [preset, fromDate, toDate, fromTime, toTime]);

  const setPreset = (value: DatePreset): void => {
    if (value === 'custom') {
      const fallbackDates = presetToDates(DEFAULT_PRESET, new Date());
      const nextFromDate = isValidDateInput(fromDate) ? fromDate : fallbackDates.fromDate;
      const nextToDate = isValidDateInput(toDate) ? toDate : fallbackDates.toDate;
      const nextFromTime = isValidTimeInput(fromTime) ? fromTime : '00:00';
      const nextToTime = isValidTimeInput(toTime) ? toTime : '23:59';
      void navigate({
        to: '/sites/$siteId/analytics',
        params: { siteId },
        search: (prev) => ({
          ...prev,
          preset: 'custom',
          from: nextFromDate,
          to: nextToDate,
          fromTime: nextFromTime,
          toTime: nextToTime,
        }),
      });
      return;
    }
    if (value === 'all') {
      void navigate({
        to: '/sites/$siteId/analytics',
        params: { siteId },
        search: (prev) => {
          const { from: _from, to: _to, fromTime: _fromTime, toTime: _toTime, ...rest } = prev;
          return { ...rest, preset: 'all' };
        },
      });
      return;
    }
    void navigate({
      to: '/sites/$siteId/analytics',
      params: { siteId },
      search: (prev) => {
        const { from: _from, to: _to, fromTime: _fromTime, toTime: _toTime, ...rest } = prev;
        return { ...rest, preset: value };
      },
    });
  };

  const applyCustomRange = (range: DateRangeInput): boolean => {
    const {
      fromDate: nextFrom,
      toDate: nextTo,
      fromTime: nextFromTime,
      toTime: nextToTime,
    } = range;
    if (
      !isValidDateInput(nextFrom) ||
      !isValidDateInput(nextTo) ||
      !isValidTimeInput(nextFromTime) ||
      !isValidTimeInput(nextToTime)
    ) {
      return false;
    }
    const candidate = buildDateRange(nextFrom, nextTo, nextFromTime, nextToTime);
    if (candidate === undefined || new Date(candidate.from) > new Date(candidate.to)) {
      return false;
    }
    void navigate({
      to: '/sites/$siteId/analytics',
      params: { siteId },
      search: (prev) => ({
        ...prev,
        preset: 'custom',
        from: nextFrom,
        to: nextTo,
        fromTime: nextFromTime,
        toTime: nextToTime,
      }),
    });
    return true;
  };

  return {
    preset,
    fromDate,
    toDate,
    fromTime,
    toTime,
    dateRange,
    setPreset,
    applyCustomRange,
  };
}
