import { createEffect, createSignal, onCleanup } from 'solid-js'
import TextArea from '~/components/TextArea'
import { useMainButton } from '~/lib/useMainButton'
import { submitFeedbackSurvey } from '~/lib/api'
import { useNavigate } from '@solidjs/router'

export default function SurveyPage() {
	const [feedback, setFeedback] = createSignal('')

	const mainButton = useMainButton()
	const navigate = useNavigate()

	const sendFeedback = async () => {
		if (!feedback()) return
		await submitFeedbackSurvey(feedback())
		window.Telegram.WebApp.CloudStorage.setItem('surveyCompleted', 'true')
		navigate('/rewards')
	}

	createEffect(() => {
		if (feedback().length > 0) {
			mainButton.enable('Submit').onClick(sendFeedback)
		} else {
			mainButton.disable('Submit').offClick(sendFeedback)
		}
	})

	onCleanup(() => {
		mainButton.hide()
		mainButton.offClick(sendFeedback)
	})

	return (
		<div class="min-h-screen bg-secondary p-4">
			<p class="mt-2 text-center text-3xl text-main">Share Your Feedback</p>
			<p class="mt-1 text-center text-sm text-hint">
				Help us improve the collaboration and talent search experience on our
				social network.
			</p>
			<TextArea
				value={feedback()}
				setValue={setFeedback}
				placeholder="Share your thoughts and suggestion, what can we improve?"
			/>
		</div>
	)
}
