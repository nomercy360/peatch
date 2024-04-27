import { createEffect, createSignal, Match, Switch } from 'solid-js';
import { QueryClient, QueryClientProvider } from '@tanstack/solid-query';
import { setToken, setUser } from './store';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: { staleTime: 1000 * 60 * 5, gcTime: 1000 * 60 * 10 },
  },
});

export default function App(props: any) {
  const [isAuthenticated, setIsAuthenticated] = createSignal(false);
  const [isLoading, setIsLoading] = createSignal(true);

  createEffect(async () => {
    const initData = window.Telegram.WebApp.initData;
    console.log('Init data:', initData);

    try {
      const resp = await fetch(
        'http://localhost:8080/auth/telegram?' + initData,
        { method: 'POST' },
      );

      const { user, token } = await resp.json();

      setUser(user);
      setToken(token);

      window.Telegram.WebApp.ready();
      window.Telegram.WebApp.expand();

      setIsAuthenticated(true);
      setIsLoading(false);
    } catch (e) {
      console.error('Failed to authenticate user:', e);
      setIsAuthenticated(false);
      setIsLoading(false);
    }
  });

  return (
    <QueryClientProvider client={queryClient}>
      <Switch>
        <Match when={isAuthenticated()}>
          <div>{props.children}</div>
        </Match>
        <Match when={!isAuthenticated() && isLoading()}>
          <div>Authenticating...</div>
        </Match>
        <Match when={!isAuthenticated() && !isLoading()}>
          <div>Failed to authenticate user</div>
        </Match>
      </Switch>
    </QueryClientProvider>
  );
}
