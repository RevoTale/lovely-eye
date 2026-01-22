import type { CodegenConfig } from '@graphql-codegen/cli';

const config: CodegenConfig = {
  schema: '../server/schema.graphqls',
  documents: 'src/**/*.graphql',
  generates: {
    './src/gql/': {
      preset: 'client',
      presetConfig: {
        fragmentMasking: true,
      },
      config: {
        strictScalars: true,
        scalars: { Time: 'string' },
        enumsAsConst: true,
        avoidOptionals: { field: true, inputValue: true, object: true },
        nonOptionalTypename: true,
        maybeValue: 'T | null',
        inputMaybeValue: 'T | null',
        useTypeImports: true,
      },
    },
  },
};

export default config;
