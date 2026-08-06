import { useQuery } from '@apollo/client/react';
import { type KeyboardEvent, useMemo, useState } from 'react';
import {
  normalizeCountryCodesPreserveOrder,
  normalizeIPInput,
} from '@/features/site-settings/ui/utils';
import { useFragment as getFragmentData } from '@/shared/api/generated/fragment-masking';
import {
  CountriesByCodeDocument,
  CountryFieldsFragmentDoc,
  GeoIpCountriesDocument,
} from '@/shared/api/generated/graphql';
import {
  COUNTRY_PAGE_OFFSET,
  COUNTRY_PAGE_SIZE,
  EMPTY_COUNT,
  EMPTY_STRING,
  MAX_COUNTRIES,
  SEARCH_MIN_LENGTH,
} from './constants';
import {
  blockedIPValues,
  buildBlockedIPEntries,
  buildNextBlockedCountries,
  countrySearchTarget,
  validateNewIP,
} from './traffic-blocking-utils';
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
  const countryQueries = useCountryQueries(
    trimmedCountrySearch,
    shouldSearchCountries,
    normalizedBlockedCountries
  );
  const blockedIPCount = blockedIPValues(blockedIPs).length;
  const countryNameLookup = useMemo(
    () =>
      new Map(
        [...countryQueries.selectedCountries, ...countryQueries.geoIPCountries].map(
          (country) => [country.code, country.name] as const
        )
      ),
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
      setCountryActionError(
        err instanceof Error ? err.message : 'Failed to update blocked countries'
      );
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
  const { data: geoIPCountriesData, loading: geoIPCountriesLoading } = useQuery(
    GeoIpCountriesDocument,
    {
      variables: { search, paging: { limit: COUNTRY_PAGE_SIZE, offset: COUNTRY_PAGE_OFFSET } },
      skip: !shouldSearch,
    }
  );
  const { data: selectedCountriesData } = useQuery(CountriesByCodeDocument, {
    variables: { codes, paging: { limit: codes.length || 1, offset: COUNTRY_PAGE_OFFSET } },
    skip: codes.length === EMPTY_COUNT,
  });
  return {
    geoIPCountries: getFragmentData(
      CountryFieldsFragmentDoc,
      geoIPCountriesData?.geoIPCountries ?? []
    ),
    geoIPCountriesLoading,
    selectedCountries: getFragmentData(
      CountryFieldsFragmentDoc,
      selectedCountriesData?.geoIPCountries ?? []
    ),
  };
}
