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
					placeholder="Name it!"
					value={editCollaboration.title}
					onInput={e => setEditCollaboration('title', e.currentTarget.value)}
				/>
				<button
					class="flex h-10 w-full items-center justify-between"
					onClick={() =>
						setEditCollaboration('is_payable', !editCollaboration.is_payable)
					}
				>
					<p class="text-sm text-main">Is it this opportunity payable?</p>
					<span
						class="size-6 rounded-lg border"
						classList={{
							'bg-button': !editCollaboration.is_payable,
							'bg-secondary': editCollaboration.is_payable,
						}}
					/>
				</button>
				<TextArea
					value={editCollaboration.description}
					setValue={d => setEditCollaboration('description', d)}
					placeholder="For example: I'm looking for a designer to participate in non-profit hackaton"
				/>
			</div>
		</FormLayout>
	)
}
