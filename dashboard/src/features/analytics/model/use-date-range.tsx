import { useNavigate, useParams, useSearch } from '@tanstack/react-router';
import { useMemo } from 'react';
import {
  DEFAULT_DATE_PRESET,
  resolveAnalyticsDateState,
} from '@/features/analytics/model/analytics-date-range';
import { ANALYTICS_ROUTE_ID } from '@/features/analytics/model/analytics-search';
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

export function useDateRange(): DateRangeState {
  const search = useSearch({ from: ANALYTICS_ROUTE_ID });
  const { siteId } = useParams({ from: ANALYTICS_ROUTE_ID });
  const navigate = useNavigate();

  const resolvedState = useMemo(() => resolveAnalyticsDateState(search, new Date()), [search]);

  const { preset, fromDate, toDate, fromTime, toTime } = resolvedState;

  const dateRange = useMemo(
    () => (preset === 'all' ? undefined : buildDateRange(fromDate, toDate, fromTime, toTime)),
    [preset, fromDate, toDate, fromTime, toTime]
  );

  const setPreset = (value: DatePreset): void => {
    if (value === 'custom') {
      const fallbackDates = presetToDates(DEFAULT_DATE_PRESET, new Date());
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
