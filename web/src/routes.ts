import type { RouteDefinition } from '@solidjs/router'

import UserProfilePage from '~/pages/users/handle'
import FeedPage from '~/pages/feed'
import PostsPage from '~/pages/posts'
import NavigationTabs from '~/components/navigation-tabs'

import UserEditPage from '~/pages/users/edit/index'
import UserEditGeneralPage from '~/pages/users/edit/general'
import UserEditBadgesPage from '~/pages/users/edit/badges'
import UserEditCreateBadgePage from '~/pages/users/edit/createBadge'
import UserEditLocationPage from '~/pages/users/edit/location'
import UserEditInterestsPage from '~/pages/users/edit/interests'
import UserEditDescriptionPage from '~/pages/users/edit/description'
import UserEditImagePage from '~/pages/users/edit/image'

import CollaborationPage from '~/pages/collaborations/id'
import CollaborationEditPage from '~/pages/collaborations/edit'
import CollaborationEditGeneralPage from '~/pages/collaborations/edit/general'
import CollaborationEditLocationPage from '~/pages/collaborations/edit/location'
import CollaborationEditCreateBadgePage from '~/pages/collaborations/edit/createBadge'
import CollaborationEditBadgesPage from '~/pages/collaborations/edit/badges'
import CollaborationEditInterestsPage from '~/pages/collaborations/edit/interests'

import PostIdPage from '~/pages/posts/id'
import PostEditPage from '~/pages/posts/edit'
import PostEditGeneralPage from '~/pages/posts/edit/general'
import PostEditLocationPage from '~/pages/posts/edit/location'

import UserCollaboratePage from '~/pages/users/collaborate'
import CollaborationCollaboratePage from '~/pages/collaborations/collaborate'
import NotFoundPage from './pages/404'

export const routes: RouteDefinition[] = [
	{
		path: '/',
		component: NavigationTabs,
		children: [
			{
				path: '/posts',
				component: PostsPage,
			},
			{
				path: '/',
				component: FeedPage,
			},
		],
	},
	{
		path: '/users/:handle',
		component: UserProfilePage,
	},
	{
		path: '/users/edit',
		component: UserEditPage,
		children: [
			{
				path: '/',
				component: UserEditGeneralPage,
			},
			{
				path: '/badges',
				component: UserEditBadgesPage,
			},
			{
				path: '/create-badge',
				component: UserEditCreateBadgePage,
			},
			{
				path: '/location',
				component: UserEditLocationPage,
			},
			{
				path: '/interests',
				component: UserEditInterestsPage,
			},
			{
				path: '/description',
				component: UserEditDescriptionPage,
			},
			{
				path: '/image',
				component: UserEditImagePage,
			},
		],
	},
	{
		path: '/collaborations/:id',
		component: CollaborationPage,
	},
	{
		path: '/collaborations/edit/:id?',
		component: CollaborationEditPage,
		children: [
			{
				path: '/',
				component: CollaborationEditGeneralPage,
			},
			{
				path: '/location',
				component: CollaborationEditLocationPage,
			},
			{
				path: '/create-badge',
				component: CollaborationEditCreateBadgePage,
			},
			{
				path: '/badges',
				component: CollaborationEditBadgesPage,
			},
			{
				path: '/interests',
				component: CollaborationEditInterestsPage,
			},
		],
	},
	{
		path: '/posts/:id',
		component: PostIdPage,
	},
	{
		path: '/posts/edit/:id?',
		component: PostEditPage,
		children: [
			{
				path: '/',
				component: PostEditGeneralPage,
			},
			{
				path: '/location',
				component: PostEditLocationPage,
			},
		],
	},
	{
		path: '/users/:handle/collaborate',
		component: UserCollaboratePage,
	},
	{
		path: 'collaborations/:id/collaborate',
		component: CollaborationCollaboratePage,
	},
	{
		path: '**',
		component: NotFoundPage,
	},
]
