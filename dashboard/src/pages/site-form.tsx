import React, { useState } from 'react';
import { useParams, useNavigate } from '@tanstack/react-router';
import { useMutation, useQuery } from '@apollo/client';
import { CREATE_SITE_MUTATION, SITE_QUERY, SITES_QUERY, DELETE_SITE_MUTATION, REGENERATE_SITE_KEY_MUTATION } from '@/graphql';
import { Button, Card, CardContent, CardDescription, CardHeader, CardTitle, Input, Label, Skeleton } from '@/components/ui';
import { Globe, ArrowLeft, Save, Trash2, RefreshCw, Copy, CheckCircle2, AlertCircle } from 'lucide-react';
import { siteDetailRoute } from '@/router';
import type { Site } from '@/generated/graphql';

export function SiteFormPage(): React.JSX.Element {
  const { siteId } = useParams({ from: siteDetailRoute.id });
  const navigate = useNavigate();
  const isNew = siteId === 'new';

  // Form state
  const [name, setName] = useState('');
  const [domain, setDomain] = useState('');
  const [copied, setCopied] = useState(false);
  const [error, setError] = useState('');

  // Queries and mutations
  const { data: siteData, loading: siteLoading } = useQuery(SITE_QUERY, {
    variables: { id: siteId },
    skip: isNew,
    onCompleted: (data) => {
      if (data?.site) {
        const site = data.site as Site;
        setName(site.name);
        setDomain(site.domain);
      }
    },
  });

  const [createSite, { loading: creating }] = useMutation(CREATE_SITE_MUTATION, {
    refetchQueries: [{ query: SITES_QUERY }],
    onCompleted: (data) => {
      if (data?.createSite) {
        navigate({ to: '/sites/$siteId', params: { siteId: data.createSite.id } });
      }
    },
    onError: (err) => {
      setError(err.message);
    },
  });

  const [deleteSite, { loading: deleting }] = useMutation(DELETE_SITE_MUTATION, {
    refetchQueries: [{ query: SITES_QUERY }],
    onCompleted: () => {
      navigate({ to: '/' });
    },
    onError: (err) => {
      setError(err.message);
    },
  });

  const [regenerateKey, { loading: regenerating }] = useMutation(REGENERATE_SITE_KEY_MUTATION, {
    onError: (err) => {
      setError(err.message);
    },
  });

  const site = siteData?.site as Site | undefined;

  const handleSubmit = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();
    setError('');

    const trimmedName = name.trim();
    const trimmedDomain = domain.trim();

    // Validate required fields
    if (!trimmedName || !trimmedDomain) {
      setError('Name and domain are required');
      return;
    }

    // Validate site name (1-100 characters, alphanumeric and common punctuation)
    if (trimmedName.length < 1 || trimmedName.length > 100) {
      setError('Site name must be between 1 and 100 characters');
      return;
    }

    // Validate domain format (basic domain validation)
    const domainRegex = /^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?(\.[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?)*$/;
    if (!domainRegex.test(trimmedDomain)) {
      setError('Please enter a valid domain (e.g., example.com)');
      return;
    }

    if (isNew) {
      await createSite({
        variables: {
          input: {
            name: trimmedName,
            domain: trimmedDomain,
          },
        },
      });
    }
  };

  const handleDelete = async (): Promise<void> => {
    if (!site) return;

    if (!window.confirm(`Are you sure you want to delete "${site.name}"? This action cannot be undone.`)) {
      return;
    }

    await deleteSite({
      variables: { id: site.id },
    });
  };

  const handleRegenerateKey = async (): Promise<void> => {
    if (!site) return;

    if (!window.confirm('Are you sure you want to regenerate the site key? The old key will stop working.')) {
      return;
    }

    await regenerateKey({
      variables: { id: site.id },
    });
  };

  const handleCopyKey = async (): Promise<void> => {
    if (!site) return;

    await navigator.clipboard.writeText(site.publicKey);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const generateTrackingScript = (): string => {
    if (!site) return '';

    // Use the base path from environment config
    const basePath = window.__ENV__?.BASE_PATH || '';
    const trackerUrl = `${window.location.origin}${basePath}/tracker.js`;

    return `<script>
  (function() {
    var script = document.createElement('script');
    script.src = '${trackerUrl}';
    script.setAttribute('data-site-key', '${site.publicKey}');
    script.async = true;
    document.head.appendChild(script);
  })();
</script>`;
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
        <Button
          variant="outline"
          size="sm"
          onClick={() => navigate({ to: '/' })}
        >
          <ArrowLeft className="h-4 w-4 mr-2" />
          Back to Sites
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

      {error ? (
        <div className="flex items-center gap-2 p-4 bg-destructive/10 border border-destructive/20 rounded-lg text-destructive">
          <AlertCircle className="h-5 w-5" />
          <p className="text-sm">{error}</p>
        </div>
      ) : null}

      <form onSubmit={handleSubmit}>
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <div className="h-8 w-8 rounded-lg bg-primary/10 flex items-center justify-center">
                <Globe className="h-4 w-4 text-primary" />
              </div>
              Site Information
            </CardTitle>
            <CardDescription>
              {isNew ? 'Enter your website details' : 'Site configuration and tracking details'}
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="space-y-2">
              <Label htmlFor="name">Site Name</Label>
              <Input
                id="name"
                placeholder="My Awesome Website"
                value={name}
                onChange={(e) => setName(e.target.value)}
                disabled={!isNew}
                required
              />
              <p className="text-xs text-muted-foreground">
                A friendly name to identify your site
              </p>
            </div>

            <div className="space-y-2">
              <Label htmlFor="domain">Domain</Label>
              <Input
                id="domain"
                placeholder="example.com"
                value={domain}
                onChange={(e) => {
                  const value = e.target.value;
                  // Automatically truncate domain to domain.com format
                  const truncated = value
                    .replace(/^https?:\/\//, '') // Remove http:// or https://
                    .replace(/^www\./, '') // Remove www.
                    .replace(/\/.*$/, '') // Remove path and trailing slashes
                    .toLowerCase() // Convert to lowercase
                    .trim();
                  setDomain(truncated);
                }}
                disabled={!isNew}
                required
              />
              <p className="text-xs text-muted-foreground">
                Your website domain (without https://)
              </p>
            </div>

            {isNew ? (
              <div className="flex gap-3 pt-4">
                <Button type="submit" disabled={creating}>
                  <Save className="h-4 w-4 mr-2" />
                  {creating ? 'Creating...' : 'Create Site'}
                </Button>
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => navigate({ to: '/' })}
                  disabled={creating}
                >
                  Cancel
                </Button>
              </div>
            ) : null}
          </CardContent>
        </Card>
      </form>

      {!isNew && site ? (
        <>
          <Card>
            <CardHeader>
              <CardTitle>Tracking Code</CardTitle>
              <CardDescription>
                Add this script to your website's HTML to start tracking
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <Label>Site Key</Label>
                <div className="flex gap-2">
                  <Input
                    value={site.publicKey}
                    readOnly
                    className="font-mono text-sm"
                  />
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    onClick={handleCopyKey}
                  >
                    {copied ? (
                      <CheckCircle2 className="h-4 w-4 text-green-500" />
                    ) : (
                      <Copy className="h-4 w-4" />
                    )}
                  </Button>
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    onClick={handleRegenerateKey}
                    disabled={regenerating}
                  >
                    <RefreshCw className={`h-4 w-4 ${regenerating ? 'animate-spin' : ''}`} />
                  </Button>
                </div>
                <p className="text-xs text-muted-foreground">
                  Keep this key secure. Use regenerate if compromised.
                </p>
              </div>

              <div className="space-y-2">
                <Label>Tracking Script</Label>
                <div className="relative">
                  <pre className="p-4 bg-muted rounded-lg overflow-x-auto text-xs border">
                    <code>{generateTrackingScript()}</code>
                  </pre>
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    className="absolute top-2 right-2"
                    onClick={async () => {
                      await navigator.clipboard.writeText(generateTrackingScript());
                      setCopied(true);
                      setTimeout(() => setCopied(false), 2000);
                    }}
                  >
                    {copied ? (
                      <CheckCircle2 className="h-4 w-4 text-green-500" />
                    ) : (
                      <Copy className="h-4 w-4" />
                    )}
                  </Button>
                </div>
                <p className="text-xs text-muted-foreground">
                  Add this script to the &lt;head&gt; section of your website
                </p>
              </div>

              <div className="pt-4">
                <Button
                  variant="outline"
                  onClick={() => navigate({ to: '/sites/$siteId', params: { siteId: site.id }, search: { view: 'analytics' } })}
                >
                  View Analytics
                </Button>
              </div>
            </CardContent>
          </Card>

          <Card className="border-destructive/50">
            <CardHeader>
              <CardTitle className="text-destructive">Danger Zone</CardTitle>
              <CardDescription>
                Irreversible actions that affect this site
              </CardDescription>
            </CardHeader>
            <CardContent>
              <Button
                variant="destructive"
                onClick={handleDelete}
                disabled={deleting}
              >
                <Trash2 className="h-4 w-4 mr-2" />
                {deleting ? 'Deleting...' : 'Delete Site'}
              </Button>
              <p className="text-xs text-muted-foreground mt-2">
                This will permanently delete all analytics data for this site
              </p>
            </CardContent>
          </Card>
        </>
      ) : null}
    </div>
  );
}
