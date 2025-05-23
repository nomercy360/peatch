import { createEffect, createSignal, For, onMount, Show } from 'solid-js'
import { useTranslations } from '~/lib/locale-context'
import { cn } from '~/lib/utils'
import { OpportunityResponse } from '~/gen'
import { Motion, Presence } from 'solid-motionone'

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
				<Presence exitBeforeEnter>
					<Show when={search()}>
						<Motion.button
							class="text-hint flex h-full items-center justify-center px-2.5 text-sm"
							onClick={() => setSearch('')}
							initial={{ opacity: 0, x: 10 }}
							animate={{ opacity: 1, x: 0 }}
							exit={{ opacity: 0, x: 10 }}
							transition={{ duration: 0.2 }}
						>
							Clear
						</Motion.button>
					</Show>
				</Presence>
			</div>
			<div class="flex h-11 w-full flex-row items-center justify-between">
				<div />
				<Motion.div
					class="text-hint flex h-11 items-center justify-center text-sm"
					animate={{ scale: [1, 1.1, 1] }}
					transition={{ duration: 0.3 }}
				>
					{Array.isArray(props.selected)
						? `${props.selected.length} / 10`
						: 'choose one'}
				</Motion.div>
			</div>
			<Motion.div
				class="flex w-full flex-row flex-wrap items-center justify-start gap-1"
				initial={{ opacity: 0 }}
				animate={{ opacity: 1 }}
				transition={{ duration: 0.3 }}
			>
				<Show when={!props.loading} fallback={<AnimatedLoader />}>
					<For each={filtered()}>
						{(op, index) => (
							<Motion.button
								onClick={() => onClick(op.id)}
								class={'flex h-[60px] w-full flex-row items-center justify-start gap-2.5 overflow-hidden rounded-2xl px-2.5'}
								initial={{ opacity: 0, x: -20 }}
								animate={{
									opacity: 1,
									x: 0,
									backgroundColor: includes(op.id) ? `#${op.color}` : 'var(--secondary)',
								}}
								transition={{
									duration: 0.3,
									delay: index() * 0.03,
									backgroundColor: { duration: 0.3 },
								}}
								press={{ scale: 0.98 }}
							>
								<Motion.div
									class="flex size-10 shrink-0 items-center justify-center rounded-full bg-border"
									animate={{
										scale: includes(op.id) ? [1, 1.2, 1] : 1,
									}}
									transition={{
										duration: 0.5,
									}}
								>
									<Motion.span
										class="material-symbols-rounded shrink-0"
										animate={{
											color: includes(op.id) ? `#${op.color}` : 'var(--foreground)',
										}}
										transition={{ duration: 0.3 }}
									>
										{String.fromCodePoint(parseInt(op.icon!, 16))}
									</Motion.span>
								</Motion.div>

								<div class="text-start">
									<Motion.p
										class={cn('text-xs font-semibold')}
										animate={{
											color: includes(op.id!) ? 'white' : 'var(--secondary-foreground)',
										}}
										transition={{ duration: 0.3 }}
									>
										{op.text}
									</Motion.p>
									<Motion.p
										class={cn('text-xs leading-tight')}
										animate={{
											color: includes(op.id!) ? 'rgba(255, 255, 255, 0.8)' : 'var(--secondary-foreground)',
										}}
										transition={{ duration: 0.3 }}
									>
										{op.description}
									</Motion.p>
								</div>
							</Motion.button>
						)}
					</For>
				</Show>
			</Motion.div>
		</>
	)
}

function AnimatedLoader() {
	return (
		<For each={[1, 2, 3, 4, 5, 6, 7, 8, 9]}>
			{(_, index) => (
				<Motion.div
					class="bg-main flex h-[60px] w-full flex-row items-center justify-start gap-2.5 rounded-2xl px-2.5"
					initial={{ opacity: 0 }}
					animate={{ opacity: [0.3, 0.6, 0.3] }}
					transition={{
						duration: 1.5,
						repeat: Infinity,
						delay: index() * 0.1,
					}}
				/>
			)}
		</For>
	)
}
