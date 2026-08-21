import type { ApolloClient } from '@apollo/client';
import { ApolloProvider } from '@apollo/client/react';
import { RouterProvider, useRouterState } from '@tanstack/react-router';
import { Suspense, useMemo, useRef } from 'react';
import { router } from '@/app/router';
import { DashboardBoot } from '@/app/ui/dashboard-boot';
import { AuthProvider, useAuth } from '@/features/auth/model/auth-context';
import { createApolloClient } from '@/shared/api/apollo';

const InnerApp = ({ apolloClient }: { apolloClient: ApolloClient }): React.ReactNode => {
  const auth = useAuth();
  const hasResolvedRoute = useRouterState({
    router,
    select: (state) => state.resolvedLocation !== undefined,
  });
  if (auth.isLoading) return <DashboardBoot />;
  return (
    <Suspense fallback={<DashboardBoot />}>
      <RouterProvider router={router} context={{ apolloClient, auth }} />
      {hasResolvedRoute ? null : <DashboardBoot />}
    </Suspense>
  );
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
        <InnerApp apolloClient={apolloClient} />
      </AuthProvider>
    </ApolloProvider>
  );
};
