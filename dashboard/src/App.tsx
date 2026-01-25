import { ApolloProvider } from '@apollo/client/react';
import { RouterProvider } from '@tanstack/react-router';
import { useRef, useMemo } from 'react';
import { createApolloClient } from '@/lib/apollo';
import { AuthProvider, useAuth } from '@/hooks';
import { router } from '@/router';

const InnerApp = (): React.ReactNode => {
  const auth = useAuth();
  return <RouterProvider router={router} context={{ auth }} />;
}

export const App = (): React.ReactNode => {
  const authErrorHandlerRef = useRef<(() => void) | null>(null);
  const apolloClient = useMemo(() => createApolloClient(() => {
    authErrorHandlerRef.current?.();
  }), []);

  return (
    <ApolloProvider client={apolloClient}>
      <AuthProvider authErrorHandlerRef={authErrorHandlerRef}>
        <InnerApp />
      </AuthProvider>
    </ApolloProvider>
  );
}
