import {
	createEffect,
	createSignal,
	For,
	Match,
	onCleanup,
	Show,
	Switch,
} from 'solid-js'
import {
	Navigate,
	useNavigate,
	useParams,
	useSearchParams,
} from '@solidjs/router'
import {
	CDN_URL,
	fetchProfile,
	followUser,
	publishProfile,
	unfollowUser,
} from '~/lib/api'
import { setUser, store } from '~/store'
import ActionDonePopup from '~/components/ActionDonePopup'
import { useMainButton } from '~/lib/useMainButton'
import { usePopup } from '~/lib/usePopup'
import { createMutation, createQuery } from '@tanstack/solid-query'
import { queryClient } from '~/App'
import { Link } from '~/components/Link'
import { UserProfile } from '~/gen/types'
import { PeatchIcon } from '~/components/peatch-icon'

export default function UserProfilePage() {
	const mainButton = useMainButton()
	const [published, setPublished] = createSignal(false)

	const navigate = useNavigate()

	const params = useParams()
	const [searchParams] = useSearchParams()

	const username = params.handle

	const query = createQuery(() => ({
		queryKey: ['profiles', username],
		queryFn: () => fetchProfile(username),
		retry: 1,
	}))

	const followMutate = createMutation(() => ({
		mutationFn: (id: number) => followUser(id),
		onMutate: async (id: number) => {
			await queryClient.cancelQueries({ queryKey: ['profiles', username] })
			queryClient.setQueryData(['profiles', username], (old: UserProfile) => {
				if (old) {
					return {
						...old,
						is_following: true,
						followers_count: old.followers_count! + 1,
					}
				}
				return old
			})
			queryClient.invalidateQueries({ queryKey: ['followers', id.toString()] })
		},
	}))

	const unfollowMutate = createMutation(() => ({
		mutationFn: (id: number) => unfollowUser(id),
		onMutate: async (id: number) => {
			await queryClient.cancelQueries({ queryKey: ['profiles', username] })
			queryClient.setQueryData(['profiles', username], (old: UserProfile) => {
				if (old) {
					return {
						...old,
						is_following: false,
						followers_count: old.followers_count! - 1,
					}
				}
				return old
			})
			queryClient.invalidateQueries({ queryKey: ['followers', id.toString()] })
		},
	}))

	const { showAlert } = usePopup()

	const isCurrentUserProfile = store.user.username === username

	const navigateToEdit = () => {
		navigate('/users/edit', { state: { back: true } })
	}

	createEffect(async () => {
		if (searchParams.refetch) {
			await query.refetch()
			if (query.data.id === store.user.id) {
				setUser(query.data)
			}
		}
	})

	const navigateToCollaborate = async () => {
		if (!store.user.published_at) {
			showAlert(
				`Publish your profile first, so ${query.data.first_name} will see it`,
			)
		} else if (store.user.hidden_at) {
			showAlert(
				`Unhide your profile first, so ${query.data.first_name} will see it`,
			)
		} else {
			navigate(`/users/${username}/collaborate`, { state: { back: true } })
		}
	}

	const closePopup = () => {
		setPublished(false)
	}

	const publish = async () => {
		setUser({
			...store.user,
			published_at: new Date().toISOString(),
		})
		await publishProfile()
		setPublished(true)
	}

	const follow = async () => {
		if (!query.data) return
		followMutate.mutate(query.data.id)
		window.Telegram.WebApp.HapticFeedback.impactOccurred('light')
	}

	const unfollow = async () => {
		if (!query.data) return
		unfollowMutate.mutate(query.data.id)
		window.Telegram.WebApp.HapticFeedback.impactOccurred('light')
	}

	createEffect(() => {
		if (isCurrentUserProfile) {
			if (!store.user.published_at) {
				mainButton.enable('Publish')
				mainButton.onClick(publish)
			} else {
				if (published()) {
					mainButton.onClick(closePopup)
					mainButton.enable('Back to profile')
				} else {
					mainButton.enable('Edit')
					mainButton.onClick(navigateToEdit)
				}
			}
		} else if (store.user.published_at && !store.user.hidden_at) {
			mainButton.enable('Collaborate')
			mainButton.onClick(navigateToCollaborate)
		}

		onCleanup(() => {
			mainButton.offClick(navigateToCollaborate)
			mainButton.offClick(publish)
			mainButton.offClick(closePopup)
			mainButton.offClick(navigateToEdit)
		})
	})

	onCleanup(async () => {
		mainButton.hide()
	})

	function shareURL() {
		const url =
			'https://t.me/share?' +
			new URLSearchParams({
				url: 'https://t.me/peatch_bot/app?startapp=t-users-' + username,
				text: `Check out ${query.data.first_name} ${query.data.last_name}'s profile on Peatch`,
			}).toString()

		console.log('URL', url)

		window.Telegram.WebApp.openTelegramLink(url)
	}

	const [badgesExpanded, setBadgesExpanded] = createSignal(false)

	const [opportunitiesExpanded, setOpportunitiesExpanded] = createSignal(false)

	const showInfoPopup = () => {
		window.Telegram.WebApp.showAlert(
			'Your profile was hidden by our moderators. Try to make it more genuine.',
		)
	}

	return (
		<div>
			<Switch>
				<Match when={query.isLoading}>
					<Loader />
				</Match>
				<Match when={query.isError}>
					<Navigate href={'/404'} />
				</Match>
				<Match when={query.isSuccess}>
					<Switch>
						<Match when={published() && isCurrentUserProfile}>
							<ActionDonePopup
								action="Profile is under review"
								description="We will notify you once we finish moderation process"
								callToAction="There are 12 people you might be interested to collaborate with"
							/>
						</Match>
						<Match when={!query.isLoading}>
							<div class="h-fit min-h-screen bg-secondary">
								<Show when={isCurrentUserProfile && store.user.hidden_at}>
									<button
										onClick={showInfoPopup}
										class="absolute left-6 top-4 flex size-8 items-center justify-start"
									>
										<span class="material-symbols-rounded text-white">
											visibility_off
										</span>
									</button>
								</Show>
								<Show when={store.user.published_at && !store.user.hidden_at}>
									<Switch>
										<Match when={isCurrentUserProfile}>
											<Link
												class="absolute right-4 top-4 z-10 flex h-8 items-center justify-center gap-2 rounded-lg bg-white px-4 text-sm font-semibold text-[#FF8C42]"
												href="/rewards"
												state={{ back: true }}
											>
												<PeatchIcon width={16} height={16} />
												{store.user.peatch_points}
											</Link>
										</Match>
										<Match
											when={!isCurrentUserProfile && !query.data?.is_following}
										>
											<ActionButton
												disabled={unfollowMutate.isPending}
												text="Follow"
												onClick={follow}
											/>
										</Match>
										<Match
											when={!isCurrentUserProfile && query.data?.is_following}
										>
											<ActionButton
												disabled={followMutate.isPending}
												text="Unfollow"
												onClick={unfollow}
											/>
										</Match>
									</Switch>
								</Show>
								<div class="p-2">
									<img
										src={CDN_URL + '/' + query.data.avatar_url}
										alt="avatar"
										class="aspect-square size-full rounded-xl object-cover"
									/>
								</div>
								<div class="px-4 py-2.5">
									<div class="flex flex-row items-center justify-between pb-4">
										<div class="flex h-8 flex-row items-center space-x-2 text-sm font-semibold">
											<Link
												href={`/users/${query.data.id}/followers?show=following`}
												state={{ back: true }}
												class="flex h-full flex-row items-center gap-1.5 text-main"
											>
												<span>{query.data.following_count}</span>
												<span class="text-secondary">following</span>
											</Link>
											<Link
												class="flex h-full flex-row items-center gap-1.5 text-main"
												href={`/users/${query.data.id}/followers?show=followers`}
												state={{ back: true }}
											>
												{query.data.followers_count}{' '}
												<span class="text-secondary">followers</span>
											</Link>
										</div>
										<button
											class="flex h-8 flex-row items-center space-x-1.5 bg-transparent px-2.5"
											onClick={() => shareURL()}
										>
											<span class="text-sm font-semibold">Share profile</span>
											<span class="material-symbols-rounded text-[16px] text-main">
												content_copy
											</span>
										</button>
									</div>
									<p class="text-3xl text-pink">
										{query.data.first_name} {query.data.last_name}:
									</p>
									<p class="text-3xl text-main">{query.data.title}</p>
									<p class="mt-1 text-lg font-normal text-secondary">
										{query.data.description}
									</p>
									<div class="mt-5 flex flex-row flex-wrap items-center justify-start gap-1">
										<For
											each={
												badgesExpanded()
													? query.data.badges
													: query.data.badges.slice(0, 3)
											}
										>
											{badge => (
												<div
													class="flex h-10 flex-row items-center justify-center gap-[5px] rounded-2xl border px-2.5"
													style={{
														'background-color': `#${badge.color}`,
														'border-color': `#${badge.color}`,
													}}
												>
													<span class="material-symbols-rounded text-white">
														{String.fromCodePoint(parseInt(badge.icon!, 16))}
													</span>
													<p class="text-sm font-semibold text-white">
														{badge.text}
													</p>
												</div>
											)}
										</For>
									</div>
									<Show when={query.data.badges.length > 3}>
										<ExpandButton
											expanded={badgesExpanded()}
											setExpanded={setBadgesExpanded}
										/>
									</Show>
									<p class="py-4 text-3xl font-extrabold text-main">
										Available for
									</p>
									<div class="flex w-full flex-col items-center justify-start gap-1">
										<For
											each={
												opportunitiesExpanded()
													? query.data.opportunities
													: query.data.opportunities.slice(0, 3)
											}
										>
											{op => (
												<div
													class="flex h-[60px] w-full flex-row items-center justify-start gap-2.5 rounded-2xl border px-2.5"
													style={{
														'background-color': `#${op.color}`,
													}}
												>
													<div class="flex size-10 shrink-0 items-center justify-center rounded-full bg-secondary">
														<span class="material-symbols-rounded shrink-0 text-main">
															{String.fromCodePoint(parseInt(op.icon!, 16))}
														</span>
													</div>
													<div class="text-start text-white">
														<p class="text-sm font-semibold">{op.text}</p>
														<p class="text-xs leading-tight text-white/60">
															{op.description}
														</p>
													</div>
												</div>
											)}
										</For>
										<Show when={query.data.opportunities.length > 3}>
											<ExpandButton
												expanded={opportunitiesExpanded()}
												setExpanded={setOpportunitiesExpanded}
											/>
										</Show>
									</div>
								</div>
							</div>
						</Match>
					</Switch>
				</Match>
			</Switch>
		</div>
	)
}
// background: ;

const ActionButton = (props: {
	disabled: boolean
	text: string
	onClick: () => void
}) => {
	return (
		<button
			disabled={props.disabled}
			class="absolute right-4 top-4 z-10 h-9 w-[90px] rounded-xl bg-black/80 px-2.5 text-sm font-semibold text-white"
			onClick={() => props.onClick()}
		>
			{props.text}
		</button>
	)
}

const ExpandButton = (props: {
	expanded: boolean
	setExpanded: (val: boolean) => void
}) => {
	return (
		<button
			class="flex h-8 w-full items-center justify-start rounded-xl bg-transparent text-sm font-semibold text-secondary"
			onClick={() => props.setExpanded(!props.expanded)}
		>
			<span class="material-symbols-rounded text-secondary">
				{props.expanded ? 'expand_less' : 'expand_more'}
			</span>
			{props.expanded ? 'show less' : 'show more'}
		</button>
	)
}

const Loader = () => {
	return (
		<div class="flex h-screen flex-col items-start justify-start bg-secondary p-2">
			<div class="aspect-square w-full rounded-xl bg-main" />
			<div class="flex flex-col items-start justify-start p-2">
				<div class="mt-2 h-6 w-1/2 rounded bg-main" />
				<div class="mt-2 h-6 w-1/2 rounded bg-main" />
				<div class="mt-2 h-20 w-full rounded bg-main" />
				<div class="mt-4 flex w-full flex-row flex-wrap items-center justify-start gap-2">
					<div class="h-10 w-40 rounded-2xl bg-main" />
					<div class="h-10 w-32 rounded-2xl bg-main" />
					<div class="h-10 w-40 rounded-2xl bg-main" />
					<div class="h-10 w-28 rounded-2xl bg-main" />
					<div class="h-10 w-32 rounded-2xl bg-main" />
				</div>
			</div>
		</div>
	)
}
