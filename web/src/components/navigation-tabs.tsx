import { cn } from '~/lib/utils'
import { useLocation } from '@solidjs/router'
import { Link } from '~/components/link'
import { store } from '~/store'
import { Avatar, AvatarFallback, AvatarImage } from '~/components/ui/avatar'

export default function NavigationTabs(props: any) {
	const location = useLocation()

	const tabs = [
		{
			href: '/posts',
			icon: 'workspace_premium',
			name: 'Opportunities',
		},
		{
			href: '/',
			icon: 'groups_3',
			name: 'People',
		},
	]

	function getUserInitials() {
		const firstInitial = store.user?.first_name ? store.user?.first_name[0] : ''
		const lastInitial = store.user?.last_name ? store.user?.last_name[0] : ''

		return `${firstInitial}${lastInitial}`
	}

	return (
		<>
			<div
				class="grid grid-cols-3 items-start border shadow-sm h-[100px] py-4 fixed bottom-0 w-full border-t bg-background z-50"
			>
				<Link
					href={`/users/${store.user?.username}`}
					state={{ from: location.pathname }}
					class="flex items-center justify-center"
				>
					<Avatar>
						<AvatarImage src={`https://assets.peatch.io/${store.user?.avatar_url}`} />
						<AvatarFallback>
							{getUserInitials()}
						</AvatarFallback>
					</Avatar>
				</Link>
				{tabs.map(({ href, icon, name }) => (
					<Link
						href={href}
						state={{ from: location.pathname }}
						class={cn('h-12 flex items-center justify-between flex-col text-sm text-secondary-foreground', {
							'text-foreground': location.pathname === href,
						})}
					>
						<span class="material-symbols-rounded text-[32px]">{icon}</span>
						<span class="text-xs font-medium">{name}</span>
					</Link>
				))}
			</div>
			{props.children}
		</>
	)
}
