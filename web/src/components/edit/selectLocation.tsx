import { createSignal, For, Show, Suspense } from 'solid-js';
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
            locationName={
              props.city ? `${props.city}, ${props.country}` : props.country
            }
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
        <Suspense fallback={<div class="text-sm text-hint">Loading...</div>}>
          <For each={query.data!}>
            {location => (
              <LocationButton
                locationName={location.display_name}
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
  locationName: string;
}) {
  const findFlag = (code: string) => {
    const flag = countryFlags.find(
      (flag: CountryFlag) => flag.code?.toLowerCase() === code?.toLowerCase(),
    );
    return flag?.flag;
  };

  const shortenLocation = (location: string) => {
    if (location.length > 40) {
      return location.slice(0, 40) + '...';
    }
    return location;
  };

  return (
    <button
      onClick={() => props.onClick()}
      class="flex h-16 w-full flex-row items-center justify-between rounded-2xl border border-main px-2.5 text-sm text-main"
      classList={{
        'bg-secondary': !props.isActive,
        'bg-main': props.isActive,
      }}
    >
      <p class="">
        <Show
          when={props.location.address.country && props.location.address.city}
        >
          {props.location.address.city}, {props.location.address.country}
        </Show>
        <Show when={!props.location.address.city}>
          {shortenLocation(props.locationName)}
        </Show>
      </p>
      <svg
        xmlns="http://www.w3.org/2000/svg"
        viewBox={
          findFlag(props.location.address.country_code)
            ? '0 0 512 512'
            : '0 0 24 24'
        }
        class="mr-2 size-5"
        innerHTML={findFlag(props.location.address.country_code)}
      ></svg>
    </button>
  );
}
