import type { AuthContextValue } from '@/features/auth/model/auth-context';

export type AuthRouteKind = 'protected' | 'login' | 'register';
export type AuthRedirect = '/' | '/login' | '/register' | null;

type AuthRoutingState = Pick<
  AuthContextValue,
  'authMode' | 'bootstrapError' | 'isAuthenticated' | 'isLoading' | 'unauthenticatedRoute'
>;

export function getAuthRedirect(kind: AuthRouteKind, auth: AuthRoutingState): AuthRedirect {
  if (auth.isLoading || auth.bootstrapError !== null) {
    return null;
  }
  if (auth.isAuthenticated) {
    return kind === 'protected' ? null : '/';
  }
  if (kind === 'protected') {
    return auth.unauthenticatedRoute;
  }
  if (kind === 'login' && auth.authMode === 'register-only') {
    return '/register';
  }
  if (kind === 'register' && auth.authMode === 'login-only') {
    return '/login';
  }
  return null;
}
