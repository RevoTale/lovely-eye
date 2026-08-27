import { expect, test } from '@playwright/test';
import { createSite, signInAsAdmin } from './helpers/admin';
import { installGraphQLOperationController } from './helpers/graphql-gate';

test('auth, first load, mutation feedback, and retained refresh state stay coherent', async ({
  context,
  page,
}) => {
  await page.addInitScript(() => {
    const startViewTransition = document.startViewTransition?.bind(document);
    if (startViewTransition === undefined) return;

    Reflect.set(globalThis, '__lovelyEyeViewTransitionCount', 0);
    document.startViewTransition = (update) => {
      const count = Reflect.get(globalThis, '__lovelyEyeViewTransitionCount');
      Reflect.set(
        globalThis,
        '__lovelyEyeViewTransitionCount',
        typeof count === 'number' ? count + 1 : 1
      );
      return startViewTransition(update);
    };
  });
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
  await expect(page.getByRole('heading', { name: 'Welcome back' })).toBeVisible();
  await expect(page.getByRole('button', { name: 'Signing in...' })).toBeDisabled();
  await expect(page.locator('.animate-pulse')).toHaveCount(0);
  sitesLoad.release();
  await expect(page.getByRole('heading', { name: 'Add New Site' })).toBeVisible();
  const viewTransitionCount = await page.evaluate(() =>
    Reflect.get(globalThis, '__lovelyEyeViewTransitionCount')
  );
  expect(viewTransitionCount).toBe(1);

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
  await freshPage.addInitScript(() => {
    let score = 0;
    const observer = new PerformanceObserver((list) => {
      for (const entry of list.getEntries()) {
        const shift = entry as PerformanceEntry & { hadRecentInput?: boolean; value?: number };
        if (shift.hadRecentInput !== true) score += shift.value ?? 0;
      }
    });
    observer.observe({ type: 'layout-shift', buffered: true });
    Reflect.set(globalThis, '__lovelyEyeLayoutShiftScore', () => score);
    Reflect.set(globalThis, '__lovelyEyeLayoutShiftObserver', observer);
  });
  await freshPage.emulateMedia({ reducedMotion: 'reduce' });
  const freshOperations = await installGraphQLOperationController(freshPage);
  const firstSiteLoad = freshOperations.blockNext('Site');
  await freshPage.goto(siteURL);
  await firstSiteLoad.seen;
  await expect(freshPage.getByRole('heading', { name: 'Loading dashboard' })).toBeVisible();
  await expect(freshPage.locator('.animate-pulse')).toHaveCount(0);
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
  const reducedMotionDuration = await freshPage
    .locator('.animate-spin')
    .first()
    .evaluate((element) => Number.parseFloat(getComputedStyle(element).animationDuration));
  expect(reducedMotionDuration).toBeLessThanOrEqual(0.01);
  dashboardRefresh.release();
  await expect(freshPage).toHaveURL(/preset=7d/u);
  const layoutShiftScore = await freshPage.evaluate(() => {
    const readScore: unknown = Reflect.get(globalThis, '__lovelyEyeLayoutShiftScore');
    return typeof readScore === 'function' ? (readScore as () => number)() : 0;
  });
  expect(layoutShiftScore).toBeLessThan(0.1);
});

test('site switching and URL-owned analytics state survive history and refresh', async ({
  page,
}) => {
  const operations = await installGraphQLOperationController(page);
  await signInAsAdmin(page);
  await createSite(page, 'Switch Alpha', ['alpha.example']);
  const betaURL = await createSite(page, 'Switch Beta', ['beta.example', 'app.beta.example']);

  await page.goto('./');
  await expect(page.getByRole('heading', { name: 'Switch Beta' })).toBeVisible();
  await page.goto('sites');
  const alphaLoad = operations.blockNext('Site');
  await page.getByRole('link', { name: /Switch Alpha/u }).click();
  await alphaLoad.seen;
  await expect(page.getByRole('heading', { name: 'Sites' })).toBeVisible();
  await expect(page.locator('.animate-pulse')).toHaveCount(0);
  alphaLoad.release();
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
  await expect(page.getByRole('tab', { name: '7d' })).toHaveAttribute('aria-selected', 'true');
  await page.goForward();
  await expect(page).toHaveURL(/preset=30d/u);
  await page.reload();
  await expect(page).toHaveURL(`${betaURL.replace(/\?.*$/u, '')}?preset=30d`);
  await expect(page.getByRole('tab', { name: '30d' })).toHaveAttribute('aria-selected', 'true');
});

test('top pages path search is explicit, URL-owned, and removable', async ({ page }) => {
  await signInAsAdmin(page);
  await createSite(page, 'Path Search Site', ['path-search.example']);

  const pathInput = page.getByRole('searchbox', { name: 'Page path contains' });
  await pathInput.fill('  /blog  ');
  await page.getByRole('button', { name: 'Search' }).click();

  await expect(page).toHaveURL(/pagePathContains=%2Fblog/u);
  await expect(pathInput).toHaveValue('/blog');
  const activeFilter = page.getByRole('link', { name: /Page contains: \/blog/u });
  await expect(activeFilter).toBeVisible();

  const dashboardRequest = page.waitForRequest((request) => {
    if (!request.url().endsWith('/graphql')) return false;
    const payload: unknown = request.postDataJSON();
    return (
      typeof payload === 'object' &&
      payload !== null &&
      'operationName' in payload &&
      payload.operationName === 'Dashboard'
    );
  });
  await page.reload();
  const dashboardPayload: unknown = (await dashboardRequest).postDataJSON();
  expect(dashboardPayload).toMatchObject({
    variables: { filter: { pagePathContains: '/blog' } },
  });

  await activeFilter.click();
  await expect(page).not.toHaveURL(/pagePathContains/u);
  await expect(pathInput).toHaveValue('');
});

test('card pagination preserves the viewport position', async ({ page }) => {
  await page.setViewportSize({ width: 1024, height: 900 });
  await signInAsAdmin(page);
  const analyticsURL = await createSite(page, 'Pagination Scroll Site', [
    'pagination-scroll.example',
  ]);
  const settingsURL = analyticsURL.replace(/\/analytics(?:\?.*)?$/u, '/settings');

  await page.goto(settingsURL);
  const siteKey = await page.getByRole('textbox', { name: 'Site Key', exact: true }).inputValue();
  for (let index = 0; index < 6; index += 1) {
    const collectResponse = await page.request.post(
      `${new URL(settingsURL).origin}/api/collect?site_key=${encodeURIComponent(siteKey)}`,
      {
        data: JSON.stringify({ path: `/pagination-${index}` }),
        headers: {
          'Content-Type': 'text/plain;charset=UTF-8',
          Origin: 'https://pagination-scroll.example',
        },
      }
    );
    expect(collectResponse.ok()).toBe(true);
  }

  await page.goto(analyticsURL);
  const recentEventsCard = page
    .locator('[data-slot="card"]')
    .filter({ hasText: 'Recent Events' })
    .first();
  await expect(recentEventsCard.getByText('Page 1 of 2', { exact: true })).toBeVisible();
  const nextPageButton = recentEventsCard.getByRole('button', { name: 'Next' });
  const initialCardBox = await recentEventsCard.boundingBox();
  expect(initialCardBox).not.toBeNull();
  await page.evaluate(
    (scrollDelta) => window.scrollBy(0, scrollDelta),
    (initialCardBox?.y ?? 0) - 100
  );
  const visibleStartCardBox = await recentEventsCard.boundingBox();
  expect(visibleStartCardBox?.y).toBeCloseTo(100, 0);

  const scrollBefore = await page.evaluate(() => window.scrollY);
  await nextPageButton.click();
  await expect(page).toHaveURL(/eventsPage=2/u);
  await expect(recentEventsCard.getByText('Page 2 of 2', { exact: true })).toBeVisible();
  const scrollAfter = await page.evaluate(() => window.scrollY);

  expect(Math.abs(scrollAfter - scrollBefore)).toBeLessThanOrEqual(1);

  await recentEventsCard.getByRole('button', { name: 'Prev' }).click();
  await expect(recentEventsCard.getByText('Page 1 of 2', { exact: true })).toBeVisible();
  const cardBox = await recentEventsCard.boundingBox();
  expect(cardBox).not.toBeNull();
  await page.evaluate((scrollDelta) => window.scrollBy(0, scrollDelta), (cardBox?.y ?? 0) + 100);
  const positionedCardBox = await recentEventsCard.boundingBox();
  expect(positionedCardBox?.y).toBeCloseTo(-100, 0);
  expect((positionedCardBox?.y ?? 0) + (positionedCardBox?.height ?? 0)).toBeGreaterThan(0);

  await nextPageButton.click();
  await expect(recentEventsCard.getByText('Page 2 of 2', { exact: true })).toBeVisible();
  await expect
    .poll(async () => (await recentEventsCard.boundingBox())?.height ?? Number.POSITIVE_INFINITY)
    .toBeLessThan(500);
  const repositionedCardBox = await recentEventsCard.boundingBox();

  expect(repositionedCardBox?.y).toBeCloseTo(0, 0);
});
