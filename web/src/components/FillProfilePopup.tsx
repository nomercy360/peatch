import { A } from '@solidjs/router'

export default function FillProfilePopup(props: { onClose: () => void }) {
	return (
		<div class="fixed bottom-0 right-0 z-30 flex w-full flex-col items-center justify-center rounded-t-3xl bg-button px-4 py-3">
			<button
				class="flex w-full flex-row items-start justify-between gap-4 text-start"
				onClick={() => props.onClose()}
			>
				<span class="text-[24px] font-bold leading-tight text-button">
					Set up your profile to collaborate with others.
				</span>
				<span class="flex size-6 items-center justify-center rounded-full bg-white/10">
					<span class="material-symbols-rounded text-[24px] text-button">
						close
					</span>
				</span>
			</button>
			<p class="mt-2 text-xl text-button">
				It only takes 5 minutes, but it can significantly improve your
				networking. According to our data, every third user finds someone to
				collaborate with within the first three days.
			</p>
			<A
				class="mt-4 flex h-12 w-full items-center justify-center rounded-2xl bg-secondary text-center text-main"
				href="/users/edit"
			>
				Set up profile
			</A>
		</div>
	)
}
