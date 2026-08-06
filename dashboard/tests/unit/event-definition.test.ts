import assert from 'node:assert/strict';
import { describe, it } from 'node:test';
import { validateEventDefinition } from '../../src/features/event-definitions/ui/editor/use-event-definition-editor';

const field = { key: 'message', type: 'STRING' as const, required: true, maxLength: 100 };

describe('event definition validation', () => {
  it('requires an event name and non-empty unique field keys', () => {
    assert.equal(validateEventDefinition('', []), 'Event name is required.');
    assert.equal(
      validateEventDefinition('signup_error', [{ ...field, key: '' }]),
      'Field key cannot be empty.'
    );
    assert.equal(
      validateEventDefinition('signup_error', [field, { ...field }]),
      'Field keys must be unique.'
    );
    assert.equal(validateEventDefinition('signup_error', [field]), '');
  });
});
