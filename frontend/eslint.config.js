import js from '@eslint/js';
import globals from 'globals';

import reactRefresh from 'eslint-plugin-react-refresh';
import tseslint from 'typescript-eslint';
import { defineConfig, globalIgnores } from 'eslint/config';
import simpleImportSort from 'eslint-plugin-simple-import-sort';
import unusedImports from 'eslint-plugin-unused-imports';
import pluginQuery from '@tanstack/eslint-plugin-query';

export default defineConfig([
  globalIgnores(['dist']),

  {
    files: ['**/*.{ts,tsx}'],

    extends: [
      js.configs.recommended,
      ...tseslint.configs.recommended,
      reactRefresh.configs.recommended
    ],

    ignores: [
      'node_modules/**',
      '.next/**',
      'out/**',
      'build/**',
      'coverage/**',
      '.turbo/**',
      '.vercel/**',
      '.cache/**',
      'public/**',
      'tmp/**',
      'temp/**',
      'next-env.d.ts',
      '**/*.config.*',
      '**/*.d.ts',
      'src/components/ui/**'
    ],

    languageOptions: {
      ecmaVersion: 2020,
      globals: globals.browser,
      parserOptions: {
        ecmaFeatures: {
          jsx: true
        }
      }
    },

    plugins: {
      'simple-import-sort': simpleImportSort,
      'unused-imports': unusedImports,
      '@tanstack/query': pluginQuery
    },

    rules: {
      // --- import sorting ---
      'simple-import-sort/imports': 'error',
      'simple-import-sort/exports': 'error',

      // --- unused imports ---
      'unused-imports/no-unused-imports': 'error',

      // disable base rule in favor of unused-imports
      '@typescript-eslint/no-unused-vars': 'off',

      'unused-imports/no-unused-vars': [
        'warn',
        {
          vars: 'all',
          varsIgnorePattern: '^_',
          args: 'after-used',
          argsIgnorePattern: '^_'
        }
      ],

      // --- TanStack Query ---
      '@tanstack/query/exhaustive-deps': 'error'
    }
  }
]);
