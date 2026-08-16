import { format } from 'date-fns';
import { enUS } from 'date-fns/locale';
import type { DateTimeLocale, Period, TimePickerType } from './types';

type GetValidNumberConfig = { max: number; min?: number; loop?: boolean };
type GetValidArrowNumberConfig = { min: number; max: number; step: number };

export function getDateByType(date: Date | null, type: TimePickerType): string {
  if (!date) return '00';
  switch (type) {
    case 'minutes':
      return getValidMinuteOrSecond(String(date.getMinutes()));
    case 'seconds':
      return getValidMinuteOrSecond(String(date.getSeconds()));
    case 'hours':
      return getValidHour(String(date.getHours()));
    case '12hours':
      return getValid12Hour(String(display12HourValue(date.getHours())));
    default:
      return '00';
  }
}

export function getArrowByType(value: string, step: number, type: TimePickerType): string {
  switch (type) {
    case 'minutes':
    case 'seconds':
      return getValidArrowMinuteOrSecond(value, step);
    case 'hours':
      return getValidArrowHour(value, step);
    case '12hours':
      return getValidArrow12Hour(value, step);
    default:
      return '00';
  }
}

export function setDateByType(
  date: Date,
  value: string,
  type: TimePickerType,
  period?: Period
): Date {
  switch (type) {
    case 'minutes':
      return setMinutes(date, value);
    case 'seconds':
      return setSeconds(date, value);
    case 'hours':
      return setHours(date, value);
    case '12hours':
      return period ? set12Hours(date, value, period) : date;
    default:
      return date;
  }
}

export function display12HourValue(hours: number): string {
  if (hours === 0 || hours === 12) return '12';
  if (hours >= 22) return `${hours - 12}`;
  if (hours % 12 > 9) return `${hours}`;
  return `0${hours % 12}`;
}

export function genMonths(locale: DateTimeLocale): Array<{ value: number; label: string }> {
  return Array.from({ length: 12 }, (_, i) => ({
    value: i,
    label: format(new Date(2021, i), 'MMMM', { locale }),
  }));
}

export function genYears(yearRange = 50): Array<{ value: number; label: string }> {
  const today = new Date();
  return Array.from({ length: yearRange * 2 + 1 }, (_, i) => ({
    value: today.getFullYear() - yearRange + i,
    label: (today.getFullYear() - yearRange + i).toString(),
  }));
}

export function normalizedLocale(locale: Partial<DateTimeLocale> | undefined): DateTimeLocale {
  const { options, localize, formatLong } = locale ?? {};
  if (options && localize && formatLong) {
    return { ...enUS, options, localize, formatLong };
  }
  return enUS;
}

function getValidNumber(
  value: string,
  { max, min = 0, loop = false }: GetValidNumberConfig
): string {
  let numericValue = Number.parseInt(value, 10);
  if (Number.isNaN(numericValue)) return '00';
  if (!loop) {
    if (numericValue > max) numericValue = max;
    if (numericValue < min) numericValue = min;
  } else {
    if (numericValue > max) numericValue = min;
    if (numericValue < min) numericValue = max;
  }
  return numericValue.toString().padStart(2, '0');
}

function getValidHour(value: string): string {
  return /^(0[0-9]|1[0-9]|2[0-3])$/u.test(value) ? value : getValidNumber(value, { max: 23 });
}

function getValid12Hour(value: string): string {
  return /^(0[1-9]|1[0-2])$/u.test(value) ? value : getValidNumber(value, { min: 1, max: 12 });
}

function getValidMinuteOrSecond(value: string): string {
  return /^[0-5][0-9]$/u.test(value) ? value : getValidNumber(value, { max: 59 });
}

function getValidArrowNumber(value: string, config: GetValidArrowNumberConfig): string {
  const numericValue = Number.parseInt(value, 10);
  if (Number.isNaN(numericValue)) return '00';
  return getValidNumber(String(numericValue + config.step), { ...config, loop: true });
}

function getValidArrowHour(value: string, step: number): string {
  return getValidArrowNumber(value, { min: 0, max: 23, step });
}

function getValidArrow12Hour(value: string, step: number): string {
  return getValidArrowNumber(value, { min: 1, max: 12, step });
}

function getValidArrowMinuteOrSecond(value: string, step: number): string {
  return getValidArrowNumber(value, { min: 0, max: 59, step });
}

function setMinutes(date: Date, value: string): Date {
  date.setMinutes(Number.parseInt(getValidMinuteOrSecond(value), 10));
  return date;
}

function setSeconds(date: Date, value: string): Date {
  date.setSeconds(Number.parseInt(getValidMinuteOrSecond(value), 10));
  return date;
}

function setHours(date: Date, value: string): Date {
  date.setHours(Number.parseInt(getValidHour(value), 10));
  return date;
}

function set12Hours(date: Date, value: string, period: Period): Date {
  const hour = Number.parseInt(getValid12Hour(value), 10);
  date.setHours(convert12HourTo24Hour(hour, period));
  return date;
}

function convert12HourTo24Hour(hour: number, period: Period): number {
  if (period === 'PM') return hour <= 11 ? hour + 12 : hour;
  if (period === 'AM') return hour === 12 ? 0 : hour;
  return hour;
}
