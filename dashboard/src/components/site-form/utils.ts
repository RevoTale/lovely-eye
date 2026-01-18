export const normalizeDomainInput = (value: string): string => value
  .replace(/^https?:\/\//, '')
  .replace(/^www\./, '')
  .replace(/\/.*$/, '')
  .toLowerCase()
  .trim();

export const getNormalizedDomains = (values: string[]): string[] => {
  const normalized = values
    .map((domainValue) => normalizeDomainInput(domainValue))
    .filter((domainValue) => domainValue.length > 0);
  return Array.from(new Set(normalized));
};

export const normalizeIPInput = (value: string): string => value.trim();

export const getNormalizedBlockedIPs = (values: string[]): string[] => {
  const normalized = values
    .map((ipValue) => normalizeIPInput(ipValue))
    .filter((ipValue) => ipValue.length > 0);
  return Array.from(new Set(normalized));
};

export const normalizeCountryCodes = (values: string[]): string[] => {
  return Array.from(new Set(values.map((code) => code.trim().toUpperCase()).filter((code) => code.length > 0)))
    .sort((a, b) => a.localeCompare(b));
};

export const normalizeCountryCodesPreserveOrder = (values: string[]): string[] => {
  const result: string[] = [];
  const seen = new Set<string>();
  values.forEach((value) => {
    const normalized = value.trim().toUpperCase();
    if (!normalized || seen.has(normalized)) {
      return;
    }
    seen.add(normalized);
    result.push(normalized);
  });
  return result;
};
