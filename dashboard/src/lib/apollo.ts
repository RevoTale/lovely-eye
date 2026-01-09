import {
  ApolloClient,
  InMemoryCache,
  createHttpLink,
  ApolloLink,
  type NormalizedCacheObject,
} from '@apollo/client';
import { getGraphQLUrl } from '@/config';

// Helper to get cookie value
function getCookie(name: string): string | null {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) {
    return parts.pop()?.split(';').shift() || null;
  }
  return null;
}

// Middleware to add CSRF token to requests
const csrfMiddleware = new ApolloLink((operation, forward) => {
  const csrfToken = getCookie('le_csrf');
  if (csrfToken) {
    operation.setContext(({ headers = {} }) => ({
      headers: {
        ...headers,
        'X-CSRF-Token': csrfToken,
      },
    }));
  }
  return forward(operation);
});

const httpLink = createHttpLink({
  uri: getGraphQLUrl(),
  credentials: 'include', // Include cookies for auth
});

export const apolloClient = new ApolloClient<NormalizedCacheObject>({
  link: ApolloLink.from([csrfMiddleware, httpLink]),
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
