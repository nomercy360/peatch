import { FormLayout } from '~/components/edit/layout'
import { editUser, setEditUser } from '~/store'
import { useMainButton } from '~/lib/useMainButton'
import { createEffect, onCleanup, onMount } from 'solid-js'
import { useNavigate } from '@solidjs/router'
import { useTranslations } from '~/lib/locale-context'

export default function GeneralInfo() {
	const mainButton = useMainButton()
	const { t } = useTranslations()

	const navigate = useNavigate()

	const navigateNext = () => {
		navigate('/users/edit/badges', { state: { back: true } })
	}

	onMount(() => {
		mainButton.onClick(navigateNext)
		window.Telegram.WebApp.enableClosingConfirmation()
	})

	createEffect(() => {
		if (editUser.first_name && editUser.last_name && editUser.title) {
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
			title={t('pages.users.edit.general.title')}
			description={t('pages.users.edit.general.description')}
			screen={1}
			totalScreens={6}
		>
			<div class="mt-5 flex w-full flex-col items-center justify-start gap-3">
				<input
					class="text-main h-10 w-full rounded-lg bg-secondary px-2.5 placeholder:text-secondary-foreground"
					placeholder={t('pages.users.edit.general.firstName')}
					autocomplete="given-name"
					maxLength={50}
					value={editUser.first_name}
					onInput={e => setEditUser('first_name', e.currentTarget.value)}
				/>
				<input
					class="text-main h-10 w-full rounded-lg bg-secondary px-2.5 placeholder:text-secondary-foreground"
					placeholder={t('pages.users.edit.general.lastName')}
					autocomplete="family-name"
					maxLength={50}
					value={editUser.last_name}
					onInput={e => setEditUser('last_name', e.currentTarget.value)}
				/>
				<input
					class="text-main h-10 w-full rounded-lg bg-secondary px-2.5 placeholder:text-secondary-foreground"
					placeholder={t('pages.users.edit.general.jobTitle')}
					maxLength={70}
					value={editUser.title}
					onInput={e => setEditUser('title', e.currentTarget.value)}
				/>
			</div>
		</FormLayout>
	)
}
