import type { RouteDefinition } from '@solidjs/router'

import UserProfilePage from '~/pages/users/handle'
import FeedPage from '~/pages/feed'
import NavigationTabs from '~/components/navigation-tabs'

import UserEditPage from '~/pages/users/edit/index'
import UserEditGeneralPage from '~/pages/users/edit/general'
import UserEditBadgesPage from '~/pages/users/edit/badges'
import UserEditCreateBadgePage from '~/pages/users/edit/create-badge'
import UserEditLocationPage from '~/pages/users/edit/location'
import UserEditInterestsPage from '~/pages/users/edit/interests'
import UserEditDescriptionPage from '~/pages/users/edit/description'
import UserEditImagePage from '~/pages/users/edit/image'

import CollaborationPage from '~/pages/collaborations/id'
import CollaborationEditPage from '~/pages/collaborations/edit'
import CollaborationEditGeneralPage from '~/pages/collaborations/edit/general'
import CollaborationEditLocationPage from '~/pages/collaborations/edit/location'
import CollaborationEditCreateBadgePage from '~/pages/collaborations/edit/create-badge'
import CollaborationEditBadgesPage from '~/pages/collaborations/edit/badges'
import CollaborationEditInterestsPage from '~/pages/collaborations/edit/interests'

import NotFoundPage from './pages/404'
import PostsPage from '~/pages/posts'

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
    path: '**',
    component: NotFoundPage,
  },
]
