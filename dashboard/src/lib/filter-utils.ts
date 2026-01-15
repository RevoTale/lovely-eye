type SearchValue = string | string[] | undefined;

export function encodeFilterValue(value: string): string {
  return encodeURIComponent(value);
}

export function decodeFilterValue(value?: string): string | undefined {
  if (!value) {
    return value;
  }
  try {
    return decodeURIComponent(value);
  } catch {
    return value;
  }
}

export function normalizeFilterValue(value: SearchValue): string[] {
  if (!value) {
    return [];
  }
  if (Array.isArray(value)) {
    return value
      .map(item => decodeFilterValue(item))
      .filter((item): item is string => Boolean(item));
  }
  const decoded = decodeFilterValue(value);
  return decoded ? [decoded] : [];
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
  return next.length > 0 ? next : undefined;
}

export function updateFilterSearch<T extends Record<string, unknown>>(
  prev: T,
  key: keyof T,
  value: string[] | undefined
): T {
  const { [key]: _ignored, ...rest } = prev;
  if (!value || value.length === 0) {
    return rest as T;
  }
  return { ...(rest as T), [key]: value };
}
