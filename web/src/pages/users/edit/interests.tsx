import { FormLayout } from '~/components/edit/layout'
import { useMainButton } from '~/lib/useMainButton'
import { useNavigate } from '@solidjs/router'
import { createEffect, onCleanup } from 'solid-js'
import { editUser, setEditUser } from '~/store'
import { fetchOpportunities } from '~/lib/api'
import { SelectOpportunity } from '~/components/edit/select-opp'
import { useTranslations } from '~/lib/locale-context'
import { useQuery } from '@tanstack/solid-query'

export default function SelectOpportunities() {
	const mainButton = useMainButton()
	const { t } = useTranslations()

	const navigate = useNavigate()

	const navigateNext = () => {
		navigate('/users/edit/location', { state: { back: true } })
	}

	const fetchOpportunityQuery = useQuery(() => ({
		queryKey: ['opportunities'],
		queryFn: () => fetchOpportunities(),
	}))

	mainButton.onClick(navigateNext)

	createEffect(() => {
		if (editUser.opportunity_ids.length) {
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
			title={t('pages.users.edit.interests.title')}
			description={t('pages.users.edit.interests.description')}
			screen={3}
			totalScreens={6}
		>
			<SelectOpportunity
				selected={editUser.opportunity_ids}
				setSelected={b => setEditUser('opportunity_ids', b as any)}
				opportunities={fetchOpportunityQuery.data}
				loading={fetchOpportunityQuery.isLoading}
			/>
		</FormLayout>
	)
}
