import { FormLayout } from '~/components/edit/layout'
import { editUser, setEditUser } from '~/store'
import { useMainButton } from '~/lib/useMainButton'
import { createEffect, onCleanup, onMount } from 'solid-js'
import { useNavigate } from '@solidjs/router'

export default function GeneralInfo() {
	const mainButton = useMainButton()

	const navigate = useNavigate()

	const navigateNext = () => {
		navigate('/users/edit/badges', { state: { back: true } })
	}

	onMount(() => {
		mainButton.onClick(navigateNext)
		window.Telegram.WebApp.enableClosingConfirmation()
	})

	createEffect(() => {
		if (editUser.first_name && editUser.last_name && editUser.title) {
			mainButton.enable('Next')
		} else {
			mainButton.disable('Next')
		}
	})

	onCleanup(() => {
		mainButton.offClick(navigateNext)
	})

	return (
		<FormLayout
			title="Introduce yourself"
			description="It will appears in your profile card, everyone will see it"
			screen={1}
			totalScreens={6}
		>
			<div class="mt-5 flex w-full flex-col items-center justify-start gap-3">
				<input
					class="h-10 w-full rounded-lg bg-secondary px-2.5 text-main placeholder:text-secondary-foreground"
					placeholder="First Name"
					autocomplete="given-name"
					maxLength={50}
					value={editUser.first_name}
					onInput={e => setEditUser('first_name', e.currentTarget.value)}
				/>
				<input
					class="h-10 w-full rounded-lg bg-secondary px-2.5 text-main placeholder:text-secondary-foreground"
					placeholder="Last Name"
					autocomplete="family-name"
					maxLength={50}
					value={editUser.last_name}
					onInput={e => setEditUser('last_name', e.currentTarget.value)}
				/>
				<input
					class="h-10 w-full rounded-lg bg-secondary px-2.5 text-main placeholder:text-secondary-foreground"
					placeholder="Title"
					maxLength={70}
					value={editUser.title}
					onInput={e => setEditUser('title', e.currentTarget.value)}
				/>
			</div>
		</FormLayout>
	)
}
