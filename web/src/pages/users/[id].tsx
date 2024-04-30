import { useButtons } from '../../hooks/useBackButton';
import { createEffect, For, Match, onCleanup, Show, Suspense, Switch } from 'solid-js';
import { useNavigate, useParams } from '@solidjs/router';
import { CDN_URL, fetchProfile, followUser, hideProfile, publishProfile, showProfile, unfollowUser } from '../../api';
import { createQuery } from '@tanstack/solid-query';
import { setFollowing, setUser, store } from '../../store';

export default function UserProfile() {
  const { mainButton, backButton } = useButtons();

  const navigate = useNavigate();
  const params = useParams();

  const userId = params.id;

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

  const query = createQuery(() => ({
    queryKey: ['profiles', userId],
    queryFn: () => fetchProfile(Number(userId)),
  }));

  const isCurrentUserProfile = store.user.id === Number(userId);

  const publish = async () => {
    setUser({
      ...store.user,
      published_at: new Date().toISOString(),
    });
    await publishProfile();
  };

  const hide = async () => {
    setUser({
      ...store.user,
      hidden_at: new Date().toISOString(),
    });
    await hideProfile();
  };

  const show = async () => {
    setUser({
      ...store.user,
      hidden_at: undefined,
    });
    await showProfile();
  };

  const follow = async () => {
    setFollowing([...store.following, Number(userId)]);
    await followUser(Number(userId));
  };

  const unfollow = async () => {
    setFollowing(store.following.filter(id => id !== Number(userId)));
    await unfollowUser(Number(userId));
  };

  const collaborate = async () => {
    navigate(`/users/${userId}/collaborate`);
  };

  const pushToEdit = () => {
    navigate(`/users/edit`);
  };

  createEffect(() => {
    if (isCurrentUserProfile) {
      if (!store.user.published_at) {
        mainButton.offClick(pushToEdit);
        mainButton.setVisible('Publish');
        mainButton.onClick(publish);
      } else {
        mainButton.offClick(publish);
        mainButton.setVisible('Edit');
        mainButton.onClick(pushToEdit);
      }
    } else {
      mainButton.offClick(pushToEdit);
      mainButton.setVisible('Collaborate');
      mainButton.onClick(collaborate);
    }
  });

  onCleanup(() => {
    mainButton.hide();
    mainButton.offClick(collaborate);
    mainButton.offClick(publish);
    mainButton.offClick(pushToEdit);
  });

  createEffect(() => {
    console.error('User Published at: ', store.user.published_at);
  });

  return (
    <div>
      <Suspense fallback={<div>Loading...</div>}>
        <Show when={query.data}>
          <div class="min-h-screen">
            <Switch>
              <Match when={isCurrentUserProfile && store.user.hidden_at}>
                <ActionButton text="Show" onClick={show} />
              </Match>
              <Match when={isCurrentUserProfile && !store.user.hidden_at}>
                <ActionButton text="Hide" onClick={hide} />
              </Match>
              <Match
                when={
                  !isCurrentUserProfile &&
                  !store.following.includes(Number(userId))
                }
              >
                <ActionButton text="Follow" onClick={follow} />
              </Match>
              <Match
                when={
                  !isCurrentUserProfile &&
                  store.following.includes(Number(userId))
                }
              >
                <ActionButton text="Unfollow" onClick={unfollow} />
              </Match>
            </Switch>
            <div class="image-container">
              <img
                src={CDN_URL + '/' + query.data.avatar_url}
                alt="avatar"
                class="aspect-square size-full object-cover"
              />
            </div>
            <div class="px-4 py-2.5">
              <p class="text-3xl text-pink">
                {query.data.first_name} {query.data.last_name}:
              </p>
              <p class="text-3xl text-black">{query.data.title}</p>
              <p class="text-lg font-normal"> {query.data.description}</p>
              <div class="mt-5 flex flex-row flex-wrap items-center justify-start gap-1">
                <For each={query.data.badges}>
                  {badge => (
                    <div
                      class="flex h-10 flex-row items-center justify-center gap-[5px] rounded-2xl border border-peatch-stroke px-2.5"
                      style={{
                        'background-color': `#${badge.color}`,
                        'border-color': `#${badge.color}`,
                      }}
                    >
                      <span class="material-symbols-rounded text-white">
                        {String.fromCodePoint(parseInt(badge.icon!, 16))}
                      </span>
                      <p class="text-sm font-semibold text-white">
                        {badge.text}
                      </p>
                    </div>
                  )}
                </For>
              </div>
              <div class="mt-5 flex w-full flex-col items-center justify-start gap-1">
                <For each={query.data.opportunities}>
                  {op => (
                    <div
                      class="flex h-[60px] w-full flex-row items-center justify-start gap-2.5 rounded-2xl border border-peatch-stroke px-2.5"
                      style={{
                        'background-color': `#${op.color}`,
                      }}
                    >
                      <div class="flex size-10 items-center justify-center rounded-full bg-white">
                        <span class="material-symbols-rounded text-black">
                          {String.fromCodePoint(parseInt(op.icon!, 16))}
                        </span>
                      </div>
                      <div class="text-start text-white">
                        <p class="text-sm font-semibold">{op.text}</p>
                        <p class="text-xs text-white/60">{op.description}</p>
                      </div>
                    </div>
                  )}
                </For>
              </div>
            </div>
          </div>
        </Show>
      </Suspense>
    </div>
  );
}
// background: ;

const ActionButton = (props: { text: string; onClick: () => void }) => {
  return (
    <button
      class="absolute left-4 top-4 z-10 h-8 w-20 rounded-lg px-2.5 text-white"
      onClick={props.onClick}
      style={{ background: 'rgba(255, 255, 255, 0.20)' }}
    >
      {props.text}
    </button>
  );
};
