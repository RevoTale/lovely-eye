const DOMAIN_PATTERN =
  /^[a-z0-9](?:[a-z0-9\-]{0,61}[a-z0-9])?(?:\.[a-z0-9](?:[a-z0-9\-]{0,61}[a-z0-9])?)*$/v;

export function normalizeDomain(value: string): string {
  return value
    .trim()
    .toLowerCase()
    .replace(/^https?:\/\//v, '')
    .replace(/^www\./v, '')
    .replace(/\/.*$/v, '');
}

export function normalizeDomains(values: string[]): string[] {
  return Array.from(new Set(values.map(normalizeDomain).filter((value) => value !== '')));
}

export function validateDomains(values: string[]): { domains: string[]; error: string | null } {
  const domains = normalizeDomains(values);
  if (domains.length === 0) {
    return { domains, error: 'At least one domain is required' };
  }
  if (!domains.every((domain) => DOMAIN_PATTERN.test(domain))) {
    return { domains: [], error: 'Please enter valid domains (e.g., example.com)' };
  }
  return { domains, error: null };
}
