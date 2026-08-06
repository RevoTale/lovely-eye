export interface BlockedIPEntry {
  id: string;
  value: string;
}

export interface TrafficBlockingCardProps {
  initialBlockedIPs: string[];
  initialBlockedCountries: string[];
  savingBlockedIPs: boolean;
  savingBlockedCountries: boolean;
  geoIPReady: boolean;
  onUpdateBlockedIPs: (blockedIPs: string[]) => Promise<void>;
  onUpdateBlockedCountries: (blockedCountries: string[]) => Promise<void>;
}
