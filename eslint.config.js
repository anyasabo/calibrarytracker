import js from '@eslint/js';
import ts from 'typescript-eslint';
import svelte from 'eslint-plugin-svelte';
import svelteParser from 'svelte-eslint-parser';
import a11y from 'eslint-plugin-svelte';
import globals from 'globals';

export default ts.config(
	js.configs.recommended,
	...ts.configs.recommended,
	...svelte.configs['flat/recommended'],
	{
		languageOptions: {
			globals: {
				...globals.browser,
				...globals.node
			}
		}
	},
	{
		files: ['**/*.svelte', '**/*.svelte.ts', '**/*.svelte.js'],
		languageOptions: {
			parser: svelteParser,
			parserOptions: {
				parser: ts.parser
			}
		}
	},
	{
		rules: {
			// Enforce explicit types on function signatures
			'@typescript-eslint/explicit-function-return-type': [
				'warn',
				{ allowExpressions: true, allowTypedFunctionExpressions: true }
			],
			// Ban 'any' — forces contributors to think about types
			'@typescript-eslint/no-explicit-any': 'error',
			// Catch unused variables (prefix with _ to intentionally ignore)
			'@typescript-eslint/no-unused-vars': ['error', { argsIgnorePattern: '^_' }],
			// No console.log in production code (use console.warn/error for real issues)
			'no-console': ['warn', { allow: ['warn', 'error'] }],
			// Disabled: we use adapter-static with base: '' so resolve() adds no value.
			// Re-enable this if we deploy under a sub-path (e.g., /calibrarytracker/).
			'svelte/no-navigation-without-resolve': 'off'
		}
	},
	{
		files: ['scripts/check-data-coordinates.mjs'],
		rules: {
			// This is an explicit CLI/reporting script, so console output is intentional.
			'no-console': 'off',
			// JS utility script; explicit TS return annotations are not meaningful here.
			'@typescript-eslint/explicit-function-return-type': 'off'
		}
	},
	{
		ignores: [
			'build/',
			'.svelte-kit/',
			'.playwright-pages/',
			'playwright-report/',
			'test-results/',
			'node_modules/',
			'data/',
			'cmd/',
			'internal/',
			'*.config.js',
			'*.config.ts'
		]
	}
);
