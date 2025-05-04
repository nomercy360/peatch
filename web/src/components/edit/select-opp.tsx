import { createEffect, createSignal, For, onMount, Show } from 'solid-js'
import { useTranslations } from '~/lib/locale-context'
import { cn } from '~/lib/utils'
import { OpportunityResponse } from '~/gen'

export function SelectOpportunity(props: {
	selected: string[] | string
	setSelected: (selected: string[] | string) => void
	opportunities: OpportunityResponse[]
	loading: boolean
}) {
	const [filtered, setFiltered] = createSignal<OpportunityResponse[]>([])
	const [search, setSearch] = createSignal('')
	const { t } = useTranslations()

	const onClick = (oppId?: string) => {
		if (!oppId) return
		if (Array.isArray(props.selected)) {
			if (props.selected.includes(oppId)) {
				props.setSelected(props.selected.filter(b => b !== oppId))
			} else if (props.selected.length < 10) {
				props.setSelected([...props.selected, oppId])
			}
		} else {
			props.setSelected(oppId)
		}
	}

	createEffect(() => {
		if (props.opportunities && props.opportunities.length > 0) {
			setFiltered(
				props.opportunities.filter(
					op =>
						op.text?.toLowerCase().includes(search().toLowerCase()) ||
						op.description?.toLowerCase().includes(search().toLowerCase()),
				),
			)
		}
	})

	const includes = (oppId?: string) => {
		if (!oppId) return false
		if (Array.isArray(props.selected)) {
			return props.selected.includes(oppId)
		}
		return props.selected === oppId
	}

	onMount(() => {
		if (Array.isArray(props.selected)) {
			setFiltered(props.opportunities)
		}
	})

	return (
		<>
			<div class="mt-5 flex h-10 w-full flex-row items-center justify-between rounded-lg bg-secondary px-2.5">
				<input
					class="text-main placeholder:text-hint h-10 w-full bg-transparent focus:outline-none"
					placeholder={t('pages.collaborations.edit.interests.searchPlaceholder')}
					type="text"
					onInput={e => setSearch(e.currentTarget.value)}
					value={search()}
				/>
				<Show when={search()}>
					<button
						class="text-hint flex h-full items-center justify-center px-2.5 text-sm"
						onClick={() => setSearch('')}
					>
						Clear
					</button>
				</Show>
			</div>
			<div class="flex h-11 w-full flex-row items-center justify-between">
				<div />
				<div class="text-hint flex h-11 items-center justify-center text-sm">
					{Array.isArray(props.selected)
						? `${props.selected.length} / 10`
						: 'choose one'}
				</div>
			</div>
			<div class="flex w-full flex-row flex-wrap items-center justify-start gap-1">
				<Show when={!props.loading} fallback={<Loader />}>
					<For each={filtered()}>
						{op => (
							<button
								onClick={() => onClick(op.id)}
								class={'flex h-[60px] w-full flex-row items-center justify-start gap-2.5 rounded-2xl px-2.5'}
								style={{ 'background-color': includes(op.id) ? `#${op.color}` : 'var(--secondary)' }}
							>
								<div class="flex size-10 shrink-0 items-center justify-center rounded-full bg-border">
									<span class="material-symbols-rounded shrink-0">
										{String.fromCodePoint(parseInt(op.icon!, 16))}
									</span>
								</div>

								<div class="text-start">
									<p class={cn('text-xs font-semibold', includes(op.id!) ? 'text-white' : 'text-secondary-foreground')}>
										{op.text}
									</p>
									<p
										class={cn('text-xs leading-tight', includes(op.id!) ? 'text-white/80' : 'text-secondary-foreground')}>
										{op.description}
									</p>
								</div>
							</button>
						)}
					</For>
				</Show>
			</div>
		</>
	)
}

function Loader() {
	return (
		<For each={[1, 2, 3, 4, 5, 6, 7, 8, 9]}>
			{() => (
				<div class="bg-main flex h-[60px] w-full flex-row items-center justify-start gap-2.5 rounded-2xl px-2.5" />
			)}
		</For>
	)
}
