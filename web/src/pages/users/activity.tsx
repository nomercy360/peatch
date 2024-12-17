import { For, Match, Show, Switch } from 'solid-js'
import { fetchActivity } from '~/lib/api'
import { createQuery } from '@tanstack/solid-query'
import { Link } from '~/components/link'

export function timeSince(dateString: string) {
	const date = new Date(dateString)
	const now = new Date()
	//@ts-ignore
	const seconds = Math.floor((now - date) / 1000)
	const minutes = Math.floor(seconds / 60)
	const hours = Math.floor(minutes / 60)
	const days = Math.floor(hours / 24)
	const months = Math.floor(days / 30)
	const years = Math.floor(days / 365)

	if (years > 0) {
		return years + 'y ago'
	} else if (months > 0) {
		return months + 'mo ago'
	} else if (days > 0) {
		return days + 'd ago'
	} else if (hours > 0) {
		return hours + 'h ago'
	} else if (minutes > 0) {
		return minutes + 'm ago'
	} else if (seconds <= 10) {
		return 'just now'
	} else {
		return seconds + 's ago'
	}
}

export default function UserProfilePage() {
	const query = createQuery(() => ({
		queryKey: ['activity'],
		queryFn: () => fetchActivity(),
		retry: 1,
	}))

	return (
		<div class="flex flex-col items-center justify-start px-4 py-5">
			<NotificationIcon />
			<h1 class="mt-2 max-w-[285px] text-3xl text-main">Notifications</h1>
			<Switch>
				<Match when={query.isLoading}>
					<Loader />
				</Match>
				<Match when={query.isSuccess}>
					<div class="mt-8 grid gap-4">
						<For each={query.data}>
							{activity => <ActivityCard activity={activity} />}
						</For>
					</div>
				</Match>
			</Switch>
		</div>
	)
}

function ResolveText(props: {
	type: string
	content: string | null
	contentID: number | null
}) {
	switch (props.type) {
		case 'post_like':
			return (
				<span>
					liked your post{' '}
					<Link class="text-main" href={`/posts/${props.contentID}`}>
						{props.content}
					</Link>
				</span>
			)
		case 'collab_like':
			return (
				<span>
					liked your collaboration{' '}
					<Link class="text-main" href={`/collaborations/${props.contentID}`}>
						{props.content}
					</Link>
				</span>
			)
		case 'user_like':
			return <span>liked your profile</span>
		case 'collab_request':
			return (
				<span>
					responded to your collaboration{' '}
					<Link class="text-main" href={`/collaborations/${props.contentID}`}>
						{props.content}
					</Link>
				</span>
			)
		case 'follow':
			return <span>started following you</span>
		case 'user_collab':
			return <span>sent you a collaboration request</span>
	}
}

function ActivityCard(props: { activity: any }) {
	return (
		<div class="flex flex-row items-center justify-start gap-4">
			<div class="flex-shrink-0">
				<div class="flex size-10 items-center justify-center rounded-full">
					<Switch>
						<Match
							when={
								props.activity.activity_type === 'post_like' ||
								props.activity.activity_type === 'collab_like' ||
								props.activity.activity_type === 'user_like'
							}
						>
							<span class="no-fill material-symbols-rounded">favorite</span>
						</Match>
						<Match
							when={
								props.activity.activity_type === 'collab_request' ||
								props.activity.activity_type === 'user_collab'
							}
						>
							<span class="no-fill material-symbols-rounded">send</span>
						</Match>
						<Match when={props.activity.activity_type === 'follow'}>
							<span class="no-fill material-symbols-rounded">person_add</span>
						</Match>
					</Switch>
				</div>
			</div>
			<div class="flex-1">
				<div class="mt-1 text-sm text-hint">
					<Link
						href={`/users/${props.activity.actor_username}`}
						state={{ from: '/users/activity' }}
						class="font-medium text-main"
					>
						@{props.activity.actor_username}{' '}
					</Link>
					<ResolveText
						type={props.activity.activity_type}
						content={props.activity.content}
						contentID={props.activity.content_id}
					/>
					<Show when={props.activity.message}>
						<span> "{props.activity.message}"</span>
					</Show>
				</div>
				<div class="text-xs text-secondary-foreground">
					{timeSince(props.activity.timestamp)}
				</div>
			</div>
		</div>
	)
}

const Loader = () => {
	return (
		<div class="w-full flex h-screen flex-col items-start justify-start p-2">
			<div class="mt-2 h-12 w-full rounded bg-secondary" />
			<div class="mt-2 h-12 w-full rounded bg-secondary" />
			<div class="mt-2 h-12 w-full rounded bg-secondary" />
			<div class="mt-2 h-12 w-full rounded bg-secondary" />
			<div class="mt-2 h-12 w-full rounded bg-secondary" />
			<div class="mt-2 h-12 w-full rounded bg-secondary" />
		</div>
	)
}

export const NotificationIcon = (props: any) => {
	return (
		<svg
			width="48"
			height="46"
			viewBox="0 0 48 46"
			fill="none"
			xmlns="http://www.w3.org/2000/svg"
			{...props}
		>
			<mask
				id="mask0_4780_2058"
				style="mask-type:alpha"
				maskUnits="userSpaceOnUse"
				x="0"
				y="0"
				width="48"
				height="46"
			>
				<rect width="48" height="46" fill="#D9D9D9" />
			</mask>
			<g mask="url(#mask0_4780_2058)">
				<path
					d="M11.3996 35.7106C10.8896 35.7106 10.4621 35.5366 10.1171 35.1885C9.77211 34.8405 9.59961 34.4093 9.59961 33.8948C9.59961 33.3803 9.77211 32.9491 10.1171 32.6011C10.4621 32.253 10.8896 32.079 11.3996 32.079H13.1996V19.3685C13.1996 16.8566 13.9496 14.6247 15.4496 12.6728C16.9496 10.7208 18.8996 9.44216 21.2996 8.8369V7.56585C21.2996 6.80927 21.5621 6.16618 22.0871 5.63657C22.6121 5.10697 23.2496 4.84216 23.9996 4.84216C24.7496 4.84216 25.3871 5.10697 25.9121 5.63657C26.4371 6.16618 26.6996 6.80927 26.6996 7.56585V8.8369C29.0996 9.44216 31.0496 10.7208 32.5496 12.6728C34.0496 14.6247 34.7996 16.8566 34.7996 19.3685V32.079H36.5996C37.1096 32.079 37.5371 32.253 37.8821 32.6011C38.2271 32.9491 38.3996 33.3803 38.3996 33.8948C38.3996 34.4093 38.2271 34.8405 37.8821 35.1885C37.5371 35.5366 37.1096 35.7106 36.5996 35.7106H11.3996ZM23.9996 41.158C23.0096 41.158 22.1621 40.8024 21.4571 40.0912C20.7521 39.38 20.3996 38.5251 20.3996 37.5264H27.5996C27.5996 38.5251 27.2471 39.38 26.5421 40.0912C25.8371 40.8024 24.9896 41.158 23.9996 41.158Z"
					fill="#FF8C42"
				/>
				<path
					d="M11.3996 35.7106C10.8896 35.7106 10.4621 35.5366 10.1171 35.1885C9.77211 34.8405 9.59961 34.4093 9.59961 33.8948C9.59961 33.3803 9.77211 32.9491 10.1171 32.6011C10.4621 32.253 10.8896 32.079 11.3996 32.079H13.1996V19.3685C13.1996 16.8566 13.9496 14.6247 15.4496 12.6728C16.9496 10.7208 18.8996 9.44216 21.2996 8.8369V7.56585C21.2996 6.80927 21.5621 6.16618 22.0871 5.63657C22.6121 5.10697 23.2496 4.84216 23.9996 4.84216C24.7496 4.84216 25.3871 5.10697 25.9121 5.63657C26.4371 6.16618 26.6996 6.80927 26.6996 7.56585V8.8369C29.0996 9.44216 31.0496 10.7208 32.5496 12.6728C34.0496 14.6247 34.7996 16.8566 34.7996 19.3685V32.079H36.5996C37.1096 32.079 37.5371 32.253 37.8821 32.6011C38.2271 32.9491 38.3996 33.3803 38.3996 33.8948C38.3996 34.4093 38.2271 34.8405 37.8821 35.1885C37.5371 35.5366 37.1096 35.7106 36.5996 35.7106H11.3996ZM23.9996 41.158C23.0096 41.158 22.1621 40.8024 21.4571 40.0912C20.7521 39.38 20.3996 38.5251 20.3996 37.5264H27.5996C27.5996 38.5251 27.2471 39.38 26.5421 40.0912C25.8371 40.8024 24.9896 41.158 23.9996 41.158Z"
					fill="url(#paint0_radial_4780_2058)"
				/>
				<path
					d="M11.3996 35.7106C10.8896 35.7106 10.4621 35.5366 10.1171 35.1885C9.77211 34.8405 9.59961 34.4093 9.59961 33.8948C9.59961 33.3803 9.77211 32.9491 10.1171 32.6011C10.4621 32.253 10.8896 32.079 11.3996 32.079H13.1996V19.3685C13.1996 16.8566 13.9496 14.6247 15.4496 12.6728C16.9496 10.7208 18.8996 9.44216 21.2996 8.8369V7.56585C21.2996 6.80927 21.5621 6.16618 22.0871 5.63657C22.6121 5.10697 23.2496 4.84216 23.9996 4.84216C24.7496 4.84216 25.3871 5.10697 25.9121 5.63657C26.4371 6.16618 26.6996 6.80927 26.6996 7.56585V8.8369C29.0996 9.44216 31.0496 10.7208 32.5496 12.6728C34.0496 14.6247 34.7996 16.8566 34.7996 19.3685V32.079H36.5996C37.1096 32.079 37.5371 32.253 37.8821 32.6011C38.2271 32.9491 38.3996 33.3803 38.3996 33.8948C38.3996 34.4093 38.2271 34.8405 37.8821 35.1885C37.5371 35.5366 37.1096 35.7106 36.5996 35.7106H11.3996ZM23.9996 41.158C23.0096 41.158 22.1621 40.8024 21.4571 40.0912C20.7521 39.38 20.3996 38.5251 20.3996 37.5264H27.5996C27.5996 38.5251 27.2471 39.38 26.5421 40.0912C25.8371 40.8024 24.9896 41.158 23.9996 41.158Z"
					fill="url(#paint1_radial_4780_2058)"
				/>
				<path
					d="M11.3996 35.7106C10.8896 35.7106 10.4621 35.5366 10.1171 35.1885C9.77211 34.8405 9.59961 34.4093 9.59961 33.8948C9.59961 33.3803 9.77211 32.9491 10.1171 32.6011C10.4621 32.253 10.8896 32.079 11.3996 32.079H13.1996V19.3685C13.1996 16.8566 13.9496 14.6247 15.4496 12.6728C16.9496 10.7208 18.8996 9.44216 21.2996 8.8369V7.56585C21.2996 6.80927 21.5621 6.16618 22.0871 5.63657C22.6121 5.10697 23.2496 4.84216 23.9996 4.84216C24.7496 4.84216 25.3871 5.10697 25.9121 5.63657C26.4371 6.16618 26.6996 6.80927 26.6996 7.56585V8.8369C29.0996 9.44216 31.0496 10.7208 32.5496 12.6728C34.0496 14.6247 34.7996 16.8566 34.7996 19.3685V32.079H36.5996C37.1096 32.079 37.5371 32.253 37.8821 32.6011C38.2271 32.9491 38.3996 33.3803 38.3996 33.8948C38.3996 34.4093 38.2271 34.8405 37.8821 35.1885C37.5371 35.5366 37.1096 35.7106 36.5996 35.7106H11.3996ZM23.9996 41.158C23.0096 41.158 22.1621 40.8024 21.4571 40.0912C20.7521 39.38 20.3996 38.5251 20.3996 37.5264H27.5996C27.5996 38.5251 27.2471 39.38 26.5421 40.0912C25.8371 40.8024 24.9896 41.158 23.9996 41.158Z"
					fill="url(#paint2_radial_4780_2058)"
					fill-opacity="0.6"
				/>
			</g>
			<defs>
				<radialGradient
					id="paint0_radial_4780_2058"
					cx="0"
					cy="0"
					r="1"
					gradientUnits="userSpaceOnUse"
					gradientTransform="translate(33.8739 34.4136) rotate(-155.082) scale(17.239 20.3059)"
				>
					<stop stop-color="#F35D28" />
					<stop offset="1" stop-color="#F35D28" stop-opacity="0" />
				</radialGradient>
				<radialGradient
					id="paint1_radial_4780_2058"
					cx="0"
					cy="0"
					r="1"
					gradientUnits="userSpaceOnUse"
					gradientTransform="translate(17.211 41.158) rotate(-64.5537) scale(24.4177 21.4736)"
				>
					<stop stop-color="#FFD67E" />
					<stop offset="1" stop-color="#FFD77F" stop-opacity="0" />
				</radialGradient>
				<radialGradient
					id="paint2_radial_4780_2058"
					cx="0"
					cy="0"
					r="1"
					gradientUnits="userSpaceOnUse"
					gradientTransform="translate(28.3196 17.8121) rotate(98.3114) scale(29.8853 19.1519)"
				>
					<stop stop-color="white" />
					<stop offset="0.489583" stop-color="white" stop-opacity="0" />
				</radialGradient>
			</defs>
		</svg>
	)
}
