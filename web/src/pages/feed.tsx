import {
	createEffect,
	createSignal,
	For,
	onMount,
	Show,
	Suspense,
} from 'solid-js'
import { Link } from '~/components/link'
import BadgeList from '~/components/badge-list'
import useDebounce from '~/lib/useDebounce'
import { store } from '~/store'
import FillProfilePopup from '~/components/fill-profile-popup'
import { LocationBadge } from '~/components/location-badge'
import { useTranslations } from '~/lib/locale-context'
import { useInfiniteQuery } from '@tanstack/solid-query'
import { verificationStatus, UserProfileResponse } from '~/gen'
import { fetchUsers } from '~/lib/api'
import { useNavigation } from '~/lib/useNavigation'


export const [search, setSearch] = createSignal('')

export default function FeedPage() {
	const { t } = useTranslations()
	const navigation = useNavigation()

	const updateSearch = useDebounce(setSearch, 350)

	const query = useInfiniteQuery(() => ({
		queryKey: ['users', search()],
		queryFn: fetchUsers,
		getNextPageParam: (lastPage) => lastPage.nextPage,
		initialPageParam: 1,
	}))

	const [scroll, setScroll] = createSignal(0)
	const [profilePopup, setProfilePopup] = createSignal(false)
	const [communityPopup, setCommunityPopup] = createSignal(false)
	const [isLoadingMore, setIsLoadingMore] = createSignal(false)

	const loadMoreUsers = () => {
		if (query.hasNextPage && !query.isFetchingNextPage && !isLoadingMore()) {
			setIsLoadingMore(true)
			query.fetchNextPage().finally(() => setIsLoadingMore(false))
		}
	}

	createEffect(() => {
		const onScroll = () => {
			setScroll(window.scrollY)

			const feedElement = document.getElementById('feed')
			if (feedElement) {
				const { scrollTop, scrollHeight, clientHeight } = feedElement
				if (scrollHeight - scrollTop - clientHeight < 300) {
					loadMoreUsers()
				}
			}
		}

		const feedElement = document.getElementById('feed')
		if (feedElement) {
			feedElement.addEventListener('scroll', onScroll)
			return () => feedElement.removeEventListener('scroll', onScroll)
		}
	})

	onMount(() => {
		window.Telegram.WebApp.CloudStorage.getItem(
			'profilePopup',
			updateProfilePopup,
		)
		window.Telegram.WebApp.CloudStorage.getItem(
			'communityPopup',
			updateCommunityPopup,
		)

		window.Telegram.WebApp.disableClosingConfirmation()
		// window.Telegram.WebApp.CloudStorage.removeItem('profilePopup')
		// window.Telegram.WebApp.CloudStorage.removeItem('communityPopup')
	})

	const closePopup = (name: string) => {
		switch (name) {
			case 'profilePopup':
				setProfilePopup(false)
				break
			case 'communityPopup':
				setCommunityPopup(false)
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

	const allUsers = () => {
		if (!query.data) return []
		return query.data.pages.flatMap(page => page.data)
	}

	return (
		<div class="flex h-screen flex-col overflow-hidden">
			<div class="flex w-full flex-shrink-0 flex-col items-center justify-between space-y-4 border-b p-4">
				<Show
					when={store.user.verification_status == verificationStatus.VerificationStatusUnverified && profilePopup()}>
					<FillProfilePopup onClose={() => closePopup('profilePopup')} />
				</Show>
				<Show
					when={communityPopup() && store.user.verification_status == verificationStatus.VerificationStatusVerified}>
					<OpenCommunityPopup onClose={() => closePopup('communityPopup')} />
				</Show>
				<div class="relative flex h-10 w-full flex-row items-center justify-center rounded-lg bg-secondary">
					<input
						class="h-full w-full bg-transparent px-2.5 placeholder:text-secondary-foreground"
						placeholder={t('common.search.people')}
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
			<div class="flex h-full w-full flex-shrink-0 flex-col overflow-y-auto pb-20" id="feed">
				<Suspense fallback={<ListPlaceholder />}>
					<For each={allUsers()}>
						{(user, _) => (
							<div>
								<UserCard user={user} scroll={scroll()} />
								<div class="h-px w-full bg-border" />
							</div>
						)}
					</For>

					<Show when={query.isFetchingNextPage}>
						<div class="flex justify-center p-4">
							<div class="h-10 w-10 animate-spin rounded-full border-4 border-primary border-t-transparent"></div>
						</div>
					</Show>

					<Show when={!query.hasNextPage && allUsers().length > 0}>
						<div class="p-4 text-center text-secondary-foreground">
							{t('common.search.noMoreResults')}
						</div>
					</Show>

					<Show when={allUsers().length === 0 && !query.isLoading}>
						<div class="p-4 text-center text-secondary-foreground">
							{t('common.search.noResults')}
						</div>
					</Show>
				</Suspense>
			</div>
		</div>
	)
}

type UserCardProps = {
	user: UserProfileResponse
	scroll: number
}

const UserCard = (props: UserCardProps) => {
	const shortenDescription = (description: string) => {
		if (description.length <= 120) return description
		return description.slice(0, 120) + '...'
	}
	return (
		<Link
			class="flex flex-col items-start px-4 pb-5 pt-4 text-start"
			href={`/users/${props.user.id}`}
			state={{ from: '/' }}
		>
			<img
				class="size-10 rounded-xl object-cover"
				src={`https://assets.peatch.io/cdn-cgi/image/width=100/${props.user.avatar_url}`}
				loading="lazy"
				alt="User Avatar"
			/>
			<p class="mt-3 text-3xl font-semibold capitalize text-primary">{props.user.first_name?.trimEnd()}:</p>
			<p class="text-3xl capitalize">{props.user.title}</p>
			<p class="mt-2 text-sm text-secondary-foreground">
				{shortenDescription(props.user.description!)}
			</p>
			<LocationBadge
				country={props.user.location?.country_name}
				city={props.user.location?.name}
				countryCode={props.user.location?.country_code}
			/>
			<Show when={props.user.badges && props.user.badges.length > 0}>
				<BadgeList badges={props.user.badges || []} position="start" />
			</Show>
		</Link>
	)
}

const OpenCommunityPopup = (props: { onClose: () => void }) => {
	return (
		<div class="relative w-full rounded-xl bg-secondary p-3 text-center">
			<button
				class="absolute right-4 top-4 flex size-6 items-center justify-center rounded-full bg-background"
				onClick={() => props.onClose()}
			>
					<span class="material-symbols-rounded text-[20px] text-secondary-foreground">
						close
					</span>
			</button>
			<div class="text-green flex items-center justify-center gap-1 text-2xl font-extrabold">
					<span class="material-symbols-rounded text-[36px] text-green-400">
						maps_ugc
					</span>
				Join community
			</div>
			<p class="mt-2 text-base font-normal text-secondary-foreground">
				To talk with founders and users. Discuss and solve problems together
			</p>
			<button
				class="mt-4 flex h-10 w-full items-center justify-center rounded-xl bg-primary text-sm font-semibold"
				onClick={() =>
					window.Telegram.WebApp.openTelegramLink(
						'https://t.me/peatch_community',
					)
				}
			>
				Open Peatch Community
			</button>
		</div>
	)
}

export const ListPlaceholder = () => {
	return (
		<div class="flex flex-col items-start justify-start gap-4 px-4 py-2.5">
			<div class="h-52 w-full rounded-2xl bg-secondary" />
			<div class="h-64 w-full rounded-2xl bg-secondary" />
			<div class="h-48 w-full rounded-2xl bg-secondary" />
			<div class="h-56 w-full rounded-2xl bg-secondary" />
		</div>
	)
}
