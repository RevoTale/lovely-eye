interface RuntimeConfig {
  BASE_PATH: string;
  API_URL: string;
  GRAPHQL_URL: string;
}

const defaultConfig: RuntimeConfig = {
  BASE_PATH: '/',
  API_URL: '/api',
  GRAPHQL_URL: '/graphql',
};

function getConfig(): RuntimeConfig {
  const env = window.__ENV__ ?? {};
  return {
    BASE_PATH: env.BASE_PATH ?? defaultConfig.BASE_PATH,
    API_URL: env.API_URL ?? defaultConfig.API_URL,
    GRAPHQL_URL: env.GRAPHQL_URL ?? defaultConfig.GRAPHQL_URL,
  };
}

export const config = getConfig();

export function getBasePath(): string {
  // Remove trailing slash for consistency
  return config.BASE_PATH.replace(/\/$/, '') || '/';
}

export function getApiUrl(): string {
  return config.API_URL;
}

export function getGraphQLUrl(): string {
  return config.GRAPHQL_URL;
}
