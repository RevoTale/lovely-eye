import { useMutation, useQuery } from '@apollo/client/react';
import { Link, useNavigate, useParams } from '@tanstack/react-router';
import { ArrowLeft } from 'lucide-react';
import { EventDefinitionsSection } from '@/features/event-definitions/event-definitions-section';
import { CountryTrackingSection } from '@/features/site-settings/ui/country-tracking-section';
import { DangerZoneSection } from '@/features/site-settings/ui/danger-zone-section';
import { SiteInfoCard } from '@/features/site-settings/ui/site-info-card';
import { TrackingCodeSection } from '@/features/site-settings/ui/tracking-code-section';
import { TrafficBlockingSection } from '@/features/site-settings/ui/traffic-blocking-section';
import {
  GeoIpStatusDocument,
  GeoIpStatusFieldsFragmentDoc,
  type SiteDetailsFieldsFragment,
  SiteDetailsFieldsFragmentDoc,
  SiteDocument,
  UpdateSiteDocument,
  type UpdateSiteMutation,
  type UpdateSiteMutationVariables,
} from '@/shared/api/generated/graphql';
import { readFragment } from '@/shared/api/read-fragment';
import { buttonVariants } from '@/shared/ui/button';
import { Card, CardContent, CardHeader } from '@/shared/ui/card';
import { Skeleton } from '@/shared/ui/skeleton';

type SiteDetails = SiteDetailsFieldsFragment;

const SITE_SETTINGS_ROUTE_ID = '/_auth/sites/$siteId/settings' as const;

const SiteFormLoading = (): React.ReactNode => (
  <div className='mx-auto w-full max-w-3xl space-y-6'>
    <Skeleton className='h-8 w-64' />
    <Card>
      <CardHeader>
        <Skeleton className='h-6 w-32' />
      </CardHeader>
      <CardContent className='space-y-4'>
        <Skeleton className='h-10 w-full' />
        <Skeleton className='h-10 w-full' />
      </CardContent>
    </Card>
  </div>
);

export const SiteSettingsScreen = (): React.ReactNode => {
  const { siteId } = useParams({ from: SITE_SETTINGS_ROUTE_ID });
  const navigate = useNavigate();
  const {
    data: siteData,
    error: siteError,
    loading: siteLoading,
  } = useQuery(SiteDocument, {
    variables: { id: siteId },
  });

  const { data: geoIPData } = useQuery(GeoIpStatusDocument);

  const [updateSite, { loading: updating }] = useMutation<
    UpdateSiteMutation,
    UpdateSiteMutationVariables
  >(UpdateSiteDocument);

  const siteDataValue = siteData?.site;
  const geoIPStatusValue = geoIPData?.geoIPStatus;
  const site: SiteDetails | undefined =
    siteDataValue === null || siteDataValue === undefined
      ? undefined
      : readFragment(SiteDetailsFieldsFragmentDoc, siteDataValue);
  const geoIPStatus =
    geoIPStatusValue === undefined
      ? undefined
      : readFragment(GeoIpStatusFieldsFragmentDoc, geoIPStatusValue);
  const handleDomainsSave = async (nameValue: string, domainsValue: string[]): Promise<void> => {
    if (site === undefined) return;
    await updateSite({
      variables: {
        id: site.id,
        input: {
          name: nameValue,
          domains: domainsValue,
          trackCountry: null,
          blockedIPs: null,
          blockedCountries: null,
        },
      },
    });
  };

  if (siteLoading) {
    return <SiteFormLoading />;
  }
  if (siteError !== undefined) {
    return <p className='text-destructive'>Error loading site settings: {siteError.message}</p>;
  }
  if (site === undefined) {
    return <p className='text-muted-foreground'>Site not found.</p>;
  }

  return (
    <div className='mx-auto w-full max-w-3xl space-y-6'>
      <div className='flex items-center gap-4'>
        <Link to='/sites' className={buttonVariants({ variant: 'outline', size: 'sm' })}>
          <ArrowLeft className='h-4 w-4 mr-2' />
          Back to Sites
        </Link>
      </div>

      <div>
        <h1 className='text-3xl font-bold tracking-tight'>{site.name}</h1>
        <p className='text-muted-foreground mt-1'>View and manage site settings</p>
      </div>

      <SiteInfoCard
        initialName={site.name}
        initialDomains={site.domains}
        updating={updating}
        onSaveDomains={async (newName, newDomains) => {
          await handleDomainsSave(newName, newDomains);
        }}
      />

      <TrackingCodeSection
        siteId={site.id}
        publicKey={site.publicKey}
        onViewAnalytics={() => {
          void navigate({
            to: '/sites/$siteId/analytics',
            params: { siteId: site.id },
          });
        }}
      />

      <CountryTrackingSection
        siteId={site.id}
        siteName={site.name}
        initialTrackCountry={site.trackCountry}
        geoIPStatus={geoIPStatus}
      />

      <TrafficBlockingSection
        siteId={site.id}
        siteName={site.name}
        initialBlockedIPs={site.blockedIPs}
        initialBlockedCountries={site.blockedCountries}
        geoIPReady={geoIPStatus?.state === 'READY'}
      />

      <EventDefinitionsSection siteId={site.id} />

      <DangerZoneSection
        siteId={site.id}
        onDeleted={() => {
          void navigate({ to: '/' });
        }}
      />
    </div>
  );
};
