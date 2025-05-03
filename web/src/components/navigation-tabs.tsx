import { cn } from '~/lib/utils'
import { useLocation } from '@solidjs/router'
import { Link } from '~/components/link'
import { store } from '~/store'
import { useMainButton } from '~/lib/useMainButton'
import { onMount } from 'solid-js'
import { useTranslations } from '~/lib/locale-context'

export default function NavigationTabs(props: any) {
	const location = useLocation()
	const mainButton = useMainButton()
	const { t } = useTranslations()

	const tabs = [
		{
			href: '/posts',
			icon: 'local_fire_department',
			name: t('common.tabs.posts'),
		},
		{
			href: '/',
			icon: 'group',
			name: t('common.tabs.network'),
		},
		{
			href: '/collaborations/edit',
			icon: 'edit_note',
			name: t('common.tabs.collaborations'),
		},
	]

	onMount(() => {
		mainButton.hide()
	})

	return (
		<>
			<div
				class="fixed bottom-0 z-50 flex h-[100px] w-full items-center justify-between space-x-10 border border-t bg-background px-5 shadow-sm">
				<div class="flex items-center">
					<Link
						href={store.user.first_name && store.user.description ? `/users/${store.user?.username}` : '/users/edit'}
						state={{ from: location.pathname }}
						class="flex items-center justify-center"
					>
						<img
							src={`https://assets.peatch.io/cdn-cgi/image/width=100/${store.user?.avatar_url}`}
							alt="User Avatar"
							class="h-10 w-10 rounded-full object-cover"
							onError={(e) => {
								const target = e.target as HTMLImageElement
								target.src = '/fallback-avatar.svg'
							}}
						/>
					</Link>
				</div>
				<div class="flex flex-1 justify-evenly">
					{tabs.map((props) => (
						<Link
							href={props.href}
							state={{ from: location.pathname }}
							class={cn('flex h-12 flex-col items-center justify-between text-sm text-secondary-foreground', {
								'text-foreground': location.pathname === props.href,
							})}
						>
							<span class="material-symbols-rounded text-[32px]">{props.icon}</span>
							<span class="text-xs font-medium">{props.name}</span>
						</Link>
					))}
				</div>
			</div>
			{props.children}
		</>
	)
}
