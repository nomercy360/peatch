import { useButtons } from '../../hooks/useBackButton';
import { createEffect, createSignal, For, onCleanup, Show, Suspense } from 'solid-js';
import { User } from '../../../gen';
import Badge from '../../components/badge';
import { useNavigate } from '@solidjs/router';
import { fetchUsers } from '../../api';
import { createQuery } from '@tanstack/solid-query';
import useDebounce from '../../hooks/useDebounce';

export default function Index() {
  const { backButton } = useButtons();
  const [search, setSearch] = createSignal('');

  const navigate = useNavigate();

  const back = () => {
    navigate('/');
  };

  createEffect(() => {
    backButton.setVisible(true);
    backButton.onClick(back);
  });

  onCleanup(() => {
    backButton.setVisible(false);
    backButton.offClick(back);
  });

  const updateSearch = useDebounce(setSearch, 300);

  const query = createQuery(() => ({
    queryKey: ['profiles', search()],
    queryFn: () => fetchUsers(search()),
  }));

  return (
    <div>
      <div class="px-4 py-2.5">
        <input
          class="h-10 w-full rounded-lg bg-peatch-bg px-2.5 text-black placeholder:text-gray"
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

  const badgeSlice = props.user.badges?.slice(0, 5);

  return (
    <a
      class="flex flex-col items-start px-4 pb-5 pt-4 text-start"
      href={`/users/${props.user.id}`}
    >
      <img
        class="size-10 rounded-2xl object-cover"
        src={imgUrl}
        alt="User Avatar"
      />
      <p class="mt-3 text-3xl text-blue">{props.user.first_name}:</p>
      <p class="text-3xl">{props.user.title}</p>
      <p class="mt-2 text-sm text-gray">
        {shortenDescription(props.user.description!)}
      </p>
      <div class="mt-3 flex w-full flex-row flex-wrap items-center justify-start gap-1">
        <For each={badgeSlice}>
          {badge => (
            <Badge icon={badge.icon!} name={badge.text!} color={badge.color!} />
          )}
        </For>
        <Show when={props.user.badges?.length || 0 > 5}>
          <div class="flex h-5 flex-row items-center justify-center rounded bg-black px-2.5 text-xs text-white">
            + {props.user.badges?.length || 0 - 5} more
          </div>
        </Show>
      </div>
    </a>
  );
};

const UserListPlaceholder = () => {
  return (
    <div class="flex flex-col items-start justify-start gap-4 px-4 py-2.5">
      <div class="h-48 w-full rounded-2xl bg-peatch-stroke"></div>
      <div class="h-48 w-full rounded-2xl bg-peatch-stroke"></div>
      <div class="h-48 w-full rounded-2xl bg-peatch-stroke"></div>
      <div class="h-48 w-full rounded-2xl bg-peatch-stroke"></div>
    </div>
  );
};
