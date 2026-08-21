interface RuntimeConfig {
  BASE_PATH: string;
  GRAPHQL_URL: string;
}

function getConfig(): RuntimeConfig {
  const { __ENV__: runtimeEnv } = window;
  if (runtimeEnv === undefined) {
    throw new Error('Runtime environment configuration is missing.');
  }
  const { BASE_PATH, GRAPHQL_URL } = runtimeEnv;
  if (BASE_PATH === undefined || GRAPHQL_URL === undefined || GRAPHQL_URL === '') {
    throw new Error('Incomplete runtime environment configuration.');
  }
  return { BASE_PATH, GRAPHQL_URL };
}

const config = getConfig();

export function getBasePath(): string {
  const trimmedBasePath = config.BASE_PATH.replace(/\/$/v, '');
  return trimmedBasePath === '' ? '/' : trimmedBasePath;
}

export function getGraphQLUrl(): string {
  return config.GRAPHQL_URL;
}
