
import { useMemo, useState, type KeyboardEvent, type ReactElement } from 'react';
import { useQuery } from '@apollo/client/react';
import { GeoIpCountriesDocument, TrafficBlockingCountryFieldsFragmentDoc } from '@/gql/graphql';
import { useFragment as getFragmentData } from '@/gql/fragment-masking';
import { Badge, Button, Card, CardContent, CardDescription, CardHeader, CardTitle, Input, Label } from '@/components/ui';
import { Loader2, Shield, X } from 'lucide-react';
import { getNormalizedBlockedIPs, normalizeCountryCodesPreserveOrder, normalizeIPInput } from '@/components/site-form/utils';

interface BlockedIPEntry {
  id: string;
  value: string;
}

interface TrafficBlockingCardProps {
  initialBlockedIPs: string[];
  initialBlockedCountries: string[];
  savingBlockedIPs: boolean;
  savingBlockedCountries: boolean;
  geoIPReady: boolean;
  onUpdateBlockedIPs: (blockedIPs: string[]) => Promise<void>;
  onUpdateBlockedCountries: (blockedCountries: string[]) => Promise<void>;
}

export const TrafficBlockingCard = ({
  initialBlockedIPs,
  initialBlockedCountries,
  savingBlockedIPs,
  savingBlockedCountries,
  geoIPReady,
  onUpdateBlockedIPs,
  onUpdateBlockedCountries,
}: TrafficBlockingCardProps): ReactElement => {
  const EMPTY_COUNT = 0;
  const EMPTY_STRING = '';
  const FIRST_INDEX = 0;
  const ID_OFFSET = 1;
  const IPV4_PARTS_COUNT = 4;
  const IPV4_MAX_VALUE = 255;
  const MAX_COUNTRY_MATCHES = 8;
  const MAX_COUNTRIES = 250;
  const MAX_IPS = 500;
const SEARCH_MIN_LENGTH = 2;
const COUNTRY_PAGE_SIZE = 100;
const COUNTRY_PAGE_OFFSET = 0;
  const SEARCH_SINGLE_MATCH_COUNT = 1;

  const buildBlockedIPEntries = (values: string[]): BlockedIPEntry[] =>
    values.map((ip, index) => ({ id: String(index + ID_OFFSET), value: ip }));

  const blockedIPs = useMemo(
    () => buildBlockedIPEntries(initialBlockedIPs),
    [initialBlockedIPs]
  );
  const blockedCountries = useMemo(
    () => initialBlockedCountries,
    [initialBlockedCountries]
  );
  const [countrySearch, setCountrySearch] = useState('');
  const [ipActionError, setIpActionError] = useState('');
  const [countryActionError, setCountryActionError] = useState('');
  const [newIPValue, setNewIPValue] = useState('');
  const [newIPError, setNewIPError] = useState('');

  const trimmedCountrySearch = countrySearch.trim();
  const shouldSearchCountries = geoIPReady && trimmedCountrySearch.length >= SEARCH_MIN_LENGTH;
  const { data: geoIPCountriesData, loading: geoIPCountriesLoading } = useQuery(GeoIpCountriesDocument, {
    variables: {
      search: trimmedCountrySearch,
      paging: {
        limit: COUNTRY_PAGE_SIZE,
        offset: COUNTRY_PAGE_OFFSET,
      },
    },
    skip: !shouldSearchCountries,
  });

  const geoIPCountries = getFragmentData(
    TrafficBlockingCountryFieldsFragmentDoc,
    geoIPCountriesData?.geoIPCountries ?? []
  );
  const blockedIPCount = useMemo(
    () => getNormalizedBlockedIPs(blockedIPs.map(({ value }) => value)).length,
    [blockedIPs]
  );
  const normalizedBlockedCountries = useMemo(
    () => normalizeCountryCodesPreserveOrder(blockedCountries),
    [blockedCountries]
  );
  const { length: blockedCountryCount } = normalizedBlockedCountries;
  const countryNameLookup = useMemo(
    () => new Map(geoIPCountries.map((country) => [country.code, country.name] as const)),
    [geoIPCountries]
  );
  const isUpdating = savingBlockedIPs || savingBlockedCountries;

  const availableCountries = useMemo(() => {
    const selected = new Set(normalizedBlockedCountries);
    return geoIPCountries.filter((country) => !selected.has(country.code));
  }, [geoIPCountries, normalizedBlockedCountries]);
  const matchingCountries = useMemo(
    () => (shouldSearchCountries ? availableCountries : []),
    [availableCountries, shouldSearchCountries]
  );

  const buildNextBlockedCountries = (code: string): string[] => {
    const trimmed = code.trim();
    if (trimmed === EMPTY_STRING) return blockedCountries;
    const normalized = trimmed.toUpperCase();
    const existing = new Set(blockedCountries.map((value) => value.trim().toUpperCase()));
    if (existing.has(normalized)) {
      return blockedCountries;
    }
    return [...blockedCountries, trimmed];
  };

  const handleAddBlockedCountry = async (code: string): Promise<void> => {
    if (isUpdating) return;
    const next = buildNextBlockedCountries(code);
    if (next.length === blockedCountries.length) {
      return;
    }
    if (next.length > MAX_COUNTRIES) {
      setCountryActionError('Blocked country list can include up to 250 entries');
      return;
    }
    setCountryActionError('');
    try {
      await onUpdateBlockedCountries(next);
    } catch (err) {
      setCountryActionError(err instanceof Error ? err.message : 'Failed to update blocked countries');
    }
  };

  const handleRemoveBlockedCountry = async (code: string): Promise<void> => {
    if (isUpdating) return;
    const normalized = code.trim().toUpperCase();
    const next = blockedCountries.filter((value) => value.trim().toUpperCase() !== normalized);
    if (next.length === blockedCountries.length) {
      return;
    }
    setCountryActionError('');
    try {
      await onUpdateBlockedCountries(next);
    } catch (err) {
      setCountryActionError(err instanceof Error ? err.message : 'Failed to update blocked countries');
    }
  };

  const isValidIP = (value: string): boolean => {
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

  const handleAddIP = async (): Promise<void> => {
    if (savingBlockedIPs) return;
    const trimmed = normalizeIPInput(newIPValue);
    if (trimmed === EMPTY_STRING) {
      setNewIPError('Enter a valid IP before saving.');
      return;
    }
    if (!isValidIP(trimmed)) {
      setNewIPError('Enter a valid IP address.');
      return;
    }
    const existing = getNormalizedBlockedIPs(blockedIPs.map(({ value }) => value));
    if (existing.includes(trimmed)) {
      setNewIPError('That IP is already blocked.');
      return;
    }
    const next = [...existing, trimmed];
    if (next.length > MAX_IPS) {
      setNewIPError('Blocked IP list can include up to 500 entries');
      return;
    }
    setIpActionError('');
    setNewIPError('');
    try {
      await onUpdateBlockedIPs(next);
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
    if (next.length === blockedIPs.length) {
      return;
    }
    setIpActionError('');
    try {
      await onUpdateBlockedIPs(getNormalizedBlockedIPs(next.map(({ value: blockedValue }) => blockedValue)));
    } catch (err) {
      setIpActionError(err instanceof Error ? err.message : 'Failed to remove blocked IP');
    }
  };

  const handleCountrySearchKeyDown = (event: KeyboardEvent<HTMLInputElement>): void => {
    if (event.key !== 'Enter') return;
    if (!shouldSearchCountries) return;
    const query = trimmedCountrySearch.toLowerCase();
    const exactMatch = matchingCountries.find(
      (country) => country.code.toLowerCase() === query || country.name.toLowerCase() === query,
    );
    const singleMatch =
      matchingCountries.length === SEARCH_SINGLE_MATCH_COUNT
        ? matchingCountries[FIRST_INDEX]
        : undefined;
    const target = exactMatch ?? singleMatch;
    if (target?.code !== undefined && target.code !== EMPTY_STRING) {
      const { code } = target;
      void (async () => {
        await handleAddBlockedCountry(code);
        setCountrySearch('');
      })();
      event.preventDefault();
    }
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <div className="h-8 w-8 rounded-lg bg-primary/10 flex items-center justify-center">
            <Shield className="h-4 w-4 text-primary" />
          </div>
          Traffic Blocking
        </CardTitle>
        <CardDescription>
          Block specific IPs and countries from being tracked
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="space-y-3">
          <div className="flex items-center justify-between">
            <Label>Blocked IPs</Label>
            <span className="text-xs text-muted-foreground">
              {blockedIPCount}/{MAX_IPS}
            </span>
          </div>
          <div className="space-y-2">
            {blockedIPs.length === EMPTY_COUNT ? (
              <span className="text-xs text-muted-foreground">
                No blocked IPs yet.
              </span>
            ) : (
              blockedIPs.map((entry) => (
                <div key={entry.id} className="flex items-center gap-2">
                  <Input
                    value={entry.value}
                    readOnly
                  />
                  <Button
                    type="button"
                    variant="outline"
                    size="icon"
                    disabled={savingBlockedIPs}
                    onClick={() => {
                      void handleRemoveBlockedIP(entry.value);
                    }}
                    aria-label="Remove blocked IP"
                  >
                    <X className="h-4 w-4" />
                  </Button>
                </div>
              ))
            )}
            <div className="space-y-1">
              <div className="flex items-center gap-2">
                <Input
                  placeholder="203.0.113.10"
                  value={newIPValue}
                  onChange={(e) => {
                    const { currentTarget } = e;
                    const { value } = currentTarget;
                    setNewIPValue(normalizeIPInput(value));
                    if (newIPError !== EMPTY_STRING) {
                      setNewIPError(EMPTY_STRING);
                    }
                  }}
                  disabled={savingBlockedIPs}
                />
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => {
                    void handleAddIP();
                  }}
                  disabled={savingBlockedIPs || blockedIPCount >= MAX_IPS}
                >
                  {savingBlockedIPs ? (
                    <Loader2 className="h-4 w-4 animate-spin" />
                  ) : (
                    'Add IP'
                  )}
                </Button>
              </div>
              {newIPError === EMPTY_STRING ? null : (
                <span className="text-xs text-destructive">{newIPError}</span>
              )}
            </div>
          </div>
          {savingBlockedIPs ? (
            <p className="text-xs text-muted-foreground">
              Updating blocked IPs...
            </p>
          ) : null}
          {ipActionError === EMPTY_STRING ? null : (
            <p className="text-xs text-destructive">
              {ipActionError}
            </p>
          )}
        </div>

        <div className="space-y-3">
          <div className="flex items-center justify-between">
            <Label>Blocked Countries</Label>
            <span className="text-xs text-muted-foreground">
              {blockedCountryCount}/{MAX_COUNTRIES}
            </span>
          </div>
          <div className="flex flex-wrap gap-2">
            {blockedCountries.length === EMPTY_COUNT ? (
              <span className="text-xs text-muted-foreground">
                No blocked countries yet.
              </span>
            ) : (
              blockedCountries.map((code) => {
                const normalizedCode = code.trim().toUpperCase();
                const displayName =
                  countryNameLookup.get(normalizedCode) ?? code;
                const showCode = displayName.trim().toUpperCase() !== normalizedCode;
                return (
                  <Badge key={code} variant="secondary" className="flex items-center gap-2">
                    <span>{displayName}</span>
                    {showCode ? (
                      <span className="text-xs text-muted-foreground">
                        {code}
                      </span>
                    ) : null}
                    <Button
                      type="button"
                      variant="ghost"
                      size="icon"
                      className="h-5 w-5"
                      disabled={savingBlockedCountries}
                      onClick={() => {
                        void handleRemoveBlockedCountry(code);
                      }}
                      aria-label={`Remove ${code}`}
                    >
                      <X className="h-3 w-3" />
                    </Button>
                </Badge>
                );
              })
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="country-search">Search countries</Label>
            <Input
              id="country-search"
              placeholder="Start typing a country name or code"
              value={countrySearch}
              onChange={(e) => {
                const { currentTarget } = e;
                const { value } = currentTarget;
                setCountrySearch(value);
              }}
              onKeyDown={handleCountrySearchKeyDown}
              disabled={!geoIPReady}
            />
            {geoIPReady &&
            trimmedCountrySearch.length > EMPTY_COUNT &&
            trimmedCountrySearch.length < SEARCH_MIN_LENGTH ? (
              <p className="text-xs text-muted-foreground">
                Type at least 2 characters to search.
              </p>
            ) : null}
          </div>

          <div className="rounded-lg border bg-muted/30 p-2 space-y-1">
            {geoIPCountriesLoading ? (
              <p className="text-xs text-muted-foreground">Searching...</p>
            ) : matchingCountries.length === EMPTY_COUNT ? (
              <p className="text-xs text-muted-foreground">
                {trimmedCountrySearch.length >= SEARCH_MIN_LENGTH
                  ? `No matches for "${trimmedCountrySearch}".`
                  : 'Search results will appear here.'}
              </p>
            ) : (
              matchingCountries.slice(FIRST_INDEX, MAX_COUNTRY_MATCHES).map((country) => {
                const { code: countryCode, name: countryName } = country;
                if (countryCode === EMPTY_STRING) return null;
                return (
                  <button
                    key={countryCode}
                    type="button"
                    className="flex w-full items-center justify-between rounded-md px-2 py-1 text-left text-sm hover:bg-accent"
                    disabled={savingBlockedCountries}
                    onClick={() => {
                      void (async () => {
                        await handleAddBlockedCountry(countryCode);
                        setCountrySearch('');
                      })();
                    }}
                  >
                    <span>{countryName}</span>
                    <span className="text-xs text-muted-foreground">{countryCode}</span>
                  </button>
                );
              })
            )}
          </div>

          {geoIPReady ? null : (
            <p className="text-xs text-muted-foreground">
              GeoIP database is not ready. Download it to manage country blocking.
            </p>
          )}
          {savingBlockedCountries ? (
            <p className="text-xs text-muted-foreground">
              Updating blocked countries...
            </p>
          ) : null}
          {countryActionError === EMPTY_STRING ? null : (
            <p className="text-xs text-destructive">
              {countryActionError}
            </p>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
