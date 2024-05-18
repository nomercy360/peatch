import { lazy } from 'solid-js';
import type { RouteDefinition } from '@solidjs/router';

import HomePage from '~/pages/index';

export const routes: RouteDefinition[] = [
	{
		path: '/',
		component: HomePage,
	},
	{
		path: '/users',
		component: lazy(() => import('~/pages/users')),
	},
	{
		path: '/users/:handle',
    component: lazy(() => import('~/pages/users/handle')),
  },
  {
    path: '/users/:id/followers',
    component: lazy(() => import('~/pages/users/followers')),
	},
	{
		path: '/users/edit',
		component: lazy(() => import('~/pages/users/edit/index')),
		children: [
			{
				path: '/',
				component: lazy(() => import('~/pages/users/edit/general')),
			},
			{
				path: '/badges',
				component: lazy(() => import('~/pages/users/edit/badges')),
			},
			{
				path: '/create-badge',
				component: lazy(() => import('~/pages/users/edit/createBadge')),
			},
			{
				path: '/location',
				component: lazy(() => import('~/pages/users/edit/location')),
			},
			{
				path: '/interests',
				component: lazy(() => import('~/pages/users/edit/interests')),
			},
			{
				path: '/description',
				component: lazy(() => import('~/pages/users/edit/description')),
			},
			{
				path: '/image',
				component: lazy(() => import('~/pages/users/edit/image')),
			},
		],
	},
	{
		path: '/collaborations',
		component: lazy(() => import('~/pages/collaborations')),
	},
	{
		path: '/collaborations/:id',
    component: lazy(() => import('~/pages/collaborations/id')),
	},
	{
		path: '/collaborations/edit/:id?',
		component: lazy(() => import('~/pages/collaborations/edit')),
		children: [
			{
				path: '/',
				component: lazy(() => import('~/pages/collaborations/edit/general')),
			},
			{
				path: '/location',
				component: lazy(() => import('~/pages/collaborations/edit/location')),
			},
			{
				path: '/create-badge',
				component: lazy(
					() => import('~/pages/collaborations/edit/createBadge'),
				),
			},
			{
				path: '/badges',
				component: lazy(() => import('~/pages/collaborations/edit/badges')),
			},
			{
				path: '/interests',
				component: lazy(() => import('~/pages/collaborations/edit/interests')),
			},
		],
	},
	{
		path: '/users/:handle/collaborate',
		component: lazy(() => import('~/pages/users/collaborate')),
	},
	{
		path: 'collaborations/:id/collaborate',
		component: lazy(() => import('~/pages/collaborations/collaborate')),
	},
	{
    path: '/users/shuffle',
    component: lazy(() => import('~/pages/users/shuffle')),
  },
  {
		path: '**',
		component: lazy(() => import('./pages/404')),
	},
]
