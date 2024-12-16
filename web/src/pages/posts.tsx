import {
	createEffect,
	createSignal,
	For, Match,
	onCleanup,
	onMount,
	Show,
	Suspense, Switch,
} from 'solid-js'
import { Collaboration, Post, UserProfile } from '~/gen/types'
import { CDN_URL, fetchFeed, fetchUsers, likeContent, unlikeContent } from '~/lib/api'
import { Link } from '~/components/Link'
import useDebounce from '~/lib/useDebounce'
import { createMutation, createQuery } from '@tanstack/solid-query'
import { useMainButton } from '~/lib/useMainButton'
import { useNavigate } from '@solidjs/router'
import { UserCardSmall } from '~/pages/posts/id'
import { LocationBadge } from '~/components/location-badge'
import { queryClient } from '~/App'
import { HeartIcon, ListPlaceholder } from '~/pages/feed'

export const [search, setSearch] = createSignal('')

export default function PostsPage() {
	const updateSearch = useDebounce(setSearch, 350)

	const mainButton = useMainButton()
	const navigate = useNavigate()

	const query = createQuery(() => ({
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
		<div class="flex h-screen flex-col">
			<div class="flex w-full flex-shrink-0 flex-col items-center justify-between space-y-4 border-b p-4">
				<div class="relative flex h-10 w-full flex-row items-center justify-center rounded-lg bg-secondary">
					<input
						class="h-full w-full bg-transparent px-2.5 placeholder:text-secondary-foreground"
						placeholder="Search people or collaborations"
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
			<div class="bg-secondary flex h-full w-full flex-shrink-0 flex-col overflow-y-auto pb-20">
				<Suspense fallback={<ListPlaceholder />}>
					<For each={query.data}>
						{(data, i) => (
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
			<LikeButton
				liked={props.post.is_liked!}
				likes={props.post.likes_count!}
				id={props.post.id!}
				type="post"
			/>
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
						src={CDN_URL + '/' + props.collab.user?.avatar_url}
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
			<LikeButton
				liked={props.collab.is_liked!}
				likes={props.collab.likes_count!}
				id={props.collab.id!}
				type="collaboration"
			/>
		</Link>
	)
}

const LikeButton = (props: {
	id: number
	liked: boolean
	likes: number
	type: 'user' | 'collaboration' | 'post'
}) => {
	const handleMutate = async (userId: number) => {
		await queryClient.cancelQueries({ type: 'active' })
		queryClient.setQueryData(['users', search()], (old: any[]) =>
			old.map(item => {
				if (item.id === userId) {
					return {
						...item,
						is_liked: !item.is_liked,
						likes_count: item.is_liked
							? item.likes_count - 1
							: item.likes_count + 1,
					}
				}
				return item
			}),
		)
		if (search()) {
			queryClient.invalidateQueries({ queryKey: ['users', ''] })
		}
	}

	const likeMutate = createMutation(() => ({
		mutationFn: (id: number) => likeContent(id, props.type),
		onMutate: id => handleMutate(id),
	}))

	const mutateUnLike = createMutation(() => ({
		mutationFn: (id: number) => unlikeContent(id, props.type),
		onMutate: id => handleMutate(id),
	}))

	const handleClick = (e: Event) => {
		e.preventDefault()
		if (!props.liked) {
			likeMutate.mutate(props.id)
		} else {
			mutateUnLike.mutate(props.id)
		}
		window.Telegram.WebApp.HapticFeedback.selectionChanged()
	}

	return (
		<button
			class="mt-2 flex items-center justify-start rounded-xl text-sm font-semibold"
			onClick={(e: Event) => handleClick(e)}
		>
			<Show
				when={!props.liked}
				fallback={<HeartIcon class="size-6 shrink-0" />}
			>
				<span class="material-symbols-rounded no-fill text-[24px]">
					favorite
				</span>
			</Show>
			<Show when={props.likes > 0}>
				<span class="ml-1 font-semibold">{props.likes}</span>
			</Show>
		</button>
	)
}
