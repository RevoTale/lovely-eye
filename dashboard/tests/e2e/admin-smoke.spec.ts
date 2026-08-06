import { expect, test } from '@playwright/test';

test('an administrator can sign in to the real application', async ({ page }) => {
  await page.goto('login');

  await expect(page.getByRole('heading', { name: 'Welcome back' })).toBeVisible();
  await page.getByLabel('Username').fill('e2e-admin');
  await page.getByLabel('Password').fill('e2e-password');
  await page.getByRole('button', { name: 'Sign in' }).click();

  await expect(page.getByRole('heading', { name: 'Add New Site' })).toBeVisible();
  await expect(page.getByRole('button', { name: 'Create Site', exact: true })).toBeVisible();
});
