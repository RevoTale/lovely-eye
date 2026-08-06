import { type FormEvent, useRef, useState } from 'react';
import { normalizeDomain, validateDomains } from '@/shared/lib/domains';

export interface DomainField {
  id: number;
  value: string;
}

interface CreateSiteFormOptions {
  onCreate: (name: string, domains: string[]) => Promise<void>;
}

export function useCreateSiteForm({ onCreate }: CreateSiteFormOptions) {
  const nextDomainId = useRef(2);
  const [domains, setDomains] = useState<DomainField[]>([{ id: 1, value: '' }]);
  const [name, setName] = useState('');
  const [error, setError] = useState<string | null>(null);

  const addDomain = (): void => {
    const id = nextDomainId.current;
    nextDomainId.current += 1;
    setDomains((current) => [...current, { id, value: '' }]);
    setError(null);
  };

  const removeDomain = (id: number): void => {
    setDomains((current) => current.filter((domain) => domain.id !== id));
    setError(null);
  };

  const changeDomain = (id: number, value: string): void => {
    const previousPrimary = domains[0]?.value ?? '';
    const normalized = normalizeDomain(value);
    setDomains((current) =>
      current.map((domain) => (domain.id === id ? { ...domain, value: normalized } : domain))
    );
    if (id === domains[0]?.id && (name.trim() === '' || name.trim() === previousPrimary)) {
      setName(normalized);
    }
    setError(null);
  };

  const changeName = (value: string): void => {
    setName(value);
    setError(null);
  };

  const submit = async (event: FormEvent): Promise<void> => {
    event.preventDefault();
    const normalizedName = name.trim();
    if (normalizedName === '' || normalizedName.length > 100) {
      setError('Site name must be between 1 and 100 characters');
      return;
    }
    const result = validateDomains(domains.map(({ value }) => value));
    if (result.error !== null) {
      setError(result.error);
      return;
    }
    setError(null);
    try {
      await onCreate(normalizedName, result.domains);
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : 'Failed to create site');
    }
  };

  return {
    addDomain,
    changeDomain,
    domains,
    error,
    name,
    removeDomain,
    setName: changeName,
    submit,
  };
}
