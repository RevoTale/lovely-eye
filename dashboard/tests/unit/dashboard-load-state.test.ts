import assert from 'node:assert/strict';
import { describe, it } from 'node:test';
import { resolveDashboardLoadState } from '../../src/features/analytics/model/dashboard-load-state';

describe('resolveDashboardLoadState', () => {
  it('keeps previous data visible during a background refresh', () => {
    const previousData = { visitors: 42 };

    assert.deepEqual(resolveDashboardLoadState(undefined, previousData, true), {
      data: previousData,
      state: 'refreshing',
    });
  });

  it('uses an initial state only when no displayable data exists', () => {
    assert.deepEqual(resolveDashboardLoadState(undefined, undefined, true), {
      data: undefined,
      state: 'initial',
    });
  });
});
