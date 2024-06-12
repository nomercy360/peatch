import { Show } from 'solid-js'
import countryFlags from '~/assets/countries.json'
import { CountryFlag } from '~/components/edit/selectLocation'

export const LocationBadge = ({
	country,
	city,
	countryCode,
}: {
	country: string
	city: string
	countryCode: string
}) => {
	const findFlag = countryFlags.find(
		(flag: CountryFlag) => flag.code === countryCode,
	)

	return (
		<div class="flex h-5 flex-row items-center justify-center gap-[5px] rounded bg-button px-2.5 text-xs font-semibold text-button">
			<Show when={findFlag}>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					viewBox={findFlag!.viewBox}
					class="z-10 size-3.5"
					innerHTML={findFlag!.flag}
				/>
			</Show>
			<p class="text-secondary-bg text-xs font-semibold">
				{city ? `${city}, ${country}` : country}
			</p>
		</div>
	)
}
