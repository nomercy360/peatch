import { For, Show } from 'solid-js'
import Badge from '~/components/badge'
import { BadgeResponse } from '~/gen'

export default function BadgeList(props: {
	badges: BadgeResponse[]
	position: 'center' | 'start'
	children?: any
}) {
	const badgeSlice = props.badges!.slice(0, 5)

	return (
		<div
			class="mt-2 flex w-full flex-row flex-wrap items-center justify-start gap-1"
			classList={{
				'justify-center': props.position === 'center',
			}}
		>
			{props.children}
			<For each={badgeSlice}>
				{badge => (
					<Badge icon={badge.icon!} name={badge.text!} color={badge.color!} />
				)}
			</For>
			<Show when={props.badges.length > 5}>
				<div
					class="flex h-5 flex-row items-center justify-center rounded px-1.5 text-xs font-semibold text-muted-foreground">
					+ {props.badges!.length - 5} more
				</div>
			</Show>
		</div>
	)
}
