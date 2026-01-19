const EMPTY_COUNT = 0;
const EMPTY_STRING = '';

export const normalizeDomainInput = (value: string): string => value
  .replace(/^https?:\/\//v, '')
  .replace(/^www\./v, '')
  .replace(/\/.*$/v, '')
  .toLowerCase()
  .trim();

export const getNormalizedDomains = (values: string[]): string[] => {
  const normalized = values
    .map((domainValue) => normalizeDomainInput(domainValue))
    .filter((domainValue) => domainValue.length > EMPTY_COUNT);
  return Array.from(new Set(normalized));
};

export const normalizeIPInput = (value: string): string => value.trim();

export const getNormalizedBlockedIPs = (values: string[]): string[] => {
  const normalized = values
    .map((ipValue) => normalizeIPInput(ipValue))
    .filter((ipValue) => ipValue.length > EMPTY_COUNT);
  return Array.from(new Set(normalized));
};

export const normalizeCountryCodes = (values: string[]): string[] =>
  Array.from(
    new Set(
      values
        .map((code) => code.trim().toUpperCase())
        .filter((code) => code.length > EMPTY_COUNT)
    )
  ).sort((a, b) => a.localeCompare(b));

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
