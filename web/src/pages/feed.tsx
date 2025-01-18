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
import { fetchUsers, likeContent, unlikeContent } from '~/lib/api'
import { Link } from '~/components/link'
import BadgeList from '~/components/BadgeList'
import useDebounce from '~/lib/useDebounce'
import { createMutation, createQuery } from '@tanstack/solid-query'
import { store } from '~/store'
import FillProfilePopup from '~/components/FillProfilePopup'
import { useMainButton } from '~/lib/useMainButton'
import { useNavigate } from '@solidjs/router'
import { LocationBadge } from '~/components/location-badge'
import { queryClient } from '~/App'
import NavigationTabs from '~/components/navigation-tabs'

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

	const imgUrl = `https://assets.peatch.io/${user.avatar_url}`

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
			<Show when={user.badges && user.badges.length > 0}>
				<BadgeList badges={user.badges!} position="start">
					<LocationBadge
						country={user.country!}
						city={user.city!}
						countryCode={user.country_code!}
					/>
				</BadgeList>
			</Show>
			<LikeButton
				liked={user.is_liked!}
				likes={user.likes_count!}
				id={user.id!}
			/>
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

const LikeButton = (props: {
	id: number
	liked: boolean
	likes: number
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
		mutationFn: (id: number) => likeContent(id, 'user'),
		onMutate: id => handleMutate(id),
	}))

	const mutateUnLike = createMutation(() => ({
		mutationFn: (id: number) => unlikeContent(id, 'user'),
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

export function HeartIcon(props: any) {
	return (
		<svg
			{...props}
			width="24"
			height="24"
			viewBox="0 0 24 24"
			fill="none"
			xmlns="http://www.w3.org/2000/svg"
		>
			<mask
				id="mask0_4780_1925"
				style="mask-type:alpha"
				maskUnits="userSpaceOnUse"
				x="0"
				y="0"
				width="24"
				height="24"
			>
				<rect width="24" height="24" fill="#D9D9D9" />
			</mask>
			<g mask="url(#mask0_4780_1925)">
				<path
					d="M12 20.3249C11.7667 20.3249 11.5292 20.2832 11.2875 20.1999C11.0458 20.1166 10.8333 19.9832 10.65 19.7999L8.925 18.2249C7.15833 16.6082 5.5625 15.0041 4.1375 13.4124C2.7125 11.8207 2 10.0666 2 8.1499C2 6.58324 2.525 5.2749 3.575 4.2249C4.625 3.1749 5.93333 2.6499 7.5 2.6499C8.38333 2.6499 9.21667 2.8374 10 3.2124C10.7833 3.5874 11.45 4.0999 12 4.7499C12.55 4.0999 13.2167 3.5874 14 3.2124C14.7833 2.8374 15.6167 2.6499 16.5 2.6499C18.0667 2.6499 19.375 3.1749 20.425 4.2249C21.475 5.2749 22 6.58324 22 8.1499C22 10.0666 21.2917 11.8249 19.875 13.4249C18.4583 15.0249 16.85 16.6332 15.05 18.2499L13.35 19.7999C13.1667 19.9832 12.9542 20.1166 12.7125 20.1999C12.4708 20.2832 12.2333 20.3249 12 20.3249Z"
					fill="#FF8C42"
				/>
				<path
					d="M12 20.3249C11.7667 20.3249 11.5292 20.2832 11.2875 20.1999C11.0458 20.1166 10.8333 19.9832 10.65 19.7999L8.925 18.2249C7.15833 16.6082 5.5625 15.0041 4.1375 13.4124C2.7125 11.8207 2 10.0666 2 8.1499C2 6.58324 2.525 5.2749 3.575 4.2249C4.625 3.1749 5.93333 2.6499 7.5 2.6499C8.38333 2.6499 9.21667 2.8374 10 3.2124C10.7833 3.5874 11.45 4.0999 12 4.7499C12.55 4.0999 13.2167 3.5874 14 3.2124C14.7833 2.8374 15.6167 2.6499 16.5 2.6499C18.0667 2.6499 19.375 3.1749 20.425 4.2249C21.475 5.2749 22 6.58324 22 8.1499C22 10.0666 21.2917 11.8249 19.875 13.4249C18.4583 15.0249 16.85 16.6332 15.05 18.2499L13.35 19.7999C13.1667 19.9832 12.9542 20.1166 12.7125 20.1999C12.4708 20.2832 12.2333 20.3249 12 20.3249Z"
					fill="url(#paint0_radial_4780_1925)"
				/>
				<path
					d="M12 20.3249C11.7667 20.3249 11.5292 20.2832 11.2875 20.1999C11.0458 20.1166 10.8333 19.9832 10.65 19.7999L8.925 18.2249C7.15833 16.6082 5.5625 15.0041 4.1375 13.4124C2.7125 11.8207 2 10.0666 2 8.1499C2 6.58324 2.525 5.2749 3.575 4.2249C4.625 3.1749 5.93333 2.6499 7.5 2.6499C8.38333 2.6499 9.21667 2.8374 10 3.2124C10.7833 3.5874 11.45 4.0999 12 4.7499C12.55 4.0999 13.2167 3.5874 14 3.2124C14.7833 2.8374 15.6167 2.6499 16.5 2.6499C18.0667 2.6499 19.375 3.1749 20.425 4.2249C21.475 5.2749 22 6.58324 22 8.1499C22 10.0666 21.2917 11.8249 19.875 13.4249C18.4583 15.0249 16.85 16.6332 15.05 18.2499L13.35 19.7999C13.1667 19.9832 12.9542 20.1166 12.7125 20.1999C12.4708 20.2832 12.2333 20.3249 12 20.3249Z"
					fill="url(#paint1_radial_4780_1925)"
				/>
				<path
					d="M12 20.3249C11.7667 20.3249 11.5292 20.2832 11.2875 20.1999C11.0458 20.1166 10.8333 19.9832 10.65 19.7999L8.925 18.2249C7.15833 16.6082 5.5625 15.0041 4.1375 13.4124C2.7125 11.8207 2 10.0666 2 8.1499C2 6.58324 2.525 5.2749 3.575 4.2249C4.625 3.1749 5.93333 2.6499 7.5 2.6499C8.38333 2.6499 9.21667 2.8374 10 3.2124C10.7833 3.5874 11.45 4.0999 12 4.7499C12.55 4.0999 13.2167 3.5874 14 3.2124C14.7833 2.8374 15.6167 2.6499 16.5 2.6499C18.0667 2.6499 19.375 3.1749 20.425 4.2249C21.475 5.2749 22 6.58324 22 8.1499C22 10.0666 21.2917 11.8249 19.875 13.4249C18.4583 15.0249 16.85 16.6332 15.05 18.2499L13.35 19.7999C13.1667 19.9832 12.9542 20.1166 12.7125 20.1999C12.4708 20.2832 12.2333 20.3249 12 20.3249Z"
					fill="url(#paint2_radial_4780_1925)"
					fill-opacity="0.6"
				/>
			</g>
			<defs>
				<radialGradient
					id="paint0_radial_4780_1925"
					cx="0"
					cy="0"
					r="1"
					gradientUnits="userSpaceOnUse"
					gradientTransform="translate(18.8571 17.0424) rotate(-161.965) scale(11.4181 10.362)"
				>
					<stop stop-color="#F35D28" />
					<stop offset="1" stop-color="#F35D28" stop-opacity="0" />
				</radialGradient>
				<radialGradient
					id="paint1_radial_4780_1925"
					cx="0"
					cy="0"
					r="1"
					gradientUnits="userSpaceOnUse"
					gradientTransform="translate(7.28571 20.3249) rotate(-55.8264) scale(12.9708 13.6629)"
				>
					<stop stop-color="#FFD67E" />
					<stop offset="1" stop-color="#FFD77F" stop-opacity="0" />
				</radialGradient>
				<radialGradient
					id="paint2_radial_4780_1925"
					cx="0"
					cy="0"
					r="1"
					gradientUnits="userSpaceOnUse"
					gradientTransform="translate(15 8.9624) rotate(101.774) scale(14.7018 13.1583)"
				>
					<stop stop-color="white" />
					<stop offset="0.489583" stop-color="white" stop-opacity="0" />
				</radialGradient>
			</defs>
		</svg>
	)
}
