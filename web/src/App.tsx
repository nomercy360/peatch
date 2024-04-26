import { createEffect } from 'solid-js';
import { QueryClient, QueryClientProvider } from '@tanstack/solid-query';

const queryClient = new QueryClient();

export default function App(props: any) {
  createEffect(() => {
    window.Telegram.WebApp.ready();
    window.Telegram.WebApp.expand();

    console.log(window.Telegram.WebApp.initData);
  });

  return (
    <QueryClientProvider client={queryClient}>
      <>{props.children}</>
    </QueryClientProvider>
  );
}
