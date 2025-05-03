import eslint from '@eslint/js'
import * as tseslint from 'typescript-eslint'
import prettierConfig from 'eslint-config-prettier'
import solidPlugin from 'eslint-plugin-solid'
import tailwindPlugin from 'eslint-plugin-tailwindcss'

export default tseslint.config(
	eslint.configs.recommended,
	...tseslint.configs.recommended,
	{
		files: ['**/*.ts', '**/*.tsx'],
		ignores: [
			'node_modules/**',
			'dist/**',
			'public/**',
			'**/*.d.ts',
		],
		plugins: {
			solid: solidPlugin,
			tailwindcss: tailwindPlugin,
		},
		languageOptions: {
			parser: tseslint.parser,
			parserOptions: {
				ecmaVersion: 'latest',
				sourceType: 'module',
			},
			globals: {
				window: 'readonly',
				document: 'readonly',
				console: 'readonly',
				setTimeout: 'readonly',
				clearTimeout: 'readonly',
				setInterval: 'readonly',
				clearInterval: 'readonly',
				fetch: 'readonly',
				URL: 'readonly',
				URLSearchParams: 'readonly',
				FileReader: 'readonly',
				File: 'readonly',
				AbortController: 'readonly',
				queueMicrotask: 'readonly',
				structuredClone: 'readonly',
				Event: 'readonly',
			},
		},
		rules: {
			...prettierConfig.rules,
			'@typescript-eslint/no-explicit-any': 'off',
			'@typescript-eslint/no-unused-vars': ['error', {
				argsIgnorePattern: '^_',
				varsIgnorePattern: '^_',
			}],
			'@typescript-eslint/no-unused-expressions': 'off',
			'tailwindcss/classnames-order': 'warn',
			'tailwindcss/no-custom-classname': ['warn', {
				whitelist: [
					'text-main',
					'text-peatch-green',
					'text-green',
					'text-link',
					'text-hint',
					'text-button',
					'text-secondary-bg',
					'bg-main',
					'placeholder:text-hint',
				],
			}],
			'tailwindcss/no-contradicting-classname': 'error',
			'solid/reactivity': 'warn',
			'solid/no-destructure': 'warn',
			'solid/jsx-no-undef': 'error',
			'@typescript-eslint/ban-ts-comment': ['error', {
				'ts-expect-error': 'allow-with-description',
				'ts-ignore': false,
				'ts-nocheck': false,
				'ts-check': false,
			}],
			'no-empty': ['error', { 'allowEmptyCatch': true }],
			'no-constant-binary-expression': 'error',
			'@typescript-eslint/no-non-null-asserted-optional-chain': 'error',
		},
		settings: {
			tailwindcss: {
				callees: ['cn', 'clsx', 'twMerge'],
				config: 'tailwind.config.cjs',
			},
		},
	},
)
