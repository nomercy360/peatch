import { SelectBadge } from '~/components/edit/selectBadge'
import { FormLayout } from '~/components/edit/layout'
import { useMainButton } from '~/lib/useMainButton'
import { useNavigate, useParams, useSearchParams } from '@solidjs/router'
import { createEffect, createSignal, onCleanup } from 'solid-js'
import { editCollaboration, setEditCollaboration } from '~/store'
import { fetchBadges } from '~/lib/api'
import { Badge } from '~/gen'
import { useTranslations } from '~/lib/locale-context'
import { useQuery } from '@tanstack/solid-query'

export default function SelectBadges() {
	const mainButton = useMainButton()

	const [badgeSearch, setBadgeSearch] = createSignal('')
	const { t } = useTranslations()

	const fetchBadgeQuery = useQuery(() => ({
		queryKey: ['badges'],
		// then push selected to the top
		queryFn: () =>
			fetchBadges().then(badges => {
				const selected = editCollaboration.badge_ids
				return [
					...selected.map((id: number) => badges.find((b: Badge) => b.id === id)),
					...badges.filter((b: Badge) => !selected.includes(b.id!)),
				]
			}),
	}))

	const navigate = useNavigate()
	const idPath = useParams().id ? '/' + useParams().id : ''

	const [searchParams] = useSearchParams()

	const navigateNext = () => {
		navigate(`/collaborations/edit${idPath}/interests`, {
			state: { back: true },
		})
	}

	const navigateCreateBadge = () => {
		navigate(
			`/collaborations/edit${idPath}/create-badge?badge_name=` + badgeSearch(),
		)
	}

	createEffect(async () => {
		if (searchParams.refetch) {
			await fetchBadgeQuery.refetch()
		}
	})

	mainButton.onClick(navigateNext)

	createEffect(() => {
		if (editCollaboration.badge_ids && editCollaboration.badge_ids.length > 0) {
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
			title={t('pages.collaborations.edit.badges.title')}
			description={t('pages.collaborations.edit.badges.description')}
			screen={2}
			totalScreens={6}
		>
			<SelectBadge
				selected={editCollaboration.badge_ids}
				setSelected={b => setEditCollaboration('badge_ids', b)}
				onCreateBadgeButtonClick={navigateCreateBadge}
				search={badgeSearch()}
				setSearch={setBadgeSearch}
				badges={fetchBadgeQuery.data!}
			/>
		</FormLayout>
	)
}
