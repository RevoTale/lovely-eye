import React from 'react';
import { useQuery } from '@apollo/client/react';
import { GeoIpCountriesDocument } from '@/gql/graphql';
import { Badge, Button, Card, CardContent, CardDescription, CardHeader, CardTitle, Input, Label } from '@/components/ui';
import { Loader2, Shield, X } from 'lucide-react';
import { getNormalizedBlockedIPs, normalizeCountryCodesPreserveOrder, normalizeIPInput } from '@/components/site-form/utils';

interface BlockedIPEntry {
  id: string;
  value: string;
}

interface TrafficBlockingCardProps {
  siteId?: string;
  initialBlockedIPs: string[];
  initialBlockedCountries: string[];
  savingBlockedIPs: boolean;
  savingBlockedCountries: boolean;
  geoIPReady: boolean;
  onUpdateBlockedIPs: (blockedIPs: string[]) => Promise<void>;
  onUpdateBlockedCountries: (blockedCountries: string[]) => Promise<void>;
}

export function TrafficBlockingCard({
  siteId,
  initialBlockedIPs,
  initialBlockedCountries,
  savingBlockedIPs,
  savingBlockedCountries,
  geoIPReady,
  onUpdateBlockedIPs,
  onUpdateBlockedCountries,
}: TrafficBlockingCardProps): React.JSX.Element {
  const maxIPs = 500;
  const maxCountries = 250;
  const [blockedIPs, setBlockedIPs] = React.useState<BlockedIPEntry[]>(() => {
    return initialBlockedIPs.map((ip, index) => ({ id: String(index + 1), value: ip }));
  });
  const [blockedCountries, setBlockedCountries] = React.useState<string[]>(initialBlockedCountries);
  const [blockedCountryNames, setBlockedCountryNames] = React.useState<Record<string, string>>({});
  const [countrySearch, setCountrySearch] = React.useState('');
  const [ipActionError, setIpActionError] = React.useState('');
  const [countryActionError, setCountryActionError] = React.useState('');
  const [newIPValue, setNewIPValue] = React.useState('');
  const [newIPError, setNewIPError] = React.useState('');

  React.useEffect(() => {
    setBlockedIPs(initialBlockedIPs.map((ip, index) => ({ id: String(index + 1), value: ip })));
    setBlockedCountries(initialBlockedCountries);
    setCountrySearch('');
    setIpActionError('');
    setCountryActionError('');
    setNewIPValue('');
    setNewIPError('');
  }, [initialBlockedCountries, initialBlockedIPs, siteId]);

  const trimmedCountrySearch = countrySearch.trim();
  const shouldSearchCountries = geoIPReady && trimmedCountrySearch.length >= 2;
  const { data: geoIPCountriesData, loading: geoIPCountriesLoading } = useQuery(GeoIpCountriesDocument, {
    variables: { search: trimmedCountrySearch },
    skip: !shouldSearchCountries,
  });

  const geoIPCountries = React.useMemo(() => geoIPCountriesData?.geoIPCountries ?? [], [geoIPCountriesData]);
  const blockedIPCount = React.useMemo(() => getNormalizedBlockedIPs(blockedIPs.map((entry) => entry.value)).length, [blockedIPs]);
  const normalizedBlockedCountries = React.useMemo(() => normalizeCountryCodesPreserveOrder(blockedCountries), [blockedCountries]);
  const blockedCountryCount = normalizedBlockedCountries.length;
  const countryNameLookup = React.useMemo(() => {
    return new Map(geoIPCountries.map((country) => [country.code, country.name] as const));
  }, [geoIPCountries]);
  const isUpdating = savingBlockedIPs || savingBlockedCountries;

  const availableCountries = React.useMemo(() => {
    const selected = new Set(normalizedBlockedCountries);
    return geoIPCountries.filter((country) => !selected.has(country.code));
  }, [geoIPCountries, normalizedBlockedCountries]);
  const matchingCountries = React.useMemo(() => {
    if (!shouldSearchCountries) {
      return [];
    }
    return availableCountries;
  }, [availableCountries, shouldSearchCountries]);

  React.useEffect(() => {
    if (geoIPCountries.length === 0) {
      return;
    }
    setBlockedCountryNames((prev) => {
      const next = { ...prev };
      for (const country of geoIPCountries) {
        const code = country.code;
        const name = country.name;
        if (code && name && next[code] === undefined) {
          next[code] = name;
        }
      }
      return next;
    });
  }, [geoIPCountries]);

  const buildNextBlockedCountries = (code: string): string[] => {
    const trimmed = code.trim();
    if (!trimmed) return blockedCountries;
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
    if (next.length > maxCountries) {
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
    if (!value) return false;
    if (value.includes(':')) {
      return /^([0-9a-f]{0,4}:){2,7}[0-9a-f]{0,4}$/i.test(value);
    }
    const parts = value.split('.');
    if (parts.length !== 4) return false;
    return parts.every((part) => {
      if (!/^\d{1,3}$/.test(part)) return false;
      const num = Number(part);
      return num >= 0 && num <= 255;
    });
  };

  const handleAddIP = async (): Promise<void> => {
    if (savingBlockedIPs) return;
    const trimmed = normalizeIPInput(newIPValue);
    if (!trimmed) {
      setNewIPError('Enter a valid IP before saving.');
      return;
    }
    if (!isValidIP(trimmed)) {
      setNewIPError('Enter a valid IP address.');
      return;
    }
    const existing = getNormalizedBlockedIPs(blockedIPs.map((item) => item.value));
    if (existing.includes(trimmed)) {
      setNewIPError('That IP is already blocked.');
      return;
    }
    const next = [...existing, trimmed];
    if (next.length > maxIPs) {
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
    if (!normalized) return;
    const next = blockedIPs.filter((item) => item.value.trim() !== normalized);
    if (next.length === blockedIPs.length) {
      return;
    }
    setIpActionError('');
    try {
      await onUpdateBlockedIPs(getNormalizedBlockedIPs(next.map((item) => item.value)));
    } catch (err) {
      setIpActionError(err instanceof Error ? err.message : 'Failed to remove blocked IP');
    }
  };

  const handleCountrySearchKeyDown = (event: React.KeyboardEvent<HTMLInputElement>): void => {
    if (event.key !== 'Enter') return;
    if (!shouldSearchCountries) return;
    const query = trimmedCountrySearch.toLowerCase();
    const exactMatch = matchingCountries.find(
      (country) => country.code.toLowerCase() === query || country.name.toLowerCase() === query,
    );
    const singleMatch = matchingCountries.length === 1 ? matchingCountries[0] : undefined;
    const target = exactMatch ?? singleMatch;
    if (target?.code) {
      const code = target.code;
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
              {blockedIPCount}/{maxIPs}
            </span>
          </div>
          <div className="space-y-2">
            {blockedIPs.length === 0 ? (
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
                    setNewIPValue(normalizeIPInput(e.target.value));
                    if (newIPError) {
                      setNewIPError('');
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
                  disabled={savingBlockedIPs || blockedIPCount >= maxIPs}
                >
                  {savingBlockedIPs ? (
                    <Loader2 className="h-4 w-4 animate-spin" />
                  ) : (
                    'Add IP'
                  )}
                </Button>
              </div>
              {newIPError ? (
                <span className="text-xs text-destructive">{newIPError}</span>
              ) : null}
            </div>
          </div>
          {savingBlockedIPs ? (
            <p className="text-xs text-muted-foreground">
              Updating blocked IPs...
            </p>
          ) : null}
          {ipActionError ? (
            <p className="text-xs text-destructive">
              {ipActionError}
            </p>
          ) : null}
        </div>

        <div className="space-y-3">
          <div className="flex items-center justify-between">
            <Label>Blocked Countries</Label>
            <span className="text-xs text-muted-foreground">
              {blockedCountryCount}/{maxCountries}
            </span>
          </div>
          <div className="flex flex-wrap gap-2">
            {blockedCountries.length === 0 ? (
              <span className="text-xs text-muted-foreground">
                No blocked countries yet.
              </span>
            ) : (
              blockedCountries.map((code) => {
                const normalizedCode = code.trim().toUpperCase();
                const displayName = countryNameLookup.get(normalizedCode)
                  ?? blockedCountryNames[normalizedCode]
                  ?? code;
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
                setCountrySearch(e.target.value);
              }}
              onKeyDown={handleCountrySearchKeyDown}
              disabled={!geoIPReady}
            />
            {geoIPReady && trimmedCountrySearch.length > 0 && trimmedCountrySearch.length < 2 ? (
              <p className="text-xs text-muted-foreground">
                Type at least 2 characters to search.
              </p>
            ) : null}
          </div>

          <div className="rounded-lg border bg-muted/30 p-2 space-y-1">
            {geoIPCountriesLoading ? (
              <p className="text-xs text-muted-foreground">Searching...</p>
            ) : matchingCountries.length === 0 ? (
              <p className="text-xs text-muted-foreground">
                {trimmedCountrySearch.length >= 2 ? `No matches for "${trimmedCountrySearch}".` : 'Search results will appear here.'}
              </p>
            ) : (
              matchingCountries.slice(0, 8).map((country) => {
                const countryCode = country.code;
                const countryName = country.name;
                if (!countryCode) return null;
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

          {!geoIPReady ? (
            <p className="text-xs text-muted-foreground">
              GeoIP database is not ready. Download it to manage country blocking.
            </p>
          ) : null}
          {savingBlockedCountries ? (
            <p className="text-xs text-muted-foreground">
              Updating blocked countries...
            </p>
          ) : null}
          {countryActionError ? (
            <p className="text-xs text-destructive">
              {countryActionError}
            </p>
          ) : null}
        </div>
      </CardContent>
    </Card>
  );
}
