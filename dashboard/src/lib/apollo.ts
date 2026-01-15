import {
  ApolloClient,
  InMemoryCache,
  createHttpLink,
  type NormalizedCacheObject,
} from '@apollo/client';
import { getGraphQLUrl } from '@/config';

const httpLink = createHttpLink({
  uri: getGraphQLUrl(),
  credentials: 'include',
});

export const apolloClient = new ApolloClient<NormalizedCacheObject>({
  link: httpLink,
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
