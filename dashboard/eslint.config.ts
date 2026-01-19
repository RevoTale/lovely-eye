import love from 'eslint-config-love';
import reactHooks from 'eslint-plugin-react-hooks';
import globals from 'globals';

export default [
  // Ignore patterns
  {
    ignores: [
      'dist/**',
      'node_modules/**',
      'src/gql/**',
      'src/components/ui/**',
      'src/vite-env.d.ts',
      'public/**',
      '*.config.js',
      '*.config.ts',
    ],
  },

  // Base eslint-config-love with TypeScript
  {
    ...love,
    files: ['**/*.ts', '**/*.tsx'],
    languageOptions: {
      ...love.languageOptions,
      globals: {
        ...globals.browser,
        ...globals.es2020,
      },
    },
    rules: {
      ...love.rules,
      'complexity': ['error', { max: 25 }],
    },
  },

  // React Hooks
  {
    files: ['**/*.tsx'],
    plugins: {
      'react-hooks': reactHooks,
    },
    rules: {
      ...reactHooks.configs.recommended.rules,
    },
  },
];
