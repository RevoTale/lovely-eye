import { expect, test } from '@playwright/test';
import { createSite, signInAsAdmin } from './helpers/admin';
import { installGraphQLOperationController } from './helpers/graphql-gate';

test('event definitions preserve explicit mutation feedback across create, edit, and delete', async ({
  page,
}) => {
  const operations = await installGraphQLOperationController(page);
  await signInAsAdmin(page);
  await createSite(page, 'Event Definitions Site', ['events.example']);
  await page.getByRole('link', { name: 'Settings' }).click();
  await expect(page.getByText('Event Definitions', { exact: true })).toBeVisible();

  await page.getByRole('button', { name: 'New event name' }).click();
  await page.getByLabel('Event Name').fill('signup_completed');
  const createMutation = operations.blockNext('UpsertEventDefinition');
  await page.getByRole('button', { name: 'Save Definition' }).click();
  await createMutation.seen;
  await expect(page.getByRole('button', { name: 'Saving...' })).toBeDisabled();
  await expect(page.getByLabel('Event Name')).toHaveValue('signup_completed');
  createMutation.release();
  await expect(page.getByText('signup_completed', { exact: true })).toBeVisible();

  await page.getByRole('button', { name: 'Edit' }).click();
  await page.getByRole('button', { name: 'Add field' }).click();
  await page.getByLabel('Field key 1').fill('plan');
  await page.getByRole('button', { name: 'Update Definition' }).click();
  await expect(page.getByText('1 field', { exact: true })).toBeVisible();

  await page.getByRole('button', { name: 'Delete', exact: true }).click();
  const deleteMutation = operations.blockNext('DeleteEventDefinition');
  await page.getByRole('button', { name: 'Confirm delete' }).click();
  await deleteMutation.seen;
  await expect(page.getByRole('button', { name: 'Confirm delete' })).toBeDisabled();
  deleteMutation.release();
  await expect(page.getByText('signup_completed', { exact: true })).toHaveCount(0);
  await expect(page.getByText('No event definitions yet.')).toBeVisible();
});
