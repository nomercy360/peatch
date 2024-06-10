import {
	createEffect,
	createSignal,
	For,
	Match,
	onCleanup,
	onMount,
	Show,
	Suspense,
	Switch,
} from 'solid-js'
import { Collaboration, User } from '~/gen/types'
import { CDN_URL, fetchFeed } from '~/lib/api'
import { Link } from '~/components/Link'
import BadgeList from '~/components/BadgeList'
import useDebounce from '~/lib/useDebounce'
import { createQuery } from '@tanstack/solid-query'
import { store } from '~/store'
import FillProfilePopup from '~/components/FillProfilePopup'
import { useMainButton } from '~/lib/useMainButton'
import { useNavigate } from '@solidjs/router'

export default function FeedPage() {
	const [search, setSearch] = createSignal('')

	const updateSearch = useDebounce(setSearch, 350)

	const mainButton = useMainButton()
	const navigate = useNavigate()

	const query = createQuery(() => ({
		queryKey: ['feed', search()],
		queryFn: () => fetchFeed(search()),
	}))

	const [scroll, setScroll] = createSignal(0)

	const [profilePopup, setProfilePopup] = createSignal(false)
	const [communityPopup, setCommunityPopup] = createSignal(false)
	const [rewardsPopup, setRewardsPopup] = createSignal(true)

	createEffect(() => {
		const onScroll = () => setScroll(window.scrollY)
		window.addEventListener('scroll', onScroll)
		return () => window.removeEventListener('scroll', onScroll)
	})

	const toCreateCollab = () => {
		navigate('/collaborations/edit')
	}

	onMount(() => {
		window.Telegram.WebApp.CloudStorage.getItem(
			'profilePopup',
			updateProfilePopup,
		)
		window.Telegram.WebApp.CloudStorage.getItem(
			'communityPopup',
			updateCommunityPopup,
		)
		window.Telegram.WebApp.CloudStorage.getItem(
			'rewardsPopup',
			updateRewardsPopup,
		)

		if (store.user.published_at) {
			mainButton.enable('Post to Peatch').onClick(toCreateCollab)
		}

		window.Telegram.WebApp.disableClosingConfirmation()
		// window.Telegram.WebApp.CloudStorage.removeItem('profilePopup')
	})

	const getUserLink = () => {
		if (store.user.first_name && store.user.description) {
			return '/users/' + store.user?.username
		} else {
			return '/users/edit'
		}
	}

	const closePopup = (name: string) => {
		switch (name) {
			case 'profilePopup':
				setProfilePopup(false)
				break
			case 'communityPopup':
				setCommunityPopup(false)
				break
			case 'rewardsPopup':
				setRewardsPopup(false)
				break
		}
		window.Telegram.WebApp.CloudStorage.setItem(name, 'closed')
	}

	const updateProfilePopup = (err: unknown, value: unknown) => {
		setProfilePopup(value !== 'closed')
	}

	const updateCommunityPopup = (err: unknown, value: unknown) => {
		setCommunityPopup(value !== 'closed')
	}

	const updateRewardsPopup = (err: unknown, value: unknown) => {
		setRewardsPopup(value !== 'closed')
	}

	onCleanup(() => {
		mainButton.hide()
		mainButton.offClick(toCreateCollab)
	})

	return (
		<div class="min-h-screen bg-secondary pb-56 pt-16">
			<Show when={!store.user.published_at && profilePopup()}>
				<FillProfilePopup onClose={() => closePopup('profilePopup')} />
			</Show>
			<Show when={communityPopup()}>
				<OpenCommunityPopup onClose={() => closePopup('communityPopup')} />
			</Show>
			<Show when={!communityPopup() && rewardsPopup()}>
				<RewardsPopup onClose={() => closePopup('rewardsPopup')} />
			</Show>
			<div class="fixed top-0 z-30 flex w-full flex-row items-center justify-between space-x-4 border-b bg-secondary p-4">
				<div class="relative flex h-10 w-full flex-row items-center justify-center rounded-lg bg-main">
					<input
						class="h-full w-full bg-transparent px-2.5 text-main placeholder:text-hint"
						placeholder="Search people or collaborations"
						type="text"
						value={search()}
						onInput={e => updateSearch(e.currentTarget.value)}
					/>
					<Show when={search()}>
						<button
							class="absolute right-2.5 flex size-5 shrink-0 items-center justify-center rounded-full bg-neutral-400"
							onClick={() => setSearch('')}
						>
							<span class="material-symbols-rounded text-[20px] text-button">
								close
							</span>
						</button>
					</Show>
				</div>
				<Link
					class="flex shrink-0 flex-row items-center justify-between"
					href={getUserLink()}
				>
					<Switch>
						<Match when={store.user.avatar_url}>
							<img
								class="size-10 rounded-xl border border-main object-cover"
								src={CDN_URL + '/' + store.user.avatar_url}
								alt="User Avatar"
							/>
						</Match>
						<Match when={!store.user.avatar_url}>
							<div class="flex size-10 items-center justify-center rounded-xl border-2 border-main bg-main">
								<span class="material-symbols-rounded text-peatch-main">
									account_circle
								</span>
							</div>
						</Match>
					</Switch>
				</Link>
			</div>
			<Suspense fallback={<ListPlaceholder />}>
				<For each={query.data}>
					{(data, i) => (
						<>
							<Switch fallback={<div />}>
								<Match when={data.type === 'user'}>
									<UserCard user={data.data as User} scroll={scroll()} />
								</Match>
								<Match when={data.type === 'collaboration'}>
									<CollaborationCard
										collab={data.data as Collaboration}
										scroll={scroll()}
									/>
								</Match>
							</Switch>
							<div class="mt-5 h-px w-full bg-border" />
						</>
					)}
				</For>
			</Suspense>
		</div>
	)
}

const UserCard = (props: { user: User; scroll: number }) => {
	const shortenDescription = (description: string) => {
		if (description.length <= 120) return description
		return description.slice(0, 120) + '...'
	}

	const imgUrl = `https://assets.peatch.io/${props.user.avatar_url}`

	return (
		<Link
			class="flex flex-col items-start bg-secondary px-4 pb-5 pt-4 text-start"
			href={`/users/${props.user.username}`}
			state={{ from: '/', scroll: props.scroll }}
		>
			<img
				class="size-11 rounded-xl object-cover"
				src={imgUrl}
				loading="lazy"
				alt="User Avatar"
			/>
			<p class="mt-3 text-3xl text-blue">{props.user.first_name}:</p>
			<p class="text-3xl text-main">{props.user.title}</p>
			<p class="mt-2 text-sm text-hint">
				{shortenDescription(props.user.description!)}
			</p>
			<Show when={props.user.badges && props.user.badges.length > 0}>
				<BadgeList
					badges={props.user.badges!}
					position="start"
					city={props.user.city!}
					country={props.user.country!}
					countryCode={props.user.country_code!}
				/>
			</Show>
		</Link>
	)
}

const OpenCommunityPopup = (props: { onClose: () => void }) => {
	return (
		<div class="w-full p-4">
			<div class="relative rounded-2xl bg-main p-4 text-center">
				<span class="material-symbols-rounded text-[48px] text-green">
					maps_ugc
				</span>
				<button
					class="absolute right-4 top-4 flex size-6 items-center justify-center rounded-full bg-neutral-200"
					onClick={props.onClose}
				>
					<span class="material-symbols-rounded text-[24px] text-button">
						close
					</span>
				</button>
				<p class="text-3xl font-extrabold text-green">Join community</p>
				<p class="mt-2 text-xl font-normal text-main">
					to talk with founders and users. Discuss and solve problems together
				</p>
				<button
					class="mt-4 flex h-12 w-full items-center justify-center rounded-xl bg-secondary text-main"
					onClick={() =>
						window.Telegram.WebApp.openTelegramLink(
							'https://t.me/peatch_community',
						)
					}
				>
					Open Peatch Community
				</button>
			</div>
		</div>
	)
}

const RewardsPopup = (props: { onClose: () => void }) => {
	return (
		<div class="w-full p-4">
			<div class="relative rounded-2xl bg-main p-4 text-center">
				<span class="material-symbols-rounded text-[48px] text-pink">
					emoji_events
				</span>
				<button
					class="absolute right-4 top-4 flex size-6 items-center justify-center rounded-full bg-neutral-200"
					onClick={props.onClose}
				>
					<span class="material-symbols-rounded text-[24px] text-button">
						close
					</span>
				</button>
				<p class="text-3xl font-extrabold text-pink">Peatch Rewards</p>
				<p class="mt-2 text-xl font-normal text-main">
					Learn more about our internal currency and how to earn it
				</p>
				<Link
					href={'/rewards'}
					class="mt-4 flex h-12 w-full items-center justify-center rounded-xl bg-secondary text-main"
				>
					Show me
				</Link>
			</div>
		</div>
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
			class="flex flex-col items-start px-4 pt-4 text-start"
			href={`/collaborations/${props.collab.id}`}
			state={{ from: '/', scroll: props.scroll }}
		>
			<div
				class="flex size-10 flex-row items-center justify-center rounded-full"
				style={{ 'background-color': `#${props.collab.opportunity?.color}` }}
			>
				<span class="material-symbols-rounded text-[20px] text-white">
					{String.fromCodePoint(parseInt(props.collab.opportunity?.icon!, 16))}
				</span>
			</div>
			<p
				class="mt-3 text-3xl"
				style={{ color: `#${props.collab.opportunity?.color}` }}
			>
				{props.collab.opportunity?.text}:
			</p>
			<p class="text-3xl text-main">{props.collab.title}</p>
			<p class="mt-2 text-sm text-hint">
				{shortenDescription(props.collab.description!)}
			</p>
		</Link>
	)
}

const ListPlaceholder = () => {
	return (
		<div class="flex flex-col items-start justify-start gap-4 px-4 py-2.5">
			<div class="h-52 w-full rounded-2xl bg-main" />
			<div class="h-64 w-full rounded-2xl bg-main" />
			<div class="h-48 w-full rounded-2xl bg-main" />
			<div class="h-56 w-full rounded-2xl bg-main" />
		</div>
	)
}
