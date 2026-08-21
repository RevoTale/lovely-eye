import { expect, type Page, test } from '@playwright/test';
import { createSite } from './helpers/admin';

const recordGraphQLOperations = (
  page: Page
): { counts: Map<string, number>; variables: Map<string, string[]> } => {
  const operationCounts = new Map<string, number>();
  const operationVariables = new Map<string, string[]>();
  page.on('request', (request) => {
    if (!request.url().endsWith('/graphql')) return;
    const rawBody = request.postData();
    if (rawBody === null) return;
    try {
      const payload: unknown = JSON.parse(rawBody);
      if (
        typeof payload === 'object' &&
        payload !== null &&
        'operationName' in payload &&
        typeof payload.operationName === 'string'
      ) {
        const variables =
          'variables' in payload && payload.variables !== undefined
            ? JSON.stringify(payload.variables)
            : '{}';
        operationCounts.set(
          payload.operationName,
          (operationCounts.get(payload.operationName) ?? 0) + 1
        );
        operationVariables.set(payload.operationName, [
          ...(operationVariables.get(payload.operationName) ?? []),
          variables,
        ]);
      }
    } catch {
      return;
    }
  });
  return { counts: operationCounts, variables: operationVariables };
};

test('records critical browser readiness and GraphQL fan-out', async ({ context, page }) => {
  await page.goto('login');
  await page.getByLabel('Username').fill('e2e-admin');
  await page.getByLabel('Password').fill('e2e-password');
  const loginStartedAt = performance.now();
  await page.getByRole('button', { name: 'Sign in' }).click();
  await expect(page.getByRole('button', { name: 'Open user menu' })).toBeVisible();
  const loginReadyMilliseconds = performance.now() - loginStartedAt;

  await createSite(page, 'Performance Baseline', ['performance.example']);
  const dashboardURL = page.url();
  const dashboardPage = await context.newPage();
  const { counts: operationCounts, variables: operationVariables } =
    recordGraphQLOperations(dashboardPage);
  const dashboardStartedAt = performance.now();
  await dashboardPage.goto(dashboardURL);
  await expect(dashboardPage.getByText('Total Visitors', { exact: true })).toBeVisible();
  const dashboardReadyMilliseconds = performance.now() - dashboardStartedAt;
  const navigation = await dashboardPage.evaluate(() => {
    const [entry] = performance.getEntriesByType('navigation');
    if (!(entry instanceof PerformanceNavigationTiming)) return null;
    return {
      responseEndMilliseconds: entry.responseEnd,
      domContentLoadedMilliseconds: entry.domContentLoadedEventEnd,
      loadEventMilliseconds: entry.loadEventEnd,
    };
  });

  const queryCount = [...operationCounts.values()].reduce((total, count) => total + count, 0);
  expect(queryCount).toBeGreaterThan(0);
  expect(queryCount).toBeLessThanOrEqual(8);
  process.stdout.write(
    `PERFORMANCE_BASELINE ${JSON.stringify({
      loginReadyMilliseconds: Math.round(loginReadyMilliseconds),
      dashboardReadyMilliseconds: Math.round(dashboardReadyMilliseconds),
      navigation,
      graphQLOperations: Object.fromEntries(
        [...operationCounts].sort(([left], [right]) => left.localeCompare(right))
      ),
      graphQLVariables: Object.fromEntries(
        [...operationVariables].sort(([left], [right]) => left.localeCompare(right))
      ),
      graphQLQueryCount: queryCount,
    })}\n`
  );
});
