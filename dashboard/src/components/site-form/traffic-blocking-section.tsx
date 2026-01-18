import React from 'react';
import { useMutation } from '@apollo/client/react';
import { UpdateSiteDocument, SiteDocument } from '@/gql/graphql';
import { TrafficBlockingCard } from '@/components/site-form/traffic-blocking-card';

interface TrafficBlockingSectionProps {
  siteId: string;
  siteName: string;
  initialBlockedIPs: string[];
  initialBlockedCountries: string[];
  geoIPReady: boolean;
}

export function TrafficBlockingSection({
  siteId,
  siteName,
  initialBlockedIPs,
  initialBlockedCountries,
  geoIPReady,
}: TrafficBlockingSectionProps): React.JSX.Element {
  const [updateBlockedIPs, { loading: savingBlockedIPs }] = useMutation(UpdateSiteDocument);
  const [updateBlockedCountries, { loading: savingBlockedCountries }] = useMutation(UpdateSiteDocument);

  const handleUpdateBlockedIPs = async (blockedIPs: string[]): Promise<void> => {
    await updateBlockedIPs({
      variables: {
        id: siteId,
        input: {
          name: siteName,
          blockedIPs,
        },
      },
      refetchQueries: [{ query: SiteDocument, variables: { id: siteId } }],
      awaitRefetchQueries: true,
    });
  };

  const handleUpdateBlockedCountries = async (blockedCountries: string[]): Promise<void> => {
    await updateBlockedCountries({
      variables: {
        id: siteId,
        input: {
          name: siteName,
          blockedCountries,
        },
      },
      refetchQueries: [{ query: SiteDocument, variables: { id: siteId } }],
      awaitRefetchQueries: true,
    });
  };

  return (
    <TrafficBlockingCard
      siteId={siteId}
      initialBlockedIPs={initialBlockedIPs}
      initialBlockedCountries={initialBlockedCountries}
      savingBlockedIPs={savingBlockedIPs}
      savingBlockedCountries={savingBlockedCountries}
      geoIPReady={geoIPReady}
      onUpdateBlockedIPs={handleUpdateBlockedIPs}
      onUpdateBlockedCountries={handleUpdateBlockedCountries}
    />
  );
}
