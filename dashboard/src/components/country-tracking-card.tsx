import React from 'react';
import { Badge, Button, Card, CardContent, CardDescription, CardHeader, CardTitle, Checkbox, Label } from '@/components/ui';

interface GeoIPStatus {
  state: string;
  dbPath?: string | null;
  source?: string | null;
  lastError?: string | null;
}

interface CountryTrackingCardProps {
  trackCountry: boolean;
  updating: boolean;
  refreshing: boolean;
  geoIPStatus?: GeoIPStatus | null;
  onToggle: (enabled: boolean) => void;
  onRetry: () => void;
}

export function CountryTrackingCard({
  trackCountry,
  updating,
  refreshing,
  geoIPStatus,
  onToggle,
  onRetry,
}: CountryTrackingCardProps): React.JSX.Element {
  const geoIPState = geoIPStatus?.state ?? 'disabled';
  const statusMessage = (() => {
    switch (geoIPState) {
      case 'downloading':
        return 'Downloading GeoIP database...';
      case 'missing':
        return 'GeoIP database not available yet.';
      case 'error':
        return 'GeoIP download failed. Use Retry to attempt again.';
      case 'ready':
        return 'GeoIP database is ready.';
      default:
        return 'GeoIP downloads are disabled.';
    }
  })();

  const geoIPBadgeVariant = (): 'default' | 'secondary' | 'outline' | 'destructive' => {
    switch (geoIPState) {
      case 'ready':
        return 'default';
      case 'downloading':
        return 'secondary';
      case 'missing':
        return 'outline';
      case 'error':
        return 'destructive';
      default:
        return 'outline';
    }
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Country Tracking</CardTitle>
        <CardDescription>
          Enable country-level analytics for this site
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="flex items-center gap-3">
          <Checkbox
            id="track-country"
            checked={trackCountry}
            onCheckedChange={(value) => {
              onToggle(value === true);
            }}
            disabled={updating}
          />
          <Label htmlFor="track-country" className="text-sm font-medium">
            Track visitor country
          </Label>
        </div>

        <div className="flex items-center gap-2 flex-wrap text-sm">
          <span className="text-muted-foreground">GeoIP database:</span>
          <Badge variant={geoIPBadgeVariant()} className="uppercase tracking-wide text-[10px]">
            {geoIPState}
          </Badge>
          {geoIPStatus?.source ? (
            <span className="text-xs text-muted-foreground">
              source: {geoIPStatus.source}
            </span>
          ) : null}
        </div>

        {geoIPStatus?.dbPath ? (
          <p className="text-xs text-muted-foreground">
            Path: <span className="font-mono">{geoIPStatus.dbPath}</span>
          </p>
        ) : null}

        {geoIPStatus?.lastError ? (
          <p className="text-xs text-destructive">
            {geoIPStatus.lastError}
          </p>
        ) : null}

        <div className="flex gap-2 items-center">
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={onRetry}
            disabled={refreshing}
          >
            {refreshing ? 'Retrying...' : 'Retry download'}
          </Button>
          {trackCountry ? (
            <span className="text-xs text-muted-foreground">
              Country tracking requires a GeoLite2 database to be downloaded.
            </span>
          ) : null}
        </div>

        <p className="text-xs text-muted-foreground">{statusMessage}</p>
      </CardContent>
    </Card>
  );
}
