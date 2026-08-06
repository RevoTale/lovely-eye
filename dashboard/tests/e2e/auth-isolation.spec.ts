import { expect, type Page, test } from '@playwright/test';

const { TEST_APP_PORT = '4173' } = process.env;
const APP_ORIGIN = `http://127.0.0.1:${TEST_APP_PORT}`;

const signIn = async (page: Page, loginURL: string): Promise<void> => {
  await page.goto(loginURL);
  await expect(page.getByRole('heading', { name: 'Welcome back' })).toBeVisible();
  await page.getByLabel('Username').fill('e2e-admin');
  await page.getByLabel('Password').fill('e2e-password');
  await page.getByRole('button', { name: 'Sign in' }).click();
  await expect(page.getByRole('button', { name: 'Open user menu' })).toBeVisible();
};

test('same-origin instances retain independent authentication sessions', async ({
  context,
  page,
}) => {
  const secondPage = await context.newPage();

  await signIn(page, `${APP_ORIGIN}/instance-a/login`);
  await signIn(secondPage, `${APP_ORIGIN}/instance-b/login`);

  const authCookies = (await context.cookies()).filter(({ name }) => name.startsWith('le_'));
  expect(authCookies.filter(({ path }) => path === '/instance-a')).toHaveLength(2);
  expect(authCookies.filter(({ path }) => path === '/instance-b')).toHaveLength(2);

  await page.getByRole('button', { name: 'Open user menu' }).click();
  await page.getByRole('menuitem', { name: 'Log out' }).click();
  await expect(page.getByRole('heading', { name: 'Welcome back' })).toBeVisible();

  await secondPage.reload();
  await expect(secondPage.getByRole('button', { name: 'Open user menu' })).toBeVisible();
  const remainingAuthCookies = (await context.cookies()).filter(({ name }) =>
    name.startsWith('le_')
  );
  expect(remainingAuthCookies.filter(({ path }) => path === '/instance-a')).toHaveLength(0);
  expect(remainingAuthCookies.filter(({ path }) => path === '/instance-b')).toHaveLength(2);
});
