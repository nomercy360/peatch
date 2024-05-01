import { retrieveLaunchParams, setDebug } from '@tma.js/sdk';
import { DisplayGate, SDKProvider, useLaunchParams } from '@tma.js/sdk-solid';
import { Component, createEffect, createSignal, Match, Switch } from 'solid-js';

import { App } from '~/components/App';
import { QueryClient, QueryClientProvider } from '@tanstack/solid-query';
import { API_BASE_URL } from '~/api';
import { setToken, setUser } from '~/store';

const Err: Component<{ error: unknown }> = props => (
  <div>
    <p>An error occurred while initializing the SDK</p>
    <blockquote>
      <code>
        {props.error instanceof Error
          ? props.error.message
          : JSON.stringify(props.error)}
      </code>
    </blockquote>
  </div>
);

const Loading: Component = () => <div>Application is loading</div>;

const queryClient = new QueryClient({
  defaultOptions: {
    queries: { staleTime: 1000 * 60 * 5, gcTime: 1000 * 60 * 10 },
  },
});

const { initDataRaw } = retrieveLaunchParams();

export const Root: Component = () => {
  if (useLaunchParams().startParam === 'debug') {
    setDebug(true);
    import('eruda').then(lib => lib.default.init());
  }

  const [isAuthenticated, setIsAuthenticated] = createSignal(false);
  const [isLoading, setIsLoading] = createSignal(true);

  createEffect(async () => {
    try {
      const resp = await fetch(`${API_BASE_URL}/auth/telegram?` + initDataRaw, {
        method: 'POST',
      });

      const { user, token } = await resp.json();

      setUser(user);
      setToken(token);

      setIsAuthenticated(true);
      setIsLoading(false);
    } catch (e) {
      console.error('Failed to authenticate user:', e);
      setIsAuthenticated(false);
      setIsLoading(false);
    }
  });

  return (
    <SDKProvider
      options={{ acceptCustomStyles: true, cssVars: true, complete: true }}
    >
      <DisplayGate error={Err} loading={Loading} initial={Loading}>
        <QueryClientProvider client={queryClient}>
          <Switch>
            <Match when={isLoading()}>Loading...</Match>
            <Match when={!isAuthenticated()}>Not authenticated</Match>
            <Match when={isAuthenticated()}>
              <App />
            </Match>
          </Switch>
        </QueryClientProvider>
      </DisplayGate>
    </SDKProvider>
  );
};
