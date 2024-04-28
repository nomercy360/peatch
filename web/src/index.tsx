/* @refresh reload */
import { render } from 'solid-js/web';

import './index.css';
import { Route, Router } from '@solidjs/router';
import { lazy } from 'solid-js';
import App from './App';

const Users = lazy(() => import('./pages/users'));
const Home = lazy(() => import('./pages'));
const User = lazy(() => import('./pages/users/[id]'));
const EditUser = lazy(() => import('./pages/collaborations/create'));
const CreateCollaboration = lazy(() => import('./pages/collaborations/create'));

const root = document.getElementById('root');

if (import.meta.env.DEV && !(root instanceof HTMLElement)) {
  throw new Error(
    'Root element not found. Did you forget to add it to your index.html? Or maybe the id attribute got misspelled?',
  );
}

const NotFound = () => <div>Not Found</div>;

render(
  () => (
    <Router root={App}>
      <Route path="/" component={Home} />
      <Route path="/users" component={Users} />
      <Route path="/users/:id" component={User} />
      <Route path="/users/edit" component={EditUser} />
      <Route path="/collaborations/create" component={CreateCollaboration} />
      <Route path="/*all" component={NotFound} />
    </Router>
  ),
  root!,
);
