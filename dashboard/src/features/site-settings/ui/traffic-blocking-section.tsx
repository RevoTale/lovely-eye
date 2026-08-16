import { useMutation } from '@apollo/client/react';
import { TrafficBlockingCard } from '@/features/site-settings/ui/traffic-blocking-card';
import { UpdateSiteDocument } from '@/shared/api/generated/graphql';

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
  const [updateBlockedCountries, { loading: savingBlockedCountries }] =
    useMutation(UpdateSiteDocument);

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
};
