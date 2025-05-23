import {
	createEffect,
	createSignal,
	For,
	Match,
	Show,
	Suspense,
	Switch,
} from 'solid-js'
import { useTranslations } from '~/lib/locale-context'
import { BadgeResponse } from '~/gen'
import { Motion, Presence } from 'solid-motionone'

export function SelectBadge(props: {
	selected: string[]
	setSelected: (selected: string[]) => void
	search: string
	setSearch: (search: string) => void
	badges: BadgeResponse[]
	onCreateBadgeButtonClick: () => void
}) {
	const { t } = useTranslations()
	const [filteredBadges, setFilteredBadges] = createSignal(props.badges)

	const onBadgeClick = (badgeId?: string) => {
		if (!badgeId) return
		if (props.selected.includes(badgeId!)) {
			props.setSelected(props.selected.filter(b => b !== badgeId))
		} else if (props.selected.length < 10) {
			props.setSelected([...props.selected, badgeId!])
		}
	}

	createEffect(() => {
		if (props.badges && props.badges.length > 0) {
			setFilteredBadges(
				props.badges.filter(badge =>
					badge.text?.toLowerCase().includes(props.search.toLowerCase()),
				),
			)
		}
	})

	return (
		<>
			<div class="mt-5 flex h-10 w-full flex-row items-center justify-between rounded-lg bg-secondary px-2.5">
				<input
					autofocus
					autocomplete="off"
					autocorrect="off"
					autocapitalize="off"
					spellcheck={false}
					maxLength={50}
					class="text-main h-10 w-full bg-transparent placeholder:text-secondary-foreground focus:outline-none"
					placeholder={t('pages.collaborations.edit.badges.searchPlaceholder')}
					type="text"
					onInput={e => props.setSearch(e.currentTarget.value)}
					value={props.search}
				/>
				<Presence exitBeforeEnter>
					<Show when={props.search}>
						<Motion.button
							class="flex h-10 items-center justify-center px-2.5 text-sm text-secondary-foreground"
							onClick={() => props.setSearch('')}
							initial={{ opacity: 0, x: 10 }}
							animate={{ opacity: 1, x: 0 }}
							exit={{ opacity: 0, x: 10 }}
							transition={{ duration: 0.2 }}
						>
							{t('common.buttons.clear')}
						</Motion.button>
					</Show>
				</Presence>
			</div>
			<div class="flex h-11 w-full flex-row items-center justify-between">
				<Switch>
					<Match when={filteredBadges()?.length > 0}>
						<div />
						<Motion.div
							class="flex h-11 items-center justify-center text-sm text-secondary-foreground"
							animate={{ scale: [1, 1.1, 1] }}
							transition={{ duration: 0.3 }}
						>
							{props.selected.length} / 10
						</Motion.div>
					</Match>
					<Match when={filteredBadges()?.length === 0}>
						<Motion.button
							class="size-full text-start text-sm text-secondary-foreground"
							onClick={() =>
								props.selected.length < 10 && props.onCreateBadgeButtonClick()
							}
							hover={{ x: 5 }}
							transition={{ duration: 0.2 }}
						>
							Can't find such thing. <span class="text-accent">Create it</span>
						</Motion.button>
						<p class="text-nowrap text-sm text-secondary-foreground">
							{props.selected.length} of 10
						</p>
					</Match>
				</Switch>
			</div>
			<Motion.div
				class="flex min-h-fit w-full flex-row flex-wrap items-center justify-start gap-1"
				initial={{ opacity: 0 }}
				animate={{ opacity: 1 }}
				transition={{ duration: 0.3 }}
			>
				<Suspense fallback={<div>Loading...</div>}>
					<For each={filteredBadges()}>
						{(badge, index) => (
							<Motion.button
								onClick={() => onBadgeClick(badge.id!)}
								class="flex h-10 flex-row items-center justify-center gap-[5px] overflow-hidden rounded-2xl px-2.5"
								initial={{ opacity: 0, scale: 0.8 }}
								animate={{
									opacity: 1,
									scale: 1,
									backgroundColor: props.selected.includes(badge.id!) ? `#${badge.color}` : 'var(--secondary)',
								}}
								transition={{
									duration: 0.2,
									delay: index() * 0.02,
									backgroundColor: { duration: 0.3 },
								}}
								press={{ scale: 0.98 }}
							>
								<Motion.span
									class="material-symbols-rounded"
									animate={{
										color: props.selected.includes(badge.id!) ? 'white' : `#${badge.color}`,
									}}
									transition={{ color: { duration: 0.3 } }}
								>
									{String.fromCodePoint(parseInt(badge.icon!, 16))}
								</Motion.span>
								<Motion.p
									class="text-sm font-semibold"
									animate={{
										color: props.selected.includes(badge.id!) ? 'white' : 'var(--foreground)',
									}}
									transition={{ duration: 0.3 }}
								>
									{badge.text}
								</Motion.p>
							</Motion.button>
						)}
					</For>
				</Suspense>
			</Motion.div>
		</>
	)
}
