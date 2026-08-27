import { expect, type Page, test } from '@playwright/test';
import { createSite, signInAsAdmin } from './helpers/admin';

const readThemeTokens = async (page: Page): Promise<Record<string, string>> =>
  page.evaluate(() => {
    const styles = getComputedStyle(document.documentElement);
    return Object.fromEntries(
      [
        '--background',
        '--card',
        '--primary',
        '--chart-1',
        '--chart-2',
        '--chart-3',
        '--destructive-foreground',
        '--radius',
        '--theme-font-sans',
        '--theme-font-mono',
      ].map((token) => [token, styles.getPropertyValue(token).trim()])
    );
  });

test('Zen Inspired semantic tokens switch between light and dark without external fonts', async ({
  page,
}) => {
  await page.emulateMedia({ colorScheme: 'light' });
  await signInAsAdmin(page);

  await expect
    .poll(() => readThemeTokens(page))
    .toMatchObject({
      '--background': '#e9e4d8',
      '--card': '#f4efe4',
      '--primary': '#2e2e2e',
      '--chart-1': '#c2410c',
      '--chart-2': '#1d4ed8',
      '--chart-3': '#047857',
      '--destructive-foreground': '#fff',
      '--radius': '.5rem',
      '--theme-font-sans':
        'ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif',
      '--theme-font-mono':
        'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace',
    });

  await page.getByRole('button', { name: 'Toggle theme' }).click();
  await expect(page.locator('html')).toHaveClass(/dark/u);
  await expect
    .poll(() => readThemeTokens(page))
    .toMatchObject({
      '--background': '#141414',
      '--card': '#1c1c1c',
      '--primary': '#d1cfc0',
      '--chart-1': '#fb923c',
      '--chart-2': '#60a5fa',
      '--chart-3': '#34d399',
      '--destructive-foreground': '#fff',
    });
});

test('analytics chart labels and progress indicators use subdued dark-theme colors', async ({
  page,
}) => {
  await page.emulateMedia({ colorScheme: 'dark' });
  await signInAsAdmin(page);
  const analyticsURL = await createSite(page, 'Dark Analytics Theme Site', [
    'dark-analytics-theme.example',
  ]);
  const settingsURL = analyticsURL.replace(/\/analytics(?:\?.*)?$/u, '/settings');

  await page.goto(settingsURL);
  const siteKey = await page.getByRole('textbox', { name: 'Site Key', exact: true }).inputValue();
  const collectResponse = await page.request.post(
    `${new URL(settingsURL).origin}/api/collect?site_key=${encodeURIComponent(siteKey)}`,
    {
      data: JSON.stringify({ path: '/dark-theme' }),
      headers: {
        'Content-Type': 'text/plain;charset=UTF-8',
        Origin: 'https://dark-analytics-theme.example',
      },
    }
  );
  expect(collectResponse.ok()).toBe(true);

  await page.route('**/graphql', async (route) => {
    const payload: unknown = route.request().postDataJSON();
    if (
      typeof payload === 'object' &&
      payload !== null &&
      'operationName' in payload &&
      payload.operationName === 'ChartData'
    ) {
      await route.fulfill({
        contentType: 'application/json',
        body: JSON.stringify({
          data: {
            dashboard: {
              __typename: 'DashboardStats',
              dailyStats: [
                {
                  __typename: 'DailyStats',
                  date: '2026-08-26',
                  visitors: 1,
                  pageViews: 1,
                  sessions: 1,
                },
                {
                  __typename: 'DailyStats',
                  date: '2026-08-27',
                  visitors: 1,
                  pageViews: 1,
                  sessions: 1,
                },
              ],
            },
          },
        }),
      });
      return;
    }
    await route.continue();
  });

  await page.goto(analyticsURL);
  const xAxisLabel = page.getByText('Aug 26', { exact: true });
  await expect(xAxisLabel).toBeVisible();
  await expect(page.locator('[data-slot="progress-indicator"]').first()).toBeVisible();

  const xAxisLabelFill = await xAxisLabel.evaluate((element) => getComputedStyle(element).fill);
  const progressColor = await page
    .locator('[data-slot="progress-indicator"]')
    .first()
    .evaluate((element) => getComputedStyle(element).backgroundColor);

  expect(xAxisLabelFill).toBe('rgb(142, 138, 131)');
  expect(progressColor).toBe('rgb(116, 134, 158)');
});
