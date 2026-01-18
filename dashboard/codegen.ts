import type { CodegenConfig } from '@graphql-codegen/cli';

const config: CodegenConfig = {
  schema: '../server/schema.graphqls',
  documents: 'src/**/*.graphql',
  generates: {
    './src/gql/': {
      preset: 'client',
      presetConfig: {
        fragmentMasking: false,
      },
      config: {
        strictScalars: true,
        scalars: { Time: 'string' },
        enumsAsTypes: true,
        avoidOptionals: { field: true, inputValue: false, object: true },
        nonOptionalTypename: true,
        useTypeImports: true,
      },
    },
  },
};

export default config;
