import { expect, test } from '@playwright/test';
import { createSite, signInAsAdmin } from './helpers/admin';
import { installGraphQLOperationController } from './helpers/graphql-gate';

test('site settings preserve neutral loading, multi-domain mutation, and destructive-action state', async ({
  context,
  page,
}) => {
  await signInAsAdmin(page);
  const analyticsURL = await createSite(page, 'Settings Contract Site', [
    'settings.example',
    'app.settings.example',
  ]);
  const settingsURL = analyticsURL.replace(/\/analytics(?:\?.*)?$/u, '/settings');

  const settingsPage = await context.newPage();
  const operations = await installGraphQLOperationController(settingsPage);
  const firstLoad = operations.blockNext('Site');
  await settingsPage.goto(settingsURL);
  await firstLoad.seen;
  await expect(settingsPage.getByRole('heading', { name: 'Loading dashboard' })).toBeVisible();
  await expect(settingsPage.locator('.animate-pulse')).toHaveCount(0);
  await expect(settingsPage.getByRole('heading', { name: 'Settings Contract Site' })).toHaveCount(
    0
  );
  firstLoad.release();
  await expect(settingsPage.getByRole('heading', { name: 'Settings Contract Site' })).toBeVisible();

  await settingsPage.locator('#domain-1').fill('admin.settings.example');
  await settingsPage.getByRole('button', { name: 'Add domain' }).click();
  await settingsPage.locator('#domain-2').fill('status.settings.example');
  const updateMutation = operations.blockNext('UpdateSite');
  await settingsPage.getByRole('button', { name: 'Save Domains' }).click();
  await updateMutation.seen;
  await expect(settingsPage.getByRole('button', { name: 'Saving...' })).toBeDisabled();
  await expect(settingsPage.locator('#domain-1')).toHaveValue('admin.settings.example');
  await expect(settingsPage.locator('#domain-2')).toHaveValue('status.settings.example');
  updateMutation.release();
  await expect(settingsPage.getByRole('button', { name: 'Save Domains' })).toBeDisabled();

  await settingsPage.reload();
  await expect(settingsPage.locator('#domain-1')).toHaveValue('admin.settings.example');
  await expect(settingsPage.locator('#domain-2')).toHaveValue('status.settings.example');

  await settingsPage.getByRole('button', { name: 'Delete Site' }).click();
  const deleteMutation = operations.blockNext('DeleteSite');
  await settingsPage.getByRole('button', { name: 'Confirm Delete' }).click();
  await deleteMutation.seen;
  await expect(settingsPage.getByRole('button', { name: 'Deleting...' })).toBeDisabled();
  deleteMutation.release();
  await expect(settingsPage).not.toHaveURL(settingsURL);
  await settingsPage.goto(settingsURL);
  await expect(settingsPage.getByText(/site not found/iu)).toBeVisible();
  await settingsPage.goto(analyticsURL);
  await expect(settingsPage.getByText(/site not found/iu)).toBeVisible();
});
