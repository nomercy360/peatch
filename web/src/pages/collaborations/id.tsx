import {
	createEffect,
	createSignal,
	For,
	Match,
	onCleanup,
	Suspense,
	Switch,
} from 'solid-js'
import { useNavigate, useParams, useSearchParams } from '@solidjs/router'
import { fetchCollaboration } from '~/lib/api'
import { store } from '~/store'
import { useMainButton } from '~/lib/useMainButton'
import { useQuery } from '@tanstack/solid-query'

export default function Collaboration() {
	const mainButton = useMainButton()

	const [isCurrentUserCollab, setIsCurrentUserCollab] = createSignal(false)

	const navigate = useNavigate()
	const params = useParams()
	const [searchParams] = useSearchParams()

	const collabId = params.id

	const query = useQuery(() => ({
		queryKey: ['collaborations', collabId],
		queryFn: () => fetchCollaboration(collabId),
	}))

	createEffect(async () => {
		if (searchParams.refetch) {
			await query.refetch()
		}
	})

	createEffect(() => {
		if (query.data?.id) {
			setIsCurrentUserCollab(store.user.id === query.data.user.id)
		}
	})

	const navigateToEdit = () => {
		navigate('/collaborations/edit/' + collabId, {
			state: { from: '/collaborations/' + collabId },
		})
	}

	createEffect(() => {
		if (isCurrentUserCollab()) {
			mainButton.enable('Edit')
			mainButton.onClick(navigateToEdit)
		}
		onCleanup(() => {
			mainButton.offClick(navigateToEdit)
		})
	})

	onCleanup(async () => {
		mainButton.hide()
	})

	return (
		<Suspense fallback={<Loader />}>
			<Switch>
				<Match when={!query.isLoading}>
					<div class="h-fit min-h-screen bg-secondary">
						<Switch>
							<Match when={isCurrentUserCollab() && !query.data.published_at}>
								<ActionButton text="Edit" onClick={navigateToEdit} />
							</Match>
						</Switch>
						<div
							class="flex w-full flex-col items-start justify-start px-4 pb-5 pt-4"
							style={{
								'background-color': `#${query.data.opportunity.color}`,
							}}
						>
							<span class="material-symbols-rounded text-[48px] text-white">
								{String.fromCodePoint(
									parseInt(query.data.opportunity.icon!, 16),
								)}
							</span>
							<p class="text-3xl text-white">{query.data.opportunity.text}:</p>
							<p class="text-3xl text-white">{query.data.title}:</p>
							<div class="mt-4 flex w-full flex-row items-center justify-start gap-2">
								<img
									class="size-11 rounded-xl object-cover"
									src={`https://assets.peatch.io/cdn-cgi/image/width=100/${query.data.user?.avatar_url}`}
									alt="User Avatar"
								/>
								<div>
									<p class="text-sm font-bold text-white">
										{query.data.user?.first_name} {query.data.user?.last_name}:
									</p>
									<p class="text-sm text-white">{query.data.user?.title}</p>
								</div>
							</div>
						</div>
						<div class="px-4 py-2.5">
							<p class="mt-1 text-start text-sm font-normal text-secondary-foreground">
								{query.data.description}
							</p>
							<div class="mt-5 flex flex-row flex-wrap items-center justify-start gap-1">
								<For each={query.data.badges}>
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
						</div>
					</div>
				</Match>
			</Switch>
		</Suspense>
	)
}

const ActionButton = (props: { text: string; onClick: () => void }) => {
	return (
		<button
			class="absolute right-4 top-4 z-10 h-9 w-[90px] rounded-xl bg-black/80 px-2.5 text-white"
			onClick={() => props.onClick()}
		>
			{props.text}
		</button>
	)
}

const Loader = () => {
	return (
		<div class="flex h-screen flex-col items-start justify-start bg-secondary">
			<div class="bg-main h-[260px] w-full" />
			<div class="flex flex-col items-start justify-start p-4">
				<div class="bg-main h-36 w-full rounded" />
				<div class="mt-4 flex w-full flex-row flex-wrap items-center justify-start gap-2">
					<div class="bg-main h-10 w-40 rounded-2xl" />
					<div class="bg-main h-10 w-32 rounded-2xl" />
					<div class="bg-main h-10 w-36 rounded-2xl" />
					<div class="bg-main h-10 w-24 rounded-2xl" />
					<div class="bg-main h-10 w-40 rounded-2xl" />
					<div class="bg-main h-10 w-28 rounded-2xl" />
					<div class="bg-main h-10 w-32 rounded-2xl" />
				</div>
			</div>
		</div>
	)
}
