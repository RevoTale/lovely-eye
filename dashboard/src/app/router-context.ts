import type { AuthContextValue } from '@/features/auth/model/auth-context';

export interface RouterContext {
  auth: AuthContextValue;
}

const noopAuthAction = async (): Promise<void> => {
  await Promise.resolve();
};

export function createInitialContext(): RouterContext {
  return {
    auth: {
      user: null,
      isLoading: true,
      isAuthenticated: false,
      authMode: 'login-only',
      bootstrapError: null,
      unauthenticatedRoute: '/login',
      login: noopAuthAction,
      register: noopAuthAction,
      logout: noopAuthAction,
    },
  };
}
