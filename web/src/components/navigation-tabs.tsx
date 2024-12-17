import { cn } from '~/lib/utils'
import { useLocation } from '@solidjs/router'
import { Link } from '~/components/link'
import { store } from '~/store'

export default function NavigationTabs(props: any) {
	const location = useLocation()

	const tabs = [
		{
			href: '/posts',
			icon: 'workspace_premium',
			name: 'Opportunities',
		},
		{
			href: '/home',
			icon: 'home',
			name: 'Home',
		},
		{
			href: '/',
			icon: 'groups_3',
			name: 'People',
		},

		{
			href: '/friends',
			icon: 'group',
			name: 'Friends',
		},
	]

	return (
		<>
			<div
				class="grid grid-cols-5 items-center border shadow-sm h-[72px] pb-2 fixed bottom-0 w-full border-t bg-background z-50"
			>
				<Link
					href={`/users/${store.user?.username}`}
					state={{ from: location.pathname }}
					class="flex items-center justify-center"
				>
					<img
						src={`https://assets.peatch.io/${store.user?.avatar_url}`}
						alt="User avatar"
						class="size-9 rounded-full object-cover"
					/>
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
						<span class="text-xs">{name}</span>
					</Link>
				))}
			</div>
			{props.children}
		</>
	)
}
