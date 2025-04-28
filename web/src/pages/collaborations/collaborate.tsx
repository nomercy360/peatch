import {
	CDN_URL,
	createCollaborationRequest,
	fetchCollaboration,
	findCollaborationRequest,
} from '~/lib/api'
import { useNavigate, useParams } from '@solidjs/router'
import { store } from '~/store'
import {
	createEffect,
	createResource,
	createSignal,
	Match,
	onCleanup,
	Show,
	Switch,
} from 'solid-js'

import { usePopup } from '~/lib/usePopup'
import { useMainButton } from '~/lib/useMainButton'
import ActionDonePopup from '~/components/action-done-popup'
import TextArea from '~/components/text-area'
import { Badge } from '~/gen/types'
import { createQuery } from '@tanstack/solid-query'

export default function Collaborate() {
	const params = useParams()
	const collabId = Number(params.id)

	const [created, setCreated] = createSignal(false)
	const mainButton = useMainButton()
	const { showAlert } = usePopup()
	const navigate = useNavigate()

	const backToCollab = () => {
		navigate(`/collaborations/${collabId}`, {
			state: { from: '/collaborations' },
		})
	}

	const [message, setMessage] = createSignal('')

	const query = createQuery(() => ({
		queryKey: ['collaborations', collabId],
		queryFn: () => fetchCollaboration(collabId),
	}))

	const [existedRequest] = createResource(async () => {
		try {
			return await findCollaborationRequest(collabId)
		} catch (e: unknown) {
			if ((e as { status: number }).status === 404) {
				return null
			}
		}
	})

	const postCollaboration = async () => {
		if (!store.user.published_at) {
			showAlert('You must publish your profile first')
			return
		}
		try {
			await createCollaborationRequest(collabId, message())
			setCreated(true)
		} catch (e) {
			console.error(e)
		}
	}

	createEffect(() => {
		if (created() || existedRequest()) {
			mainButton.offClick(postCollaboration)
			mainButton.onClick(backToCollab)
			mainButton.enable('Back to collaboration')
		} else if (!existedRequest.loading && !existedRequest()) {
			mainButton.onClick(postCollaboration)
			if (message() !== '') {
				mainButton.enable('Send message')
			} else {
				mainButton.disable('Send message')
			}
		}

		onCleanup(() => {
			mainButton.offClick(postCollaboration)
			mainButton.offClick(backToCollab)
		})
	})

	return (
		<Switch>
			<Match when={created()}>
				<ActionDonePopup
					action="Message sent"
					description={`Once ${query.data.user.first_name} accepts your invitation, we'll share your contacts`}
					callToAction={'There are 12 more collaborations like this'}
				/>
			</Match>
			<Match when={existedRequest.loading && !query.data}>
				<div />
			</Match>
			<Match when={!existedRequest.loading && query.data}>
				<Show when={existedRequest()}>
					<ActionDonePopup
						action="Message sent"
						description={`Once ${query.data.user.first_name} accepts your invitation, we'll share your contacts`}
						callToAction={'There are 12 more collaborations like this'}
					/>
				</Show>
				<Show when={!existedRequest()}>
					<div class="flex flex-col items-center justify-center bg-secondary p-4">
						<div class="mb-4 mt-1 flex flex-col items-center justify-center text-center">
							<p class="max-w-[220px] text-3xl text-main">Express interest</p>
							<p class="mt-2 text-sm text-secondary">
								Say hello, ask question and let {query.data.user.first_name}{' '}
								know you’re interested
							</p>
							<div class="my-5 flex w-full flex-row items-center justify-center">
								<img
									class="z-10 size-24 rounded-3xl border-2 border-secondary object-cover object-center"
									src={CDN_URL + '/' + store.user.avatar_url}
									alt="User Avatar"
								/>
								<img
									class="-ml-4 size-24 rounded-3xl border-2 border-secondary object-cover object-center"
									src={CDN_URL + '/' + query.data.user.avatar_url}
									alt="User Avatar"
								/>
							</div>
							<p class="text-secondary">
								{query.data.user.first_name} is looking for a{' '}
								{query.data.badges.map((b: Badge) => b.text).join(', ')}
							</p>
						</div>
						<TextArea
							value={message()}
							setValue={(value: string) => setMessage(value)}
							placeholder={`Hi, ${query.data.user.first_name}! My name is ${store.user.first_name}. I'm interested in collaborating with you on your ${query.data.opportunity.text}. If you're too, please reach out to me.`}
						/>
					</div>
				</Show>
			</Match>
		</Switch>
	)
}
