import { defineConfig, devices } from '@playwright/test';

const APP_HOST = '127.0.0.1';
const {
  TEST_APP_PORT = '4173',
  TEST_BASE_PATH = '/',
  TEST_MULTI_INSTANCE = 'false',
  TEST_PERFORMANCE = 'false',
} = process.env;
const NORMALIZED_BASE_PATH =
  TEST_BASE_PATH === '/' ? '' : `/${TEST_BASE_PATH.split('/').filter(Boolean).join('/')}`;
const APP_ORIGIN = `http://${APP_HOST}:${TEST_APP_PORT}`;
const APP_URL = `${APP_ORIGIN}${NORMALIZED_BASE_PATH}/`;
const IS_CI = Object.hasOwn(process.env, 'CI');

export default defineConfig({
  testDir: './tests/e2e',
  testIgnore:
    TEST_MULTI_INSTANCE === 'true'
      ? /admin-smoke|runtime-base-path|admin-state|performance-baseline/u
      : TEST_PERFORMANCE === 'true'
        ? /admin-smoke|runtime-base-path|admin-state|auth-isolation/u
        : /auth-isolation|performance-baseline/u,
  fullyParallel: false,
  workers: 1,
  forbidOnly: IS_CI,
  retries: IS_CI ? 2 : 0,
  reporter: IS_CI ? 'github' : 'list',
  use: {
    baseURL: APP_URL,
    trace: 'retain-on-failure',
  },
  webServer: {
    command: './scripts/start-test-app.sh',
    url: `${APP_ORIGIN}/health`,
    reuseExistingServer: false,
    timeout: 120_000,
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
});
