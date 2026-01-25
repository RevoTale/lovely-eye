
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

export const TrafficBlockingSection = ({
  siteId,
  siteName,
  initialBlockedIPs,
  initialBlockedCountries,
  geoIPReady,
}: TrafficBlockingSectionProps): React.ReactNode => {
  const [updateBlockedIPs, { loading: savingBlockedIPs }] = useMutation(UpdateSiteDocument);
  const [updateBlockedCountries, { loading: savingBlockedCountries }] = useMutation(UpdateSiteDocument);

  const handleUpdateBlockedIPs = async (blockedIPs: string[]): Promise<void> => {
    await updateBlockedIPs({
      variables: {
        id: siteId,
        input: {
          name: siteName,
          blockedIPs,
          trackCountry: null,
          domains: null,
          blockedCountries: null,
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
          trackCountry: null,
          domains: null,
          blockedIPs: null,
        },
      },
      refetchQueries: [{ query: SiteDocument, variables: { id: siteId } }],
      awaitRefetchQueries: true,
    });
  };

  return (
    <TrafficBlockingCard
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
