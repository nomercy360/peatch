import { FormLayout } from '~/components/edit/layout'
import { editCollaboration, setEditCollaboration } from '~/store'
import { useMainButton } from '~/lib/useMainButton'
import { createEffect, onCleanup } from 'solid-js'
import { useLocation, useNavigate } from '@solidjs/router'
import TextArea from '~/components/text-area'
import { useTranslations } from '~/lib/locale-context'

export default function GeneralInfo() {
	const mainButton = useMainButton()

	const navigate = useNavigate()

	const { t } = useTranslations()

	const path = useLocation().pathname

	const navigateNext = () => {
		navigate(path + '/badges')
	}

	mainButton.onClick(navigateNext)

	createEffect(() => {
		if (editCollaboration.title && editCollaboration.description) {
			mainButton.enable(t('common.buttons.next'))
		} else {
			mainButton.disable(t('common.buttons.next'))
		}
	})

	onCleanup(() => {
		mainButton.offClick(navigateNext)
	})

	return (
		<FormLayout
			title={t('pages.collaborations.edit.general.title')}
			description={t('pages.collaborations.edit.general.description')}
			screen={1}
			totalScreens={4}
		>
			<div class="mt-5 flex w-full flex-col items-center justify-start gap-3">
				<input
					maxLength={70}
					class="text-main placeholder:text-hint h-10 w-full rounded-lg bg-secondary px-2.5"
					placeholder={t('pages.collaborations.edit.general.titlePlaceholder')}
					value={editCollaboration.title}
					onInput={e => setEditCollaboration('title', e.currentTarget.value)}
				/>
				<CheckBoxInput
					text={t('pages.collaborations.edit.general.checkboxPlaceholder')}
					checked={editCollaboration.is_payable || false}
					setChecked={v => setEditCollaboration('is_payable', v)}
				/>
				<TextArea
					value={editCollaboration.description}
					setValue={d => setEditCollaboration('description', d)}
					placeholder={t('pages.collaborations.edit.general.descriptionPlaceholder')}
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
		<label class="group flex h-12 w-full cursor-pointer items-center justify-between px-1">
			<p class="text-secondary-foreground">{props.text}</p>
			<input
				type="checkbox"
				class="sr-only"
				checked={props.checked}
				onChange={() => props.setChecked(!props.checked)}
			/>
			<span class="flex size-7 items-center justify-center rounded-lg border-2">
				<svg xmlns="http://www.w3.org/2000/svg"
						 width="24"
						 height="24"
						 class="size-5 scale-0 text-accent opacity-0 group-has-[:checked]:scale-100 group-has-[:checked]:opacity-100"
						 viewBox="0 0 24 24"
						 fill="none"
						 stroke="currentColor"
						 stroke-width="2"
						 stroke-linecap="round"
						 stroke-linejoin="round"
				><path d="M20 6 9 17l-5-5" /></svg>
			</span>
		</label>
	)
}
