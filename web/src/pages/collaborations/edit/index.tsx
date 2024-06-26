import { RouteSectionProps, useNavigate, useParams } from '@solidjs/router'
import { fetchCollaboration } from '~/lib/api'
import { setEditCollaboration, setEditCollaborationId } from '~/store'
import { createEffect, createResource, Show } from 'solid-js'

export default function EditCollaboration(props: RouteSectionProps) {
	const params = useParams()

	if (!params.id) {
		setEditCollaboration({
			badge_ids: [],
			city: '',
			country: '',
			country_code: '',
			description: '',
			is_payable: false,
			opportunity_id: 0,
			title: '',
		})

		const navigate = useNavigate()
		navigate('/collaborations/edit')
		return <div>{props.children}</div>
	} else {
		const [collaboration, _] = createResource(async () => {
			return await fetchCollaboration(Number(params.id))
		})

		createEffect(() => {
			if (!collaboration.loading) {
				setEditCollaboration({
					badge_ids: collaboration().badges.map(
						(badge: { id: number }) => badge.id,
					),
					city: collaboration().city,
					country: collaboration().country,
					country_code: collaboration().country_code,
					description: collaboration().description,
					is_payable: collaboration().is_payable,
					opportunity_id: collaboration().opportunity.id,
					title: collaboration().title,
				})

				setEditCollaborationId(collaboration().id)
			}
		})

		return (
			<Show when={!collaboration.loading}>
				<div>{props.children}</div>
			</Show>
		)
	}
}
