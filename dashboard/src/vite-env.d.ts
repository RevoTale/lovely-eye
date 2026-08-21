/// <reference types="vite/client" />

interface RuntimeConfig {
  BASE_PATH: string;
  GRAPHQL_URL: string;
}

declare global {
  interface Window {
    __ENV__?: Partial<RuntimeConfig>;
  }
}

export {};
