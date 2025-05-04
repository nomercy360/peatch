import { createEffect, createSignal, Match, Switch } from 'solid-js'
import { setToken, setUser } from './store'
import { API_BASE_URL } from '~/lib/api'
import { NavigationProvider } from './lib/useNavigation'
import { useNavigate } from '@solidjs/router'
import { QueryClient, QueryClientProvider } from '@tanstack/solid-query'
import Toast from '~/components/toast'
import { LocaleContextProvider } from '~/lib/locale-context'

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

function transformStartParam(startParam?: string): string | null {
	if (!startParam) return null;

	if (startParam.startsWith('u_')) {
		const id = startParam.slice('u_'.length);
		return `/users/${id}`;
	} else if (startParam.startsWith('c_')) {
		const id = startParam.slice('c_'.length);
		return `/collaborations/${id}`;
	} else {
		return null;
	}
}

export default function App(props: any) {
	const [isAuthenticated, setIsAuthenticated] = createSignal(false)
	const [isLoading, setIsLoading] = createSignal(true)

	const navigate = useNavigate()

	createEffect(() => {
		const authenticate = async () => {
			const initData = window.Telegram.WebApp.initData

			console.log('WEBAPP:', window.Telegram)

			try {
				const resp = await fetch(`${API_BASE_URL}/auth/telegram`, {
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
					},
					body: JSON.stringify({ query: initData }),
				})

				const { user, token } = await resp.json()

				setUser(user)
				setToken(token)

				window.Telegram.WebApp.ready()
				window.Telegram.WebApp.expand()
				window.Telegram.WebApp.disableVerticalSwipes()

				setIsAuthenticated(true)
				setIsLoading(false)

				const startapp = window.Telegram.WebApp.initDataUnsafe.start_param
				const redirectUrl = transformStartParam(startapp)

				if (redirectUrl) {
					navigate(redirectUrl)
					return
				}
			} catch (e) {
				console.error('Failed to authenticate user:', e)
				setIsAuthenticated(false)
				setIsLoading(false)
			}
		}

		authenticate()
	})

	return (
		<LocaleContextProvider>
			<NavigationProvider>
				<QueryClientProvider client={queryClient}>
					<Switch>
						<Match when={isAuthenticated()}>
							<div>{props.children}</div>
						</Match>
						<Match when={!isAuthenticated() && isLoading()}>
							<div class="h-screen w-full flex-col items-start justify-center" />
						</Match>
						<Match when={!isAuthenticated() && !isLoading()}>
							<div class="text-main h-screen min-h-screen w-full flex-col items-start justify-center text-3xl">
								Something went wrong. Please try again later.
							</div>
						</Match>
					</Switch>
					<Toast />
				</QueryClientProvider>
			</NavigationProvider>
		</LocaleContextProvider>
	)
}
