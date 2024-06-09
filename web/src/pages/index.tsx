import {
	createEffect,
	createSignal,
	For,
	Match,
	Show,
	Suspense,
	Switch,
} from 'solid-js'
import { store } from '~/store'
import { CDN_URL, fetchPreview } from '~/lib/api'
import FillProfilePopup from '~/components/FillProfilePopup'
import { Link } from '~/components/Link'
import { useMainButton } from '~/lib/useMainButton'
import { createQuery } from '@tanstack/solid-query'
import { A } from '@solidjs/router'

export default function Index() {
	const [profilePopup, setProfilePopup] = createSignal(false)

	const previewQuery = createQuery(() => ({
		queryKey: ['preview'],
		queryFn: () => fetchPreview(),
	}))

	const getUserLink = () => {
		if (store.user.first_name && store.user.last_name) {
			return '/users/' + store.user?.username
		} else {
			return '/users/edit'
		}
	}

	const mainButton = useMainButton()

	mainButton.hide()

	const closePopup = () => {
		setProfilePopup(false)
		window.Telegram.WebApp.CloudStorage.setItem('profilePopup', 'closed')
	}

	const updateProfilePopup = (err: unknown, value: unknown) => {
		setProfilePopup(value !== 'closed')
	}

	createEffect(() => {
		window.Telegram.WebApp.CloudStorage.getItem(
			'profilePopup',
			updateProfilePopup,
		)
		// window.Telegram.WebApp.CloudStorage.removeItem('profilePopup');
	})

	return (
		<div class="flex min-h-screen flex-col bg-secondary px-4">
			<Show when={!store.user.published_at && profilePopup()}>
				<FillProfilePopup onClose={closePopup} />
			</Show>
			<Link
				class="flex flex-row items-center justify-between py-4"
				href={getUserLink()}
			>
				<p class="text-3xl text-main">
					Bonsoir, {store.user?.first_name || store.user?.username}!
				</p>
				<Switch>
					<Match when={store.user.avatar_url}>
						<img
							class="size-11 rounded-xl border object-cover object-center"
							src={CDN_URL + '/' + store.user.avatar_url}
							alt="User Avatar"
						/>
					</Match>
					<Match when={!store.user.avatar_url}>
						<div class="flex size-11 items-center justify-center rounded-xl border-2 bg-main">
							<span class="material-symbols-rounded text-peatch-main">
								account_circle
							</span>
						</div>
					</Match>
				</Switch>
			</Link>
			<Link class="flex flex-col items-start justify-start py-4" href="/users">
				<div class="flex w-full flex-row items-center justify-start">
					<Suspense fallback={<ImagesLoader />}>
						<For each={previewQuery.data}>
							{(image, idx) => (
								<img
									src={image}
									alt="User Avatar"
									class="-ml-1 size-11 rounded-xl border object-cover object-center"
									classList={{
										'ml-0': idx() === 0,
										'z-20': idx() === 0,
										'z-10': idx() === 1,
									}}
								/>
							)}
						</For>
					</Suspense>
				</div>
				<div class="flex flex-row items-center justify-between">
					<p class="mt-2 text-3xl text-main">
						<span class="text-accent">Explore people</span> you may like to
						collaborate
					</p>
					<span class="material-symbols-rounded text-[48px] text-pink">
						maps_ugc
					</span>
				</div>
				<p class="mt-1.5 text-sm text-hint">
					Figma Wizards, Consultants, Founders, and more
				</p>
			</Link>
			<div class="h-px w-full bg-main" />
			<Link
				class="flex flex-col items-start justify-start py-4"
				href="/collaborations"
			>
				<div class="flex w-full flex-row items-center justify-start">
					<div class="z-20 flex size-11 flex-col items-center justify-center rounded-xl border bg-orange">
						<span class="material-symbols-rounded text-white">
							self_improvement
						</span>
					</div>
					<div class="z-10 -ml-1 flex size-11 flex-col items-center justify-center rounded-xl border bg-red">
						<span class="material-symbols-rounded text-white">wine_bar</span>
					</div>
					<div class="-ml-1 flex size-11 flex-col items-center justify-center rounded-xl border bg-blue">
						<span class="material-symbols-rounded text-white">
							directions_run
						</span>
					</div>
				</div>
				<div class="flex flex-row items-start justify-between">
					<p class="mt-2 text-3xl text-main">
						<span class="text-pink">Find collaborations</span> that you may be
						interested to join
					</p>
					<span class="material-symbols-rounded text-[48px] text-red">
						arrow_circle_right
					</span>
				</div>
				<p class="mt-1.5 text-sm text-hint">
					Yoga practice, Running, Grabbing a coffee, and more
				</p>
			</Link>
			<div class="h-px w-full bg-main" />
			<button
				class="flex flex-col items-start justify-start py-4"
				onClick={() =>
					window.Telegram.WebApp.openTelegramLink(
						'https://t.me/peatch_community',
					)
				}
			>
				<div class="flex flex-row items-start justify-between text-start">
					<p class="mt-2 text-3xl text-main">
						<span class="text-green">Join community</span> to talk with founders
						and users. Discuss and solve problems together
					</p>
					<span class="material-symbols-rounded text-[48px] text-green">
						forum
					</span>
				</div>
			</button>
		</div>
	)
}

const ImagesLoader = () => {
	return (
		<div class="flex w-full flex-row items-center justify-start">
			<For each={[1, 2, 3] as number[]}>
				{(image, idx) => (
					<div
						class="-ml-1 size-11 rounded-xl border bg-hint"
						classList={{
							'ml-0': idx() === 0,
							'z-20': idx() === 0,
							'z-10': idx() === 1,
						}}
					/>
				)}
			</For>
		</div>
	)
}

const ShufflePopup = () => {
	return (
		<div class="mb-4 flex flex-col items-center justify-center overflow-hidden rounded-lg bg-main py-4 text-center">
			<p class="mt-4 text-3xl text-main">Shuffle Raffle</p>
			<p class="mb-6 text-sm font-semibold text-main">
				Discover people with common interests
			</p>
			<div class="mt-3 flex flex-row items-center justify-center gap-2.5">
				<ShuffleBadge icon="wine_bar" text="Wine Lovers" color="#EF5DA8" />
				<ShuffleBadge icon="fitness_center" text="Gym Rat" color="#F9A826" />
				<ShuffleBadge icon="music_note" text="Meloman" color="#2D9CDB" />
			</div>
			<div class="mt-3 flex flex-row items-center justify-center gap-2.5">
				<ShuffleBadge icon="palette" text="Design" color="#6D214F" />
				<ShuffleBadge
					icon="translate"
					text="Language Exchange"
					color="#F2994A"
				/>
				<ShuffleBadge icon="sports_bar" text="Friends" color="#F2C94C" />
			</div>
			<A
				class="mb-2 mt-8 flex h-8 w-20 items-center justify-center rounded-lg bg-button text-sm font-semibold text-button"
				href={'/users/shuffle'}
			>
				Explore
			</A>
		</div>
	)
}

const ShuffleBadge = (props: { icon: string; text: string; color: string }) => {
	return (
		<div class="flex h-8 flex-row items-center justify-center gap-1.5 text-nowrap rounded-2xl bg-secondary px-3">
			<span class="material-symbols-rounded" style={{ color: props.color }}>
				{props.icon}
			</span>
			<p class="text-sm font-semibold text-main">{props.text}</p>
		</div>
	)
}
