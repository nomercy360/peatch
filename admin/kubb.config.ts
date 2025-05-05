import { defineConfig } from '@kubb/core'
import { pluginTs } from '@kubb/plugin-ts'
import { pluginOas } from '@kubb/plugin-oas'
import { pluginSolidQuery } from '@kubb/plugin-solid-query'

export default defineConfig(() => {
	return {
		root: '.',
		input: {
			path: '../docs/swagger.yaml',
		},
		output: {
			path: './src/gen',
		},
		plugins: [
			pluginOas(),
			pluginTs({
				output: {
					path: './types',
				},
			}),
		],
	}
})
