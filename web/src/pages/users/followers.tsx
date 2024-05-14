import { createSignal, For, Match, Show, Switch } from 'solid-js'
import { useParams } from '@solidjs/router'
import {
	CDN_URL,
	fetchFollowers,
	fetchFollowing,
	followUser,
	unfollowUser,
} from '~/lib/api'
import { createMutation, createQuery } from '@tanstack/solid-query'
import { queryClient } from '~/App'
import { Link } from '~/components/Link'
import { store } from '~/store'
import { UserProfileShort } from '~/gen/types'

export default function UserFollowersPage() {
	const [showFollowing, setShowFollowing] = createSignal(true)

	const params = useParams()

	const following = createQuery(() => ({
		queryKey: ['following', params.id],
		queryFn: () => fetchFollowing(Number(params.id)),
	}))

	const followers = createQuery(() => ({
		queryKey: ['followers', params.id],
		queryFn: () => fetchFollowers(Number(params.id)),
	}))

	const handleMutate = async (id: number, isFollowing: boolean) => {
		await queryClient.cancelQueries({ type: 'active' })
		queryClient.setQueryData(
			['following', params.id],
			(old: UserProfileShort[]) =>
				old.map(user => {
					if (user.id === id) {
						return { ...user, is_following: isFollowing }
					}
					return user
				}),
		)
		queryClient.setQueryData(
			['followers', params.id],
			(old: UserProfileShort[]) =>
				old.map(user => {
					if (user.id === id) {
						return { ...user, is_following: isFollowing }
					}
					return user
				}),
		)
		queryClient.invalidateQueries({ queryKey: ['following', store.user.id] })
		queryClient.invalidateQueries({
			queryKey: ['profiles', store.user.username],
		})
	}

	const followMutate = createMutation(() => ({
		mutationFn: (id: number) => followUser(id),
		onMutate: (id: number) => handleMutate(id, true),
	}))

	const unfollowMutate = createMutation(() => ({
		mutationFn: (id: number) => unfollowUser(id),
		onMutate: (id: number) => handleMutate(id, false),
	}))

	const follow = async (e: Event, userID: number) => {
		e.preventDefault()
		followMutate.mutate(userID)
		window.Telegram.WebApp.HapticFeedback.impactOccurred('light')
	}

	const unfollow = async (e: Event, userID: number) => {
		e.preventDefault()
		unfollowMutate.mutate(userID)
		window.Telegram.WebApp.HapticFeedback.impactOccurred('light')
	}

	return (
		<div>
			<ul
				class="grid w-full grid-cols-2 text-center text-sm font-medium"
				id="default-tab"
				role="tablist"
			>
				<li class="me-2 flex items-center justify-center" role="presentation">
					<button
						onClick={() => setShowFollowing(false)}
						class="flex h-12 w-3/5 items-center justify-center px-4 text-sm font-medium text-main transition-all duration-300 ease-in-out"
						classList={{
							'border-b-4 border-accent': !showFollowing(),
							'border-b-4 border-transparent': showFollowing(),
						}}
						id="feed"
						role="tab"
						aria-controls="feed"
						aria-selected="false"
					>
						Followers
					</button>
				</li>
				<li class="me-2 flex items-center justify-center" role="presentation">
					<button
						onClick={() => setShowFollowing(true)}
						class="flex h-12 w-3/5 items-center justify-center px-4 text-sm font-medium text-main transition-all duration-300 ease-in-out"
						classList={{
							'border-b-4 border-accent': showFollowing(),
							'border-b-4 border-transparent': !showFollowing(),
						}}
						id="posts-tab"
						role="tab"
						aria-controls="posts"
						aria-selected="false"
					>
						Following
					</button>
				</li>
			</ul>
			<Switch>
				<Match when={showFollowing()}>
					<Show when={!following.isLoading} fallback={<Loader />}>
						<div class="grid w-full space-y-4 p-4">
							<For each={following.data}>
								{following => (
									<UserCardSmall
										user={following}
										onFollow={follow}
										onUnfollow={unfollow}
									/>
								)}
							</For>
						</div>
					</Show>
				</Match>
				<Match when={!showFollowing()}>
					<Show when={!followers.isLoading} fallback={<Loader />}>
						<div class="grid w-full space-y-4 p-4">
							<For each={followers.data}>
								{follower => (
									<UserCardSmall
										user={follower}
										onFollow={follow}
										onUnfollow={unfollow}
									/>
								)}
							</For>
						</div>
					</Show>
				</Match>
			</Switch>
		</div>
	)
}

const Loader = () => {
	return (
		<div class="flex h-screen flex-col items-start justify-start space-y-4 bg-secondary p-2">
			<div class="h-14 w-full rounded-2xl bg-main" />
			<div class="h-14 w-full rounded-2xl bg-main" />
			<div class="h-14 w-full rounded-2xl bg-main" />
			<div class="h-14 w-full rounded-2xl bg-main" />
		</div>
	)
}

function UserCardSmall(props: {
	user: UserProfileShort
	onFollow: (e: Event, id: number) => void
	onUnfollow: (e: Event, id: number) => void
}) {
	return (
		<Link
			href={'/users/' + props.user.username}
			class="flex w-full flex-row items-center justify-between"
			state={{ back: true }}
		>
			<div class="flex flex-row items-center justify-start space-x-4">
				<img
					src={CDN_URL + '/' + props.user.avatar_url}
					alt={props.user.username}
					class="size-11 shrink-0 rounded-xl object-cover"
				/>
				<div class="grid">
					<p class="font-semibold text-main">
						{props.user.last_name
							? props.user.first_name + ' ' + props.user.last_name
							: props.user.username}
					</p>
					<p class="text-sm text-secondary">{props.user.title}</p>
				</div>
			</div>
			<Switch>
				<Match
					when={props.user.is_following && props.user.id !== store.user.id}
				>
					<button
						class="rounded-full bg-main px-4 py-2 text-main"
						onClick={(e: Event) => props.onUnfollow(e, props.user.id!)}
					>
						Followed
					</button>
				</Match>
				<Match
					when={!props.user.is_following && props.user.id !== store.user.id}
				>
					<button
						class="rounded-full bg-accent px-4 py-2 text-main"
						onClick={(e: Event) => props.onFollow(e, props.user.id!)}
					>
						Follow
					</button>
				</Match>
			</Switch>
		</Link>
	)
}
