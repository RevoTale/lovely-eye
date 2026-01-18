import React from 'react';
import { useMutation } from '@apollo/client';
import { CountryTrackingCard } from '@/components/country-tracking-card';
import { REFRESH_GEOIP_MUTATION, UPDATE_SITE_MUTATION } from '@/graphql';

interface GeoIPStatus {
  state: string;
  dbPath?: string | null;
  source?: string | null;
  lastError?: string | null;
}

interface CountryTrackingSectionProps {
  siteId?: string | undefined;
  siteName?: string | undefined;
  initialTrackCountry: boolean;
  geoIPStatus?: GeoIPStatus | null | undefined;
}

export function CountryTrackingSection({
  siteId,
  siteName,
  initialTrackCountry,
  geoIPStatus,
}: CountryTrackingSectionProps): React.JSX.Element {
  const [trackCountry, setTrackCountry] = React.useState(initialTrackCountry);
  const [actionError, setActionError] = React.useState('');
  const [updateSite, { loading: updating }] = useMutation(UPDATE_SITE_MUTATION);
  const [refreshGeoIP, { loading: refreshing }] = useMutation(REFRESH_GEOIP_MUTATION);

  React.useEffect(() => {
    setTrackCountry(initialTrackCountry);
  }, [initialTrackCountry, siteId]);

  const handleToggle = async (enabled: boolean): Promise<void> => {
    if (!siteId || !siteName) return;
    const previous = trackCountry;
    setTrackCountry(enabled);
    setActionError('');
    try {
      await updateSite({
        variables: {
          id: siteId,
          input: {
            name: siteName,
            trackCountry: enabled,
          },
        },
      });
    } catch {
      setTrackCountry(previous);
      setActionError('Failed to update country tracking.');
    }
  };

  const handleRetry = async (): Promise<void> => {
    if (!siteId) return;
    setActionError('');
    try {
      await refreshGeoIP();
    } catch {
      setActionError('Failed to refresh GeoIP database.');
    }
  };

  const resolvedGeoIPStatus = geoIPStatus ?? null;

  return (
    <div className="space-y-2">
      <CountryTrackingCard
        trackCountry={trackCountry}
        updating={updating}
        refreshing={refreshing}
        geoIPStatus={resolvedGeoIPStatus}
        onToggle={(enabled) => {
          void handleToggle(enabled);
        }}
        onRetry={() => {
          void handleRetry();
        }}
      />
      {actionError ? (
        <p className="text-xs text-destructive">
          {actionError}
        </p>
      ) : null}
    </div>
  );
}
