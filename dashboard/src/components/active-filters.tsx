import { useQuery } from '@apollo/client/react';
import ActiveFilterChip, { type ActiveFilterField } from '@/components/active-filter-chip';
import { Badge } from '@/components/ui';
import { CountryFieldsFragmentDoc, CountriesByCodeDocument } from '@/gql/graphql';
import { useFragment as getFragmentData } from '@/gql/fragment-masking';
import { normalizeFilterValue } from '@/lib/filter-utils';
import { Link } from '@/router';

interface FilterSearch {
  referrer?: string | string[] | undefined;
  browser?: string | string[] | undefined;
  device?: string | string[] | undefined;
  os?: string | string[] | undefined;
  page?: string | string[] | undefined;
  country?: string | string[] | undefined;
  eventName?: string | string[] | undefined;
  eventPath?: string | string[] | undefined;
}

interface ActiveFiltersProps {
  siteId: string;
  search: FilterSearch;
}

const EMPTY_COUNT = 0;

export const ActiveFilters = ({ siteId, search }: ActiveFiltersProps): React.ReactNode => {
  const referrers = normalizeFilterValue(search.referrer);
  const browsers = normalizeFilterValue(search.browser);
  const devices = normalizeFilterValue(search.device);
  const operatingSystems = normalizeFilterValue(search.os);
  const pages = normalizeFilterValue(search.page);
  const countries = normalizeFilterValue(search.country);
  const normalizedCountryCodes = Array.from(
    new Set(
      countries
        .map((country) => country.trim().toUpperCase())
        .filter((country) => country.length > EMPTY_COUNT)
    )
  );
  const eventNames = normalizeFilterValue(search.eventName);
  const eventPaths = normalizeFilterValue(search.eventPath);
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
  const countryNameLookup = new Map(lookedUpCountries.map((country) => [country.code, country.name] as const));
  const hasFilters =
    referrers.length > EMPTY_COUNT ||
    browsers.length > EMPTY_COUNT ||
    devices.length > EMPTY_COUNT ||
    operatingSystems.length > EMPTY_COUNT ||
    pages.length > EMPTY_COUNT ||
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
    field: ActiveFilterField;
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
    <div className="flex items-center gap-2 flex-wrap">
      <span className="text-sm text-muted-foreground">Filtered by:</span>
      {filterGroups.flatMap(({ displayValue, field, label, values }) =>
        values.map((value) => (
          <ActiveFilterChip key={`${field}-${value}`} field={field} label={label} siteId={siteId} value={value}>
            {displayValue?.(value) ?? value}
          </ActiveFilterChip>
        ))
      )}
      <Link to="/sites/$siteId" params={{ siteId }} search={{}}>
        <Badge variant="outline" className="cursor-pointer hover:bg-accent text-xs">
          Clear all
        </Badge>
      </Link>
    </div>
  );
}
