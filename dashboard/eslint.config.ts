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
      // Practical adjustments
      '@typescript-eslint/no-magic-numbers': 'off',
      '@typescript-eslint/explicit-function-return-type': 'off',
      '@typescript-eslint/prefer-destructuring': 'off',
      '@typescript-eslint/strict-boolean-expressions': 'off',
      '@typescript-eslint/max-params': 'off',
      '@typescript-eslint/prefer-nullish-coalescing': 'off',
      '@typescript-eslint/promise-function-async': 'off',
      '@typescript-eslint/only-throw-error': 'off',
      '@typescript-eslint/triple-slash-reference': 'off',
      'complexity': ['error', { max: 25 }],
      'require-unicode-regexp': 'off',
      'no-negated-condition': 'off',
      'no-alert': 'off',
      'arrow-body-style': 'off',
      'prefer-named-capture-group': 'off',
      'import/no-absolute-path': 'off',
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
