import { getNormalizedBlockedIPs } from '@/features/site-settings/ui/utils';
import {
  EMPTY_COUNT,
  EMPTY_STRING,
  FIRST_INDEX,
  ID_OFFSET,
  IPV4_MAX_VALUE,
  IPV4_PARTS_COUNT,
  MAX_IPS,
  SEARCH_SINGLE_MATCH_COUNT,
} from './constants';
import type { BlockedIPEntry } from './types';

export const buildBlockedIPEntries = (values: string[]): BlockedIPEntry[] =>
  values.map((ip, index) => ({ id: String(index + ID_OFFSET), value: ip }));

export const blockedIPValues = (entries: BlockedIPEntry[]): string[] =>
  getNormalizedBlockedIPs(entries.map(({ value }) => value));

export const isValidIP = (value: string): boolean => {
  if (value === EMPTY_STRING) return false;
  if (value.includes(':')) {
    return /^(?:[0-9a-f]{0,4}:){2,7}[0-9a-f]{0,4}$/iv.test(value);
  }
  const parts = value.split('.');
  if (parts.length !== IPV4_PARTS_COUNT) return false;
  return parts.every((part) => {
    if (!/^\d{1,3}$/v.test(part)) return false;
    const num = Number(part);
    return num >= EMPTY_COUNT && num <= IPV4_MAX_VALUE;
  });
};

export const buildNextBlockedCountries = (code: string, countries: string[]): string[] => {
  const trimmed = code.trim();
  if (trimmed === EMPTY_STRING) return countries;
  const normalized = trimmed.toUpperCase();
  const existing = new Set(countries.map((value) => value.trim().toUpperCase()));
  return existing.has(normalized) ? countries : [...countries, trimmed];
};

export const validateNewIP = (
  value: string,
  blockedIPs: BlockedIPEntry[],
  count: number
): string => {
  if (value === EMPTY_STRING) return 'Enter a valid IP before saving.';
  if (!isValidIP(value)) return 'Enter a valid IP address.';
  if (blockedIPValues(blockedIPs).includes(value)) return 'That IP is already blocked.';
  if (count + 1 > MAX_IPS) return 'Blocked IP list can include up to 500 entries';
  return '';
};

export const countrySearchTarget = <T extends { code: string; name: string }>(
  queryRaw: string,
  countries: T[]
): T | undefined => {
  const query = queryRaw.toLowerCase();
  const exactMatch = countries.find(
    (country) => country.code.toLowerCase() === query || country.name.toLowerCase() === query
  );
  return (
    exactMatch ??
    (countries.length === SEARCH_SINGLE_MATCH_COUNT ? countries[FIRST_INDEX] : undefined)
  );
};
