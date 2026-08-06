import { ApolloClient, ApolloLink, InMemoryCache } from '@apollo/client';
import { CombinedGraphQLErrors, ServerError } from '@apollo/client/errors';
import { ErrorLink } from '@apollo/client/link/error';
import { HttpLink } from '@apollo/client/link/http';
import { getGraphQLErrorCode } from '@/shared/api/errors';
import { getGraphQLUrl } from '@/shared/config/runtime';

const HTTP_UNAUTHORIZED = 401;

const httpLink = new HttpLink({
  uri: getGraphQLUrl(),
  credentials: 'include',
});

type AuthErrorHandler = () => void;

export function createApolloClient(onAuthError?: AuthErrorHandler): ApolloClient {
  const errorLink = new ErrorLink(({ error }) => {
    const hasAuthNetworkError = ServerError.is(error) && error.statusCode === HTTP_UNAUTHORIZED;
    const hasAuthGraphQLError =
      CombinedGraphQLErrors.is(error) &&
      error.errors.some(
        (graphQLError) => getGraphQLErrorCode(graphQLError.extensions) === 'UNAUTHENTICATED'
      );

    if (hasAuthNetworkError || hasAuthGraphQLError) {
      onAuthError?.();
    }
  });

  return new ApolloClient({
    link: ApolloLink.from([errorLink, httpLink]),
    cache: new InMemoryCache({
      typePolicies: {
        Country: { keyFields: ['code'] },
        // Operational status has no entity identity; each owning query replaces its embedded value.
        GeoIPStatus: { keyFields: false },
      },
    }),
    defaultOptions: {
      watchQuery: {
        fetchPolicy: 'cache-and-network',
      },
      query: {
        fetchPolicy: 'network-only',
      },
    },
  });
}
