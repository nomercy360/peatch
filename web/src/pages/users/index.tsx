import {
  createEffect,
  createSignal,
  For,
  onCleanup,
  Show,
  Suspense,
} from 'solid-js';
import { User } from '../../../gen';
import { fetchUsers } from '~/api';
import { createQuery } from '@tanstack/solid-query';
import useDebounce from '~/hooks/useDebounce';
import Badge from '~/components/Badge';
import { Link } from '~/components/Link';
import { useLocation } from '@solidjs/router';

export default function Index() {
  const [search, setSearch] = createSignal('');

  const updateSearch = useDebounce(setSearch, 300);

  const query = createQuery(() => ({
    queryKey: ['profiles', search()],
    queryFn: () => fetchUsers(search()),
  }));

  return (
    <div class="pb-52 bg-peatch-secondary min-h-screen">
      <div class="px-4 py-2.5">
        <input
          class="h-10 w-full rounded-lg bg-peatch-main px-2.5 text-main placeholder:text-hint"
          placeholder="Search for profiles"
          type="text"
          value={search()}
          onInput={e => updateSearch(e.currentTarget.value)}
        />
      </div>
      <Suspense fallback={<UserListPlaceholder />}>
        <For each={query.data}>{profile => <UserCard user={profile} />}
        </For>
      </Suspense>
    </div>
  );
}

const UserCard = (props: { user: User }) => {
  const shortenDescription = (description: string) => {
    if (description.length <= 120) return description;
    return description.slice(0, 120) + '...';
  };

  const imgUrl = `https://assets.peatch.io/${props.user.avatar_url}`;

  const badgeSlice = props.user.badges?.slice(0, 5);

  return (
    <Link
      class="flex flex-col items-start px-4 pb-5 pt-4 text-start bg-peatch-secondary"
      href={`/users/${props.user.id}`}
      state={{ from: '/users' }}
    >
      <img
        class="size-10 rounded-2xl object-cover"
        src={imgUrl}
        alt="User Avatar"
      />
      <p class="mt-3 text-3xl text-blue">{props.user.first_name}:</p>
      <p class="text-3xl text-main">{props.user.title}</p>
      <p class="mt-2 text-sm text-hint">
        {shortenDescription(props.user.description!)}
      </p>
      <div class="mt-3 flex w-full flex-row flex-wrap items-center justify-start gap-1">
        <For each={badgeSlice}>
          {badge => (
            <Badge icon={badge.icon!} name={badge.text!} color={badge.color!} />
          )}
        </For>
        <Show when={props.user.badges && props.user.badges.length > 5}>
          <div class="flex h-5 flex-row items-center justify-center rounded bg-black px-2.5 text-xs text-white">
            + {props.user.badges!.length - 5} more
          </div>
        </Show>
      </div>
      <div class="h-px bg-peatch-main w-full mt-5"></div>
    </Link>
  );
};

const UserListPlaceholder = () => {
  return (
    <div class="flex flex-col items-start justify-start gap-4 px-4 py-2.5">
      <div class="h-56 w-full rounded-2xl bg-peatch-main"></div>
      <div class="h-56 w-full rounded-2xl bg-peatch-main"></div>
      <div class="h-56 w-full rounded-2xl bg-peatch-main"></div>
      <div class="h-56 w-full rounded-2xl bg-peatch-main"></div>
    </div>
  );
};
