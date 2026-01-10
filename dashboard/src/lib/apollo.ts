import {
  ApolloClient,
  InMemoryCache,
  createHttpLink,
  type NormalizedCacheObject,
} from '@apollo/client';
import { getGraphQLUrl } from '@/config';

// Create HTTP link with credentials to send cookies
// Auth uses HttpOnly + Secure cookies with SameSite=Strict/Lax
// No CSRF tokens needed - modern browser security handles it
const httpLink = createHttpLink({
  uri: getGraphQLUrl(),
  credentials: 'include', // Include cookies for auth
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
