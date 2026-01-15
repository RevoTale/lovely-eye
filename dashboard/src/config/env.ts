interface RuntimeConfig {
  BASE_PATH: string;
  API_URL: string;
  GRAPHQL_URL: string;
}

function getConfig(): RuntimeConfig {
  const env = window.__ENV__
  if (!env) {
    throw new Error('Runtime environment configuration is missing.');
  }
  if (env.BASE_PATH === undefined || !env.API_URL || !env.GRAPHQL_URL) {
    throw new Error('Incomplete runtime environment configuration.');
  }
  return {
    BASE_PATH: env.BASE_PATH,
    API_URL: env.API_URL,
    GRAPHQL_URL: env.GRAPHQL_URL,
  };
}

export const config = getConfig();

export function getBasePath(): string {
  return config.BASE_PATH.replace(/\/$/, '') || '/';
}

export function getApiUrl(): string {
  return config.API_URL;
}

export function getGraphQLUrl(): string {
  return config.GRAPHQL_URL;
}
