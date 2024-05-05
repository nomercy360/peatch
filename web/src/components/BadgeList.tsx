import { For, Show } from 'solid-js';
import Badge from '~/components/Badge';
import { Badge as UserBadge } from '../../gen';
import countryFlags from '~/assets/countries.json';
import { CountryFlag } from '~/components/edit/selectLocation';

export default function BadgeList(props: {
  badges: UserBadge[];
  position: 'center' | 'start';
  country: string;
  city: string;
  countryCode: string;
}) {
  const badgeSlice = props.badges!.slice(0, 5);

  const findFlag = countryFlags.find(
    (flag: CountryFlag) => flag.code === props.countryCode,
  );

  return (
    <div
      class="mt-3 flex w-full flex-row flex-wrap items-center justify-start gap-1"
      classList={{
        'justify-center': props.position === 'center',
      }}
    >
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
          {props.city ? `${props.city}, ${props.country}` : props.country}
        </p>
      </div>
      <For each={badgeSlice}>
        {badge => (
          <Badge icon={badge.icon!} name={badge.text!} color={badge.color!} />
        )}
      </For>
      <Show when={props.badges.length > 5}>
        <div class="flex h-5 flex-row items-center justify-center rounded bg-black px-2.5 text-xs font-semibold text-white">
          + {props.badges!.length - 5} more
        </div>
      </Show>
    </div>
  );
}
