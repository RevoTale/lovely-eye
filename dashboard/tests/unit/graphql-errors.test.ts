import assert from 'node:assert/strict';
import { describe, it } from 'node:test';
import { getGraphQLErrorCode } from '../../src/shared/api/errors';

describe('GraphQL error categories', () => {
  it('accepts only documented machine-readable codes', () => {
    assert.equal(getGraphQLErrorCode({ code: 'UNAUTHENTICATED' }), 'UNAUTHENTICATED');
    assert.equal(getGraphQLErrorCode({ code: 'FORBIDDEN' }), 'FORBIDDEN');
    assert.equal(getGraphQLErrorCode({ code: 'SOMETHING_NEW' }), null);
    assert.equal(getGraphQLErrorCode(null), null);
  });
});
