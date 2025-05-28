import { For } from 'solid-js'
import type { ComponentProps } from 'solid-js'
import { A } from '@solidjs/router'
import { IconUsers, IconUserCog, IconTarget, IconHandshake, IconMapPin, IconAward } from '~/components/icons'


const navItems = [
	{
		title: 'Users',
		url: '/users',
		icon: IconUsers,
	},
	{
		title: 'Badges',
		url: '/badges',
		icon: IconAward,
	},
	{
		title: 'Cities',
		url: '/cities',
		icon: IconMapPin,
	},
	{
		title: 'Collaborations',
		url: '/collaborations',
		icon: IconHandshake,
	},
	{
		title: 'Opportunities',
		url: '/opportunities',
		icon: IconTarget,
	},
	{
		title: 'Admins',
		url: '/admins',
		icon: IconUserCog,
	},
]

export function AppSidebar(props: ComponentProps<'nav'>) {
	return (
		<nav class="w-16 h-full bg-background text-foreground border-r" {...props}>
			<div class="p-2">
				<ul class="space-y-2">
					<For each={navItems}>
						{(item) => (
							<li>
								<A
									href={item.url}
									class="size-10 flex items-center justify-center rounded p-2 hover:bg-accent hover:text-accent-foreground transition-colors"
								>
									<item.icon class="size-4" />
								</A>
							</li>
						)}
					</For>
				</ul>
			</div>
		</nav>
	)
}
