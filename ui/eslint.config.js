// eslint.config.js
import typescript from '@typescript-eslint/eslint-plugin';
import tsParser from '@typescript-eslint/parser';

export default [
  {
    files: ['src/**/*.{ts,tsx,js,jsx}'],
    languageOptions: {
      parser: tsParser,
      parserOptions: {
        ecmaVersion: 'latest',
        sourceType: 'module',
        project: './tsconfig.json', // Only needed if you use types like `@typescript-eslint/explicit-module-boundary-types`
      },
    },
    plugins: {
      '@typescript-eslint': typescript,
    },
    rules: {
      // --- JavaScript / TypeScript basic rules ---
      'no-unused-vars': 'warn',
      'no-console': 'warn',

      // --- TypeScript specific rules ---
      '@typescript-eslint/explicit-function-return-type': 'off',
      '@typescript-eslint/no-explicit-any': 'warn',
      '@typescript-eslint/no-unused-vars': ['warn', { argsIgnorePattern: '^_' }],
      
      // --- Style rules ---
      'semi': ['error', 'always'],
      'quotes': ['error', 'single'],
      'indent': ['error', 2],

      // --- Example Prettier-like rules ---
      'comma-dangle': ['error', 'always-multiline'],
    },
  },
];
