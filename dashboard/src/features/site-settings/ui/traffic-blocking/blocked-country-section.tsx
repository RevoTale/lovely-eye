import type { FunctionComponent, KeyboardEventHandler } from 'react';
import type { CountryFieldsFragment } from '@/shared/api/generated/graphql';
import { Input } from '@/shared/ui/input';
import { Label } from '@/shared/ui/label';
import BlockedCountryBadge from './blocked-country-badge';
import { EMPTY_COUNT, EMPTY_STRING, MAX_COUNTRIES, SEARCH_MIN_LENGTH } from './constants';
import CountrySearchResults from './country-search-results';

interface BlockedCountrySectionProps {
  blockedCountryCount: number;
  blockedCountries: string[];
  countryActionError: string;
  countryNameLookup: Map<string, string>;
  countrySearch: string;
  geoIPCountriesLoading: boolean;
  geoIPReady: boolean;
  matchingCountries: CountryFieldsFragment[];
  normalizedBlockedCountries: string[];
  savingBlockedCountries: boolean;
  shouldSearchCountries: boolean;
  trimmedCountrySearch: string;
  onAddBlockedCountry: (code: string) => Promise<void>;
  onCountrySearchKeyDown: KeyboardEventHandler<HTMLInputElement>;
  onRemoveBlockedCountry: (code: string) => void;
  setCountrySearch: (value: string) => void;
}

const BlockedCountrySection: FunctionComponent<BlockedCountrySectionProps> = ({
  blockedCountryCount,
  blockedCountries,
  countryActionError,
  countryNameLookup,
  countrySearch,
  geoIPCountriesLoading,
  geoIPReady,
  matchingCountries,
  normalizedBlockedCountries,
  onAddBlockedCountry,
  onCountrySearchKeyDown,
  onRemoveBlockedCountry,
  savingBlockedCountries,
  setCountrySearch,
  shouldSearchCountries,
  trimmedCountrySearch,
}) => (
  <div className='space-y-3'>
    <div className='flex items-center justify-between'>
      <Label>Blocked Countries</Label>
      <span className='text-xs text-muted-foreground'>
        {blockedCountryCount}/{MAX_COUNTRIES}
      </span>
    </div>
    <div className='flex flex-wrap gap-2'>
      {blockedCountries.length === EMPTY_COUNT ? (
        <span className='text-xs text-muted-foreground'>No blocked countries yet.</span>
      ) : (
        normalizedBlockedCountries.map((code) => (
          <BlockedCountryBadge
            key={code}
            code={code}
            countryNameLookup={countryNameLookup}
            onRemoveBlockedCountry={onRemoveBlockedCountry}
            savingBlockedCountries={savingBlockedCountries}
          />
        ))
      )}
    </div>
    <div className='space-y-2'>
      <Label htmlFor='country-search'>Search countries</Label>
      <Input
        id='country-search'
        placeholder='Start typing a country name or code'
        value={countrySearch}
        onChange={(event) => setCountrySearch(event.currentTarget.value)}
        onKeyDown={onCountrySearchKeyDown}
        disabled={!geoIPReady}
      />
      {geoIPReady &&
      trimmedCountrySearch.length > EMPTY_COUNT &&
      trimmedCountrySearch.length < SEARCH_MIN_LENGTH ? (
        <p className='text-xs text-muted-foreground'>Type at least 2 characters to search.</p>
      ) : null}
    </div>
    <CountrySearchResults
      geoIPCountriesLoading={geoIPCountriesLoading}
      matchingCountries={matchingCountries}
      onAddBlockedCountry={onAddBlockedCountry}
      setCountrySearch={setCountrySearch}
      shouldSearchCountries={shouldSearchCountries}
      trimmedCountrySearch={trimmedCountrySearch}
    />
    {geoIPReady ? null : (
      <p className='text-xs text-muted-foreground'>
        GeoIP database is not ready. Download it to manage country blocking.
      </p>
    )}
    {savingBlockedCountries ? (
      <p className='text-xs text-muted-foreground'>Updating blocked countries...</p>
    ) : null}
    {countryActionError === EMPTY_STRING ? null : (
      <p className='text-xs text-destructive'>{countryActionError}</p>
    )}
  </div>
);

export default BlockedCountrySection;
