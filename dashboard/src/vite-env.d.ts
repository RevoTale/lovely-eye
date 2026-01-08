/// <reference types="vite/client" />

interface RuntimeConfig {
  BASE_PATH: string;
  API_URL: string;
  GRAPHQL_URL: string;
}

declare global {
  interface Window {
    __ENV__?: Partial<RuntimeConfig>;
  }
}

export {};
