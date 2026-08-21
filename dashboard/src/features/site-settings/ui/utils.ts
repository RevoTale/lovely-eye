const EMPTY_STRING = '';
const EMPTY_COUNT = 0;

export const normalizeIPInput = (value: string): string => value.trim();

export const getNormalizedBlockedIPs = (values: string[]): string[] => {
  const normalized = values
    .map((ipValue) => normalizeIPInput(ipValue))
    .filter((ipValue) => ipValue.length > EMPTY_COUNT);
  return Array.from(new Set(normalized));
};

export const normalizeCountryCodesPreserveOrder = (values: string[]): string[] => {
  const result: string[] = [];
  const seen = new Set<string>();
  values.forEach((value) => {
    const normalized = value.trim().toUpperCase();
    if (normalized === EMPTY_STRING || seen.has(normalized)) {
      return;
    }
    seen.add(normalized);
    result.push(normalized);
  });
  return result;
};
