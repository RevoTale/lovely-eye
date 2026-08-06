import { expect, test } from '@playwright/test';
import { signInAsAdmin } from './helpers/admin';

test('site creation recovers after an invalid domain is corrected', async ({ page }) => {
  await signInAsAdmin(page);
  await page.goto('sites/new');

  await page.getByLabel('Domains').fill('not a web address');
  await page.getByLabel('Site Name').fill('Recovered Site');
  await page.getByRole('button', { name: 'Create Site' }).click();
  await expect(page.getByText('Please enter valid domains (e.g., example.com)')).toBeVisible();

  await page.getByLabel('Domains').fill('recovered.example');
  await expect(page.getByText('Please enter valid domains (e.g., example.com)')).toHaveCount(0);
  await page.getByRole('button', { name: 'Create Site' }).click();

  await expect(page.getByRole('heading', { name: 'Recovered Site' })).toBeVisible();
  await expect(page.getByText('recovered.example', { exact: true })).toBeVisible();
});
