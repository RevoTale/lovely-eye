import assert from 'node:assert/strict';
import { describe, it } from 'node:test';
import { resolveAnalyticsDateRange } from '../../src/features/analytics/model/analytics-date-range';
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
      pagePathContains: '  /blog  ',
      preset: 'invalid',
      from: 'not-a-date',
    });

    assert.equal(parsed.eventsPage, 2);
    assert.deepEqual(parsed.referrer, ['example.com']);
    assert.deepEqual(parsed.browser, ['Firefox', 'Safari']);
    assert.equal(parsed.pagePathContains, '/blog');
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

  it('resolves URL-owned ranges identically for loaders and screens', () => {
    const reference = new Date('2026-08-21T12:00:00Z');
    const custom = resolveAnalyticsDateRange(
      {
        preset: 'custom',
        from: '2026-08-01',
        to: '2026-08-02',
        fromTime: '10:15',
        toTime: '11:45',
      },
      reference
    );

    assert.deepEqual(custom, {
      from: new Date('2026-08-01T10:15:00').toISOString(),
      to: new Date('2026-08-02T11:45:00').toISOString(),
    });
    assert.equal(resolveAnalyticsDateRange({ preset: 'all' }, reference), undefined);
  });
});
