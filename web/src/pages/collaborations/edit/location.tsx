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
import { CreateCollaboration } from '~/gen'

export default function SelectBadges() {
	const mainButton = useMainButton()
	const { t } = useTranslations()

	const navigate = useNavigate()

	const createCollab = async () => {
		const req = {
			opportunity_id: editCollaboration.opportunity_id,
			title: editCollaboration.title,
			description: editCollaboration.description,
			badge_ids: editCollaboration.badge_ids,
			is_payable: editCollaboration.is_payable,
			location_id: editCollaboration.location ? editCollaboration.location.id : null,
		} as CreateCollaboration

		const created = await createCollaboration(req)
		navigate('/collaborations/' + created.id)
	}

	const editCollab = async () => {
		const req = {
			opportunity_id: editCollaboration.opportunity_id,
			title: editCollaboration.title,
			description: editCollaboration.description,
			badge_ids: editCollaboration.badge_ids,
			is_payable: editCollaboration.is_payable,
			location_id: editCollaboration.location ? editCollaboration.location.id : null,
		} as CreateCollaboration

		await updateCollaboration(editCollaborationId(), req)
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
