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
import { Link } from '~/components/link'
import { UserProfile } from '~/gen/types'
import { PeatchIcon } from '~/components/peatch-icon'
import { NotificationIcon } from '~/pages/users/activity'

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
			window.Telegram.WebApp.showConfirm(
				`Publish your profile first, so ${query.data.first_name} will see it`,
				(ok: boolean) =>
					ok && navigate('/users/edit', { state: { back: true } }),
			)
		} else if (store.user.hidden_at) {
			showAlert('Your profile is hidden by our moderators')
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
		} else {
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
			'https://t.me/share/url?' +
			new URLSearchParams({
				url: 'https://t.me/peatch_bot/app?startapp=t-users-' + username,
			}).toString() +
			`&text=Check out ${query.data.first_name} ${query.data.last_name}'s profile on Peatch! 🌟`

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
							<div class="h-fit min-h-screen p-2 bg-secondary">
								<Show when={isCurrentUserProfile && store.user.hidden_at}>
									<button
										onClick={showInfoPopup}
										class="absolute left-4 top-4 flex size-8 items-center justify-center rounded-lg bg-secondary"
									>
										<span class="material-symbols-rounded">
											visibility_off
										</span>
									</button>
								</Show>
								<Show when={store.user.published_at && !store.user.hidden_at}>
									<Switch>
										<Match when={isCurrentUserProfile}>
											<Link
												href="/users/activity"
												state={{ back: true }}
												class="absolute left-4 top-4 z-10 flex size-8 items-center justify-center rounded-lg bg-secondary"
											>
												<NotificationIcon width={20} height={19} />
											</Link>
											<div
												class="absolute right-4 top-4 z-10 flex h-8 items-center justify-center gap-2 rounded-lg bg-secondary px-4 text-sm font-semibold text-orange"
											>
												<PeatchIcon width={16} height={16} />
												{store.user.peatch_points}
											</div>
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
								<div class="w-full aspect-square bg-cover bg-center relative rounded-xl"
										 style={{ 'background-image': `url(${CDN_URL + '/' + query.data.avatar_url})` }}>
									<div
										class="flex flex-row items-center justify-between absolute bottom-0 left-0 w-full rounded-b-xl h-10 px-4 bg-gradient-to-t from-background">
										<button
											class="flex h-8 flex-row items-center space-x-1 px-2.5 bg-secondary rounded-2xl"
											onClick={() => shareURL()}
										>
											<span class="material-symbols-rounded text-[16px]">
												ios_share
											</span>
											<span class="text-sm">
												Share
											</span>
										</button>
										<div class="flex h-8 flex-row items-center space-x-2 text-sm font-semibold">
											<Link
												href={`/users/${query.data.id}/followers?show=following`}
												state={{ back: true }}
												class="text-primary-foreground flex h-full flex-row items-center gap-1.5"
											>
												<span>{query.data.following_count}</span>
												<span class="text-primary-foreground opacity-70 font-normal">following</span>
											</Link>
											<Link
												class="text-primary-foreground flex h-full flex-row items-center gap-1.5"
												href={`/users/${query.data.id}/followers?show=followers`}
												state={{ back: true }}
											>
												{query.data.followers_count}{' '}
												<span class="text-primary-foreground opacity-70 font-normal">followers</span>
											</Link>
										</div>
									</div>
								</div>
								<div class="px-4 py-2.5">
									<p class="capitalize text-3xl text-primary font-semibold">
										{query.data.first_name} {query.data.last_name}:
									</p>
									<p class="text-3xl capitalize">{query.data.title}</p>
									<p class="mt-1 text-sm font-normal text-secondary-foreground">
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
									<p class="py-4 text-3xl font-extrabold">
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
														<span class="material-symbols-rounded shrink-0">
															{String.fromCodePoint(parseInt(op.icon!, 16))}
														</span>
													</div>
													<div class="text-start">
														<p class="text-sm font-semibold text-white">{op.text}</p>
														<p class="text-xs leading-tight text-white/80">
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
			class="flex h-8 w-full items-center justify-start rounded-xl bg-transparent text-sm font-semibold text-secondary-foreground"
			onClick={() => props.setExpanded(!props.expanded)}
		>
			<span class="material-symbols-rounded text-secondary-foreground">
				{props.expanded ? 'expand_less' : 'expand_more'}
			</span>
			{props.expanded ? 'show less' : 'show more'}
		</button>
	)
}

const Loader = () => {
	return (
		<div class="flex min-h-screen flex-col items-start justify-start bg-secondary p-2">
			<div class="aspect-square w-full rounded-xl bg-background" />
			<div class="flex flex-col items-start justify-start p-2">
				<div class="mt-2 h-6 w-1/2 rounded bg-background" />
				<div class="mt-2 h-6 w-1/2 rounded bg-background" />
				<div class="mt-2 h-20 w-full rounded bg-background" />
				<div class="mt-4 flex w-full flex-row flex-wrap items-center justify-start gap-2">
					<div class="h-10 w-40 rounded-2xl bg-background" />
					<div class="h-10 w-32 rounded-2xl bg-background" />
					<div class="h-10 w-40 rounded-2xl bg-background" />
					<div class="h-10 w-28 rounded-2xl bg-background" />
					<div class="h-10 w-32 rounded-2xl bg-background" />
				</div>
			</div>
		</div>
	)
}
