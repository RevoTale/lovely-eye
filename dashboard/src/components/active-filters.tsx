
import { useQuery } from '@apollo/client/react';
import { Badge } from '@/components/ui';
import { CountryFieldsFragmentDoc, CountriesByCodeDocument } from '@/gql/graphql';
import { useFragment as getFragmentData } from '@/gql/fragment-masking';
import { Link } from '@/router';
import { normalizeFilterValue, removeFilterValue, updateFilterSearch } from '@/lib/filter-utils';

interface FilterSearch {
  referrer?: string | string[] | undefined;
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

  return (
    <div className="flex items-center gap-2 flex-wrap">
      <span className="text-sm text-muted-foreground">Filtered by:</span>
      {referrers.map((referrer) => (
        <Link
          key={`referrer-${referrer}`}
          to="/sites/$siteId"
          params={{ siteId }}
          search={(prev) => ({
            ...updateFilterSearch(prev, 'referrer', removeFilterValue(prev.referrer, referrer)),
          })}
        >
          <Badge variant="secondary" className="flex items-center gap-1 cursor-pointer hover:bg-secondary/80">
            <span className="text-xs">Referrer: {referrer}</span>
            <span className="ml-1 text-xs">×</span>
          </Badge>
        </Link>
      ))}
      {devices.map((device) => (
        <Link
          key={`device-${device}`}
          to="/sites/$siteId"
          params={{ siteId }}
          search={(prev) => ({
            ...updateFilterSearch(prev, 'device', removeFilterValue(prev.device, device)),
          })}
        >
          <Badge variant="secondary" className="flex items-center gap-1 cursor-pointer hover:bg-secondary/80">
            <span className="text-xs">Device: {device}</span>
            <span className="ml-1 text-xs">×</span>
          </Badge>
        </Link>
      ))}
      {operatingSystems.map((os) => (
        <Link
          key={`os-${os}`}
          to="/sites/$siteId"
          params={{ siteId }}
          search={(prev) => ({
            ...updateFilterSearch(prev, 'os', removeFilterValue(prev.os, os)),
          })}
        >
          <Badge variant="secondary" className="flex items-center gap-1 cursor-pointer hover:bg-secondary/80">
            <span className="text-xs">OS: {os}</span>
            <span className="ml-1 text-xs">×</span>
          </Badge>
        </Link>
      ))}
      {pages.map((page) => (
        <Link
          key={`page-${page}`}
          to="/sites/$siteId"
          params={{ siteId }}
          search={(prev) => ({
            ...updateFilterSearch(prev, 'page', removeFilterValue(prev.page, page)),
          })}
        >
          <Badge variant="secondary" className="flex items-center gap-1 cursor-pointer hover:bg-secondary/80">
            <span className="text-xs">Page: {page}</span>
            <span className="ml-1 text-xs">×</span>
          </Badge>
        </Link>
      ))}
      {countries.map((country) => (
        <Link
          key={`country-${country}`}
          to="/sites/$siteId"
          params={{ siteId }}
          search={(prev) => ({
            ...updateFilterSearch(prev, 'country', removeFilterValue(prev.country, country)),
          })}
        >
          <Badge variant="secondary" className="flex items-center gap-1 cursor-pointer hover:bg-secondary/80">
            <span className="text-xs">Country: {getCountryDisplayName(country)}</span>
            <span className="ml-1 text-xs">×</span>
          </Badge>
        </Link>
      ))}
      {eventNames.map((eventName) => (
        <Link
          key={`event-name-${eventName}`}
          to="/sites/$siteId"
          params={{ siteId }}
          search={(prev) => ({
            ...updateFilterSearch(prev, 'eventName', removeFilterValue(prev.eventName, eventName)),
          })}
        >
          <Badge variant="secondary" className="flex items-center gap-1 cursor-pointer hover:bg-secondary/80">
            <span className="text-xs">Event: {eventName}</span>
            <span className="ml-1 text-xs">×</span>
          </Badge>
        </Link>
      ))}
      {eventPaths.map((eventPath) => (
        <Link
          key={`event-path-${eventPath}`}
          to="/sites/$siteId"
          params={{ siteId }}
          search={(prev) => ({
            ...updateFilterSearch(prev, 'eventPath', removeFilterValue(prev.eventPath, eventPath)),
          })}
        >
          <Badge variant="secondary" className="flex items-center gap-1 cursor-pointer hover:bg-secondary/80">
            <span className="text-xs">Event Path: {eventPath}</span>
            <span className="ml-1 text-xs">×</span>
          </Badge>
        </Link>
      ))}
      <Link to="/sites/$siteId" params={{ siteId }} search={{}}>
        <Badge variant="outline" className="cursor-pointer hover:bg-accent text-xs">
          Clear all
        </Badge>
      </Link>
    </div>
  );
}
