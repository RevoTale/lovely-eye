import { expect, test } from '@playwright/test';

const { TEST_BASE_PATH = '/' } = process.env;
const NORMALIZED_BASE_PATH =
  TEST_BASE_PATH === '/' ? '' : `/${TEST_BASE_PATH.split('/').filter(Boolean).join('/')}`;
const EXPECTED_COOKIE_PATH = NORMALIZED_BASE_PATH === '' ? '/' : NORMALIZED_BASE_PATH;

test('runtime base path owns assets, routes, GraphQL, and auth cookies', async ({
  context,
  page,
  request,
}) => {
  await page.goto('login');
  await expect(page.getByRole('heading', { name: 'Welcome back' })).toBeVisible();

  const configResponse = await request.get(`${NORMALIZED_BASE_PATH}/config.js`);
  expect(configResponse.ok()).toBe(true);
  expect(await configResponse.text()).toContain(`BASE_PATH: '${NORMALIZED_BASE_PATH}'`);

  const logoSource = await page.getByAltText('Lovely Eye').getAttribute('src');
  expect(logoSource).not.toBeNull();
  const logoURL = new URL(logoSource ?? '', await page.evaluate(() => document.baseURI));
  expect(logoURL.pathname).toBe(`${NORMALIZED_BASE_PATH}/favicon.svg`);
  expect((await request.get(logoURL.toString())).ok()).toBe(true);

  await page.getByLabel('Username').fill('e2e-admin');
  await page.getByLabel('Password').fill('e2e-password');
  await page.getByRole('button', { name: 'Sign in' }).click();
  await expect(page.getByRole('button', { name: 'Open user menu' })).toBeVisible();

  const authCookies = (await context.cookies()).filter(({ name }) => name.startsWith('le_'));
  expect(authCookies).toHaveLength(2);
  expect(authCookies.every(({ path }) => path === EXPECTED_COOKIE_PATH)).toBe(true);

  await page.goto('sites/new');
  await expect(page.getByRole('heading', { name: 'Add New Site' })).toBeVisible();
  await page.reload();
  await expect(page.getByRole('heading', { name: 'Add New Site' })).toBeVisible();

  const moduleSource = await page.locator('script[type="module"][src]').last().getAttribute('src');
  expect(moduleSource).not.toBeNull();
  const documentBaseURL = await page.evaluate(() => document.baseURI);
  const moduleURL = new URL(moduleSource ?? '', documentBaseURL);
  expect(moduleURL.pathname.startsWith(`${NORMALIZED_BASE_PATH}/assets/`)).toBe(true);
  expect((await request.get(moduleURL.toString())).ok()).toBe(true);

  if (NORMALIZED_BASE_PATH !== '') {
    const unprefixedAssetURL = new URL(
      moduleURL.pathname.slice(NORMALIZED_BASE_PATH.length),
      page.url()
    );
    expect((await request.get(unprefixedAssetURL.toString())).status()).toBe(404);
  }

  await page.getByRole('button', { name: 'Open user menu' }).click();
  await page.getByRole('menuitem', { name: 'Log out' }).click();
  await expect(page.getByRole('heading', { name: 'Welcome back' })).toBeVisible();
});
