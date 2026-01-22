import love from 'eslint-config-love';
import reactHooks from 'eslint-plugin-react-hooks';
import reactPlugin from 'eslint-plugin-react';
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
      complexity: ['error', 25],
    },
    settings: {
      react: {
        version: 'detect',
      },
    },
  },
  reactPlugin.configs.flat['recommended'], // This is not a plugin object, but a shareable config object
  reactPlugin.configs.flat['jsx-runtime'],

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
  {
    files: ['src/router.tsx'],
    rules: {
      '@typescript-eslint/only-throw-error': 'off',
    },
  },
  {
    files: ['src/App.tsx'],
    rules: {
      'react-hooks/refs': 'off',
    },
  },
];
