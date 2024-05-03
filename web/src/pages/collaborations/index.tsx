import { createEffect, createSignal, For, onCleanup, Suspense } from 'solid-js';
import { Collaboration } from '../../../gen';
import { CDN_URL, fetchCollaborations } from '~/api';
import { createQuery } from '@tanstack/solid-query';
import useDebounce from '~/hooks/useDebounce';
import { Link } from '~/components/Link';
import { useMainButton } from '~/hooks/useMainButton';
import { useNavigate } from '@solidjs/router';
import { store } from '~/store';

export default function Index() {
  const [search, setSearch] = createSignal('');

  const updateSearch = useDebounce(setSearch, 300);

  const query = createQuery(() => ({
    queryKey: ['collaborations', search()],
    queryFn: () => fetchCollaborations(search()),
  }));

  const navigate = useNavigate();

  const mainButton = useMainButton();

  const pushToCreate = () => {
    navigate('/collaborations/edit');
  };

  createEffect(() => {
    if (store.user.published_at !== null) {
      mainButton.enable('Post a collaboration');
      mainButton.onClick(pushToCreate);
    }
  });

  onCleanup(() => {
    mainButton.offClick(pushToCreate);
  });

  return (
    <div class="min-h-screen bg-secondary pb-52">
      <div class="px-4 py-2.5">
        <input
          class="h-10 w-full rounded-lg bg-main px-2.5 text-main placeholder:text-hint"
          placeholder="Search collaborations by type or keyword"
          type="text"
          value={search()}
          onInput={e => updateSearch(e.currentTarget.value)}
        />
      </div>
      <Suspense fallback={<CollabListPlaceholder />}>
        <For each={query.data}>
          {collab => <CollaborationCard collab={collab} />}
        </For>
      </Suspense>
    </div>
  );
}

const CollaborationCard = (props: { collab: Collaboration }) => {
  const shortenDescription = (description: string) => {
    if (description.length <= 160) return description;
    return description.slice(0, 160) + '...';
  };

  return (
    <Link
      class="flex flex-col items-start px-4 pt-4 text-start "
      href={`/collaborations/${props.collab.id}`}
      state={{ from: '/collaborations' }}
    >
      <p class="mt-3 text-3xl text-blue">{props.collab.opportunity?.text}:</p>
      <p class="text-3xl text-main">{props.collab.title}</p>
      <p class="mt-2 text-sm text-hint">
        {shortenDescription(props.collab.description!)}
      </p>
      <div class="mt-4 flex w-full flex-row items-center justify-start gap-2">
        <img
          class="size-11 rounded-xl object-cover"
          src={CDN_URL + '/' + props.collab.user?.avatar_url}
          alt="User Avatar"
        />
        <div>
          <p class="text-sm font-bold text-main">
            {props.collab.user?.first_name} {props.collab.user?.last_name}:
          </p>
          <p class="text-sm text-main">{props.collab.user?.title}</p>
        </div>
      </div>
      <div class="mt-5 h-px w-full bg-main"></div>
    </Link>
  );
};

const CollabListPlaceholder = () => {
  return (
    <div class="flex flex-col items-start justify-start gap-4 px-4 py-2.5">
      <div class="h-48 w-full rounded-2xl bg-main"></div>
      <div class="h-48 w-full rounded-2xl bg-main"></div>
      <div class="h-48 w-full rounded-2xl bg-main"></div>
      <div class="h-48 w-full rounded-2xl bg-main"></div>
    </div>
  );
};
