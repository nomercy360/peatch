import { Show } from 'solid-js'
import countryFlags from '~/assets/countries.json'
import { CountryFlag } from '~/components/edit/selectLocation'

type LocationBadgeProps = {
	country: string
	city: string
	countryCode: string
}
export const LocationBadge = (props: LocationBadgeProps) => {
	const findFlag = countryFlags.find(
		(flag: CountryFlag) => flag.code === props.countryCode,
	)

	return (
		<div
			class="mt-2 flex h-5 flex-row items-center justify-center gap-[5px] rounded text-xs font-semibold text-muted-foreground ">
			<Show when={findFlag}>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					viewBox={findFlag!.viewBox}
					class="z-10 size-3.5"
					innerHTML={findFlag!.flag}
				/>
			</Show>
			<p class="text-secondary-bg text-xs font-semibold">
				{props.city ? `${props.city}, ${props.country}` : props.country}
			</p>
		</div>
	)
}
