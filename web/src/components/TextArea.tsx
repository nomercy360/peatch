import { createEffect, createSignal, Match, Switch } from 'solid-js'

export default function TextArea(props: {
	value: string
	setValue: (value: string) => void
	placeholder: string
}) {
	const [count, setCount] = createSignal(0)
	const maxLength = 500

	createEffect(() => {
		setCount(props.value.length)
	})

	return (
		<div class="relative mt-5 h-80 w-full rounded-lg bg-secondary pb-6">
			<textarea
				class="size-full resize-none bg-transparent p-2.5 text-main placeholder:text-hint focus:outline-none w-full h-full"
				placeholder={props.placeholder}
				value={props.value}
				onInput={e => props.setValue((e.target as HTMLTextAreaElement).value)}
				autocomplete="off"
				autocapitalize="off"
				spellcheck={false}
				maxLength={maxLength}
			/>
			<Switch>
				<Match when={count() > 0}>
					<div class="absolute bottom-2 left-2 text-sm text-hint">
						{count()} / {maxLength}
					</div>
				</Match>
				<Match when={count() === 0}>
					<div class="absolute bottom-2 left-2 text-sm text-hint">
						max {maxLength} characters
					</div>
				</Match>
			</Switch>
		</div>
	)
}
