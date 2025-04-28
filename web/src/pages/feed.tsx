import {
	createEffect,
	createSignal,
	For,
	onCleanup,
	onMount,
	Show,
	Suspense,
} from 'solid-js'
import { User, UserProfile } from '~/gen/types'
import { fetchUsers } from '~/lib/api'
import { Link } from '~/components/link'
import BadgeList from '~/components/BadgeList'
import useDebounce from '~/lib/useDebounce'
import { createQuery } from '@tanstack/solid-query'
import { store } from '~/store'
import FillProfilePopup from '~/components/FillProfilePopup'
import { useMainButton } from '~/lib/useMainButton'
import { useNavigate } from '@solidjs/router'
import { LocationBadge } from '~/components/location-badge'


export const [search, setSearch] = createSignal('')

export default function FeedPage() {
	const updateSearch = useDebounce(setSearch, 350)

	const mainButton = useMainButton()
	const navigate = useNavigate()

	const query = createQuery(() => ({
		queryKey: ['users', search()],
		queryFn: () => fetchUsers(search()),
	}))

	const [scroll, setScroll] = createSignal(0)

	const [profilePopup, setProfilePopup] = createSignal(false)
	const [communityPopup, setCommunityPopup] = createSignal(false)

	createEffect(() => {
		const onScroll = () => setScroll(window.scrollY)
		window.addEventListener('scroll', onScroll)
		return () => window.removeEventListener('scroll', onScroll)
	})

	const toCreateCollab = () => {
		navigate('/collaborations/edit')
	}

	const [dropDown, setDropDown] = createSignal(false)

	const closeDropDown = () => {
		setDropDown(false)
	}

	const openDropDown = () => {
		document.body.style.overflow = 'hidden'
		setDropDown(true)
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

		window.Telegram.WebApp.disableClosingConfirmation()
		// window.Telegram.WebApp.CloudStorage.removeItem('profilePopup')
		// window.Telegram.WebApp.CloudStorage.removeItem('communityPopup')
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
		}
		window.Telegram.WebApp.CloudStorage.setItem(name, 'closed')
	}

	const updateProfilePopup = (err: unknown, value: unknown) => {
		setProfilePopup(value !== 'closed')
	}

	const updateCommunityPopup = (err: unknown, value: unknown) => {
		setCommunityPopup(value !== 'closed')
	}

	onMount(() => {
		// disable scroll on body when drawer is open
		document.body.style.overflow = 'hidden'
	})

	onCleanup(() => {
		mainButton.hide()
		mainButton.offClick(toCreateCollab)
		mainButton.offClick(openDropDown)
		document.removeEventListener('click', closeDropDownOnOutsideClick)
		document.body.style.overflow = 'auto'
	})

	// if dropdown is open, every click outside of the dropdown will close
	const closeDropDownOnOutsideClick = (e: MouseEvent) => {
		if (
			dropDown() &&
			!e.composedPath().includes(document.getElementById('dropdown-menu')!)
		) {
			closeDropDown()
		}
	}

	document.addEventListener('click', closeDropDownOnOutsideClick)

	return (
		<div class="flex h-screen flex-col">
			<Show when={dropDown()}>
				<div
					class="fixed inset-0 z-50 flex h-screen w-full flex-col items-center justify-end px-4 py-2.5"
					style={{
						'background-color':
							window.Telegram.WebApp.colorScheme === 'dark'
								? 'rgba(0, 0, 0, 0.8)'
								: 'rgba(255, 255, 255, 0.8)',
					}}
				>
					<div
						id="dropdown-menu"
						class="flex w-full flex-col items-center justify-center rounded-xl bg-secondary"
					>
						<button
							class="flex h-12 w-full items-center justify-center bg-transparent"
							onClick={() => navigate('/collaborations/edit')}
						>
							New collaboration
						</button>
						<div class="h-px w-full bg-border" />
						<button
							class="flex h-12 w-full items-center justify-center bg-transparent"
							onClick={() => navigate('/posts/edit')}
						>
							New post
						</button>
					</div>
				</div>
			</Show>
			<div class="flex w-full flex-shrink-0 flex-col items-center justify-between space-y-4 border-b p-4">
				<Show when={!store.user.published_at && profilePopup()}>
					<FillProfilePopup onClose={() => closePopup('profilePopup')} />
				</Show>
				<Show when={communityPopup() && store.user.published_at}>
					<OpenCommunityPopup onClose={() => closePopup('communityPopup')} />
				</Show>
				<div class="relative flex h-10 w-full flex-row items-center justify-center rounded-lg bg-secondary">
					<input
						class="h-full w-full bg-transparent px-2.5 placeholder:text-secondary-foreground"
						placeholder="Search people"
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
			<div class="bg-secondary flex h-full w-full flex-shrink-0 flex-col overflow-y-auto pb-20" id="feed">
				<Suspense fallback={<ListPlaceholder />}>
					<For each={query.data}>
						{(user, i) => (
							<div>
								<UserCard user={user as User} scroll={scroll()} />
								<div class="h-px w-full bg-border" />
							</div>
						)}
					</For>
				</Suspense>
			</div>
		</div>
	)
}

const UserCard = (props: { user: User; scroll: number }) => {
	const shortenDescription = (description: string) => {
		if (description.length <= 120) return description
		return description.slice(0, 120) + '...'
	}

	const user = props.user as UserProfile

	const imgUrl = `https://assets.peatch.io/cdn-cgi/image/width=100/${user.avatar_url}`

	return (
		<Link
			class="flex flex-col items-start px-4 pb-5 pt-4 text-start"
			href={`/users/${user.username}`}
			state={{ from: '/', scroll: props.scroll }}
		>
			<img
				class="size-10 rounded-xl object-cover"
				src={imgUrl}
				loading="lazy"
				alt="User Avatar"
			/>
			<p class="mt-3 text-3xl text-primary font-semibold capitalize">{user.first_name?.trimEnd()}:</p>
			<p class="text-3xl capitalize">{user.title}</p>
			<p class="mt-2 text-sm text-secondary-foreground">
				{shortenDescription(user.description!)}
			</p>
			<LocationBadge
				country={user.country!}
				city={user.city!}
				countryCode={user.country_code!}
			/>
			<Show when={user.badges && user.badges.length > 0}>
				<BadgeList badges={user.badges!} position="start" />
			</Show>
		</Link>
	)
}

const OpenCommunityPopup = (props: { onClose: () => void }) => {
	return (
		<div class="w-full bg-secondary rounded-xl relative p-3 text-center">
			<button
				class="absolute right-4 top-4 flex size-6 items-center justify-center rounded-full bg-background"
				onClick={props.onClose}
			>
					<span class="material-symbols-rounded text-[20px] text-secondary-foreground">
						close
					</span>
			</button>
			<div class="flex items-center gap-1 justify-center text-2xl font-extrabold text-green">
					<span class="material-symbols-rounded text-[36px] text-green-400">
						maps_ugc
					</span>
				Join community
			</div>
			<p class="mt-2 text-base font-normal text-secondary-foreground">
				To talk with founders and users. Discuss and solve problems together
			</p>
			<button
				class="bg-primary mt-4 flex h-10 w-full items-center justify-center rounded-xl text-sm font-semibold"
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
