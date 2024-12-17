import { For } from 'solid-js'
import { fetchPreview } from '~/lib/api'
import { createQuery } from '@tanstack/solid-query'
import { Link } from '~/components/link'

export default function ActionDonePopup(props: {
	action: string
	description: string
	callToAction: string
}) {
	const previewQuery = createQuery(() => ({
		queryKey: ['preview'],
		queryFn: () => fetchPreview(),
	}))

	return (
		<div class="flex h-screen w-full flex-col items-center justify-between bg-secondary p-5 text-center">
			<div class="flex flex-col items-center justify-start">
				<span class="material-symbols-rounded text-peatch-green text-[60px] text-green">
					schedule
				</span>
				<p class="text-3xl text-main">{props.action}</p>
				<p class="mt-2 text-2xl text-secondary">{props.description}</p>
			</div>
			<div class="flex flex-col items-center justify-center">
				<div class="flex w-full flex-row items-center justify-center">
					<For each={previewQuery.data!}>
						{(image, idx) => (
							<img
								src={image}
								alt="User Avatar"
								class="-ml-1 size-11 rounded-lg border object-cover object-center"
								classList={{
									'ml-0': idx() === 0,
									'z-20': idx() === 0,
									'z-10': idx() === 1,
								}}
							/>
						)}
					</For>
				</div>
				<p class="mt-4 max-w-xs text-lg text-secondary">{props.callToAction}</p>
				<Link
					class="mt-2 flex h-12 w-full items-center justify-center text-sm font-medium text-link"
					href="/"
				>
					Show them
				</Link>
			</div>
		</div>
	)
}
