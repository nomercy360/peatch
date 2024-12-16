import { cn } from '~/lib/utils'
import { useLocation } from '@solidjs/router'
import { Link } from '~/components/Link'

export default function NavigationTabs() {
	const location = useLocation()

	const tabs = [
		{
			href: '/',
			icon: 'groups_3',
			activePath: '/',
			name: 'People',
		},
		{
			href: '/posts',
			icon: 'workspace_premium',
			activePath: '/posts',
			name: 'Opportunities',
		},
		{
			href: '/',
			icon: 'home',
			activePath: '/',
			name: 'Home',
		},
		{
			href: '/friends',
			icon: 'group',
			activePath: '/friends',
			name: 'Friends',
		},
		{
			href: '/profile',
			icon: 'person',
			activePath: '/profile',
			name: 'Profile',
		},
	]

	return (
		<div
			class="grid grid-cols-5 items-center space-x-4 border shadow-sm h-[72px] fixed bottom-0 w-full border-t bg-background z-50"
		>
			{tabs.map(({ href, icon, activePath, name }) => (
				<Link
					href={href}
					class={cn('h-12 flex items-center justify-between flex-col text-sm text-secondary-foreground', {
						'text-foreground': location.pathname === activePath,
					})}
				>
					<span class="material-symbols-rounded text-[32px]">{icon}</span>
					<span class="text-xs">{name}</span>
				</Link>
			))}
		</div>
	)
}
