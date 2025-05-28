import type { RouteDefinition } from '@solidjs/router'
import { lazy } from 'solid-js'
import { Navigate } from '@solidjs/router'
import { AdminLayout } from '~/components/AdminLayout'
import UsersPage from '~/pages/users'
import CollaborationsPage from '~/pages/collaborations'
import LoginPage from '~/pages/login'

// Lazy load pages that don't exist yet
const BadgesPage = lazy(() => import('~/pages/badges'))
const CitiesPage = lazy(() => import('~/pages/cities'))
const OpportunitiesPage = lazy(() => import('~/pages/opportunities'))
const AdminsPage = lazy(() => import('~/pages/admins'))

export const routes: RouteDefinition[] = [
	{
		path: '/',
		component: AdminLayout,
		children: [
			{
				path: '/',
				component: () => Navigate({ href: "/users" }),
			},
			{
				path: 'users',
				component: UsersPage,
			},
			{
				path: 'badges',
				component: BadgesPage,
			},
			{
				path: 'cities',
				component: CitiesPage,
			},
			{
				path: 'collaborations',
				component: CollaborationsPage,
			},
			{
				path: 'opportunities',
				component: OpportunitiesPage,
			},
			{
				path: 'admins',
				component: AdminsPage,
			},
		],
	},
	{
		path: '/login',
		component: LoginPage,
	},
]
