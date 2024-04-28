import { createSignal, For, Show, Suspense } from 'solid-js';
import { FormLayout } from './layout';
import countryFlags from '../../assets/countries.json';
import { createQuery } from '@tanstack/solid-query';
import useDebounce from '../../hooks/useDebounce';

type Location = {
  address: {
    country: string;
    city: string;
    town: string;
    country_code: string;
  };
};

type CountryFlag = {
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
    queryFn: async () => {
      const resp = await fetch(
        `https://nominatim.openstreetmap.org/search?q=${search()}&format=jsonv2&addressdetails=1`,
        {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
        },
      );

      return resp.json();
    },
    enabled: !!search(),
  }));

  const updateSearch = useDebounce(setSearch, 300);

  const onLocationClick = (location: Location) => {
    props.setCountry(location.address.country);
    props.setCity(location.address.city || location.address.town);
    props.setCountryCode(location.address.country_code);
  };

  const clearLocation = () => {
    props.setCountry('');
    props.setCity('');
    props.setCountryCode('');
  };

  return (
      <>
      <div class="mt-5 flex h-10 w-full flex-row items-center justify-between rounded-lg bg-peatch-bg px-2.5">
        <input
          class="w-full bg-transparent text-black placeholder:text-gray focus:outline-none"
          placeholder="City and country"
          type="text"
          onInput={e => updateSearch(e.currentTarget.value)}
          value={search()}
        />
        <Show when={search()}>
          <button
            class="flex h-full items-center justify-center px-2.5 text-sm text-peatch-dark-gray"
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
              address: {
                country: props.country,
                city: props.city || '',
                town: '',
                country_code: props.countryCode,
              },
            }}
          />
        </Show>
        <Suspense fallback={<div>Loading...</div>}>
          <For each={query.data!}>
            {location => (
              <LocationButton
                isActive={
                  location.address.country === props.country &&
                  location.address.city === props.city &&
                  location.address.country_code === props.countryCode
                }
                onClick={() => onLocationClick(location)}
                location={location}
              />
            )}
          </For>
        </Suspense>
      </div>
    </>
  );
}

function LocationButton(props: {
  location: Location;
  onClick: () => void;
  isActive: boolean;
}) {
  const findFlag = (code: string) => {
    const flag = countryFlags.find(
      (flag: CountryFlag) => flag.code?.toLowerCase() === code?.toLowerCase(),
    );
    return flag?.flag;
  };

  return (
    <button
      onClick={() => props.onClick()}
      class="flex h-10 w-full flex-row items-center justify-between"
      classList={{
        'bg-peatch-bg': props.isActive,
      }}
    >
      <p class="">
        {props.location.address.country}
        <Show when={props.location.address.city || props.location.address.town}>
          , {props.location.address.city || props.location.address.town}
        </Show>
      </p>
      <div class="flex size-10 items-center justify-center">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="26"
          viewBox={
            findFlag(props.location.address.country_code)
              ? '0 0 512 512'
              : '0 0 24 24'
          }
          class="size-6"
          innerHTML={findFlag(props.location.address.country_code)}
        ></svg>
      </div>
    </button>
  );
}
