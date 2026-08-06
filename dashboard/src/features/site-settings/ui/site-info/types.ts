export interface DomainEntry {
  id: string;
  value: string;
}

export interface SiteInfoCardProps {
  initialName: string;
  initialDomains: string[];
  updating: boolean;
  onSaveDomains: (name: string, domains: string[]) => Promise<void>;
}
