import { Navigate, Route } from '@solidjs/router';
import { lazy, onCleanup } from 'solid-js';

import { createNavigator } from '~/navigation/createNavigator.js';
import { createRouter } from '~/navigation/createRouter.js';
import GeneralInfo from '~/pages/users/[edit]/general';
import SelectBadges from '~/pages/users/[edit]/badges';
import CreateBadge from '~/pages/users/[edit]/createBadge';
import EditLocation from '~/pages/users/[edit]/location';
import SelectOpportunities from '~/pages/users/[edit]/interests';
import ImageUpload from '~/pages/users/[edit]/image';
import Description from '~/pages/users/[edit]/description';
import CollaborationInfo from '~/pages/collaborations/[edit]/general';
import EditCollaboration from '~/pages/collaborations/[edit]';
import CollaborationLocation from '~/pages/collaborations/[edit]/location';
import CollaborationBadges from '~/pages/collaborations/[edit]/badges';
import CollaborationOpportunity from '~/pages/collaborations/[edit]/interests';

const Users = lazy(() => import('~/pages/users'));
const Collaborations = lazy(() => import('~/pages/collaborations'));
const Collaboration = lazy(() => import('~/pages/collaborations/[id]'));
const Home = lazy(() => import('~/pages'));
const User = lazy(() => import('~/pages/users/[id]'));
const EditUser = lazy(() => import('~/pages/users/[edit]/index'));
const CreateCollaboration = lazy(() => import('~/pages/collaborations/[edit]'));
const UserCollaborate = lazy(() => import('~/pages/users/collaborate'));

export function App() {
  // Create new application navigator and attach it to the browser history, so it could modify
  // it and listen to its changes.
  const navigator = createNavigator();

  void navigator.attach();

  onCleanup(() => {
    navigator.detach();
  });

  const Router = createRouter(navigator);

  return (
    <Router>
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
      <Route path="/collaborations/edit" component={EditCollaboration}>
        <Route path="/" component={CollaborationInfo} />
        <Route path="/location" component={CollaborationLocation} />
        <Route path="/badges" component={CollaborationBadges} />
        <Route path="/interests" component={CollaborationOpportunity} />
      </Route>
      <Route path="/collaborations" component={Collaborations} />
      <Route path="/collaborations/create" component={CreateCollaboration} />
      <Route path="/collaborations/:id" component={Collaboration} />
      <Route path="/users/:id/collaborate" component={UserCollaborate} />

      <Route path="*" component={() => <Navigate href={'/'} />} />
    </Router>
  );
}
