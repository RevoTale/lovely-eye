const DATE_PART_LENGTH = 10;
const DATE_SEPARATOR = '-';
const TIME_SEPARATOR = ':';
const DATE_TIME_SEPARATOR = 'T';
const SPACE_SEPARATOR = ' ';
const YEAR_LENGTH = 4;
const MONTH_LENGTH = 2;
const DAY_LENGTH = 2;
const HOURS_LENGTH = 2;
const MINUTES_LENGTH = 2;
const SECONDS_LENGTH = 2;
const SECONDS_THRESHOLD = 1_000_000_000_000;
const MS_IN_SECOND = 1000;
const UTC_SUFFIX = ' UTC';
const UTC_MARKER = 'Z';
const TZ_OFFSET_LENGTH = 5;
const TZ_OFFSET_WITH_COLON_LENGTH = 6;
const TZ_COLON_INDEX = 3;
const NOT_FOUND = -1;
const NON_EMPTY_LENGTH = 1;
const INDEX_START = 0;
const INDEX_STEP = 1;
const CHAR_CODE_ZERO = 48;
const CHAR_CODE_NINE = 57;

const normalizeEpoch = (value: number): number => (value < SECONDS_THRESHOLD ? value * MS_IN_SECOND : value);

const isDigits = (value: string): boolean => {
  if (value.length < NON_EMPTY_LENGTH) return false;
  for (let index = INDEX_START; index < value.length; index += INDEX_STEP) {
    const charCode = value.charCodeAt(index);
    if (charCode < CHAR_CODE_ZERO || charCode > CHAR_CODE_NINE) return false;
  }
  return true;
};

const isDateOnly = (value: string): boolean => {
  if (value.length !== DATE_PART_LENGTH) return false;
  const [year, month, day] = value.split(DATE_SEPARATOR);
  return year?.length === YEAR_LENGTH && month?.length === MONTH_LENGTH && day?.length === DAY_LENGTH
    && isDigits(year) && isDigits(month) && isDigits(day);
};

const isTimePart = (value: string): boolean => {
  const [hours, minutes, secondsWithZone] = value.split(TIME_SEPARATOR);
  const seconds = secondsWithZone?.slice(INDEX_START, SECONDS_LENGTH);
  return hours?.length === HOURS_LENGTH && minutes?.length === MINUTES_LENGTH && seconds?.length === SECONDS_LENGTH
    && isDigits(hours) && isDigits(minutes) && isDigits(seconds);
};

const normalizeTimezone = (value: string): string => {
  if (value.endsWith(UTC_SUFFIX)) return `${value.slice(INDEX_START, -UTC_SUFFIX.length)}${UTC_MARKER}`;
  if (value.endsWith(UTC_MARKER)) return value;
  const signIndex = Math.max(value.lastIndexOf('+'), value.lastIndexOf('-'));
  if (signIndex === NOT_FOUND) return value;
  const offset = value.slice(signIndex);
  if (offset.length === TZ_OFFSET_LENGTH && !offset.includes(TIME_SEPARATOR)) {
    return `${value.slice(INDEX_START, signIndex)}${offset.slice(INDEX_START, TZ_COLON_INDEX)}:${offset.slice(TZ_COLON_INDEX)}`;
  }
  return offset.length === TZ_OFFSET_WITH_COLON_LENGTH ? value : value;
};

const normalizeDateTime = (value: string): string | null => {
  const separator = value.includes(DATE_TIME_SEPARATOR) ? DATE_TIME_SEPARATOR : SPACE_SEPARATOR;
  const [datePart, timePart] = value.split(separator);
  if (datePart === undefined || datePart === '' || timePart === undefined || timePart === '') return null;
  if (!isDateOnly(datePart) || !isTimePart(timePart)) return null;
  return normalizeTimezone(`${datePart}${DATE_TIME_SEPARATOR}${timePart}`);
};

export const buildTimeFormatter = (bucket: 'daily' | 'hourly'): Intl.DateTimeFormat =>
  new Intl.DateTimeFormat('en-US', bucket === 'hourly'
    ? { month: 'short', day: 'numeric', hour: 'numeric', timeZone: 'UTC' }
    : { month: 'short', day: 'numeric', timeZone: 'UTC' });

export const parseChartDate = (value: string): number | null => {
  const trimmed = value.trim();
  const numeric = Number(trimmed);
  if (!Number.isNaN(numeric)) return normalizeEpoch(numeric);
  if (isDateOnly(trimmed)) {
    const parsed = Date.parse(`${trimmed}T00:00:00Z`);
    return Number.isNaN(parsed) ? null : parsed;
  }
  const normalized = normalizeDateTime(trimmed);
  if (normalized !== null) {
    const parsed = Date.parse(normalized);
    return Number.isNaN(parsed) ? null : parsed;
  }
  const parsed = Date.parse(trimmed);
  return Number.isNaN(parsed) ? null : parsed;
};
