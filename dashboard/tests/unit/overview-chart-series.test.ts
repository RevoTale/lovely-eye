import assert from 'node:assert/strict';
import { describe, it } from 'node:test';
import { OVERVIEW_SERIES } from '../../src/features/analytics/ui/overview-chart/overview-chart-series';

describe('overview chart series', () => {
  it('differentiates every series by both color and line pattern', () => {
    assert.equal(new Set(OVERVIEW_SERIES.map((series) => series.stroke)).size, 3);
    assert.equal(new Set(OVERVIEW_SERIES.map((series) => series.strokeDasharray)).size, 3);
  });
});
