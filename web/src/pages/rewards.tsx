import {
	Component,
	ComponentProps,
	createSignal,
	Match,
	onCleanup,
	onMount,
	splitProps,
	Switch,
} from 'solid-js'
import { claimDailyReward } from '~/lib/api'
import { createMutation } from '@tanstack/solid-query'
import { setUser, store } from '~/store'
import { addToast } from '~/components/toast'
import { useNavigate } from '@solidjs/router'
import { PeatchIcon } from '~/components/peatch-icon'

export default function Rewards() {
	const dailyRewardMutate = createMutation(() => ({
		mutationFn: () => claimDailyReward(),
		retry: 0,
		onSuccess: () => {
			window.Telegram.WebApp.HapticFeedback.impactOccurred('light')
			setUser({ ...store.user, peatch_points: store.user.peatch_points! + 10 })
		},
		onError: error => {
			addToast('You have already claimed your daily reward.')
		},
	}))

	const [timeLeft, setTimeLeft] = createSignal<string | null>(null)

	const calculateTimeLeft = () => {
		if (!store.user.last_check_in) {
			setTimeLeft(null)
			return
		}

		const lastCheckIn = new Date(store.user.last_check_in)
		const nextDay = new Date(
			lastCheckIn.getFullYear(),
			lastCheckIn.getMonth(),
			lastCheckIn.getDate() + 1,
			0,
			0,
			0,
		)
		const now = new Date()
		const diff = nextDay.getTime() - now.getTime()

		if (diff <= 0) {
			setTimeLeft(null)
		} else {
			const hours = Math.floor(diff / (1000 * 60 * 60))
			const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
			setTimeLeft(`${hours}h ${minutes}m`)
		}
	}

	let interval: any

	onMount(() => {
		calculateTimeLeft()
		interval = setInterval(calculateTimeLeft, 1000)
	})

	onCleanup(() => clearInterval(interval))

	const navigate = useNavigate()

	const [surveyCompleted, setSurveyCompleted] = createSignal(false)

	function shareURL() {
		const url =
			'https://t.me/share/url?' +
			new URLSearchParams({
				url: 'https://t.me/peatch_bot?start=friend' + store.user.id,
			}).toString() +
			`&text=${store.user.first_name} invited you to join Peatch! 🍑🎉`

		window.Telegram.WebApp.openTelegramLink(url)
	}

	onMount(() => {
		window.Telegram.WebApp.CloudStorage.getItem(
			'surveyCompleted',
			(err, value) => {
				if (err) {
					console.error(err)
					return
				}

				if (value) {
					setSurveyCompleted(true)
				}
			},
		)
		// window.Telegram.WebApp.CloudStorage.removeItem('surveyCompleted')
	})

	return (
		<div class="flex min-h-screen w-full flex-col items-center justify-start bg-secondary p-3.5 text-center">
			<PeatchIcon width={48} height={48} />
			<h1 class="mt-2 max-w-[285px] text-3xl text-main">Peatch Rewards</h1>
			<p class="mb-6 mt-2 max-w-[285px] text-xl font-normal text-main">
				Earn Peatches by creating content, collaborating on projects, and
				completing various tasks
			</p>
			<RewardCard
				title="Complete Daily Check-in"
				description="Sign in to the app every day to earn Peatches"
				reward="10 peatches"
			>
				<RewardButton
					disabled={timeLeft() !== null || dailyRewardMutate.isSuccess}
					onClick={() => dailyRewardMutate.mutate()}
				>
					<Switch>
						<Match when={dailyRewardMutate.isSuccess}>Claimed</Match>
						<Match when={timeLeft()}>{timeLeft()}</Match>
						<Match when={!timeLeft()}>Claim</Match>
					</Switch>
				</RewardButton>
			</RewardCard>
			<RewardCard
				title="Refer a Friend"
				description="Invite your friends to join and earn Peatches"
				reward="100 peatches"
			>
				<RewardButton onClick={() => shareURL()}>Send Invite</RewardButton>
			</RewardCard>
			<RewardCard
				title="Complete a Survey"
				description="Share your feedback and earn Peatches"
				reward="50 peatches"
			>
				<RewardButton
					disabled={surveyCompleted()}
					onClick={() => navigate('/survey')}
				>
					{surveyCompleted() ? 'Completed' : 'Complete'}
				</RewardButton>
			</RewardCard>
			<RewardCard
				title="Create a Post"
				description="Create and publish a new collaboration project to invite others to join."
				reward="30 peatches"
			>
				<RewardButton
					onClick={() =>
						navigate('/collaborations/edit', { state: { from: '/rewards' } })
					}
				>
					Create
				</RewardButton>
			</RewardCard>
			<RewardCard
				title="Collaborate on a Project"
				description="Join a collaboration project and earn Peatches"
				reward="30 peatches"
			>
				<RewardButton disabled>Join</RewardButton>
			</RewardCard>
		</div>
	)
}

const RewardCard: Component<{
	title: string
	description: string
	reward: string
	children: any
}> = props => {
	return (
		<div class="mt-2 flex w-full flex-col items-center text-start">
			<div class="flex w-full flex-col items-start justify-start rounded-xl bg-main p-4">
				<p class="text-base font-medium text-main">{props.title}</p>
				<p class="text-sm text-secondary">{props.description}</p>
				<div class="mt-4 flex w-full flex-row items-center justify-between">
					<p class="font-bold text-accent">{props.reward}</p>
					{props.children}
				</div>
			</div>
		</div>
	)
}

const RewardButton: Component<ComponentProps<'button'>> = props => {
	const [, rest] = splitProps(props, ['children'])
	return (
		<button
			class="flex h-10 items-center justify-center rounded-xl bg-button px-2.5 text-sm font-medium text-button disabled:bg-secondary disabled:text-main"
			{...rest}
		>
			{props.children}
		</button>
	)
}
