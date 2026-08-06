import { expect, test } from '@playwright/test';
import { createSite, signInAsAdmin } from './helpers/admin';
import { installGraphQLOperationController } from './helpers/graphql-gate';

test('GeoIP failures remain actionable after enable and retry', async ({ page }) => {
  await signInAsAdmin(page);
  await createSite(page, 'GeoIP Status Site', ['geoip-status.example']);
  await page.getByRole('link', { name: 'Settings' }).click();
  await expect(page.getByRole('heading', { name: 'GeoIP Status Site' })).toBeVisible();

  await page.getByRole('checkbox', { name: 'Track visitor country' }).click();
  await expect(page.getByText('error', { exact: true })).toBeVisible();
  await expect(page.getByText(/connection refused/u)).toBeVisible();

  const operations = await installGraphQLOperationController(page);
  const refreshMutation = operations.blockNext('RefreshGeoIPDatabase');
  await page.getByRole('button', { name: 'Retry download' }).click();
  await refreshMutation.seen;
  await expect(page.getByRole('button', { name: 'Retrying...' })).toBeDisabled();
  refreshMutation.release();
  await expect(page.getByRole('button', { name: 'Retry download' })).toBeEnabled();
  await expect(page.getByText(/connection refused/u)).toBeVisible();
});
