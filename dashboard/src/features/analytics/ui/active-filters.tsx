import { useQuery } from '@apollo/client/react';
import { Link } from '@tanstack/react-router';
import {
  type AnalyticsFilterKey,
  type AnalyticsSearch,
  clearAnalyticsFilters,
  removeAnalyticsFilter,
  setPagePathSearch,
} from '@/features/analytics/model/analytics-search';
import ActiveFilterChip from '@/features/analytics/ui/active-filter-chip';
import { useFragment as getFragmentData } from '@/shared/api/generated/fragment-masking';
import { CountriesByCodeDocument, CountryFieldsFragmentDoc } from '@/shared/api/generated/graphql';
import { Badge } from '@/shared/ui/badge';

interface ActiveFiltersProps {
  siteId: string;
  search: AnalyticsSearch;
}

const EMPTY_COUNT = 0;

export const ActiveFilters = ({ siteId, search }: ActiveFiltersProps): React.ReactNode => {
  const referrers = search.referrer ?? [];
  const browsers = search.browser ?? [];
  const devices = search.device ?? [];
  const operatingSystems = search.os ?? [];
  const pages = search.page ?? [];
  const pagePathContains = search.pagePathContains;
  const countries = search.country ?? [];
  const normalizedCountryCodes = Array.from(
    new Set(
      countries
        .map((country) => country.trim().toUpperCase())
        .filter((country) => country.length > EMPTY_COUNT)
    )
  );
  const eventNames = search.eventName ?? [];
  const eventPaths = search.eventPath ?? [];
  const { data: countryLookupData } = useQuery(CountriesByCodeDocument, {
    variables: {
      codes: normalizedCountryCodes,
      paging: {
        limit: normalizedCountryCodes.length || 1,
        offset: 0,
      },
    },
    skip: normalizedCountryCodes.length === EMPTY_COUNT,
  });
  const lookedUpCountries = getFragmentData(
    CountryFieldsFragmentDoc,
    countryLookupData?.geoIPCountries ?? []
  );
  const countryNameLookup = new Map(
    lookedUpCountries.map((country) => [country.code, country.name] as const)
  );
  const hasFilters =
    referrers.length > EMPTY_COUNT ||
    browsers.length > EMPTY_COUNT ||
    devices.length > EMPTY_COUNT ||
    operatingSystems.length > EMPTY_COUNT ||
    pages.length > EMPTY_COUNT ||
    pagePathContains !== undefined ||
    countries.length > EMPTY_COUNT ||
    eventNames.length > EMPTY_COUNT ||
    eventPaths.length > EMPTY_COUNT;
  if (!hasFilters) {
    return null;
  }

  const getCountryDisplayName = (country: string): string => {
    const normalizedCountryCode = country.trim().toUpperCase();
    return countryNameLookup.get(normalizedCountryCode) ?? country;
  };
  const filterGroups: Array<{
    field: AnalyticsFilterKey;
    label: string;
    values: string[];
    displayValue?: (value: string) => string;
  }> = [
    { field: 'referrer', label: 'Referrer', values: referrers },
    { field: 'browser', label: 'Browser', values: browsers },
    { field: 'device', label: 'Device', values: devices },
    { field: 'os', label: 'OS', values: operatingSystems },
    { field: 'page', label: 'Page', values: pages },
    {
      field: 'country',
      label: 'Country',
      values: countries,
      displayValue: getCountryDisplayName,
    },
    { field: 'eventName', label: 'Event', values: eventNames },
    { field: 'eventPath', label: 'Event Path', values: eventPaths },
  ];

  return (
    <div className='flex items-center gap-2 flex-wrap'>
      <span className='text-sm text-muted-foreground'>Filtered by:</span>
      {filterGroups.flatMap(({ displayValue, field, label, values }) =>
        values.map((value) => (
          <ActiveFilterChip
            key={`${field}-${value}`}
            label={label}
            siteId={siteId}
            remove={(current) => removeAnalyticsFilter(current, field, value)}
          >
            {displayValue?.(value) ?? value}
          </ActiveFilterChip>
        ))
      )}
      {pagePathContains === undefined ? null : (
        <ActiveFilterChip
          label='Page contains'
          siteId={siteId}
          remove={(current) => setPagePathSearch(current, '')}
        >
          {pagePathContains}
        </ActiveFilterChip>
      )}
      <Link
        to='/sites/$siteId/analytics'
        params={{ siteId }}
        search={(current) => clearAnalyticsFilters(current)}
      >
        <Badge variant='outline' className='cursor-pointer hover:bg-accent text-xs'>
          Clear all
        </Badge>
      </Link>
    </div>
  );
};
