import { FormLayout } from '~/components/edit/layout'
import { useMainButton } from '~/lib/useMainButton'
import { useNavigate } from '@solidjs/router'
import { onCleanup, onMount } from 'solid-js'
import { editPost, setEditPost } from '~/store'
import SelectLocation from '~/components/edit/selectLocation'

export default function AddLocation() {
	const mainButton = useMainButton()

	const navigate = useNavigate()

	const goBack = () => {
		navigate(-1)
	}

	onMount(() => {
		mainButton.enable('Back').onClick(goBack)
	})

	onCleanup(() => {
		mainButton.offClick(goBack)
	})

	return (
		<FormLayout
			title="Any special location?"
			description="People will see it when they discover your post."
			screen={4}
			totalScreens={6}
		>
			<SelectLocation
				city={editPost.city}
				setCity={c => setEditPost('city', c)}
				country={editPost.country}
				setCountry={c => setEditPost('country', c)}
				countryCode={editPost.country_code}
				setCountryCode={c => setEditPost('country_code', c)}
			/>
		</FormLayout>
	)
}
