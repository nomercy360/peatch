import {
	createEffect,
	createSignal,
	Match,
	onCleanup,
	Show,
	Suspense,
	Switch,
} from 'solid-js'
import { useNavigate, useParams } from '@solidjs/router'
import { fetchPost } from '~/lib/api'
import { store } from '~/store'
import { createQuery } from '@tanstack/solid-query'
import { useMainButton } from '~/lib/useMainButton'
import { Link } from '~/components/link'
import { UserProfile } from '~/gen/types'

export default function Collaboration() {
	const [isCurrentUserCollab, setIsCurrentUserCollab] = createSignal(false)

	const navigate = useNavigate()
	const params = useParams()
	const postId = params.id

	const mainButton = useMainButton()

	const query = createQuery(() => ({
		queryKey: ['posts', postId],
		queryFn: () => fetchPost(Number(postId)),
	}))

	createEffect(() => {
		if (query.data?.id) {
			setIsCurrentUserCollab(store.user.id === query.data.user.id)
		}
	})

	const navigateToEdit = () => {
		navigate('/posts/edit/' + postId, {
			state: { from: '/posts/' + postId },
		})
	}

	createEffect(() => {
		if (isCurrentUserCollab() && !query.data.published_at) {
			mainButton.enable('Edit').onClick(navigateToEdit)
		}
	})

	onCleanup(() => {
		mainButton.offClick(navigateToEdit)
	})

	return (
		<Suspense fallback={<Loader />}>
			<Switch>
				<Match when={!query.isLoading}>
					<div class="h-fit min-h-screen bg-secondary">
						<div class="flex w-full flex-col items-start justify-start bg-main px-4 py-4">
							<Show when={query.data.image_url}>
								<img
									class="aspect-[4/3] w-full rounded-xl object-cover"
									src={query.data.image_url}
									alt="Collaboration Image"
								/>
							</Show>
							<p class="mb-5 mt-4 text-3xl text-main">{query.data.title}</p>
							<UserCardSmall user={query.data.user} />
						</div>
						<div class="px-4 py-2.5">
							<p class="text-lg font-normal text-main">
								{query.data.description}
							</p>
						</div>
					</div>
				</Match>
			</Switch>
		</Suspense>
	)
}

export const UserCardSmall = (props: { user: UserProfile }) => {
	return (
		<Link
			class="flex w-full flex-row items-center justify-start gap-2"
			href={'/users/' + props.user.username}
			state={{ back: true }}
		>
			<img
				class="size-10 rounded-xl object-cover"
				src={`https://assets.peatch.io/cdn-cgi/image/width=100/${props.user.avatar_url}`}
				alt="User Avatar"
			/>
			<div>
				<p class="text-sm font-bold text-main">
					{props.user.first_name} {props.user.last_name}
				</p>
				<p class="text-sm text-secondary">{props.user.title}</p>
			</div>
		</Link>
	)
}
const Loader = () => {
	return (
		<div class="flex h-screen flex-col items-start justify-start bg-secondary">
			<div class="h-[260px] w-full bg-main" />
			<div class="flex flex-col items-start justify-start p-4">
				<div class="h-36 w-full rounded bg-main" />
				<div class="mt-4 flex w-full flex-row flex-wrap items-center justify-start gap-2">
					<div class="h-10 w-40 rounded-2xl bg-main" />
					<div class="h-10 w-32 rounded-2xl bg-main" />
					<div class="h-10 w-36 rounded-2xl bg-main" />
					<div class="h-10 w-24 rounded-2xl bg-main" />
					<div class="h-10 w-40 rounded-2xl bg-main" />
					<div class="h-10 w-28 rounded-2xl bg-main" />
					<div class="h-10 w-32 rounded-2xl bg-main" />
				</div>
			</div>
		</div>
	)
}
