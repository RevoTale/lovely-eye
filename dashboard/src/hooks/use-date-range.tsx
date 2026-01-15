import { useCallback, useEffect, useMemo, useState } from 'react';
import { useNavigate, useParams, useSearch } from '@tanstack/react-router';
import {
  buildDateRange,
  isDatePreset,
  isValidDateInput,
  isValidTimeInput,
  presetToDates,
  type DatePreset,
} from '@/lib/date-range';
import { siteDetailRoute } from '@/router';

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
  const search = useSearch({ from: siteDetailRoute.id });
  const { siteId } = useParams({ from: siteDetailRoute.id });
  const navigate = useNavigate();
  const defaultPreset: DatePreset = '30d';

  const resolveSearchState = useCallback((raw: typeof search): { preset: DatePreset } & DateRangeInput => {
    const presetValue = isDatePreset(raw.preset) ? raw.preset : defaultPreset;
    if (presetValue === 'all') {
      return {
        preset: 'all',
        fromDate: '',
        toDate: '',
        fromTime: '',
        toTime: '',
      };
    }
    if (presetValue === 'custom') {
      const fromDate = typeof raw.from === 'string' ? raw.from : '';
      const toDate = typeof raw.to === 'string' ? raw.to : '';
      const fromTime = typeof raw.fromTime === 'string' ? raw.fromTime : '';
      const toTime = typeof raw.toTime === 'string' ? raw.toTime : '';
      const hasValidInputs =
        isValidDateInput(fromDate) &&
        isValidDateInput(toDate) &&
        isValidTimeInput(fromTime) &&
        isValidTimeInput(toTime);
      if (hasValidInputs) {
        const candidate = buildDateRange(fromDate, toDate, fromTime, toTime);
        if (candidate && new Date(candidate.from) <= new Date(candidate.to)) {
          return { preset: 'custom', fromDate, toDate, fromTime, toTime };
        }
      }
    }
    const presetDates = presetToDates(presetValue, new Date());
    return {
      preset: presetValue,
      fromDate: presetDates.fromDate,
      toDate: presetDates.toDate,
      fromTime: '00:00',
      toTime: '23:59',
    };
  }, [defaultPreset]);

  const [state, setState] = useState(() => resolveSearchState(search));
  const { preset, fromDate, toDate, fromTime, toTime } = state;

  useEffect(() => {
    const nextState = resolveSearchState(search);
    setState((prev) => {
      if (
        prev.preset === nextState.preset &&
        prev.fromDate === nextState.fromDate &&
        prev.toDate === nextState.toDate &&
        prev.fromTime === nextState.fromTime &&
        prev.toTime === nextState.toTime
      ) {
        return prev;
      }
      return nextState;
    });
  }, [resolveSearchState, search]);

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
      const candidate = buildDateRange(fromDate, toDate, fromTime, toTime);
      if (!candidate) {
        return;
      }
      setState((prev) => ({ ...prev, preset: 'custom' }));
      void navigate({
        to: '/sites/$siteId',
        params: { siteId },
        search: (prev) => ({
          ...prev,
          preset: 'custom',
          from: fromDate,
          to: toDate,
          fromTime,
          toTime,
        }),
      });
      return;
    }
    if (value === 'all') {
      setState({
        preset: 'all',
        fromDate: '',
        toDate: '',
        fromTime: '',
        toTime: '',
      });
      void navigate({
        to: '/sites/$siteId',
        params: { siteId },
        search: (prev) => {
          const { from, to, fromTime, toTime, ...rest } = prev as Record<string, unknown>;
          return {
            ...(rest as Record<string, unknown>),
            preset: 'all',
          };
        },
      });
      return;
    }
    const presetDates = presetToDates(value, new Date());
    setState({
      preset: value,
      fromDate: presetDates.fromDate,
      toDate: presetDates.toDate,
      fromTime: '00:00',
      toTime: '23:59',
    });
    void navigate({
      to: '/sites/$siteId',
      params: { siteId },
      search: (prev) => {
        const { from, to, fromTime, toTime, ...rest } = prev as Record<string, unknown>;
        return {
          ...(rest as Record<string, unknown>),
          preset: value,
        };
      },
    });
  };

  const applyCustomRange = (range: DateRangeInput): boolean => {
    const { fromDate: nextFrom, toDate: nextTo, fromTime: nextFromTime, toTime: nextToTime } = range;
    if (
      !isValidDateInput(nextFrom) ||
      !isValidDateInput(nextTo) ||
      !isValidTimeInput(nextFromTime) ||
      !isValidTimeInput(nextToTime)
    ) {
      return false;
    }
    const candidate = buildDateRange(nextFrom, nextTo, nextFromTime, nextToTime);
    if (!candidate || new Date(candidate.from) > new Date(candidate.to)) {
      return false;
    }
    setState({
      preset: 'custom',
      fromDate: nextFrom,
      toDate: nextTo,
      fromTime: nextFromTime,
      toTime: nextToTime,
    });
    void navigate({
      to: '/sites/$siteId',
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
