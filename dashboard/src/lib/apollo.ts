import {
  ApolloClient,
  InMemoryCache,
  createHttpLink,
  ApolloLink,
  type ApolloError,
  type NormalizedCacheObject,
} from '@apollo/client';
import { onError } from '@apollo/client/link/error';
import { getGraphQLUrl } from '@/config';

const httpLink = createHttpLink({
  uri: getGraphQLUrl(),
  credentials: 'include',
});

type AuthErrorHandler = () => void;

export const createApolloClient = (onAuthError?: AuthErrorHandler): ApolloClient<NormalizedCacheObject> => {
  const errorLink = onError((response) => {
    const error = ('error' in response ? (response.error as ApolloError | undefined) : undefined);
    const networkError = error?.networkError as { statusCode?: number; status?: number } | undefined;
    const statusCode = networkError?.statusCode ?? networkError?.status;
    const hasAuthNetworkError = statusCode === 401 || statusCode === 403;
    const graphQLErrors = error?.graphQLErrors ?? [];
    const hasAuthGraphQLError = graphQLErrors.some((graphQLError) => {
      const code = graphQLError.extensions?.code;
      return code === 'UNAUTHENTICATED' || code === 'FORBIDDEN';
    });

    if (hasAuthNetworkError || hasAuthGraphQLError) {
      onAuthError?.();
    }
  });

  return new ApolloClient<NormalizedCacheObject>({
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
