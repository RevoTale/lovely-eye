import React from 'react';
import { useParams, useNavigate } from '@tanstack/react-router';
import { useQuery, useMutation } from '@apollo/client/react';
import {
  SiteDocument,
  SitesDocument,
  CreateSiteDocument,
  UpdateSiteDocument,
  GeoIpStatusDocument,
  type SiteQuery,
} from '@/gql/graphql';
import { Button, Card, CardContent, CardHeader, Skeleton } from '@/components/ui';
import { ArrowLeft } from 'lucide-react';
import { Link, siteDetailRoute } from '@/router';
import { CountryTrackingSection } from '@/components/site-form/country-tracking-section';
import { DangerZoneSection } from '@/components/site-form/danger-zone-section';
import { EventDefinitionsSection } from '@/components/site-form/event-definitions-section';
import { SiteInfoCard } from '@/components/site-form/site-info-card';
import { TrackingCodeSection } from '@/components/site-form/tracking-code-section';
import { TrafficBlockingSection } from '@/components/site-form/traffic-blocking-section';

type SiteDetails = NonNullable<SiteQuery['site']>;

export function SiteFormPage(): React.JSX.Element {
  const { siteId } = useParams({ from: siteDetailRoute.id });
  const navigate = useNavigate();
  const isNew = siteId === 'new';
  const GEO_IP_POLL_INTERVAL_MS = 5000;
  const sitesPaging = { limit: 100, offset: 0 };

  const { data: siteData, loading: siteLoading } = useQuery(SiteDocument, {
    variables: { id: siteId },
    skip: isNew,
  });

  const { data: geoIPData } = useQuery(GeoIpStatusDocument, {
    skip: isNew,
    pollInterval: GEO_IP_POLL_INTERVAL_MS,
  });

  const [createSite, { loading: creating }] = useMutation(CreateSiteDocument, {
    refetchQueries: [{ query: SitesDocument, variables: { paging: sitesPaging } }],
    onCompleted: (data) => {
      void navigate({ to: '/sites/$siteId', params: { siteId: data.createSite.id } });
    },
  });

  const [updateSite, { loading: updating }] = useMutation(UpdateSiteDocument);

  const site: SiteDetails | undefined = siteData?.site ?? undefined;
  const geoIPStatus = geoIPData?.geoIPStatus;
  const handleSubmit = async (nameValue: string, domainsValue: string[]): Promise<void> => {
    await createSite({
      variables: {
        input: {
          name: nameValue,
          domains: domainsValue,
        },
      },
    });
  };

  const handleDomainsSave = async (nameValue: string, domainsValue: string[]): Promise<void> => {
    if (site === undefined) return;
    await updateSite({
      variables: {
        id: site.id,
        input: {
          name: nameValue,
          domains: domainsValue,
        },
      },
    });
  };

  if (siteLoading) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-8 w-64" />
        <Card>
          <CardHeader>
            <Skeleton className="h-6 w-32" />
          </CardHeader>
          <CardContent className="space-y-4">
            <Skeleton className="h-10 w-full" />
            <Skeleton className="h-10 w-full" />
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="space-y-6 max-w-3xl">
      <div className="flex items-center gap-4">
        <Button variant="outline" size="sm" asChild>
          <Link to="/">
            <ArrowLeft className="h-4 w-4 mr-2" />
            Back to Sites
          </Link>
        </Button>
      </div>

      <div>
        <h1 className="text-3xl font-bold tracking-tight">
          {isNew ? 'Add New Site' : site?.name ?? 'Site Details'}
        </h1>
        <p className="text-muted-foreground mt-1">
          {isNew ? 'Create a new site to track analytics' : 'View and manage site settings'}
        </p>
      </div>

      <SiteInfoCard
        siteId={site?.id}
        isNew={isNew}
        initialName={site?.name ?? ''}
        initialDomains={site?.domains ?? []}
        creating={creating}
        updating={updating}
        onCreate={async (newName, newDomains) => {
          await handleSubmit(newName, newDomains);
        }}
        onSaveDomains={async (newName, newDomains) => {
          await handleDomainsSave(newName, newDomains);
        }}
        onCancel={() => {
          void navigate({ to: '/' });
        }}
      />

      {!isNew && site !== undefined ? (
        <>
          <TrackingCodeSection
            siteId={site.id}
            publicKey={site.publicKey}
            onViewAnalytics={() => {
              void navigate({ to: '/sites/$siteId', params: { siteId: site.id }, search: { view: 'analytics' } });
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
            geoIPReady={geoIPStatus?.state === 'ready'}
          />

          <EventDefinitionsSection siteId={site.id} />

          <DangerZoneSection
            siteId={site.id}
            onDeleted={() => {
              void navigate({ to: '/' });
            }}
          />
        </>
      ) : null}
    </div>
  );
}
