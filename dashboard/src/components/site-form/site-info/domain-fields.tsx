import type { FunctionComponent } from 'react';
import { Plus, X } from 'lucide-react';
import { Button, Input, Label } from '@/components/ui';
import type { DomainEntry } from './types';

const FIRST_DOMAIN_INDEX = 0;
const MIN_DOMAIN_COUNT = 1;

interface DomainFieldsProps {
  domains: DomainEntry[];
  onAddDomain: () => void;
  onDomainChange: (index: number, id: string, value: string) => void;
  onRemoveDomain: (id: string) => void;
}

const DomainFields: FunctionComponent<DomainFieldsProps> = ({
  domains,
  onAddDomain,
  onDomainChange,
  onRemoveDomain,
}) => (
  <div className="space-y-2">
    <Label htmlFor="primary-domain">Domains</Label>
    <div className="space-y-2">
      {domains.map((domainEntry, index) => (
        <div key={domainEntry.id} className="flex items-center gap-2">
          <Input
            id={index === FIRST_DOMAIN_INDEX ? 'primary-domain' : `domain-${index}`}
            placeholder={index === FIRST_DOMAIN_INDEX ? 'example.com' : 'blog.example.com'}
            value={domainEntry.value}
            onChange={(event) => onDomainChange(index, domainEntry.id, event.target.value)}
            required={index === FIRST_DOMAIN_INDEX}
          />
          {domains.length > MIN_DOMAIN_COUNT ? (
            <Button
              type="button"
              variant="outline"
              size="icon"
              onClick={() => onRemoveDomain(domainEntry.id)}
              aria-label="Remove domain"
            >
              <X className="h-4 w-4" />
            </Button>
          ) : null}
        </div>
      ))}
    </div>
    <div className="flex flex-wrap items-center gap-3">
      <Button type="button" variant="outline" size="sm" onClick={onAddDomain}>
        <Plus className="h-4 w-4" />
        Add domain
      </Button>
    </div>
    <p className="text-xs text-muted-foreground">
      Add domains without https://. The first domain is treated as the primary domain.
    </p>
  </div>
);

export default DomainFields;
