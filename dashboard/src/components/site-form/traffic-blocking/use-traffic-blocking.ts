import { useMemo, useState, type KeyboardEvent } from 'react';
import { useQuery } from '@apollo/client/react';
import {
  CountriesByCodeDocument,
  CountryFieldsFragmentDoc,
  GeoIpCountriesDocument,
} from '@/gql/graphql';
import { useFragment as getFragmentData } from '@/gql/fragment-masking';
import { normalizeCountryCodesPreserveOrder, normalizeIPInput } from '@/components/site-form/utils';
import {
  COUNTRY_PAGE_OFFSET,
  COUNTRY_PAGE_SIZE,
  EMPTY_COUNT,
  EMPTY_STRING,
  FIRST_INDEX,
  MAX_COUNTRIES,
  MAX_IPS,
  SEARCH_MIN_LENGTH,
  SEARCH_SINGLE_MATCH_COUNT,
} from './constants';
import { blockedIPValues, buildBlockedIPEntries, isValidIP } from './traffic-blocking-utils';
import type { TrafficBlockingCardProps } from './types';

export function useTrafficBlocking({
  geoIPReady,
  initialBlockedCountries,
  initialBlockedIPs,
  onUpdateBlockedCountries,
  onUpdateBlockedIPs,
  savingBlockedCountries,
  savingBlockedIPs,
}: TrafficBlockingCardProps) {
  const blockedIPs = buildBlockedIPEntries(initialBlockedIPs);
  const blockedCountries = initialBlockedCountries;
  const [countrySearch, setCountrySearch] = useState('');
  const [ipActionError, setIpActionError] = useState('');
  const [countryActionError, setCountryActionError] = useState('');
  const [newIPValue, setNewIPValue] = useState('');
  const [newIPError, setNewIPError] = useState('');
  const normalizedBlockedCountries = useMemo(
    () => normalizeCountryCodesPreserveOrder(blockedCountries),
    [blockedCountries]
  );
  const trimmedCountrySearch = countrySearch.trim();
  const shouldSearchCountries = geoIPReady && trimmedCountrySearch.length >= SEARCH_MIN_LENGTH;
  const countryQueries = useCountryQueries(trimmedCountrySearch, shouldSearchCountries, normalizedBlockedCountries);
  const blockedIPCount = blockedIPValues(blockedIPs).length;
  const countryNameLookup = useMemo(
    () => new Map([...countryQueries.selectedCountries, ...countryQueries.geoIPCountries].map((country) => [country.code, country.name] as const)),
    [countryQueries.geoIPCountries, countryQueries.selectedCountries]
  );
  const availableCountries = useMemo(() => {
    const selected = new Set(normalizedBlockedCountries);
    return countryQueries.geoIPCountries.filter((country) => !selected.has(country.code));
  }, [countryQueries.geoIPCountries, normalizedBlockedCountries]);
  const matchingCountries = useMemo(
    () => (shouldSearchCountries ? availableCountries : []),
    [availableCountries, shouldSearchCountries]
  );
  const isUpdating = savingBlockedIPs || savingBlockedCountries;

  const handleAddBlockedCountry = async (code: string): Promise<void> => {
    if (isUpdating) return;
    const next = buildNextBlockedCountries(code, blockedCountries);
    if (next.length === blockedCountries.length) return;
    if (next.length > MAX_COUNTRIES) {
      setCountryActionError('Blocked country list can include up to 250 entries');
      return;
    }
    setCountryActionError('');
    await updateCountries(next);
  };

  const handleRemoveBlockedCountry = async (code: string): Promise<void> => {
    if (isUpdating) return;
    const normalized = code.trim().toUpperCase();
    const next = blockedCountries.filter((value) => value.trim().toUpperCase() !== normalized);
    if (next.length === blockedCountries.length) return;
    setCountryActionError('');
    await updateCountries(next);
  };

  const handleAddIP = async (): Promise<void> => {
    if (savingBlockedIPs) return;
    const trimmed = normalizeIPInput(newIPValue);
    const validationError = validateNewIP(trimmed, blockedIPs, blockedIPCount);
    if (validationError !== '') {
      setNewIPError(validationError);
      return;
    }
    setIpActionError('');
    setNewIPError('');
    try {
      await onUpdateBlockedIPs([...blockedIPValues(blockedIPs), trimmed]);
      setNewIPValue('');
    } catch (err) {
      setNewIPError(err instanceof Error ? err.message : 'Failed to save blocked IP');
    }
  };

  const handleRemoveBlockedIP = async (value: string): Promise<void> => {
    if (savingBlockedIPs) return;
    const normalized = value.trim();
    if (normalized === EMPTY_STRING) return;
    const next = blockedIPs.filter(({ value: blockedValue }) => blockedValue.trim() !== normalized);
    if (next.length === blockedIPs.length) return;
    setIpActionError('');
    try {
      await onUpdateBlockedIPs(blockedIPValues(next));
    } catch (err) {
      setIpActionError(err instanceof Error ? err.message : 'Failed to remove blocked IP');
    }
  };

  const handleCountrySearchKeyDown = (event: KeyboardEvent<HTMLInputElement>): void => {
    if (event.key !== 'Enter' || !shouldSearchCountries) return;
    const target = countrySearchTarget(trimmedCountrySearch, matchingCountries);
    if (target?.code === undefined || target.code === EMPTY_STRING) return;
    void handleAddBlockedCountry(target.code).then(() => setCountrySearch(''));
    event.preventDefault();
  };

  const updateCountries = async (next: string[]): Promise<void> => {
    try {
      await onUpdateBlockedCountries(next);
    } catch (err) {
      setCountryActionError(err instanceof Error ? err.message : 'Failed to update blocked countries');
    }
  };

  return {
    blockedCountries,
    blockedCountryCount: normalizedBlockedCountries.length,
    blockedIPCount,
    blockedIPs,
    countryActionError,
    countryNameLookup,
    countrySearch,
    geoIPCountriesLoading: countryQueries.geoIPCountriesLoading,
    handleAddBlockedCountry,
    handleAddIP,
    handleCountrySearchKeyDown,
    handleRemoveBlockedCountry,
    handleRemoveBlockedIP,
    ipActionError,
    matchingCountries,
    newIPError,
    newIPValue,
    normalizedBlockedCountries,
    setCountrySearch,
    setNewIPError,
    setNewIPValue,
    shouldSearchCountries,
    trimmedCountrySearch,
  };
}

function useCountryQueries(search: string, shouldSearch: boolean, codes: string[]) {
  const { data: geoIPCountriesData, loading: geoIPCountriesLoading } = useQuery(GeoIpCountriesDocument, {
    variables: { search, paging: { limit: COUNTRY_PAGE_SIZE, offset: COUNTRY_PAGE_OFFSET } },
    skip: !shouldSearch,
  });
  const { data: selectedCountriesData } = useQuery(CountriesByCodeDocument, {
    variables: { codes, paging: { limit: codes.length || 1, offset: COUNTRY_PAGE_OFFSET } },
    skip: codes.length === EMPTY_COUNT,
  });
  return {
    geoIPCountries: getFragmentData(CountryFieldsFragmentDoc, geoIPCountriesData?.geoIPCountries ?? []),
    geoIPCountriesLoading,
    selectedCountries: getFragmentData(CountryFieldsFragmentDoc, selectedCountriesData?.geoIPCountries ?? []),
  };
}

function buildNextBlockedCountries(code: string, countries: string[]): string[] {
  const trimmed = code.trim();
  if (trimmed === EMPTY_STRING) return countries;
  const normalized = trimmed.toUpperCase();
  const existing = new Set(countries.map((value) => value.trim().toUpperCase()));
  return existing.has(normalized) ? countries : [...countries, trimmed];
}

function validateNewIP(value: string, blockedIPs: ReturnType<typeof buildBlockedIPEntries>, count: number): string {
  if (value === EMPTY_STRING) return 'Enter a valid IP before saving.';
  if (!isValidIP(value)) return 'Enter a valid IP address.';
  if (blockedIPValues(blockedIPs).includes(value)) return 'That IP is already blocked.';
  if (count + 1 > MAX_IPS) return 'Blocked IP list can include up to 500 entries';
  return '';
}

function countrySearchTarget<T extends { code: string; name: string }>(queryRaw: string, countries: T[]): T | undefined {
  const query = queryRaw.toLowerCase();
  const exactMatch = countries.find(
    (country) => country.code.toLowerCase() === query || country.name.toLowerCase() === query
  );
  return exactMatch ?? (countries.length === SEARCH_SINGLE_MATCH_COUNT ? countries[FIRST_INDEX] : undefined);
}
