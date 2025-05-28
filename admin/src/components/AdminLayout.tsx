import { RouteSectionProps } from '@solidjs/router'
import { AppSidebar } from '~/components/AppSidebar'


export function AdminLayout(props: RouteSectionProps) {
	return (
		<div class='flex flex-row h-screen'>
			<AppSidebar />
			<main class="flex-1 overflow-auto">
				<div class="p-6">
					{props.children}
				</div>
			</main>
		</div>
	)
}
