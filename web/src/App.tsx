import { createEffect, createSignal, Match, onCleanup, Switch } from 'solid-js';
import { QueryClient, QueryClientProvider } from '@tanstack/solid-query';
import { setToken, setUser } from './store';
import { API_BASE_URL } from './api';
import { NavigationProvider } from './hooks/useNavigation';

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

    // test sleep
    try {
      const resp = await fetch(`${API_BASE_URL}/auth/telegram?` + initData, {
        method: 'POST',
      });

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
    <NavigationProvider>
      <QueryClientProvider client={queryClient}>
        <Switch>
          <Match when={isAuthenticated()}>
            <div>{props.children}</div>
          </Match>
          <Match when={!isAuthenticated() && isLoading()}>
            <div
              class="h-screen w-full items-start flex-col justify-center bg-main">
            </div>
          </Match>
          <Match when={!isAuthenticated() && !isLoading()}>
            <div
              class="min-h-screen h-screen w-full items-start flex-col justify-center text-3xl bg-main text-main">
            Something went wrong. Please try again later.
            </div>
          </Match>
        </Switch>
      </QueryClientProvider>
    </NavigationProvider>
  );
}
