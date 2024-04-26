import { useButtons } from '../hooks/useBackButton';
import { createEffect, createSignal, For, onCleanup, Show } from 'solid-js';
import { User } from '../../gen';
import Badge from '../components/Badge';
import { createQuery } from '@tanstack/solid-query';

const mockProfiles: User[] = [
  {
    avatar_url: '/thumb.png',
    badges: [
      {
        icon: 'self_improvement',
        name: 'Self Improvement',
        color: '#FFA500',
      },
      {
        icon: 'wine_bar',
        name: 'Sip & Chat',
        color: '#d82828',
      },
      {
        icon: 'directions_run',
        name: 'Workout Buddy',
        color: '#e7a3f6',
      },
      {
        icon: 'self_improvement',
        name: 'Reading Club',
        color: '#f33333',
      },
      {
        icon: 'wine_bar',
        name: 'Wine Bar',
        color: '#4e8eb6',
      },
      {
        icon: 'directions_run',
        name: 'Directions Run',
        color: '#bd5b8f',
      },
    ],
    chat_id: 1,
    title: 'Software Engineer',
    city: 'San Francisco',
    country: 'United States',
    country_code: 'US',
    created_at: '2021-01-01',
    description:
      '27 y.o. serial entrepreneur & product director with architecture, product design, marketing & tech development background in the US, UK, and EU. I am a founder of a few startups and',
    first_name: 'John',
    followers_count: 100,
    following_count: 200,
    id: 1,
    is_published: true,
    language: 'en',
  },
];

const fetchProfiles = async () => {
  await new Promise(resolve => setTimeout(resolve, 300));
  return mockProfiles;
};

export default function Profiles() {
  const { setBackVisibility, onBackClick, offBackClick } = useButtons();

  const query = createQuery(() => ({
    queryKey: ['todos'],
    queryFn: fetchProfiles,
  }));

  const back = () => {
    history.back();
  };

  createEffect(() => {
    setBackVisibility(true);
    onBackClick(back);
  });

  onCleanup(() => {
    setBackVisibility(false);
    offBackClick(back);
  });

  const [search, setSearch] = createSignal('');

  return (
    <div>
      <div class="px-4 py-2.5">
        <input
          class="h-10 w-full rounded-lg bg-peatch-bg px-2.5 text-black placeholder:text-gray"
          placeholder="Search for profiles"
        />
      </div>
      <For each={query.data}>{profile => <UserProfile user={profile} />}</For>
    </div>
  );
}

const UserProfile = (props: { user: User }) => {
  const shortenDescription = (description: string) => {
    return description.slice(0, 120) + '...';
  };

  const badgeSlice = props.user.badges.slice(0, 5);

  return (
    <div class="flex flex-col items-start px-4 pb-5 pt-4 text-start">
      <img
        class="size-10 rounded-2xl"
        src={props.user.avatar_url}
        alt="User Avatar"
      />
      <p class="text-blue mt-3 text-3xl">{props.user.first_name}:</p>
      <p class="text-3xl">{props.user.title}</p>
      <p class="mt-2 text-sm text-gray">
        {shortenDescription(props.user.description)}
      </p>
      <div class="mt-3 flex w-full flex-row flex-wrap items-center justify-start gap-1">
        <For each={badgeSlice}>
          {badge => (
            <Badge icon={badge.icon} name={badge.name} color={badge.color} />
          )}
        </For>
        <Show when={props.user.badges.length > 5}>
          <div class="text-xs flex h-5 flex-row items-center justify-center rounded bg-black px-2.5 text-white">
            + {props.user.badges.length - 5} more
          </div>
        </Show>
      </div>
    </div>
  );
};
