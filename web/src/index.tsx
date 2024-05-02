/* @refresh reload */
import { render } from 'solid-js/web';

import './index.css';
import { Navigate, Route, Router } from '@solidjs/router';
import { lazy } from 'solid-js';
import App from './App';

import GeneralInfo from '~/pages/users/edit/general';
import SelectBadges from '~/pages/users/edit/badges';
import CreateBadge from '~/pages/users/edit/createBadge';
import EditLocation from '~/pages/users/edit/location';
import SelectOpportunities from '~/pages/users/edit/interests';
import ImageUpload from '~/pages/users/edit/image';
import Description from '~/pages/users/edit/description';
import CollaborationInfo from '~/pages/collaborations/edit/general';
import CollaborationLocation from '~/pages/collaborations/edit/location';
import CollaborationBadges from '~/pages/collaborations/edit/badges';
import CollaborationCreateBadge from '~/pages/collaborations/edit/createBadge';
import CollaborationOpportunity from '~/pages/collaborations/edit/interests';

const Users = lazy(() => import('~/pages/users'));
const Collaborations = lazy(() => import('~/pages/collaborations'));
const Collaboration = lazy(() => import('~/pages/collaborations/[id]'));
const Home = lazy(() => import('~/pages'));
const User = lazy(() => import('~/pages/users/[id]'));
const EditUser = lazy(() => import('~/pages/users/edit/index'));
const EditCollaboration = lazy(() => import('~/pages/collaborations/edit'));
const UserCollaborate = lazy(() => import('~/pages/users/collaborate'));

const root = document.getElementById('root');

if (import.meta.env.DEV && !(root instanceof HTMLElement)) {
  throw new Error(
    'Root element not found. Did you forget to add it to your index.html? Or maybe the id attribute got misspelled?',
  );
}

render(
  () => (
    <Router root={App}>
      <Route path="/" component={Home} />
      <Route path="/users" component={Users} />
      <Route path="/users/:id" component={User} />
      <Route path="/users/edit" component={EditUser}>
        <Route path="/" component={GeneralInfo} />
        <Route path="/badges" component={SelectBadges} />
        <Route path="/create-badge" component={CreateBadge} />
        <Route path="/location" component={EditLocation} />
        <Route path="/interests" component={SelectOpportunities} />
        <Route path="/description" component={Description} />
        <Route path="/image" component={ImageUpload} />
      </Route>
      <Route path="/collaborations" component={Collaborations} />
      <Route path="/collaborations/:id" component={Collaboration} />
      <Route path="/collaborations/edit/:id?" component={EditCollaboration}>
        <Route path="/" component={CollaborationInfo} />
        <Route path="/location" component={CollaborationLocation} />
        <Route path="/create-badge" component={CollaborationCreateBadge} />
        <Route path="/badges" component={CollaborationBadges} />
        <Route path="/interests" component={CollaborationOpportunity} />
      </Route>
      <Route path="/users/:id/collaborate" component={UserCollaborate} />
      <Route path="*" component={() => <Navigate href={'/'} />} />
    </Router>
  ),
  root!,
);
