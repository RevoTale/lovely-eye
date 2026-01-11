import { ApolloProvider } from '@apollo/client';
import { RouterProvider } from '@tanstack/react-router';
import { apolloClient } from '@/lib/apollo';
import { AuthProvider, useAuth } from '@/hooks';
import { router } from '@/router';

function InnerApp(): React.JSX.Element {
  const auth = useAuth();
  return <RouterProvider router={router} context={{ auth }} />;
}

export function App(): React.JSX.Element {
  return (
    <ApolloProvider client={apolloClient}>
      <AuthProvider>
        <InnerApp />
      </AuthProvider>
    </ApolloProvider>
  );
}