import {
  ApolloClient,
  InMemoryCache,
  ApolloLink,
} from '@apollo/client';
import { HttpLink } from '@apollo/client/link/http';
import { ErrorLink } from '@apollo/client/link/error';
import { CombinedGraphQLErrors, ServerError } from '@apollo/client/errors';
import { getGraphQLUrl } from '@/config';

const HTTP_UNAUTHORIZED = 401;
const HTTP_FORBIDDEN = 403;

const httpLink = new HttpLink({
  uri: getGraphQLUrl(),
  credentials: 'include',
});

type AuthErrorHandler = () => void;

export const createApolloClient = (onAuthError?: AuthErrorHandler): ApolloClient => {
  const errorLink = new ErrorLink(({ error }) => {
    const hasAuthNetworkError =
      ServerError.is(error) &&
      (error.statusCode === HTTP_UNAUTHORIZED || error.statusCode === HTTP_FORBIDDEN);
    const hasAuthGraphQLError =
      CombinedGraphQLErrors.is(error) &&
      error.errors.some((graphQLError) => {
        const code = (graphQLError.extensions as { code?: string } | undefined)?.code;
        return code === 'UNAUTHENTICATED' || code === 'FORBIDDEN';
      });

    if (hasAuthNetworkError || hasAuthGraphQLError) {
      onAuthError?.();
    }
  });

  return new ApolloClient({
    link: ApolloLink.from([errorLink, httpLink]),
    cache: new InMemoryCache(),
    defaultOptions: {
      watchQuery: {
        fetchPolicy: 'cache-and-network',
      },
      query: {
        fetchPolicy: 'network-only',
      },
    },
  });
};
