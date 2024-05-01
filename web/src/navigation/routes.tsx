import { Component, lazy } from 'solid-js';

const Users = lazy(() => import('~/pages/users'));
const Collaborations = lazy(() => import('~/pages/collaborations'));
const Collaboration = lazy(() => import('~/pages/collaborations/[id]'));
const Home = lazy(() => import('~/pages'));
const User = lazy(() => import('~/pages/users/[id]'));
const EditUser = lazy(() => import('~/pages/users/edit'));
const CreateCollaboration = lazy(() => import('~/pages/collaborations/create'));
const UserCollaborate = lazy(() => import('~/pages/users/collaborate'));

interface Route {
  path: string;
  Component: Component;
}

export const routes: Route[] = [
  { path: '/', Component: Home },
  { path: '/users', Component: Users },
  { path: '/collaborations', Component: Collaborations },
  { path: '/collaborations/:id', Component: Collaboration },
  { path: '/users/:id', Component: User },
  { path: '/users/edit', Component: EditUser },
  { path: '/collaborations/create', Component: CreateCollaboration },
  { path: '/users/:id/collaborate', Component: UserCollaborate },
];