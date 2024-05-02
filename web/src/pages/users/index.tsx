import { createSignal, For, Show, Suspense } from 'solid-js';
import { User } from '../../../gen';
import { fetchUsers } from '~/api';
import { createQuery } from '@tanstack/solid-query';
import useDebounce from '~/hooks/useDebounce';
import { Link } from '~/components/Link';
import BadgeList from '~/components/BadgeList';

export default function Index() {
  const [search, setSearch] = createSignal('');

  const updateSearch = useDebounce(setSearch, 300);

  const query = createQuery(() => ({
    queryKey: ['profiles', search()],
    queryFn: () => fetchUsers(search()),
  }));

  return (
    <div class="min-h-screen bg-secondary pb-52">
      <div class="px-4 py-2.5">
        <input
          class="h-10 w-full rounded-lg bg-main px-2.5 text-main placeholder:text-hint"
          placeholder="Search for profiles"
          type="text"
          value={search()}
          onInput={e => updateSearch(e.currentTarget.value)}
        />
      </div>
      <Suspense fallback={<UserListPlaceholder />}>
        <For each={query.data}>{profile => <UserCard user={profile} />}</For>
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

  return (
    <Link
      class="flex flex-col items-start bg-secondary px-4 pb-5 pt-4 text-start"
      href={`/users/${props.user.id}`}
      state={{ from: '/users' }}
    >
      <img
        class="size-10 rounded-xl object-cover"
        src={imgUrl}
        alt="User Avatar"
      />
      <p class="mt-3 text-3xl text-blue">{props.user.first_name}:</p>
      <p class="text-3xl text-main">{props.user.title}</p>
      <p class="mt-2 text-sm text-hint">
        {shortenDescription(props.user.description!)}
      </p>
      <Show when={props.user.badges && props.user.badges.length > 0}>
        <BadgeList badges={props.user.badges!} position="start" />
      </Show>
      <div class="mt-5 h-px w-full bg-main"></div>
    </Link>
  );
};

const UserListPlaceholder = () => {
  return (
    <div class="flex flex-col items-start justify-start gap-4 px-4 py-2.5">
      <div class="h-56 w-full rounded-2xl bg-main"></div>
      <div class="h-56 w-full rounded-2xl bg-main"></div>
      <div class="h-56 w-full rounded-2xl bg-main"></div>
      <div class="h-56 w-full rounded-2xl bg-main"></div>
    </div>
  );
};
