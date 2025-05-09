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
				<Show when={props.search}>
					<button
						class="flex h-10 items-center justify-center px-2.5 text-sm text-secondary-foreground"
						onClick={() => props.setSearch('')}
					>
						{t('common.buttons.clear')}
					</button>
				</Show>
			</div>
			<div class="flex h-11 w-full flex-row items-center justify-between">
				<Switch>
					<Match when={filteredBadges()?.length > 0}>
						<div />
						<div class="flex h-11 items-center justify-center text-sm text-secondary-foreground">
							{props.selected.length} / 10
						</div>
					</Match>
					<Match when={filteredBadges()?.length === 0}>
						<button
							class="size-full text-start text-sm text-secondary-foreground"
							onClick={() =>
								props.selected.length < 10 && props.onCreateBadgeButtonClick()
							}
						>
							Can’t find such thing. <span class="text-accent">Create it</span>
						</button>
						<p class="text-nowrap text-sm text-secondary-foreground">
							{props.selected.length} of 10
						</p>
					</Match>
				</Switch>
			</div>
			<div class="flex min-h-fit w-full flex-row flex-wrap items-center justify-start gap-1">
				<Suspense fallback={<div>Loading...</div>}>
					<For each={filteredBadges()}>
						{badge => (
							<button
								onClick={() => onBadgeClick(badge.id!)}
								class="flex h-10 flex-row items-center justify-center gap-[5px] rounded-2xl px-2.5 transition-colors duration-200 ease-in-out"
								style={{
									'background-color': `${props.selected.includes(badge.id!) ? `#${badge.color}` : 'var(--secondary)'}`,
								}}
							>
								<span
									class="material-symbols-rounded"
									style={{
										color: `${props.selected.includes(badge.id!) ? 'white' : `#${badge.color}`}`,
									}}
								>
									{String.fromCodePoint(parseInt(badge.icon!, 16))}
								</span>
								<p
									class="text-sm font-semibold"
									classList={{
										'text-white': props.selected.includes(badge.id!),
										'text-foreground': !props.selected.includes(badge.id!),
									}}
								>
									{badge.text}
								</p>
							</button>
						)}
					</For>
				</Suspense>
			</div>
		</>
	)
}
