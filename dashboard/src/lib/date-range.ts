export type DatePreset = '7d' | '30d' | '90d' | 'custom' | 'all';

const DATE_PRESET_WEEK: DatePreset = '7d';
const DATE_PRESET_MONTH: DatePreset = '30d';
const DATE_PRESET_QUARTER: DatePreset = '90d';
const DATE_PRESET_CUSTOM: DatePreset = 'custom';
const DATE_PRESET_ALL: DatePreset = 'all';

const DAY_OFFSET_WEEK = 6;
const DAY_OFFSET_MONTH = 29;
const DAY_OFFSET_QUARTER = 89;
const MONTH_INDEX_OFFSET = 1;
const PAD_LENGTH = 2;

const DATE_REGEX = /^\d{4}-\d{2}-\d{2}$/v;
const TIME_REGEX = /^\d{2}:\d{2}$/v;

const TIME_INPUT_LENGTH = 5;
const TIME_SLICE_START = 0;
const TIME_SECONDS_SUFFIX = ':00';

const ZERO_TIME = '00:00';
const MAX_TIME = '23:59';
const ZERO_TIME_FULL = '00:00:00';
const MAX_TIME_FULL = '23:59:59';

const EMPTY_STRING = '';
const HOUR_MIN = 0;
const HOUR_MAX = 23;
const MINUTE_MIN = 0;
const MINUTE_MAX = 59;

export function isDatePreset(value: string | undefined): value is DatePreset {
  return (
    value === DATE_PRESET_WEEK ||
    value === DATE_PRESET_MONTH ||
    value === DATE_PRESET_QUARTER ||
    value === DATE_PRESET_CUSTOM ||
    value === DATE_PRESET_ALL
  );
}

export function formatDateInput(date: Date): string {
  const year = date.getFullYear();
  const month = String(date.getMonth() + MONTH_INDEX_OFFSET).padStart(PAD_LENGTH, '0');
  const day = String(date.getDate()).padStart(PAD_LENGTH, '0');
  return `${year}-${month}-${day}`;
}

export function presetToDates(
  preset: DatePreset,
  reference = new Date()
): { fromDate: string; toDate: string } {
  const end = new Date(reference);
  const start = new Date(reference);
  if (preset === DATE_PRESET_WEEK) {
    start.setDate(end.getDate() - DAY_OFFSET_WEEK);
  } else if (preset === DATE_PRESET_QUARTER) {
    start.setDate(end.getDate() - DAY_OFFSET_QUARTER);
  } else {
    start.setDate(end.getDate() - DAY_OFFSET_MONTH);
  }
  return { fromDate: formatDateInput(start), toDate: formatDateInput(end) };
}

export function isValidDateInput(value: string): boolean {
  if (value === EMPTY_STRING) return false;
  if (!DATE_REGEX.test(value)) return false;
  const parsed = new Date(`${value}T${ZERO_TIME_FULL}`);
  if (Number.isNaN(parsed.getTime())) return false;
  return formatDateInput(parsed) === value;
}

export function isValidTimeInput(value: string): boolean {
  if (value === EMPTY_STRING) return false;
  if (!TIME_REGEX.test(value)) return false;
  const [hourRaw, minuteRaw] = value.split(':');
  if (hourRaw === undefined || minuteRaw === undefined) return false;
  const hour = Number(hourRaw);
  const minute = Number(minuteRaw);
  if (Number.isNaN(hour) || Number.isNaN(minute)) return false;
  return (
    hour >= HOUR_MIN &&
    hour <= HOUR_MAX &&
    minute >= MINUTE_MIN &&
    minute <= MINUTE_MAX
  );
}

export function parseDateInput(value: string): Date | undefined {
  if (!isValidDateInput(value)) return undefined;
  const parsed = new Date(`${value}T${ZERO_TIME_FULL}`);
  if (Number.isNaN(parsed.getTime())) return undefined;
  return parsed;
}

export function normalizeTimeInput(value: string, fallback: string): string {
  if (value === EMPTY_STRING) return fallback;
  const trimmed = value.trim();
  if (trimmed === EMPTY_STRING) return fallback;
  if (trimmed.length >= TIME_INPUT_LENGTH) {
    return trimmed.slice(TIME_SLICE_START, TIME_INPUT_LENGTH);
  }
  return fallback;
}

export function buildDateRange(
  fromDate: string,
  toDate: string,
  fromTime = ZERO_TIME,
  toTime = MAX_TIME
): { from: string; to: string } | undefined {
  if (fromDate === EMPTY_STRING || toDate === EMPTY_STRING) return undefined;
  const normalizeTime = (value: string, fallback: string): string => {
    if (value === EMPTY_STRING) return fallback;
    if (value.length === TIME_INPUT_LENGTH) {
      return `${value}${TIME_SECONDS_SUFFIX}`;
    }
    return value;
  };
  const from = new Date(`${fromDate}T${normalizeTime(fromTime, ZERO_TIME_FULL)}`);
  const to = new Date(`${toDate}T${normalizeTime(toTime, MAX_TIME_FULL)}`);
  if (Number.isNaN(from.getTime()) || Number.isNaN(to.getTime())) {
    return undefined;
  }
  return { from: from.toISOString(), to: to.toISOString() };
}
