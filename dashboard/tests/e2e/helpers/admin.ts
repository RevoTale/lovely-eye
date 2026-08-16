import { expect, type Page } from '@playwright/test';

export const signInAsAdmin = async (page: Page, loginURL = 'login'): Promise<void> => {
  await page.goto(loginURL);
  await expect(page.getByRole('heading', { name: 'Welcome back' })).toBeVisible();
  await page.getByLabel('Username').fill('e2e-admin');
  await page.getByLabel('Password').fill('e2e-password');
  await page.getByRole('button', { name: 'Sign in' }).click();
  await expect(page.getByRole('button', { name: 'Open user menu' })).toBeVisible();
};

export const createSite = async (
  page: Page,
  name: string,
  domains: [string, ...string[]]
): Promise<string> => {
  await page.goto('sites/new');
  await expect(page.getByRole('heading', { name: 'Add New Site' })).toBeVisible();
  await page.getByLabel('Domains').fill(domains[0]);
  for (const [index, domain] of domains.slice(1).entries()) {
    await page.getByRole('button', { name: 'Add domain' }).click();
    await page.locator(`#domain-${index + 1}`).fill(domain);
  }
  await page.getByLabel('Site Name').fill(name);
  await page.getByRole('button', { name: 'Create Site' }).click();
  await expect(page.getByRole('heading', { name })).toBeVisible();
  return page.url();
};
