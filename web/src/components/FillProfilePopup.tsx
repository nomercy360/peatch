import { Link } from '~/components/link'

export default function FillProfilePopup(props: { onClose: () => void }) {
	return (
		<div class="w-full bg-secondary rounded-xl relative p-3 text-center">
			<button
				class="absolute right-4 top-4 flex size-6 items-center justify-center rounded-full bg-background"
				onClick={props.onClose}
			>
				<span class="material-symbols-rounded text-[20px] text-secondary-foreground">
					close
				</span>
			</button>
			<div class="flex items-center gap-1 justify-center text-2xl font-extrabold text-green">
				<span class="material-symbols-rounded text-[36px] text-blue-400">
					people
				</span>
				Set up your profile
			</div>
			<p class="mt-2 text-base font-normal text-secondary-foreground">
				Complete your profile in just 5 minutes to enhance your networking and be able to collaborate with others.
			</p>
			<Link
				class="bg-primary mt-4 flex h-10 w-full items-center justify-center rounded-xl text-sm font-semibold"
				href="/users/edit"
			>
				Set up profile
			</Link>
		</div>
	)
}
