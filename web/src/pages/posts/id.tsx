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
import { useQuery } from '@tanstack/solid-query'
import { useMainButton } from '~/lib/useMainButton'
import { Link } from '~/components/link'
import { UserProfile } from '~/gen'

export default function Collaboration() {
	const [isCurrentUserCollab, setIsCurrentUserCollab] = createSignal(false)

	const navigate = useNavigate()
	const params = useParams()
	const postId = params.id

	const mainButton = useMainButton()

	const query = useQuery(() => ({
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
						<div class="bg-main flex w-full flex-col items-start justify-start px-4 py-4">
							<Show when={query.data.image_url}>
								<img
									class="aspect-[4/3] w-full rounded-xl object-cover"
									src={query.data.image_url}
									alt="Collaboration Image"
								/>
							</Show>
							<p class="text-main mb-5 mt-4 text-3xl">{query.data.title}</p>
							<UserCardSmall user={query.data.user} />
						</div>
						<div class="px-4 py-2.5">
							<p class="text-main text-lg font-normal">
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
				<p class="text-main text-sm font-bold">
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
			<div class="bg-main h-[260px] w-full" />
			<div class="flex flex-col items-start justify-start p-4">
				<div class="bg-main h-36 w-full rounded" />
				<div class="mt-4 flex w-full flex-row flex-wrap items-center justify-start gap-2">
					<div class="bg-main h-10 w-40 rounded-2xl" />
					<div class="bg-main h-10 w-32 rounded-2xl" />
					<div class="bg-main h-10 w-36 rounded-2xl" />
					<div class="bg-main h-10 w-24 rounded-2xl" />
					<div class="bg-main h-10 w-40 rounded-2xl" />
					<div class="bg-main h-10 w-28 rounded-2xl" />
					<div class="bg-main h-10 w-32 rounded-2xl" />
				</div>
			</div>
		</div>
	)
}
