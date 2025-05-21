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
		// If we're skipping location, ensure it's empty
		if (!editCollaboration.location.id) {
			setEditCollaboration('location', {} as any)
		}
		const created = await createCollaboration(editCollaboration)
		navigate('/collaborations/' + created.id)
	}

	const editCollab = async () => {
		// For editing, we don't allow empty location through the button
		// because the main button is disabled when no location is selected
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
		if (editCollaboration.location.id) {
			mainButton.enable(t('common.buttons.chooseAndSave'))
		} else if (!editCollaborationId()) {
			// For new collaborations, allow skipping
			mainButton.enable(t('common.buttons.skip'))
		} else {
			// For editing existing collaborations
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
				initialLocation={editCollaboration.location}
				setLocation={b => setEditCollaboration('location', b)}
			/>
		</FormLayout>
	)
}
