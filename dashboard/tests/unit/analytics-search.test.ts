import assert from 'node:assert/strict';
import { describe, it } from 'node:test';
import {
  addAnalyticsFilter,
  analyticsSearchSchema,
  clearAnalyticsFilters,
  clearPagination,
  removeAnalyticsFilter,
  setAnalyticsPage,
} from '../../src/features/analytics/model/analytics-search';

describe('analytics URL state', () => {
  it('returns final typed values and drops invalid input', () => {
    const parsed = analyticsSearchSchema.parse({
      eventsPage: '2.9',
      referrer: 'example.com',
      browser: ['Firefox', 'Safari'],
      preset: 'invalid',
      from: 'not-a-date',
    });

    assert.equal(parsed.eventsPage, 2);
    assert.deepEqual(parsed.referrer, ['example.com']);
    assert.deepEqual(parsed.browser, ['Firefox', 'Safari']);
    assert.equal(parsed.preset, undefined);
    assert.equal(parsed.from, undefined);
  });

  it('updates filters and pagination without discarding unrelated state', () => {
    const initial = { preset: '7d' as const, browser: ['Firefox'], eventsPage: 3 };
    const withFilter = addAnalyticsFilter(initial, 'browser', 'Safari');
    assert.deepEqual(withFilter.browser, ['Firefox', 'Safari']);
    assert.equal(setAnalyticsPage(withFilter, 'eventsPage', 1).eventsPage, undefined);
    assert.equal(clearPagination(withFilter).eventsPage, undefined);
    assert.deepEqual(removeAnalyticsFilter(withFilter, 'browser', 'Firefox').browser, ['Safari']);
    assert.equal(clearAnalyticsFilters(withFilter).browser, undefined);
    assert.equal(clearAnalyticsFilters(withFilter).preset, '7d');
  });
});
