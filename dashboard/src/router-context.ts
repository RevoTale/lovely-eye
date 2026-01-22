import type { AuthContextType } from '@/hooks/use-auth';

export interface RouterContext {
  auth: AuthContextType;
}

const noopAuthAction = async (): Promise<void> => {
  await Promise.resolve();
};

// Helper to create initial router context (overridden at runtime via RouterProvider)
export function createInitialContext(): RouterContext {
  return {
    auth: {
      user: null,
      isLoading: true,
      isAuthenticated: false,
      login: noopAuthAction,
      register: noopAuthAction,
      logout: noopAuthAction,
    },
  };
}
