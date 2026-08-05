import type { ReactElement } from 'react';
import { Shield } from 'lucide-react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui';
import BlockedCountrySection from './traffic-blocking/blocked-country-section';
import BlockedIPSection from './traffic-blocking/blocked-ip-section';
import type { TrafficBlockingCardProps } from './traffic-blocking/types';
import { useTrafficBlocking } from './traffic-blocking/use-traffic-blocking';

export const TrafficBlockingCard = (props: TrafficBlockingCardProps): ReactElement => {
  const trafficBlocking = useTrafficBlocking(props);

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10">
            <Shield className="h-4 w-4 text-primary" />
          </div>
          Traffic Blocking
        </CardTitle>
        <CardDescription>Block specific IPs and countries from being tracked</CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        <BlockedIPSection
          blockedIPCount={trafficBlocking.blockedIPCount}
          blockedIPs={trafficBlocking.blockedIPs}
          ipActionError={trafficBlocking.ipActionError}
          newIPError={trafficBlocking.newIPError}
          newIPValue={trafficBlocking.newIPValue}
          onAddIP={() => void trafficBlocking.handleAddIP()}
          onRemoveBlockedIP={(value) => void trafficBlocking.handleRemoveBlockedIP(value)}
          savingBlockedIPs={props.savingBlockedIPs}
          setNewIPError={trafficBlocking.setNewIPError}
          setNewIPValue={trafficBlocking.setNewIPValue}
        />
        <BlockedCountrySection
          blockedCountries={trafficBlocking.blockedCountries}
          blockedCountryCount={trafficBlocking.blockedCountryCount}
          countryActionError={trafficBlocking.countryActionError}
          countryNameLookup={trafficBlocking.countryNameLookup}
          countrySearch={trafficBlocking.countrySearch}
          geoIPCountriesLoading={trafficBlocking.geoIPCountriesLoading}
          geoIPReady={props.geoIPReady}
          matchingCountries={trafficBlocking.matchingCountries}
          normalizedBlockedCountries={trafficBlocking.normalizedBlockedCountries}
          onAddBlockedCountry={trafficBlocking.handleAddBlockedCountry}
          onCountrySearchKeyDown={trafficBlocking.handleCountrySearchKeyDown}
          onRemoveBlockedCountry={(code) => void trafficBlocking.handleRemoveBlockedCountry(code)}
          savingBlockedCountries={props.savingBlockedCountries}
          setCountrySearch={trafficBlocking.setCountrySearch}
          shouldSearchCountries={trafficBlocking.shouldSearchCountries}
          trimmedCountrySearch={trafficBlocking.trimmedCountrySearch}
        />
      </CardContent>
    </Card>
  );
};
