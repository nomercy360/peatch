import {
	createEffect,
	createSignal,
	For, Match,
	onCleanup,
	onMount,
	Show,
	Suspense, Switch,
} from 'solid-js'
import { Collaboration, Post, UserProfile } from '~/gen'
import { fetchFeed } from '~/lib/api'
import { Link } from '~/components/link'
import useDebounce from '~/lib/useDebounce'
import { useQuery } from '@tanstack/solid-query'
import { useMainButton } from '~/lib/useMainButton'
import { UserCardSmall } from '~/pages/posts/id'
import { LocationBadge } from '~/components/location-badge'
import { ListPlaceholder } from '~/pages/feed'
import { useTranslations } from '~/lib/locale-context'

export const [search, setSearch] = createSignal('')

export default function PostsPage() {
	const { t } = useTranslations()
	const updateSearch = useDebounce(setSearch, 350)

	const mainButton = useMainButton()

	const query = useQuery(() => ({
		queryKey: ['posts', search()],
		queryFn: () => fetchFeed(search()),
	}))

	const [scroll, setScroll] = createSignal(0)

	createEffect(() => {
		const onScroll = () => setScroll(window.scrollY)
		window.addEventListener('scroll', onScroll)
		return () => window.removeEventListener('scroll', onScroll)
	})

	onMount(() => {
		// disable scroll on body when drawer is open
		document.body.style.overflow = 'hidden'
	})

	onCleanup(() => {
		mainButton.hide()
		document.body.style.overflow = 'auto'
	})

	return (
		<div class="flex h-screen flex-col overflow-hidden">
			<div class="flex w-full flex-shrink-0 flex-col items-center justify-between space-y-4 border-b p-4">
				<div class="relative flex h-10 w-full flex-row items-center justify-center rounded-lg bg-secondary">
					<input
						class="h-full w-full bg-transparent px-2.5 placeholder:text-secondary-foreground"
						placeholder={t('common.search.posts')}
						type="text"
						value={search()}
						onInput={e => updateSearch(e.currentTarget.value)}
					/>
					<Show when={search()}>
						<button
							class="absolute right-2.5 flex size-5 shrink-0 items-center justify-center rounded-full bg-secondary"
							onClick={() => setSearch('')}
						>
							<span class="material-symbols-rounded text-[20px] text-secondary">
								close
							</span>
						</button>
					</Show>
				</div>
			</div>
			<div class="flex h-full w-full flex-shrink-0 flex-col overflow-y-auto pb-20">
				<Suspense fallback={<ListPlaceholder />}>
					<For each={query.data}>
						{(data, _) => (
							<div>
								<Switch fallback={<div />}>
									<Match when={data.type === 'collaboration'}>
										<CollaborationCard
											collab={data.data as Collaboration}
											scroll={scroll()}
										/>
									</Match>
									<Match when={data.type === 'post'}>
										<PostCard post={data.data as Post} />
									</Match>
								</Switch>
								<div class="h-px w-full bg-border" />
							</div>
						)}
					</For>
				</Suspense>
			</div>
		</div>
	)
}

const PostCard = (props: { post: Post }) => {
	return (
		<Link
			class="flex flex-col items-start px-4 pb-5 pt-4 text-start"
			href={`/posts/${props.post.id}`}
		>
			<UserCardSmall user={props.post.user as UserProfile} />
			<p class="mt-4 text-3xl">{props.post.title}</p>
			<p class="mt-1 text-sm text-secondary-foreground">{props.post.description}</p>
			<Show when={props.post.image_url}>
				<img
					class="mt-3 aspect-[4/3] w-full rounded-xl object-cover"
					src={props.post.image_url}
					alt="Post Image"
					loading="lazy"
				/>
			</Show>
			<div class="mt-3">
				<Show when={props.post.country && props.post.city}>
					<LocationBadge
						country={props.post.country!}
						city={props.post.city!}
						countryCode={props.post.country_code!}
					/>
				</Show>
			</div>
		</Link>
	)
}

const CollaborationCard = (props: {
	collab: Collaboration
	scroll: number
}) => {
	const shortenDescription = (description: string) => {
		if (description.length <= 160) return description
		return description.slice(0, 160) + '...'
	}

	return (
		<Link
			class="flex flex-col items-start px-4 pb-5 pt-4 text-start"
			href={`/collaborations/${props.collab.id}`}
			state={{ from: '/', scroll: props.scroll }}
		>
			<div class="flex flex-row items-center justify-center">
				<Link
					href={`/users/${props.collab.user?.username}`}
					state={{ back: true, scroll: props.scroll }}
				>
					<img
						class="size-10 rounded-xl object-cover"
						src={`https://assets.peatch.io/cdn-cgi/image/width=100/${props.collab.user?.avatar_url}`}
						alt="User Avatar"
					/>
				</Link>
				<div
					class="-ml-2 flex size-10 flex-row items-center justify-center rounded-full"
					style={{ 'background-color': `#${props.collab.opportunity?.color}` }}
				>
					<span class="material-symbols-rounded text-[20px] text-white">
						{String.fromCodePoint(
							parseInt(props.collab.opportunity?.icon!, 16),
						)}
					</span>
				</div>
			</div>
			<p
				class="mt-3 text-3xl font-semibold"
				style={{ color: `#${props.collab.opportunity?.color}` }}
			>
				{props.collab.opportunity?.text}:
			</p>
			<p class="text-3xl">{props.collab.title}</p>
			<p class="mt-2 text-sm text-secondary-foreground">
				{shortenDescription(props.collab.description!)}
			</p>
		</Link>
	)
}
