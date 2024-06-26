import { FormLayout } from '~/components/edit/layout'
import { editCollaboration, setEditCollaboration } from '~/store'
import { useMainButton } from '~/lib/useMainButton'
import { createEffect, onCleanup } from 'solid-js'
import { useLocation, useNavigate } from '@solidjs/router'
import TextArea from '~/components/TextArea'

export default function GeneralInfo() {
	const mainButton = useMainButton()

	const navigate = useNavigate()
	const path = useLocation().pathname

	const navigateNext = () => {
		navigate(path + '/badges')
	}

	mainButton.onClick(navigateNext)

	createEffect(() => {
		if (editCollaboration.title && editCollaboration.description) {
			mainButton.enable('Next')
		} else {
			mainButton.disable('Next')
		}
	})

	onCleanup(() => {
		mainButton.offClick(navigateNext)
	})

	return (
		<FormLayout
			title="Describe collaboration"
			description="This will help people to understand it clearly"
			screen={1}
			totalScreens={4}
		>
			<div class="mt-5 flex w-full flex-col items-center justify-start gap-3">
				<input
					maxLength={70}
					class="h-10 w-full rounded-lg bg-main px-2.5 text-main placeholder:text-hint"
					placeholder="Looking for a product designer"
					value={editCollaboration.title}
					onInput={e => setEditCollaboration('title', e.currentTarget.value)}
				/>
				<CheckBoxInput
					text="Is it this opportunity payable?"
					checked={editCollaboration.is_payable || false}
					setChecked={v => setEditCollaboration('is_payable', v)}
				/>
				<TextArea
					value={editCollaboration.description}
					setValue={d => setEditCollaboration('description', d)}
					placeholder="For example: I'm looking for a designer to participate in non-profit hackaton"
				/>
			</div>
		</FormLayout>
	)
}

const CheckBoxInput = (props: {
	text: string
	checked: boolean
	setChecked: (value: boolean) => void
}) => {
	return (
		<label class="group flex h-10 w-full cursor-pointer items-center justify-between">
			<p class="text-secondary">{props.text}</p>
			<input
				type="checkbox"
				class="sr-only"
				checked={props.checked}
				onChange={() => props.setChecked(!props.checked)}
			/>
			<span class="flex size-7 items-center justify-center rounded-lg border">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="3"
					stroke-linecap="round"
					stroke-linejoin="round"
					class="size-5 scale-0 text-accent opacity-0 group-has-[:checked]:scale-100 group-has-[:checked]:opacity-100"
				>
					<path d="M5 12l5 5l10 -10" />
				</svg>
			</span>
		</label>
	)
}
