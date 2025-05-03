import { useMainButton } from '~/lib/useMainButton'
import { useNavigate } from '@solidjs/router'
import { createEffect, onCleanup } from 'solid-js'
import { editUser, setEditUser } from '~/store'
import TextArea from '~/components/text-area'
import { FormLayout } from '~/components/edit/layout'
import { useTranslations } from '~/lib/locale-context'

export default function Description() {
	const mainButton = useMainButton()
	const { t } = useTranslations()

	const navigate = useNavigate()

	const navigateToImageUpload = async () => {
		navigate('/users/edit/image', { state: { back: true } })
	}

	mainButton.onClick(navigateToImageUpload)

	createEffect(() => {
		if (editUser.description) {
			mainButton.enable(t('common.buttons.next'))
		} else {
			mainButton.disable(t('common.buttons.next'))
		}
	})

	onCleanup(() => {
		mainButton.offClick(navigateToImageUpload)
	})

	return (
		<FormLayout
			title={t('pages.users.edit.description.title')}
			description={t('pages.users.edit.description.description')}
			screen={5}
			totalScreens={6}
		>
			<TextArea
				value={editUser.description}
				setValue={d => setEditUser('description', d)}
				placeholder={t('pages.users.edit.description.placeholder')}
			/>
		</FormLayout>
	)
}
