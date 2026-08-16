import assert from 'node:assert/strict';
import { describe, it } from 'node:test';
import { selectInitialSiteId } from '../../src/features/sites/model/site-selection';

describe('selectInitialSiteId', () => {
  it('returns the remembered site only while it remains accessible', () => {
    const sites = [{ id: 'one' }, { id: 'two' }];
    assert.equal(selectInitialSiteId(sites, 'two'), 'two');
    assert.equal(selectInitialSiteId(sites, 'removed'), null);
  });

  it('opens the only accessible site without stored state', () => {
    assert.equal(selectInitialSiteId([{ id: 'only' }], null), 'only');
  });

  it('requires an explicit choice when several sites have no remembered selection', () => {
    assert.equal(selectInitialSiteId([{ id: 'one' }, { id: 'two' }], null), null);
  });
});
