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
			mainButton.enable('Next')
		} else {
			mainButton.disable('Next')
		}
	})

	onCleanup(() => {
		mainButton.offClick(navigateToImageUpload)
	})

	return (
		<FormLayout
			title={t('pages.users.edit.description.title')}
			description="Tell others about your backround, achievments and goals"
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
