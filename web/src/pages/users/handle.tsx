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
	fetchProfile,
	followUser,
} from '~/lib/api'
import { addToast } from '~/components/toast'
import { setUser, store } from '~/store'
import { useMainButton } from '~/lib/useMainButton'
import { queryClient } from '~/App'
import { UserProfileResponse, verificationStatus } from '~/gen'
import { useTranslations } from '~/lib/locale-context'
import { useMutation, useQuery } from '@tanstack/solid-query'

export default function UserProfilePage() {
	const mainButton = useMainButton()

	const navigate = useNavigate()

	const params = useParams()
	const [searchParams] = useSearchParams()

	const id = params.handle

	const { t } = useTranslations()

	const query = useQuery(() => ({
		queryKey: ['profiles', id],
		queryFn: () => fetchProfile(id),
		retry: 1,
	}))

	const followMutate = useMutation(() => ({
		mutationFn: (id: string) => followUser(id),
		retry: 0,
		onMutate: async (id: string) => {
			await queryClient.cancelQueries({ queryKey: ['profiles', id] })

			const previousProfile = queryClient.getQueryData(['profiles', id]) as UserProfileResponse

			queryClient.setQueryData(['profiles', id], (old: UserProfileResponse) => {
				if (old) {
					return {
						...old,
						is_following: true,
					}
				}
				return old
			})

			return { previousProfile }
		},
		onSuccess: () => {
			addToast(t('pages.users.followSuccess'), 'success')
		},
		onError: (error: any, _id: string, context?: { previousProfile?: UserProfileResponse }) => {
			if (context?.previousProfile) {
				queryClient.setQueryData(['profiles', context.previousProfile.id], context.previousProfile)
			}

			if (error.botBlocked) {
				const username = error.username
				if (username) {
					addToast(
						t('pages.users.botBlocked'),
						'warning',
						{
							text: t('pages.users.messageUser'),
							onClick: () => {
								window.Telegram.WebApp.openTelegramLink(`https://t.me/${username}`)
							},
						},
					)
				} else {
					addToast(t('pages.users.botBlocked'), 'warning')
				}
			} else {
				addToast(t('pages.users.followError'), 'error')
			}
		},
	}))

	const isCurrentUserProfile = store.user.id === id

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

	const follow = async () => {
		if (!query.data) return
		followMutate.mutate(query.data.id)
		window.Telegram.WebApp.HapticFeedback.impactOccurred('light')
	}

	createEffect(() => {
		if (isCurrentUserProfile) {
			mainButton.enable(t('common.buttons.edit'))
			mainButton.onClick(navigateToEdit)
		}
	})

	onCleanup(() => {
		mainButton.offClick(navigateToEdit)
	})


	onCleanup(async () => {
		mainButton.hide()
	})

	function shareURL() {
		const url =
			'https://t.me/share/url?' +
			new URLSearchParams({
				url: 'https://t.me/peatch_bot/app?startapp=u' + id,
			}).toString() + '&text=' +
			t('pages.users.shareURLText', { name: query.data.name })

		window.Telegram.WebApp.openTelegramLink(url)
	}

	const [badgesExpanded, setBadgesExpanded] = createSignal(false)

	const [opportunitiesExpanded, setOpportunitiesExpanded] = createSignal(false)

	const showInfoPopup = () => {
		window.Telegram.WebApp.showAlert(
			t('pages.users.verificationStatusDenied'),
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
					<div class="flex h-fit min-h-screen flex-col items-center p-2 text-center">
						<Show
							when={isCurrentUserProfile && store.user.verification_status === verificationStatus.VerificationStatusDenied}>
							<button
								onClick={showInfoPopup}
								class="absolute left-4 top-4 flex size-7 items-center justify-center rounded-lg bg-secondary"
							>
								<span class="material-symbols-rounded text-[16px] text-secondary-foreground">
									visibility_off
								</span>
							</button>
						</Show>
						<img
							alt="User Avatar"
							class="relative aspect-square w-32 rounded-xl bg-cover bg-center object-cover"
							src={`https://assets.peatch.io/cdn-cgi/image/width=400/${query.data.avatar_url}`}
						/>
						<Show when={!isCurrentUserProfile}>
							<button
								onClick={shareURL}
								class="absolute right-3 top-3 flex size-8 flex-row items-center justify-center rounded-lg border bg-secondary px-3 text-accent-foreground transition-all duration-300"
							>
								<span class="material-symbols-rounded text-[16px]">
									ios_share
								</span>
							</button>
						</Show>
						<Show when={!isCurrentUserProfile}>
							<button
								class={`mt-4 flex h-8 flex-row items-center space-x-1 rounded-2xl border px-3 transition-all duration-300 ${
									query.data.is_following
										? 'border-secondary bg-secondary text-secondary-foreground'
										: 'border-primary bg-primary text-primary-foreground'
								}`}
								onClick={() => follow()}
								disabled={query.data.is_following}
							>
								<span class="material-symbols-rounded text-[16px]">
									{query.data.is_following ? 'check' : 'waving_hand'}
								</span>
								<span class="text-sm">
									{query.data.is_following ? t('pages.users.saidHi') : t('pages.users.sayHi')}
								</span>
							</button>
						</Show>
						<div class="w-full px-4 py-2.5">
							<p class="text-3xl font-semibold capitalize text-primary">
								{query.data.name}:
							</p>
							<p class="text-3xl capitalize">{query.data.title}</p>
							<p class="mt-1 text-start text-sm font-normal text-secondary-foreground">
								{query.data.description}
							</p>
							<div class="mt-3 flex flex-row flex-wrap items-center justify-start gap-1">
								<For
									each={
										badgesExpanded()
											? query.data.badges
											: query.data.badges.slice(0, 3)
									}
								>
									{badge => (
										<div
											class="flex h-8 flex-row items-center justify-center gap-1 rounded-xl border px-2"
											style={{
												'background-color': `#${badge.color}`,
												'border-color': `#${badge.color}`,
											}}
										>
											<span class="material-symbols-rounded text-sm text-white">
												{String.fromCodePoint(parseInt(badge.icon!, 16))}
											</span>
											<p class="text-xs font-semibold text-white">
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
							<p class="pb-1 pt-3 text-start text-xl font-extrabold">
								{t('pages.users.availableFor')}
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
											class="flex h-14 w-full flex-row items-center justify-start gap-2 rounded-xl bg-secondary px-2 text-secondary-foreground"
										>
											<div class="flex size-8 shrink-0 items-center justify-center rounded-full text-white"
													 style={{ 'background-color': `#${op.color}` }}
											>
												<span class="material-symbols-rounded text-sm">
													{String.fromCodePoint(parseInt(op.icon!, 16))}
												</span>
											</div>
											<div class="text-start">
												<p class="text-xs font-semibold text-foreground">{op.text}</p>
												<p class="text-[10px] leading-tight">
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
		</div>
	)
}

const ExpandButton = (props: {
	expanded: boolean
	setExpanded: (val: boolean) => void
}) => {
	return (
		<button
			class="flex h-8 w-full items-center justify-start rounded-xl bg-transparent text-xs font-medium text-secondary-foreground"
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
