type SearchValue = string | string[] | undefined;

const EMPTY_COUNT = 0;
const EMPTY_STRING = '';

export function encodeFilterValue(value: string): string {
  return encodeURIComponent(value);
}

export function decodeFilterValue(value?: string): string | undefined {
  if (value === undefined || value === EMPTY_STRING) {
    return value;
  }
  try {
    return decodeURIComponent(value);
  } catch {
    return value;
  }
}

export function normalizeFilterValue(value: SearchValue): string[] {
  if (value === undefined || value === EMPTY_STRING) {
    return [];
  }
  if (Array.isArray(value)) {
    return value
      .map(item => decodeFilterValue(item))
      .filter((item): item is string => Boolean(item));
  }
  const decoded = decodeFilterValue(value);
  return decoded !== undefined && decoded !== EMPTY_STRING ? [decoded] : [];
}

export function addFilterValue(current: SearchValue, value: string): string[] {
  const normalized = normalizeFilterValue(current);
  if (normalized.includes(value)) {
    return normalized;
  }
  return [...normalized, value];
}

export function removeFilterValue(current: SearchValue, value: string): string[] | undefined {
  const next = normalizeFilterValue(current).filter(item => item !== value);
  return next.length > EMPTY_COUNT ? next : undefined;
}

export function updateFilterSearch<T extends Record<string, unknown>>(
  prev: T,
  key: keyof T,
  value: string[] | undefined
): T {
  // eslint-disable-next-line @typescript-eslint/no-unsafe-type-assertion -- generic record reconstruction requires assertion
  const result: T = Object.fromEntries(Object.entries(prev).filter(([k]) => k !== key)) as T;
  if (value === undefined || value.length === EMPTY_COUNT) {
    return result;
  }
  return Object.assign(result, { [key]: value });
}
