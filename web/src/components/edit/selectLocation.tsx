import { createSignal, For, Show, Suspense } from 'solid-js';
import countryFlags from '../../assets/countries.json';
import { createQuery } from '@tanstack/solid-query';
import useDebounce from '../../hooks/useDebounce';
import { searchLocations } from '~/api';

type Location = {
  country: string;
  country_code: string;
  city: string;
};

export type CountryFlag = {
  flag: string;
  code: string;
  viewBox: string;
};

export default function SelectLocation(props: {
  country: string;
  setCountry: (country: string) => void;
  city?: string;
  setCity: (city: string) => void;
  countryCode: string;
  setCountryCode: (countryCode: string) => void;
}) {
  const [search, setSearch] = createSignal('');

  const query = createQuery(() => ({
    queryKey: ['locations', search()],
    queryFn: () => searchLocations(search()),
  }));

  const updateSearch = useDebounce(setSearch, 400);

  const onLocationClick = (location: Location) => {
    props.setCountry(location.country);
    props.setCity(location.city);
    props.setCountryCode(location.country_code);
  };

  const clearLocation = () => {
    props.setCountry('');
    props.setCity('');
    props.setCountryCode('');
  };

  return (
    <>
      <div class="mt-5 flex h-10 w-full flex-row items-center justify-between rounded-lg bg-main px-2.5">
        <input
          class="w-full bg-transparent text-main placeholder:text-hint focus:outline-none"
          placeholder="City and country"
          type="text"
          onInput={e => updateSearch(e.currentTarget.value)}
          value={search()}
        />
        <Show when={search()}>
          <button
            class="flex h-full items-center justify-center px-2.5 text-sm text-hint"
            onClick={() => setSearch('')}
          >
            Clear
          </button>
        </Show>
      </div>
      <div class="mt-2.5 flex w-full flex-row flex-wrap items-center justify-start gap-1">
        <Show when={!search() && props.country && props.countryCode}>
          <LocationButton
            onClick={() => clearLocation()}
            isActive={true}
            location={{
              country: props.country,
              city: props.city || '',
              country_code: props.countryCode,
            }}
          />
        </Show>
        <Show when={search() || (!props.country && !props.city)}>
          <Suspense fallback={<div class="text-sm text-hint">Loading...</div>}>
            <For each={query.data!}>
              {location => (
                <div class="w-full">
                  <LocationButton
                    isActive={
                      location.country === props.country &&
                      location.city === props.city &&
                      location.country_code === props.countryCode
                    }
                    onClick={() => onLocationClick(location)}
                    location={location}
                  />
                  <div class="mt-2 h-px w-full bg-main" />
                </div>
              )}
            </For>
          </Suspense>
        </Show>
      </div>
    </>
  );
}

function LocationButton(props: {
  location: Location;
  onClick: () => void;
  isActive: boolean;
}) {
  const findFlag = countryFlags.find(
    (flag: CountryFlag) => flag.code === props.location.country_code,
  );

  return (
    <button
      onClick={() => props.onClick()}
      class="flex h-16 w-full flex-row items-center justify-between rounded-2xl px-2.5 text-sm text-main"
      classList={{
        'bg-secondary': !props.isActive,
        'bg-main': props.isActive,
      }}
    >
      <p class="">
        <Show when={props.location.country && props.location.city}>
          {props.location.city}, {props.location.country}
        </Show>
        <Show when={!props.location.city}>{props.location.country}</Show>
      </p>
      <Show when={findFlag}>
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox={findFlag!.viewBox}
          class="z-10 mr-2 size-5"
          innerHTML={findFlag!.flag}
        />
      </Show>
    </button>
  );
}
