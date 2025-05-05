import type { RouteDefinition } from '@solidjs/router'
import HomePage from './pages'


export const routes: RouteDefinition[] = [
	{
		path: '/',
		component: HomePage,
	},
]
