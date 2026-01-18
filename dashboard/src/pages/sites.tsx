import React from 'react';
import { useQuery } from '@apollo/client/react';
import { SitesDocument, type SitesQuery } from '@/gql/graphql';
import { Link } from '@/router';
import { Button, Card, CardContent, CardDescription, CardHeader, CardTitle, Skeleton, Badge } from '@/components/ui';
import { Plus, Globe, ExternalLink, TrendingUp } from 'lucide-react';

type SiteListItem = SitesQuery['sites'][number];

const EMPTY_COUNT = 0;
const FIRST_INDEX = 0;
const EXTRA_DOMAIN_OFFSET = 1;
const SKELETON_CARD_COUNT = 3;

export function SitesPage(): React.JSX.Element {
  const { data, loading, error } = useQuery(SitesDocument);

  if (loading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div className="space-y-2">
            <Skeleton className="h-8 w-32" />
            <Skeleton className="h-4 w-48" />
          </div>
          <Skeleton className="h-10 w-32" />
        </div>
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: SKELETON_CARD_COUNT }, (_, i) => (
            <Card key={i}>
              <CardHeader>
                <Skeleton className="h-6 w-32" />
                <Skeleton className="h-4 w-48" />
              </CardHeader>
              <CardContent>
                <Skeleton className="h-4 w-24" />
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    );
  }

  if (error !== undefined) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-destructive">Error loading sites: {error.message}</div>
      </div>
    );
  }

  const sites: SiteListItem[] = data?.sites ?? [];

  return (
    <div className="space-y-8">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Sites</h1>
          <p className="text-muted-foreground mt-1">Manage your tracked websites</p>
        </div>
        <Link to="/sites/$siteId" params={{ siteId: 'new' }}>
          <Button className="shadow-sm">
            <Plus className="mr-2 h-4 w-4" />
            Add Site
          </Button>
        </Link>
      </div>

      {sites.length === EMPTY_COUNT ? (
        <Card className="border-dashed">
          <CardContent className="flex flex-col items-center justify-center py-16">
            <div className="h-16 w-16 rounded-full bg-primary/10 flex items-center justify-center mb-4">
              <Globe className="h-8 w-8 text-primary" />
            </div>
            <h3 className="text-xl font-semibold mb-2">No sites yet</h3>
            <p className="text-muted-foreground text-center mb-6 max-w-sm">
              Add your first website to start tracking analytics and monitor visitor behavior
            </p>
            <Link to="/sites/$siteId" params={{ siteId: 'new' }}>
              <Button size="lg">
                <Plus className="mr-2 h-4 w-4" />
                Add your first site
              </Button>
            </Link>
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {sites.map((site) => (
            <Link key={site.id} to="/sites/$siteId" params={{ siteId: site.id }}>
              <Card className="group hover:shadow-lg hover:border-primary/50 transition-all cursor-pointer h-full">
                <CardHeader>
                  <div className="flex items-start justify-between">
                    <div className="h-10 w-10 rounded-lg bg-primary/10 flex items-center justify-center group-hover:bg-primary/20 transition-colors">
                      <Globe className="h-5 w-5 text-primary" />
                    </div>
                    <Badge variant="secondary" className="flex items-center gap-1">
                      <TrendingUp className="h-3 w-3" />
                      Active
                    </Badge>
                  </div>
                  <CardTitle className="mt-4 group-hover:text-primary transition-colors">
                    {site.name}
                  </CardTitle>
                  <CardDescription className="flex items-center gap-2">
                    <span>{site.domains[FIRST_INDEX] ?? ''}</span>
                    {site.domains.length > EXTRA_DOMAIN_OFFSET ? (
                      <span className="text-xs text-muted-foreground">
                        +{site.domains.length - EXTRA_DOMAIN_OFFSET} more
                      </span>
                    ) : null}
                    <ExternalLink className="h-3 w-3 opacity-0 group-hover:opacity-100 transition-opacity" />
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <p className="text-sm text-muted-foreground">
                    Added {new Date(site.createdAt).toLocaleDateString('en-US', {
                      month: 'short',
                      day: 'numeric',
                      year: 'numeric'
                    })}
                  </p>
                </CardContent>
              </Card>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
