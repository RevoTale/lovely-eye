import { useMutation } from '@apollo/client/react';
import { type ReactElement, useState } from 'react';
import { CountryTrackingCard } from '@/features/site-settings/ui/country-tracking-card';
import {
  GeoIpStatusDocument,
  type GeoIpStatusFieldsFragment,
  RefreshGeoIpDatabaseDocument,
  UpdateSiteDocument,
} from '@/shared/api/generated/graphql';

interface CountryTrackingSectionProps {
  siteId?: string | undefined;
  siteName?: string | undefined;
  initialTrackCountry: boolean;
  geoIPStatus?: GeoIpStatusFieldsFragment | null | undefined;
}

export const CountryTrackingSection = ({
  siteId,
  siteName,
  initialTrackCountry,
  geoIPStatus,
}: CountryTrackingSectionProps): ReactElement => {
  const [trackCountry, setTrackCountry] = useState(initialTrackCountry);
  const [actionError, setActionError] = useState('');
  const refreshGeoIPStatus = {
    awaitRefetchQueries: true,
    refetchQueries: [GeoIpStatusDocument],
  };
  const [updateSite, { loading: updating }] = useMutation(UpdateSiteDocument, refreshGeoIPStatus);
  const [refreshGeoIP, { loading: refreshing }] = useMutation(
    RefreshGeoIpDatabaseDocument,
    refreshGeoIPStatus
  );

  const handleToggle = async (enabled: boolean): Promise<void> => {
    if (siteId === undefined || siteId === '' || siteName === undefined || siteName === '') {
      return;
    }
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
            domains: null,
            blockedIPs: null,
            blockedCountries: null,
          },
        },
      });
    } catch {
      setTrackCountry(previous);
      setActionError('Failed to update country tracking.');
    }
  };

  const handleRetry = async (): Promise<void> => {
    if (siteId === undefined || siteId === '') return;
    setActionError('');
    try {
      await refreshGeoIP();
    } catch {
      setActionError('Failed to refresh GeoIP database.');
    }
  };

  const resolvedGeoIPStatus = geoIPStatus ?? null;

  return (
    <div className='space-y-2'>
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
      {actionError === '' ? null : <p className='text-xs text-destructive'>{actionError}</p>}
    </div>
  );
};
