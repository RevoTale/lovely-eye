import assert from 'node:assert/strict';
import { describe, test } from 'node:test';
import { type AuthRouteKind, getAuthRedirect } from '@/features/auth/model/auth-routing';

interface RoutingCase {
  name: string;
  kind: AuthRouteKind;
  auth: Parameters<typeof getAuthRedirect>[1];
  expected: ReturnType<typeof getAuthRedirect>;
}

const baseAuth: Parameters<typeof getAuthRedirect>[1] = {
  authMode: 'login-only',
  bootstrapError: null,
  isAuthenticated: false,
  isLoading: false,
  unauthenticatedRoute: '/login',
};

describe('getAuthRedirect', () => {
  const cases: RoutingCase[] = [
    {
      name: 'waits while authentication is unresolved',
      kind: 'protected',
      auth: { ...baseAuth, isLoading: true },
      expected: null,
    },
    {
      name: 'keeps bootstrap errors visible',
      kind: 'protected',
      auth: { ...baseAuth, bootstrapError: 'unavailable' },
      expected: null,
    },
    {
      name: 'sends an authenticated user away from guest screens',
      kind: 'login',
      auth: { ...baseAuth, isAuthenticated: true },
      expected: '/',
    },
    {
      name: 'uses the selected unauthenticated route for protected content',
      kind: 'protected',
      auth: { ...baseAuth, authMode: 'register-only', unauthenticatedRoute: '/register' },
      expected: '/register',
    },
    {
      name: 'requires initial registration when no user exists',
      kind: 'login',
      auth: { ...baseAuth, authMode: 'register-only', unauthenticatedRoute: '/register' },
      expected: '/register',
    },
    {
      name: 'disables registration when configuration is login-only',
      kind: 'register',
      auth: baseAuth,
      expected: '/login',
    },
  ];

  for (const routingCase of cases) {
    test(routingCase.name, () => {
      assert.equal(getAuthRedirect(routingCase.kind, routingCase.auth), routingCase.expected);
    });
  }
});
