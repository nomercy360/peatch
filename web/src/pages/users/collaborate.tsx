import {
	CDN_URL,
	createUserCollaboration,
	fetchProfile,
	findUserCollaborationRequest,
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
import TextArea from '~/components/TextArea'
import { usePopup } from '~/lib/usePopup'
import { useMainButton } from '~/lib/useMainButton'
import BadgeList from '~/components/BadgeList'
import ActionDonePopup from '~/components/ActionDonePopup'
import { createQuery } from '@tanstack/solid-query'
import { LocationBadge } from '~/components/location-badge'

export default function Collaborate() {
	const params = useParams()
	const username = params.handle

	const [created, setCreated] = createSignal(false)
	const mainButton = useMainButton()
	const { showConfirm } = usePopup()
	const navigate = useNavigate()

	const backToProfile = () => {
		navigate(`/users/${username}`, { state: { from: '/users' } })
	}

	const [message, setMessage] = createSignal('')

	const query = createQuery(() => ({
		queryKey: ['profiles', username],
		queryFn: () => fetchProfile(username),
	}))

	const [existedRequest] = createResource(async () => {
		try {
			return await findUserCollaborationRequest(username)
		} catch (e: unknown) {
			if ((e as { status: number }).status === 404) {
				return null
			}
		}
	})

	const postCollaboration = async () => {
		if (!store.user.published_at) {
			showConfirm(
				'You must publish your profile first',
				(ok: boolean) =>
					ok && navigate('/users/edit', { state: { back: true } }),
			)
			return
		}
		try {
			await createUserCollaboration(query.data.id, message())
			setCreated(true)
		} catch (e) {
			console.error(e)
		}
	}

	createEffect(() => {
		if (created() || existedRequest()) {
			mainButton.offClick(postCollaboration)
			mainButton.onClick(backToProfile)
			mainButton.enable(`Back to ${query.data.first_name}'s profile`)
		} else if (!existedRequest.loading && !existedRequest()) {
			mainButton.onClick(postCollaboration)
			if (message()) {
				mainButton.enable('Send message')
			} else {
				mainButton.disable('Send message')
			}
		}

		onCleanup(() => {
			mainButton.offClick(postCollaboration)
			mainButton.offClick(backToProfile)
		})
	})

	return (
		<Show when={query.isSuccess}>
			<Switch>
				<Match when={created()}>
					<ActionDonePopup
						action="Message sent"
						description={`Once ${query.data.first_name} accepts your invitation, we'll share your contacts`}
						callToAction={`There are 12 people with a similar profiles like ${query.data.first_name}`}
					/>
				</Match>
				<Match when={existedRequest.loading && !query.data}>
					<div />
				</Match>
				<Match when={!existedRequest.loading && query.data}>
					<Show when={existedRequest()}>
						<ActionDonePopup
							action="Message sent"
							description={`Once ${query.data.first_name} accepts your invitation, we'll share your contacts`}
							callToAction={`There are 12 people with a similar profiles like ${query.data.first_name}`}
						/>
					</Show>
					<Show when={!existedRequest()}>
						<div class="flex flex-col items-center justify-center p-4">
							<div class="mb-4 mt-1 flex flex-col items-center justify-center text-center">
								<p class="max-w-[220px] text-3xl">
									Collaborate with {query.data.first_name}
								</p>
								<div class="my-5 flex w-full flex-row items-center justify-center">
									<img
										class="z-10 size-24 rounded-3xl border-2 border-secondary object-cover object-center"
										src={CDN_URL + '/' + store.user.avatar_url}
										alt="User Avatar"
									/>
									<img
										class="-ml-4 size-24 rounded-3xl border-2 border-secondary object-cover object-center"
										src={CDN_URL + '/' + query.data.avatar_url}
										alt="User Avatar"
									/>
								</div>
								<Show when={query.data.badges && query.data?.badges.length > 0}>
									<BadgeList badges={query.data.badges!} position="center">
										<LocationBadge
											country={query.data.country!}
											city={query.data.city!}
											countryCode={query.data.country_code!}
										/>
									</BadgeList>
								</Show>
							</div>
							<TextArea
								value={message()}
								setValue={(value: string) => setMessage(value)}
								placeholder="Write a message to start collaboration"
							/>
						</div>
					</Show>
				</Match>
			</Switch>
		</Show>
	)
}
