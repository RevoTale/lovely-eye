import React from 'react';
import { useMutation } from '@apollo/client';
import { SITE_QUERY, UPDATE_SITE_MUTATION } from '@/graphql';
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
  const [updateBlockedIPs, { loading: savingBlockedIPs }] = useMutation(UPDATE_SITE_MUTATION);
  const [updateBlockedCountries, { loading: savingBlockedCountries }] = useMutation(UPDATE_SITE_MUTATION);

  const handleUpdateBlockedIPs = async (blockedIPs: string[]): Promise<void> => {
    await updateBlockedIPs({
      variables: {
        id: siteId,
        input: {
          name: siteName,
          blockedIPs,
        },
      },
      refetchQueries: [{ query: SITE_QUERY, variables: { id: siteId } }],
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
      refetchQueries: [{ query: SITE_QUERY, variables: { id: siteId } }],
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
