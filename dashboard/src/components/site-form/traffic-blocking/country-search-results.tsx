import type { FunctionComponent } from 'react';
import type { CountryFieldsFragment } from '@/gql/graphql';
import { EMPTY_COUNT, EMPTY_STRING, FIRST_INDEX, MAX_COUNTRY_MATCHES } from './constants';

interface CountrySearchResultsProps {
  geoIPCountriesLoading: boolean;
  matchingCountries: CountryFieldsFragment[];
  shouldSearchCountries: boolean;
  trimmedCountrySearch: string;
  onAddBlockedCountry: (code: string) => Promise<void>;
  setCountrySearch: (value: string) => void;
}

const CountrySearchResults: FunctionComponent<CountrySearchResultsProps> = ({
  geoIPCountriesLoading,
  matchingCountries,
  onAddBlockedCountry,
  setCountrySearch,
  shouldSearchCountries,
  trimmedCountrySearch,
}) => (
  <div className="space-y-1 rounded-lg border bg-muted/30 p-2">
    {geoIPCountriesLoading ? (
      <p className="text-xs text-muted-foreground">Searching...</p>
    ) : matchingCountries.length === EMPTY_COUNT ? (
      <p className="text-xs text-muted-foreground">
        {shouldSearchCountries ? `No matches for "${trimmedCountrySearch}".` : 'Search results will appear here.'}
      </p>
    ) : (
      matchingCountries.slice(FIRST_INDEX, MAX_COUNTRY_MATCHES).map((country) => {
        if (country.code === EMPTY_STRING) return null;
        return (
          <button
            key={country.code}
            type="button"
            className="flex w-full items-center justify-between rounded-md px-2 py-1 text-left text-sm hover:bg-accent"
            onClick={() => void onAddBlockedCountry(country.code).then(() => setCountrySearch(''))}
          >
            <span>{country.name}</span>
            <span className="text-xs text-muted-foreground">{country.code}</span>
          </button>
        );
      })
    )}
  </div>
);

export default CountrySearchResults;
