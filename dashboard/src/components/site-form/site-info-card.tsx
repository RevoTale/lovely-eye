import React from 'react';
import { Button, Card, CardContent, CardDescription, CardHeader, CardTitle, Input, Label } from '@/components/ui';
import { Globe, Loader2, Plus, Save, X } from 'lucide-react';
import { getNormalizedDomains, normalizeDomainInput } from '@/components/site-form/utils';

interface DomainEntry {
  id: string;
  value: string;
}

interface SiteInfoCardProps {
  siteId?: string | undefined;
  isNew: boolean;
  initialName: string;
  initialDomains: string[];
  creating: boolean;
  updating: boolean;
  onCreate: (name: string, domains: string[]) => Promise<void>;
  onSaveDomains: (name: string, domains: string[]) => Promise<void>;
  onCancel: () => void;
}

export function SiteInfoCard({
  siteId,
  isNew,
  initialName,
  initialDomains,
  creating,
  updating,
  onCreate,
  onSaveDomains,
  onCancel,
}: SiteInfoCardProps): React.JSX.Element {
  const [name, setName] = React.useState(initialName);
  const [formError, setFormError] = React.useState('');
  const nextDomainIdRef = React.useRef(initialDomains.length > 0 ? initialDomains.length + 1 : 2);

  const buildDomainEntries = (values: string[]): DomainEntry[] => {
    if (values.length === 0) {
      return [{ id: '1', value: '' }];
    }
    return values.map((domain, index) => ({ id: String(index + 1), value: domain }));
  };

  const [domains, setDomains] = React.useState<DomainEntry[]>(() => buildDomainEntries(initialDomains));

  React.useEffect(() => {
    setName(initialName);
    setDomains(buildDomainEntries(initialDomains));
    setFormError('');
    nextDomainIdRef.current = initialDomains.length > 0 ? initialDomains.length + 1 : 2;
  }, [initialDomains, initialName, siteId]);

  const hasDomainChanges = React.useMemo(() => {
    const currentDomains = getNormalizedDomains(domains.map((entry) => entry.value));
    const savedDomains = getNormalizedDomains(initialDomains);
    if (currentDomains.length !== savedDomains.length) return true;
    const savedSet = new Set(savedDomains);
    return currentDomains.some((domainValue) => !savedSet.has(domainValue));
  }, [domains, initialDomains]);

  const validateDomains = (): string[] | null => {
    const domainRegex = /^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?(\.[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?)*$/;
    const trimmedDomains = domains
      .map((domainEntry) => normalizeDomainInput(domainEntry.value))
      .filter((domainValue) => domainValue.length > 0);

    const uniqueDomains = Array.from(new Set(trimmedDomains));

    if (uniqueDomains.length === 0) {
      setFormError('At least one domain is required');
      return null;
    }

    if (!uniqueDomains.every((domainValue) => domainRegex.test(domainValue))) {
      setFormError('Please enter valid domains (e.g., example.com)');
      return null;
    }

    return uniqueDomains;
  };

  const handleSubmit = async (event: React.FormEvent): Promise<void> => {
    event.preventDefault();
    const trimmedName = name.trim();
    const validatedDomains = validateDomains();

    if (!trimmedName) {
      setFormError('Name is required');
      return;
    }

    if (trimmedName.length < 1 || trimmedName.length > 100) {
      setFormError('Site name must be between 1 and 100 characters');
      return;
    }

    if (!validatedDomains) {
      return;
    }

    setFormError('');
    try {
      if (isNew) {
        await onCreate(trimmedName, validatedDomains);
      } else {
        await onSaveDomains(trimmedName, validatedDomains);
      }
    } catch (err) {
      setFormError(err instanceof Error ? err.message : 'Failed to save site details');
    }
  };

  const handleDomainChange = (index: number, id: string, value: string): void => {
    const previousPrimary = domains[0]?.value ?? '';
    const normalized = normalizeDomainInput(value);
    setDomains((prev) => {
      return prev.map((entry) => entry.id === id
        ? { ...entry, value: normalized }
        : entry
      );
    });
    if (isNew && index === 0 && (name.trim() === '' || name.trim() === previousPrimary)) {
      setName(normalized);
    }
  };

  return (
    <form onSubmit={(event) => {
      void handleSubmit(event);
    }}>
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
            <Label htmlFor="primary-domain">Domains</Label>
            <div className="space-y-2">
              {domains.map((domainEntry, index) => (
                <div key={domainEntry.id} className="flex items-center gap-2">
                  <Input
                    id={index === 0 ? 'primary-domain' : `domain-${index}`}
                    placeholder={index === 0 ? 'example.com' : 'blog.example.com'}
                    value={domainEntry.value}
                    onChange={(e) => {
                      handleDomainChange(index, domainEntry.id, e.target.value);
                    }}
                    required={index === 0}
                  />
                  {domains.length > 1 ? (
                    <Button
                      type="button"
                      variant="outline"
                      size="icon"
                      onClick={() => {
                        setDomains((prev) => prev.filter((entry) => entry.id !== domainEntry.id));
                      }}
                      aria-label="Remove domain"
                    >
                      <X className="h-4 w-4" />
                    </Button>
                  ) : null}
                </div>
              ))}
            </div>
            <div className="flex flex-wrap items-center gap-3">
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={() => {
                  const nextId = String(nextDomainIdRef.current);
                  nextDomainIdRef.current += 1;
                  setDomains((prev) => [...prev, { id: nextId, value: '' }]);
                }}
              >
                <Plus className="h-4 w-4" />
                Add domain
              </Button>
            </div>
            <p className="text-xs text-muted-foreground">
              Add domains without https://. The first domain is treated as the primary domain.
            </p>
          </div>

          <div className="space-y-2">
            <Label htmlFor="name">Site Name</Label>
            <Input
              id="name"
              placeholder="My Awesome Website"
              value={name}
              onChange={(e) => {
                setName(e.target.value);
              }}
              disabled={!isNew}
              required
            />
            <p className="text-xs text-muted-foreground">
              A friendly name to identify your site
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
                onClick={onCancel}
                disabled={creating}
              >
                Cancel
              </Button>
            </div>
          ) : (
            <div className="flex gap-3 pt-4">
              <Button
                type="submit"
                disabled={updating || !hasDomainChanges}
              >
                {updating ? (
                  <>
                    <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                    Saving...
                  </>
                ) : (
                  <>
                    <Save className="h-4 w-4 mr-2" />
                    Save Domains
                  </>
                )}
              </Button>
            </div>
          )}
          {formError ? (
            <p className="text-xs text-destructive">
              {formError}
            </p>
          ) : null}
        </CardContent>
      </Card>
    </form>
  );
}
