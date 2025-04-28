import { FormLayout } from '~/components/edit/layout'
import { useMainButton } from '~/lib/useMainButton'
import { createEffect, onCleanup, onMount } from 'solid-js'
import { editUser, setEditUser } from '~/store'
import SelectLocation from '~/components/edit/selectLocation'
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
		if (editUser.country && editUser.country_code) {
			mainButton.enable('Next')
		} else {
			mainButton.disable('Next')
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
				city={editUser.city}
				setCity={c => setEditUser('city', c)}
				country={editUser.country}
				setCountry={c => setEditUser('country', c)}
				countryCode={editUser.country_code}
				setCountryCode={c => setEditUser('country_code', c)}
			/>
		</FormLayout>
	)
}
