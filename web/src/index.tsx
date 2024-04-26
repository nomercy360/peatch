/* @refresh reload */
import { render } from 'solid-js/web';

import './index.css';
import { Router, Route } from '@solidjs/router';

import Home from './pages/Home';
import Users from './pages/Profiles';
import App from './App';

const root = document.getElementById('root');

if (import.meta.env.DEV && !(root instanceof HTMLElement)) {
  throw new Error(
    'Root element not found. Did you forget to add it to your index.html? Or maybe the id attribute got misspelled?',
  );
}

function NotFound() {
  return <div>Not Found</div>;
}

render(() => (
  <Router root={App}>
    <Route path="/users" component={Users} />
    <Route path="/" component={Home} />
    <Route path="*404" component={NotFound} />
  </Router>
), root!);
