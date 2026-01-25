import love from 'eslint-config-love';
import importPlugin from 'eslint-plugin-import';
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
      'src/routeTree.gen.ts',
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
    plugins: {
      ...(love.plugins ?? {}),
      import: importPlugin,
    },
    rules: {
      ...love.rules,
      complexity: ['error', 25],
      'import/no-namespace': 'error',
      'no-restricted-imports': [
        'error',
        {
          paths: [
            {
              name: '@apollo/client',
              importNames: ['gql'],
              message: 'Use .graphql files instead of gql templates.',
            },
            {
              name: 'graphql-tag',
              message: 'Use .graphql files instead of gql templates.',
            },
          ],
        },
      ],
      'no-restricted-syntax': [
        'error',
        {
          selector: "TaggedTemplateExpression[tag.name='gql']",
          message: 'Use .graphql files instead of gql templates.',
        },
      ],
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
      react: reactPlugin,
      'react-hooks': reactHooks,
    },
    rules: {
      ...reactHooks.configs.recommended.rules,
      'react/function-component-definition': [
        'error',
        {
          namedComponents: 'arrow-function',
          unnamedComponents: 'arrow-function',
        },
      ],
      'no-restricted-syntax': [
        'error',
        {
          selector: "TaggedTemplateExpression[tag.name='gql']",
          message: 'Use .graphql files instead of gql templates.',
        },
        {
          selector:
            "VariableDeclarator[id.typeAnnotation=null][init.type='ArrowFunctionExpression'][id.name=/^[A-Z]/]",
          message:
            'Component declarations must use an explicit React.FunctionComponent type annotation.',
        },
        {
          selector:
            "VariableDeclarator[id.name=/^[A-Z]/][init.type='ArrowFunctionExpression'][init.returnType.typeAnnotation.type='TSTypeReference'][init.returnType.typeAnnotation.typeName.type='TSQualifiedName'][init.returnType.typeAnnotation.typeName.right.name='Element'][init.returnType.typeAnnotation.typeName.left.type='TSQualifiedName'][init.returnType.typeAnnotation.typeName.left.right.name='JSX'][init.returnType.typeAnnotation.typeName.left.left.name='React']",
          message: 'Use React.ReactNode for component return types instead of React.JSX.Element.',
        },
      ],
      'react/jsx-no-bind': [
        'error',
        {
          ignoreRefs: true,
          allowArrowFunctions: true,
          allowBind: false,
        },
      ],
      'react-hooks/rules-of-hooks': 'error',
      'react-hooks/exhaustive-deps': 'warn',
    },
  },
  {
    files: ['src/router.tsx'],
    rules: {
      '@typescript-eslint/only-throw-error': 'off',
    },
  },
  {
    files: ['src/routes/**/*.tsx'],
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
