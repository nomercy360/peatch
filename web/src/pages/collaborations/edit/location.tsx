import { FormLayout } from '~/components/edit/layout'
import { useMainButton } from '~/lib/useMainButton'
import { useNavigate } from '@solidjs/router'
import { createEffect, onCleanup } from 'solid-js'
import {
	editCollaboration,
	editCollaborationId,
	setEditCollaboration,
} from '~/store'
import SelectLocation from '~/components/edit/select-location'
import { createCollaboration, updateCollaboration } from '~/lib/api'
import { useTranslations } from '~/lib/locale-context'

export default function SelectBadges() {
	const mainButton = useMainButton()
	const { t } = useTranslations()

	const navigate = useNavigate()

	const createCollab = async () => {
		const created = await createCollaboration(editCollaboration)
		navigate('/collaborations/' + created.id)
	}

	const editCollab = async () => {
		await updateCollaboration(editCollaborationId(), editCollaboration)
		navigate('/collaborations/' + editCollaborationId() + '?refetch=true')
	}

	const createOrEditCollab = async () => {
		if (editCollaborationId()) {
			await editCollab()
		} else {
			await createCollab()
		}
	}

	mainButton.onClick(createOrEditCollab)

	createEffect(() => {
		if (editCollaboration.country && editCollaboration.country_code) {
			mainButton.enable(t('common.buttons.chooseAndSave'))
		} else {
			mainButton.disable(t('common.buttons.chooseAndSave'))
		}
	})

	onCleanup(() => {
		mainButton.offClick(createOrEditCollab)
	})

	return (
		<FormLayout
			title={t('pages.collaborations.edit.location.title')}
			description={t('pages.collaborations.edit.location.description')}
			screen={4}
			totalScreens={6}
		>
			<SelectLocation
				city={editCollaboration.city}
				setCity={c => setEditCollaboration('city', c)}
				country={editCollaboration.country}
				setCountry={c => setEditCollaboration('country', c)}
				countryCode={editCollaboration.country_code}
				setCountryCode={c => setEditCollaboration('country_code', c)}
			/>
		</FormLayout>
	)
}
