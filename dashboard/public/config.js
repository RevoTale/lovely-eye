// Runtime configuration - this file is replaced/modified at container startup by the Go server
// The Go server will inject the actual values based on environment variables

window.__ENV__ = {
  BASE_PATH: '/',
  API_URL: '/api',
  GRAPHQL_URL: '/graphql',
};
