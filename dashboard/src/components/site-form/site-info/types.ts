export interface DomainEntry {
  id: string;
  value: string;
}

export interface SiteInfoCardProps {
  isNew: boolean;
  initialName: string;
  initialDomains: string[];
  creating: boolean;
  updating: boolean;
  onCreate: (name: string, domains: string[]) => Promise<void>;
  onSaveDomains: (name: string, domains: string[]) => Promise<void>;
  onCancel: () => void;
}
