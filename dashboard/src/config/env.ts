interface RuntimeConfig {
  BASE_PATH: string;
  API_URL: string;
  GRAPHQL_URL: string;
}

function getConfig(): RuntimeConfig {
  const { __ENV__: runtimeEnv } = window
  if (runtimeEnv === undefined) {
    throw new Error('Runtime environment configuration is missing.');
  }
  const { BASE_PATH, API_URL, GRAPHQL_URL } = runtimeEnv
  if (
    BASE_PATH === undefined ||
    API_URL === undefined ||
    API_URL === '' ||
    GRAPHQL_URL === undefined ||
    GRAPHQL_URL === ''
  ) {
    throw new Error('Incomplete runtime environment configuration.');
  }
  return {
    BASE_PATH,
    API_URL,
    GRAPHQL_URL,
  };
}

export const config = getConfig();

export function getBasePath(): string {
  const trimmedBasePath = config.BASE_PATH.replace(/\/$/v, '');
  return trimmedBasePath === '' ? '/' : trimmedBasePath;
}

export function getApiUrl(): string {
  return config.API_URL;
}

export function getGraphQLUrl(): string {
  return config.GRAPHQL_URL;
}
