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

const DOMAIN_ID_INCREMENT = 1;
const EMPTY_COUNT = 0;
const EMPTY_STRING = '';
const FIRST_DOMAIN_ID = 1;
const FIRST_DOMAIN_INDEX = 0;
const MIN_DOMAIN_COUNT = 1;
const MIN_NAME_LENGTH = 1;
const MAX_NAME_LENGTH = 100;
const SECONDARY_DOMAIN_START = 2;

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
  const nextDomainIdRef = React.useRef(
    initialDomains.length > EMPTY_COUNT
      ? initialDomains.length + DOMAIN_ID_INCREMENT
      : SECONDARY_DOMAIN_START
  );

  const buildDomainEntries = (values: string[]): DomainEntry[] => {
    if (values.length === EMPTY_COUNT) {
      return [{ id: String(FIRST_DOMAIN_ID), value: EMPTY_STRING }];
    }
    return values.map((domain, index) => ({
      id: String(index + DOMAIN_ID_INCREMENT),
      value: domain,
    }));
  };

  const [domains, setDomains] = React.useState<DomainEntry[]>(() => buildDomainEntries(initialDomains));

  React.useEffect(() => {
    setName(initialName);
    setDomains(buildDomainEntries(initialDomains));
    setFormError('');
    nextDomainIdRef.current =
      initialDomains.length > EMPTY_COUNT
        ? initialDomains.length + DOMAIN_ID_INCREMENT
        : SECONDARY_DOMAIN_START;
  }, [initialDomains, initialName, siteId]);

  const hasDomainChanges = React.useMemo(() => {
    const currentDomains = getNormalizedDomains(domains.map((entry) => entry.value));
    const savedDomains = getNormalizedDomains(initialDomains);
    if (currentDomains.length !== savedDomains.length) return true;
    const savedSet = new Set(savedDomains);
    return currentDomains.some((domainValue) => !savedSet.has(domainValue));
  }, [domains, initialDomains]);

  const validateDomains = (): string[] | null => {
    const domainRegex =
      /^[a-z0-9](?:[a-z0-9\-]{0,61}[a-z0-9])?(?:\.[a-z0-9](?:[a-z0-9\-]{0,61}[a-z0-9])?)*$/v;
    const trimmedDomains = domains
      .map((domainEntry) => normalizeDomainInput(domainEntry.value))
      .filter((domainValue) => domainValue.length > EMPTY_COUNT);

    const uniqueDomains = Array.from(new Set(trimmedDomains));

    if (uniqueDomains.length === EMPTY_COUNT) {
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

    if (trimmedName === EMPTY_STRING) {
      setFormError('Name is required');
      return;
    }

    if (trimmedName.length < MIN_NAME_LENGTH || trimmedName.length > MAX_NAME_LENGTH) {
      setFormError('Site name must be between 1 and 100 characters');
      return;
    }

    if (validatedDomains === null) {
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
    const previousPrimary = domains[FIRST_DOMAIN_INDEX]?.value ?? EMPTY_STRING;
    const normalized = normalizeDomainInput(value);
    setDomains((prev) =>
      prev.map((entry) => (entry.id === id ? { ...entry, value: normalized } : entry))
    );
    const trimmedName = name.trim();
    if (
      isNew &&
      index === FIRST_DOMAIN_INDEX &&
      (trimmedName === EMPTY_STRING || trimmedName === previousPrimary)
    ) {
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
                    id={index === FIRST_DOMAIN_INDEX ? 'primary-domain' : `domain-${index}`}
                    placeholder={index === FIRST_DOMAIN_INDEX ? 'example.com' : 'blog.example.com'}
                    value={domainEntry.value}
                    onChange={(e) => {
                      handleDomainChange(index, domainEntry.id, e.target.value);
                    }}
                    required={index === FIRST_DOMAIN_INDEX}
                  />
                  {domains.length > MIN_DOMAIN_COUNT ? (
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
                  nextDomainIdRef.current += DOMAIN_ID_INCREMENT;
                  setDomains((prev) => [...prev, { id: nextId, value: EMPTY_STRING }]);
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
          {formError === EMPTY_STRING ? null : (
            <p className="text-xs text-destructive">
              {formError}
            </p>
          )}
        </CardContent>
      </Card>
    </form>
  );
}
