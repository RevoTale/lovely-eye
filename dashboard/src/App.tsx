import { ApolloProvider } from '@apollo/client';
import { RouterProvider } from '@tanstack/react-router';
import { useMemo, useRef } from 'react';
import { createApolloClient } from '@/lib/apollo';
import { AuthProvider, useAuth } from '@/hooks';
import { router } from '@/router';

function InnerApp(): React.JSX.Element {
  const auth = useAuth();
  return <RouterProvider router={router} context={{ auth }} />;
}

export function App(): React.JSX.Element {
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
