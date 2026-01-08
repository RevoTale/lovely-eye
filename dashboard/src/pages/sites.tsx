import React from 'react';
import { useQuery } from '@apollo/client';
import { SITES_QUERY } from '@/graphql';
import { Link } from '@/router';
import { Button, Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui';
import { Plus, Globe, ExternalLink } from 'lucide-react';
import type { Site } from '@/generated/graphql';

export function SitesPage(): React.JSX.Element {
  const { data, loading, error } = useQuery(SITES_QUERY);

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-muted-foreground">Loading sites...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        <div className="text-destructive">Error loading sites: {error.message}</div>
      </div>
    );
  }

  const sites = (data?.sites ?? []) as Site[];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Sites</h1>
          <p className="text-muted-foreground">Manage your tracked websites</p>
        </div>
        <Link to="/sites/$siteId" params={{ siteId: 'new' }}>
          <Button>
            <Plus className="mr-2 h-4 w-4" />
            Add Site
          </Button>
        </Link>
      </div>

      {sites.length === 0 ? (
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-12">
            <Globe className="h-12 w-12 text-muted-foreground mb-4" />
            <h3 className="text-lg font-semibold mb-2">No sites yet</h3>
            <p className="text-muted-foreground text-center mb-4">
              Add your first website to start tracking analytics
            </p>
            <Link to="/sites/$siteId" params={{ siteId: 'new' }}>
              <Button>
                <Plus className="mr-2 h-4 w-4" />
                Add your first site
              </Button>
            </Link>
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {sites.map((site) => (
            <Link key={site.id} to="/sites/$siteId" params={{ siteId: site.id }}>
              <Card className="hover:border-primary/50 transition-colors cursor-pointer">
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Globe className="h-5 w-5" />
                    {site.name}
                  </CardTitle>
                  <CardDescription className="flex items-center gap-1">
                    {site.domain}
                    <ExternalLink className="h-3 w-3" />
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <p className="text-sm text-muted-foreground">
                    Added {new Date(site.createdAt).toLocaleDateString()}
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
