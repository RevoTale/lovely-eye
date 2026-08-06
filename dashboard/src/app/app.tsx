import { ApolloProvider } from '@apollo/client/react';
import { RouterProvider } from '@tanstack/react-router';
import { useMemo, useRef } from 'react';
import { router } from '@/app/router';
import { AuthProvider, useAuth } from '@/features/auth/model/auth-context';
import { AuthShell } from '@/features/auth/ui/auth-shell';
import { createApolloClient } from '@/shared/api/apollo';

const InnerApp = (): React.ReactNode => {
  const auth = useAuth();
  if (auth.isLoading) {
    return (
      <AuthShell title='Loading dashboard' description='Checking authentication status.'>
        <p className='text-center text-sm text-muted-foreground'>Please wait...</p>
      </AuthShell>
    );
  }
  return <RouterProvider router={router} context={{ auth }} />;
};

export const App = (): React.ReactNode => {
  const authErrorHandlerRef = useRef<(() => void) | null>(null);
  const apolloClient = useMemo(
    () =>
      createApolloClient(() => {
        authErrorHandlerRef.current?.();
      }),
    []
  );

  return (
    <ApolloProvider client={apolloClient}>
      <AuthProvider authErrorHandlerRef={authErrorHandlerRef}>
        <InnerApp />
      </AuthProvider>
    </ApolloProvider>
  );
};
