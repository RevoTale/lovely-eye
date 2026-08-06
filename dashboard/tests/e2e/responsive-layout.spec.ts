import { expect, type Page, test } from '@playwright/test';
import { createSite, signInAsAdmin } from './helpers/admin';

const MOBILE_VIEWPORT = { width: 320, height: 800 } as const;
const TABLET_VIEWPORT = { width: 768, height: 900 } as const;
const DESKTOP_VIEWPORT = { width: 1024, height: 900 } as const;
const WIDE_VIEWPORT = { width: 1440, height: 1000 } as const;

const expectNoHorizontalOverflow = async (page: Page): Promise<void> => {
  const dimensions = await page.evaluate(() => ({
    viewportWidth: document.documentElement.clientWidth,
    contentWidth: document.documentElement.scrollWidth,
  }));
  expect(dimensions.contentWidth).toBeLessThanOrEqual(dimensions.viewportWidth);
};

test('Add New Site keeps mobile gutters and a centered wide layout', async ({ page }) => {
  await page.setViewportSize(MOBILE_VIEWPORT);
  await signInAsAdmin(page);
  await page.goto('sites/new');

  const headingBox = await page.getByRole('heading', { name: 'Add New Site' }).boundingBox();
  const mobileFormBox = await page.locator('main form').boundingBox();
  expect(headingBox).not.toBeNull();
  expect(mobileFormBox).not.toBeNull();
  expect(headingBox?.x).toBeGreaterThanOrEqual(16);
  expect((mobileFormBox?.x ?? 0) + (mobileFormBox?.width ?? 0)).toBeLessThanOrEqual(
    MOBILE_VIEWPORT.width - 16
  );
  await expectNoHorizontalOverflow(page);

  await page.setViewportSize(WIDE_VIEWPORT);
  const formBox = await page.locator('main form').boundingBox();
  expect(formBox).not.toBeNull();
  const leftGutter = formBox?.x ?? 0;
  const rightGutter = WIDE_VIEWPORT.width - leftGutter - (formBox?.width ?? 0);
  expect(Math.abs(leftGutter - rightGutter)).toBeLessThanOrEqual(1);
});

test('critical admin compositions fit mobile and wide viewports', async ({ page }) => {
  test.setTimeout(60_000);
  await page.setViewportSize(MOBILE_VIEWPORT);
  await signInAsAdmin(page);
  const analyticsURL = await createSite(page, 'Responsive Contract Site', [
    'primary.responsive-contract.example',
    'secondary.responsive-contract.example',
  ]);
  const settingsURL = analyticsURL.replace(/\/analytics(?:\?.*)?$/u, '/settings');
  await page.goto(settingsURL);
  const siteKey = await page.getByRole('textbox', { name: 'Site Key', exact: true }).inputValue();
  const collectResponse = await page.request.post(
    `${new URL(settingsURL).origin}/api/collect?site_key=${encodeURIComponent(siteKey)}`,
    {
      data: JSON.stringify({ path: '/responsive-contract' }),
      headers: {
        'Content-Type': 'text/plain;charset=UTF-8',
        Origin: 'https://primary.responsive-contract.example',
      },
    }
  );
  expect(collectResponse.ok()).toBe(true);

  const routes = ['sites', 'sites/new', analyticsURL, settingsURL];
  for (const viewport of [MOBILE_VIEWPORT, TABLET_VIEWPORT, DESKTOP_VIEWPORT, WIDE_VIEWPORT]) {
    await page.setViewportSize(viewport);
    for (const route of routes) {
      await page.goto(route);
      if (route === analyticsURL) {
        await expect(page.getByText('Total Visitors', { exact: true })).toBeVisible();
      }
      await expectNoHorizontalOverflow(page);
    }
  }

  await page.setViewportSize(MOBILE_VIEWPORT);
  await page.goto(analyticsURL);
  await expect(page.getByText('Total Visitors', { exact: true })).toBeVisible();
  await page.getByRole('tab', { name: 'Custom' }).click();
  await expectNoHorizontalOverflow(page);
  await page.locator('[data-slot="popover-trigger"]').first().click();
  await expectNoHorizontalOverflow(page);

  await page.goto(settingsURL);
  await page.getByRole('button', { name: 'Regenerate site key' }).click();
  await expectNoHorizontalOverflow(page);
  await page.getByRole('button', { name: 'New event name' }).click();
  await page.getByRole('button', { name: 'Add field' }).click();
  await expectNoHorizontalOverflow(page);
  await page.getByRole('button', { name: 'Show snippet' }).click();
  await expectNoHorizontalOverflow(page);
});
