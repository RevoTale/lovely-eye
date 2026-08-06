import { type FormEvent, useMemo, useRef, useState } from 'react';
import { normalizeDomain, normalizeDomains, validateDomains } from '@/shared/lib/domains';
import type { DomainEntry, SiteInfoCardProps } from './types';

const DOMAIN_ID_INCREMENT = 1;
const EMPTY_COUNT = 0;
const EMPTY_STRING = '';
const FIRST_DOMAIN_ID = 1;
const MAX_NAME_LENGTH = 100;
const MIN_NAME_LENGTH = 1;
const SECONDARY_DOMAIN_START = 2;

export function useSiteInfoForm({
  initialDomains,
  initialName,
  onSaveDomains,
}: Pick<SiteInfoCardProps, 'initialDomains' | 'initialName' | 'onSaveDomains'>) {
  const name = initialName;
  const [formError, setFormError] = useState('');
  const nextDomainIdRef = useRef(
    initialDomains.length > EMPTY_COUNT
      ? initialDomains.length + DOMAIN_ID_INCREMENT
      : SECONDARY_DOMAIN_START
  );
  const [domains, setDomains] = useState<DomainEntry[]>(() => buildDomainEntries(initialDomains));
  const hasDomainChanges = useMemo(
    () => domainsChanged(domains, initialDomains),
    [domains, initialDomains]
  );

  const addDomain = (): void => {
    const nextId = String(nextDomainIdRef.current);
    nextDomainIdRef.current += DOMAIN_ID_INCREMENT;
    setDomains((prev) => [...prev, { id: nextId, value: EMPTY_STRING }]);
  };

  const removeDomain = (id: string): void => {
    setDomains((prev) => prev.filter((entry) => entry.id !== id));
  };

  const handleDomainChange = (_index: number, id: string, value: string): void => {
    const normalized = normalizeDomain(value);
    setDomains((prev) =>
      prev.map((entry) => (entry.id === id ? { ...entry, value: normalized } : entry))
    );
  };

  const handleSubmit = async (event: FormEvent): Promise<void> => {
    event.preventDefault();
    const trimmedName = name.trim();
    const domainValidation = validateDomains(domains.map((entry) => entry.value));
    const validationError = validateForm(trimmedName, domainValidation.error);
    if (validationError !== '') {
      setFormError(validationError);
      return;
    }
    setFormError('');
    try {
      await onSaveDomains(trimmedName, domainValidation.domains);
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
  const currentDomains = normalizeDomains(domains.map((entry) => entry.value));
  const savedDomains = normalizeDomains(initialDomains);
  if (currentDomains.length !== savedDomains.length) return true;
  const savedSet = new Set(savedDomains);
  return currentDomains.some((domainValue) => !savedSet.has(domainValue));
}

function validateForm(name: string, domainError: string | null): string {
  if (name === EMPTY_STRING) return 'Name is required';
  if (name.length < MIN_NAME_LENGTH || name.length > MAX_NAME_LENGTH) {
    return 'Site name must be between 1 and 100 characters';
  }
  if (domainError !== null) return domainError;
  return '';
}
