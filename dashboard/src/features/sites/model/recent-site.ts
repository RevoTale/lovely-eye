import { getBasePath } from '@/shared/config/runtime';

const STORAGE_KEY_PREFIX = 'lovely-eye:last-site';

function storageKey(): string {
  return `${STORAGE_KEY_PREFIX}:${getBasePath()}`;
}

export function getRememberedSiteId(): string | null {
  try {
    return localStorage.getItem(storageKey());
  } catch {
    return null;
  }
}

export function rememberSite(siteId: string): void {
  try {
    localStorage.setItem(storageKey(), siteId);
  } catch {
    // Storage may be unavailable in privacy modes; navigation remains correct without persistence.
  }
}
