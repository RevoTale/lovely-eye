import assert from 'node:assert/strict';
import { describe, it } from 'node:test';
import { getSafeReferrerHref } from '../../src/features/analytics/model/referrer-utils';

describe('referrer links', () => {
  it('allows only absolute HTTP(S) URLs', () => {
    assert.equal(getSafeReferrerHref('https://example.com/path'), 'https://example.com/path');
    assert.equal(getSafeReferrerHref('http://example.com'), 'http://example.com/');
    assert.equal(getSafeReferrerHref('javascript:alert(1)'), null);
    assert.equal(getSafeReferrerHref('data:text/html,unsafe'), null);
    assert.equal(getSafeReferrerHref('//example.com/path'), null);
    assert.equal(getSafeReferrerHref('not a URL'), null);
  });
});
