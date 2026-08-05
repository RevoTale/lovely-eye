import { useMemo, useRef, useState, type FormEvent } from 'react';
import { getNormalizedDomains, normalizeDomainInput } from '@/components/site-form/utils';
import type { DomainEntry, SiteInfoCardProps } from './types';

const DOMAIN_ID_INCREMENT = 1;
const EMPTY_COUNT = 0;
const EMPTY_STRING = '';
const FIRST_DOMAIN_ID = 1;
const FIRST_DOMAIN_INDEX = 0;
const MAX_NAME_LENGTH = 100;
const MIN_NAME_LENGTH = 1;
const SECONDARY_DOMAIN_START = 2;

export function useSiteInfoForm({
  initialDomains,
  initialName,
  isNew,
  onCreate,
  onSaveDomains,
}: Pick<SiteInfoCardProps, 'initialDomains' | 'initialName' | 'isNew' | 'onCreate' | 'onSaveDomains'>) {
  const [name, setName] = useState(initialName);
  const [formError, setFormError] = useState('');
  const nextDomainIdRef = useRef(
    initialDomains.length > EMPTY_COUNT
      ? initialDomains.length + DOMAIN_ID_INCREMENT
      : SECONDARY_DOMAIN_START
  );
  const [domains, setDomains] = useState<DomainEntry[]>(() => buildDomainEntries(initialDomains));
  const hasDomainChanges = useMemo(() => domainsChanged(domains, initialDomains), [domains, initialDomains]);

  const addDomain = (): void => {
    const nextId = String(nextDomainIdRef.current);
    nextDomainIdRef.current += DOMAIN_ID_INCREMENT;
    setDomains((prev) => [...prev, { id: nextId, value: EMPTY_STRING }]);
  };

  const removeDomain = (id: string): void => {
    setDomains((prev) => prev.filter((entry) => entry.id !== id));
  };

  const handleDomainChange = (index: number, id: string, value: string): void => {
    const previousPrimary = domains[FIRST_DOMAIN_INDEX]?.value ?? EMPTY_STRING;
    const normalized = normalizeDomainInput(value);
    setDomains((prev) => prev.map((entry) => (entry.id === id ? { ...entry, value: normalized } : entry)));
    const trimmedName = name.trim();
    if (
      isNew &&
      index === FIRST_DOMAIN_INDEX &&
      (trimmedName === EMPTY_STRING || trimmedName === previousPrimary)
    ) {
      setName(normalized);
    }
  };

  const handleSubmit = async (event: FormEvent): Promise<void> => {
    event.preventDefault();
    const trimmedName = name.trim();
    const domainValidation = validateDomains(domains);
    const validationError = validateForm(trimmedName, domainValidation.error);
    if (validationError !== '') {
      setFormError(validationError);
      return;
    }
    setFormError('');
    try {
      if (isNew) {
        await onCreate(trimmedName, domainValidation.domains);
      } else {
        await onSaveDomains(trimmedName, domainValidation.domains);
      }
    } catch (err) {
      setFormError(err instanceof Error ? err.message : 'Failed to save site details');
    }
  };

  return {
    addDomain,
    domains,
    formError,
    handleDomainChange,
    handleSubmit,
    hasDomainChanges,
    name,
    removeDomain,
    setName,
  };
}

function buildDomainEntries(values: string[]): DomainEntry[] {
  if (values.length === EMPTY_COUNT) return [{ id: String(FIRST_DOMAIN_ID), value: EMPTY_STRING }];
  return values.map((domain, index) => ({
    id: String(index + DOMAIN_ID_INCREMENT),
    value: domain,
  }));
}

function domainsChanged(domains: DomainEntry[], initialDomains: string[]): boolean {
  const currentDomains = getNormalizedDomains(domains.map((entry) => entry.value));
  const savedDomains = getNormalizedDomains(initialDomains);
  if (currentDomains.length !== savedDomains.length) return true;
  const savedSet = new Set(savedDomains);
  return currentDomains.some((domainValue) => !savedSet.has(domainValue));
}

function validateDomains(domains: DomainEntry[]): { domains: string[]; error: string } {
  const domainRegex =
    /^[a-z0-9](?:[a-z0-9\-]{0,61}[a-z0-9])?(?:\.[a-z0-9](?:[a-z0-9\-]{0,61}[a-z0-9])?)*$/v;
  const uniqueDomains = Array.from(
    new Set(
      domains
        .map((domainEntry) => normalizeDomainInput(domainEntry.value))
        .filter((domainValue) => domainValue.length > EMPTY_COUNT)
    )
  );
  if (uniqueDomains.length === EMPTY_COUNT) {
    return { domains: [], error: 'At least one domain is required' };
  }
  if (!uniqueDomains.every((domainValue) => domainRegex.test(domainValue))) {
    return { domains: [], error: 'Please enter valid domains (e.g., example.com)' };
  }
  return { domains: uniqueDomains, error: '' };
}

function validateForm(name: string, domainError: string): string {
  if (name === EMPTY_STRING) return 'Name is required';
  if (name.length < MIN_NAME_LENGTH || name.length > MAX_NAME_LENGTH) {
    return 'Site name must be between 1 and 100 characters';
  }
  if (domainError !== '') return domainError;
  return '';
}
