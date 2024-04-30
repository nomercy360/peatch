import { useButtons } from '../../hooks/useBackButton';
import {
  createEffect,
  createSignal,
  For,
  onCleanup,
  Show,
  Suspense,
} from 'solid-js';
import { Collaboration } from '../../../gen';
import { useNavigate } from '@solidjs/router';
import { CDN_URL, fetchCollaborations } from '../../api';
import { createQuery } from '@tanstack/solid-query';
import useDebounce from '../../hooks/useDebounce';
import Badge from '../../components/Badge';

export default function Index() {
  const [search, setSearch] = createSignal('');
  const { backButton } = useButtons();

  const navigate = useNavigate();

  const back = () => {
    navigate('/');
  };

  createEffect(() => {
    backButton.setVisible();
    backButton.onClick(back);
  });

  onCleanup(() => {
    backButton.hide();
    backButton.offClick(back);
  });

  const updateSearch = useDebounce(setSearch, 300);

  const query = createQuery(() => ({
    queryKey: ['collaborations', search()],
    queryFn: () => fetchCollaborations(search()),
  }));

  return (
    <div>
      <div class="px-4 py-2.5">
        <input
          class="h-10 w-full rounded-lg bg-peatch-bg px-2.5 text-black placeholder:text-gray"
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
    <a
      class="flex flex-col items-start px-4 pt-4 text-start"
      href={`/collaborations/${props.collab.id}`}
    >
      <p class="mt-3 text-3xl text-blue">{props.collab.opportunity?.text}:</p>
      <p class="text-3xl">{props.collab.title}</p>
      <p class="mt-2 text-sm text-gray">
        {shortenDescription(props.collab.description!)}
      </p>
      <div class="mt-4 gap-2 flex w-full flex-row items-center justify-start">
        <img
          class="size-10 rounded-2xl object-cover"
          src={CDN_URL + '/' + props.collab.user?.avatar_url}
          alt="User Avatar"
        />
        <div>
          <p class="text-sm font-bold text-black">
            {props.collab.user?.first_name} {props.collab.user?.last_name}:
          </p>
          <p class="text-sm text-black">{props.collab.user?.title}</p>
        </div>
      </div>
      <div class="h-px bg-peatch-stroke w-full mt-5"></div>
    </a>
  );
};

const CollabListPlaceholder = () => {
  return (
    <div class="flex flex-col items-start justify-start gap-4 px-4 py-2.5">
      <div class="h-48 w-full rounded-2xl bg-peatch-stroke"></div>
      <div class="h-48 w-full rounded-2xl bg-peatch-stroke"></div>
      <div class="h-48 w-full rounded-2xl bg-peatch-stroke"></div>
      <div class="h-48 w-full rounded-2xl bg-peatch-stroke"></div>
    </div>
  );
};
