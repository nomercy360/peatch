import { createSignal, For, Show, Suspense } from 'solid-js'
import countryFlags from '../../assets/countries.json'
import useDebounce from '../../lib/useDebounce'
import { searchLocations } from '~/lib/api'
import { useQuery } from '@tanstack/solid-query'
import { CityResponse } from '~/gen'

export type CountryFlag = {
	flag: string
	code: string
	viewBox: string
}

export default function SelectLocation(props: {
	setLocation: (city: CityResponse) => void
	initialLocation?: CityResponse
}) {
	const [search, setSearch] = createSignal('')

	const query = useQuery(() => ({
		queryKey: ['locations', search()],
		queryFn: () => searchLocations(search()),
	}))

	const updateSearch = useDebounce(setSearch, 400)

	const onLocationClick = (location: CityResponse) => {
		props.setLocation(location)
	}

	const clearLocation = () => {
		props.setLocation({} as CityResponse)
		setSearch('')
	}

	return (
		<>
			<div class="mt-5 flex h-10 w-full flex-row items-center justify-between rounded-lg bg-secondary px-2.5">
				<input
					class="text-main placeholder:text-hint w-full bg-transparent focus:outline-none"
					placeholder="City and country"
					type="text"
					onInput={e => updateSearch(e.currentTarget.value)}
					value={search()}
				/>
				<Show when={search()}>
					<button
						class="text-hint flex h-full items-center justify-center px-2.5 text-sm"
						onClick={() => setSearch('')}
					>
						Clear
					</button>
				</Show>
			</div>
			<div class="mt-2.5 flex w-full flex-row flex-wrap items-center justify-start gap-1">
				<Show when={!search() && props.initialLocation}>
					<LocationButton
						onClick={() => clearLocation()}
						isActive={true}
						location={props.initialLocation!}
					/>
				</Show>
				<Show when={search() || !props.initialLocation}>
					<Suspense fallback={<div class="text-hint text-sm">Loading...</div>}>
						<For each={query.data!}>
							{location => (
								<LocationButton
									isActive={props.initialLocation?.id === location.id}
									onClick={() => onLocationClick(location)}
									location={location}
								/>
							)}
						</For>
					</Suspense>
				</Show>
			</div>
		</>
	)
}

function LocationButton(props: {
	location: CityResponse
	onClick: () => void
	isActive: boolean
}) {
	const findFlag = countryFlags.find(
		(flag: CountryFlag) => flag.code === props.location.country_code,
	)

	return (
		<button
			onClick={() => props.onClick()}
			class="flex h-12 w-full flex-row items-center justify-start space-x-3 rounded-xl px-4 text-sm"
			classList={{
				'bg-secondary': props.isActive,
				'bg-border': props.isActive,
			}}
		>
			<Show when={findFlag}>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					viewBox={findFlag!.viewBox}
					class="z-10 size-7 flex-shrink-0"
					innerHTML={findFlag!.flag}
				/>
			</Show>
			<div class="flex flex-col justify-start text-start">
				<p class="text-sm">{props.location.name}</p>
				<p class="text-xs text-secondary-foreground">
					{props.location.country_name}
				</p>
			</div>
		</button>
	)
}
