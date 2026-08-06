import { expect, test } from '@playwright/test';
import { createSite, signInAsAdmin } from './helpers/admin';
import { installGraphQLOperationController } from './helpers/graphql-gate';

test('auth, first load, mutation feedback, and retained refresh state stay coherent', async ({
  context,
  page,
}) => {
  const operations = await installGraphQLOperationController(page);
  const authResolution = operations.blockNext('Me');
  await page.goto('login');
  await authResolution.seen;
  await expect(page.getByRole('heading', { name: 'Loading dashboard' })).toBeVisible();
  await expect(page.getByRole('heading', { name: 'Welcome back' })).toHaveCount(0);
  authResolution.release();
  await expect(page.getByRole('heading', { name: 'Welcome back' })).toBeVisible();

  const sitesLoad = operations.blockNext('Sites');
  await page.getByLabel('Username').fill('e2e-admin');
  await page.getByLabel('Password').fill('e2e-password');
  await page.getByRole('button', { name: 'Sign in' }).click();
  await sitesLoad.seen;
  await expect(page.locator('.animate-pulse').first()).toBeVisible();
  sitesLoad.release();
  await expect(page.getByRole('heading', { name: 'Add New Site' })).toBeVisible();

  await page.goto('sites/new');
  await page.getByLabel('Domains').fill('primary.example');
  await page.getByRole('button', { name: 'Add domain' }).click();
  await page.locator('#domain-1').fill('secondary.primary.example');
  await page.getByLabel('Site Name').fill('State Contract Site');

  const createMutation = operations.blockNext('CreateSite');
  await page.getByRole('button', { name: 'Create Site' }).click();
  await createMutation.seen;
  await expect(page.getByRole('button', { name: 'Creating...' })).toBeDisabled();
  await expect(page.getByLabel('Domains')).toHaveValue('primary.example');
  await expect(page.locator('#domain-1')).toHaveValue('secondary.primary.example');
  createMutation.release();
  await expect(page.getByRole('heading', { name: 'State Contract Site' })).toBeVisible();
  await expect(page.getByText('primary.example · secondary.primary.example')).toBeVisible();
  await page.goto('./');
  await expect(page.getByRole('heading', { name: 'State Contract Site' })).toBeVisible();

  const siteURL = page.url();
  const freshPage = await context.newPage();
  const freshOperations = await installGraphQLOperationController(freshPage);
  const firstSiteLoad = freshOperations.blockNext('Site');
  await freshPage.goto(siteURL);
  await firstSiteLoad.seen;
  await expect(freshPage.locator('.animate-pulse').first()).toBeVisible();
  firstSiteLoad.release();
  const siteHeading = freshPage.getByRole('heading', { name: 'State Contract Site' });
  await expect(siteHeading).toBeVisible();

  const headingBeforeRefresh = await siteHeading.boundingBox();
  expect(headingBeforeRefresh).not.toBeNull();
  const dashboardRefresh = freshOperations.blockNext('Dashboard');
  await freshPage.getByRole('tab', { name: '7d' }).click();
  await dashboardRefresh.seen;
  await expect(siteHeading).toBeVisible();
  const headingDuringRefresh = await siteHeading.boundingBox();
  expect(headingDuringRefresh?.y).toBe(headingBeforeRefresh?.y);
  expect(headingDuringRefresh?.height).toBe(headingBeforeRefresh?.height);
  await expect(freshPage.getByText('Refreshing', { exact: true }).first()).toBeVisible();
  dashboardRefresh.release();
  await expect(freshPage).toHaveURL(/preset=7d/u);
});

test('site switching and URL-owned analytics state survive history and refresh', async ({
  page,
}) => {
  await signInAsAdmin(page);
  await createSite(page, 'Switch Alpha', ['alpha.example']);
  const betaURL = await createSite(page, 'Switch Beta', ['beta.example', 'app.beta.example']);

  await page.goto('./');
  await expect(page.getByRole('heading', { name: 'Switch Beta' })).toBeVisible();
  await page.goto('sites');
  await page.getByRole('link', { name: /Switch Alpha/u }).click();
  await expect(page.getByRole('heading', { name: 'Switch Alpha' })).toBeVisible();
  await page.getByRole('link', { name: 'Sites' }).click();
  await page.getByRole('link', { name: /Switch Beta/u }).click();
  await expect(page.getByRole('heading', { name: 'Switch Beta' })).toBeVisible();

  await page.getByRole('tab', { name: '7d' }).click();
  await expect(page).toHaveURL(/preset=7d/u);
  await page.getByRole('tab', { name: '30d' }).click();
  await expect(page).toHaveURL(/preset=30d/u);
  await page.goBack();
  await expect(page).toHaveURL(/preset=7d/u);
  await expect(page.getByRole('tab', { name: '7d' })).toHaveAttribute('data-state', 'active');
  await page.goForward();
  await expect(page).toHaveURL(/preset=30d/u);
  await page.reload();
  await expect(page).toHaveURL(`${betaURL.replace(/\?.*$/u, '')}?preset=30d`);
  await expect(page.getByRole('tab', { name: '30d' })).toHaveAttribute('data-state', 'active');
});
