import { FormLayout } from '~/components/edit/layout'
import { useMainButton } from '~/lib/useMainButton'
import { createEffect, onCleanup, onMount } from 'solid-js'
import { editUser, setEditUser } from '~/store'
import SelectLocation from '~/components/edit/select-location'
import { useNavigate } from '@solidjs/router'
import { useTranslations } from '~/lib/locale-context'

export default function SelectBadges() {
	const mainButton = useMainButton()
	const { t } = useTranslations()

	const navigate = useNavigate()

	const navigateToDescription = async () => {
		navigate('/users/edit/description', { state: { back: true } })
	}

	onMount(() => {
		mainButton.onClick(navigateToDescription)
	})

	createEffect(() => {
		if (editUser.location.id) {
			mainButton.enable(t('common.buttons.next'))
		} else {
			mainButton.disable(t('common.buttons.next'))
		}
	})

	onCleanup(() => {
		mainButton.offClick(navigateToDescription)
	})

	return (
		<FormLayout
			title={t('pages.users.edit.location.title')}
			description={t('pages.users.edit.location.description')}
			screen={4}
			totalScreens={6}
		>
			<SelectLocation
				initialLocation={editUser.location}
				setLocation={b => setEditUser('location', b)}
			/>
		</FormLayout>
	)
}
