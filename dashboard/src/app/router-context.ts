import type { ApolloClient } from '@apollo/client';
import type { AuthContextValue } from '@/features/auth/model/auth-context';

export interface RouterContext {
  auth: AuthContextValue;
  apolloClient: ApolloClient | null;
}

const noopAuthAction = async (): Promise<void> => {
  await Promise.resolve();
};

export function createInitialContext(): RouterContext {
  return {
    apolloClient: null,
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
