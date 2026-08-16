import type { GeoIpStatusFieldsFragment } from '@/shared/api/generated/graphql';
import { Badge } from '@/shared/ui/badge';
import { Button } from '@/shared/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card';
import { Checkbox } from '@/shared/ui/checkbox';
import { Label } from '@/shared/ui/label';

interface CountryTrackingCardProps {
  trackCountry: boolean;
  updating: boolean;
  refreshing: boolean;
  geoIPStatus?: GeoIpStatusFieldsFragment | null;
  onToggle: (enabled: boolean) => void;
  onRetry: () => void;
}

export const CountryTrackingCard = ({
  trackCountry,
  updating,
  refreshing,
  geoIPStatus,
  onToggle,
  onRetry,
}: CountryTrackingCardProps): React.ReactNode => {
  const geoIPState = geoIPStatus?.state ?? 'DISABLED';
  const geoIPSource = geoIPStatus?.source;
  const geoIPDbPath = geoIPStatus?.dbPath;
  const geoIPLastError = geoIPStatus?.lastError;
  const statusMessage = (() => {
    switch (geoIPState) {
      case 'DOWNLOADING':
        return 'Downloading GeoIP database...';
      case 'MISSING':
        return 'GeoIP database not available yet.';
      case 'ERROR':
        return 'GeoIP download failed. Use Retry to attempt again.';
      case 'READY':
        return 'GeoIP database is ready.';
      default:
        return 'GeoIP downloads are disabled.';
    }
  })();

  const geoIPBadgeVariant = (): 'default' | 'secondary' | 'outline' | 'destructive' => {
    switch (geoIPState) {
      case 'READY':
        return 'default';
      case 'DOWNLOADING':
        return 'secondary';
      case 'MISSING':
        return 'outline';
      case 'ERROR':
        return 'destructive';
      default:
        return 'outline';
    }
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Country Tracking</CardTitle>
        <CardDescription>Enable country-level analytics for this site</CardDescription>
      </CardHeader>
      <CardContent className='space-y-4'>
        <div className='flex items-center gap-3'>
          <Checkbox
            id='track-country'
            checked={trackCountry}
            onCheckedChange={(value) => {
              onToggle(value === true);
            }}
            disabled={updating}
          />
          <Label htmlFor='track-country' className='text-sm font-medium'>
            Track visitor country
          </Label>
        </div>

        <div className='flex items-center gap-2 flex-wrap text-sm'>
          <span className='text-muted-foreground'>GeoIP database:</span>
          <Badge variant={geoIPBadgeVariant()} className='uppercase tracking-wide text-[10px]'>
            {geoIPState.toLowerCase()}
          </Badge>
          {geoIPSource !== null && geoIPSource !== undefined && geoIPSource !== '' ? (
            <span className='text-xs text-muted-foreground'>source: {geoIPSource}</span>
          ) : null}
        </div>

        {geoIPDbPath !== null && geoIPDbPath !== undefined && geoIPDbPath !== '' ? (
          <p className='break-all text-xs text-muted-foreground'>
            Path: <span className='font-mono'>{geoIPDbPath}</span>
          </p>
        ) : null}

        {geoIPLastError !== null && geoIPLastError !== undefined && geoIPLastError !== '' ? (
          <p className='break-words text-xs text-destructive'>{geoIPLastError}</p>
        ) : null}

        <div className='flex flex-col items-start gap-2 sm:flex-row sm:items-center'>
          <Button type='button' variant='outline' size='sm' onClick={onRetry} disabled={refreshing}>
            {refreshing ? 'Retrying...' : 'Retry download'}
          </Button>
          {trackCountry ? (
            <span className='text-xs text-muted-foreground'>
              Country tracking requires a GeoLite2 database to be downloaded.
            </span>
          ) : null}
        </div>

        <p className='text-xs text-muted-foreground'>{statusMessage}</p>
      </CardContent>
    </Card>
  );
};
