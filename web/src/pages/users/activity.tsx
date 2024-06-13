import { For, Match, Show, Switch } from 'solid-js'
import { fetchActivity } from '~/lib/api'
import { createQuery } from '@tanstack/solid-query'
import { Link } from '~/components/Link'

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
			<span class="material-symbols-rounded text-[48px] text-orange">
				notifications
			</span>
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
			return <span>'sent you a collaboration request'</span>
	}
}

function ActivityCard(props: { activity: any }) {
	return (
		<div class="flex flex-row items-center justify-start gap-4">
			<div class="flex-shrink-0">
				<div class="flex size-10 items-center justify-center rounded-full bg-main text-main">
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
				<div class="text-sm text-hint">
					{timeSince(props.activity.timestamp)}
				</div>
			</div>
		</div>
	)
}

const Loader = () => {
	return (
		<div class="flex h-screen flex-col items-start justify-start bg-secondary p-2">
			<div class="aspect-square w-full rounded-xl bg-main" />
			<div class="flex flex-col items-start justify-start p-2">
				<div class="mt-2 h-6 w-full rounded bg-main" />
				<div class="mt-2 h-6 w-full rounded bg-main" />
			</div>
		</div>
	)
}
