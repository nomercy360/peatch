import { createEffect, createSignal, Show } from 'solid-js'
import { QueryClient, QueryClientProvider } from '@tanstack/solid-query'
import { checkAuth } from './lib/api'

export const queryClient = new QueryClient({
	defaultOptions: {
		queries: {
			retry: 2,
			staleTime: 1000 * 60 * 5, // 5 minutes
			gcTime: 1000 * 60 * 5, // 5 minutes
		},
		mutations: {
			retry: 2,
		},
	},
})

export default function App(props: { children?: any }) {
	const [isLoading, setIsLoading] = createSignal(true)

	createEffect(async () => {
		// Skip auth check on login page
		if (window.location.pathname === '/login') {
			setIsLoading(false)
			return
		}

		try {
			const isValid = await checkAuth()
			setIsLoading(false)

			// Redirect to login if not authenticated
			if (!isValid) {
				window.location.href = '/login'
			}
		} catch (e) {
			console.error('Failed to check authentication:', e)
			setIsLoading(false)
			window.location.href = '/login'
		}
	})

	return (
		<QueryClientProvider client={queryClient}>
			<Show when={!isLoading()} fallback={
				<div class="min-h-screen w-full flex items-center justify-center">
					<div>Loading...</div>
				</div>
			}>
				{props.children}
			</Show>
		</QueryClientProvider>
	)
}
