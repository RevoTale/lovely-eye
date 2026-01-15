export type DatePreset = '7d' | '30d' | '90d' | 'custom' | 'all';

export function isDatePreset(value: string | undefined): value is DatePreset {
  return value === '7d' || value === '30d' || value === '90d' || value === 'custom' || value === 'all';
}

export function formatDateInput(date: Date): string {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  return `${year}-${month}-${day}`;
}

export function presetToDates(preset: DatePreset, reference = new Date()): { fromDate: string; toDate: string } {
  const end = new Date(reference);
  const start = new Date(reference);
  if (preset === '7d') {
    start.setDate(end.getDate() - 6);
  } else if (preset === '90d') {
    start.setDate(end.getDate() - 89);
  } else {
    start.setDate(end.getDate() - 29);
  }
  return { fromDate: formatDateInput(start), toDate: formatDateInput(end) };
}

export function isValidDateInput(value: string): boolean {
  if (!value) return false;
  if (!/^\d{4}-\d{2}-\d{2}$/.test(value)) return false;
  const parsed = new Date(`${value}T00:00:00`);
  if (Number.isNaN(parsed.getTime())) return false;
  return formatDateInput(parsed) === value;
}

export function isValidTimeInput(value: string): boolean {
  if (!value) return false;
  if (!/^\d{2}:\d{2}$/.test(value)) return false;
  const [hourRaw, minuteRaw] = value.split(':');
  if (hourRaw === undefined || minuteRaw === undefined) return false;
  const hour = Number(hourRaw);
  const minute = Number(minuteRaw);
  if (Number.isNaN(hour) || Number.isNaN(minute)) return false;
  return hour >= 0 && hour <= 23 && minute >= 0 && minute <= 59;
}

export function parseDateInput(value: string): Date | undefined {
  if (!isValidDateInput(value)) return undefined;
  const parsed = new Date(`${value}T00:00:00`);
  if (Number.isNaN(parsed.getTime())) return undefined;
  return parsed;
}

export function normalizeTimeInput(value: string, fallback: string): string {
  if (!value) return fallback;
  const trimmed = value.trim();
  if (!trimmed) return fallback;
  if (trimmed.length >= 5) {
    return trimmed.slice(0, 5);
  }
  return fallback;
}

export function buildDateRange(
  fromDate: string,
  toDate: string,
  fromTime = '00:00',
  toTime = '23:59'
): { from: string; to: string } | undefined {
  if (!fromDate || !toDate) return undefined;
  const normalizeTime = (value: string, fallback: string): string => {
    if (!value) return fallback;
    if (value.length === 5) {
      return `${value}:00`;
    }
    return value;
  };
  const from = new Date(`${fromDate}T${normalizeTime(fromTime, '00:00:00')}`);
  const to = new Date(`${toDate}T${normalizeTime(toTime, '23:59:59')}`);
  if (Number.isNaN(from.getTime()) || Number.isNaN(to.getTime())) {
    return undefined;
  }
  return { from: from.toISOString(), to: to.toISOString() };
}
