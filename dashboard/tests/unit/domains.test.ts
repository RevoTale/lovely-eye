import assert from 'node:assert/strict';
import { describe, it } from 'node:test';
import { normalizeDomains, validateDomains } from '../../src/shared/lib/domains';

describe('domain normalization', () => {
  it('normalizes aliases and removes duplicates without losing order', () => {
    assert.deepEqual(
      normalizeDomains(['HTTPS://WWW.Example.com/path', 'app.example.com', 'example.com']),
      ['example.com', 'app.example.com']
    );
  });

  it('requires at least one valid domain', () => {
    assert.equal(validateDomains(['']).error, 'At least one domain is required');
    assert.equal(
      validateDomains(['not a domain']).error,
      'Please enter valid domains (e.g., example.com)'
    );
    assert.deepEqual(validateDomains(['example.com']).domains, ['example.com']);
  });
});
