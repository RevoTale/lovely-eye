import { expect, type Page, test } from '@playwright/test';
import { signInAsAdmin } from './helpers/admin';

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
