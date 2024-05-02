import { createEffect, createSignal, For, Match, onCleanup, Suspense, Switch } from 'solid-js';
import { useNavigate, useParams, useSearchParams } from '@solidjs/router';
import { CDN_URL, fetchProfile, followUser, hideProfile, publishProfile, showProfile, unfollowUser } from '~/api';
import { createQuery } from '@tanstack/solid-query';
import { setFollowing, setUser, store } from '~/store';
import ActionDonePopup from '../../components/ActionDonePopup';
import { useMainButton } from '~/hooks/useMainButton';
import { useNavigation } from '~/hooks/useNavigation';

export default function UserProfile() {
  const mainButton = useMainButton();
  const [published, setPublished] = createSignal(false);

  const navigate = useNavigate();

  const { navigateBack } = useNavigation();

  const params = useParams();
  const [searchParams, _] = useSearchParams();

  const userId = params.id;

  const isCurrentUserProfile = store.user.id === Number(userId);

  const query = createQuery(() => ({
    queryKey: ['profiles', userId],
    queryFn: () => fetchProfile(Number(userId)),
  }));

  createEffect(async () => {
    if (searchParams.refetch) {
      await query.refetch();
      if (query.data.id === store.user.id) {
        setUser(query.data);
      }
    }
  });

  const navigateToCollaborate = async () => {
    if (store.user.published_at && !store.user.hidden_at) {
      navigate(`/users/${userId}/collaborate`);
    }
  };

  const navigateToEdit = () => {
    navigate('/users/edit', { state: { back: true } });
  };

  const publish = async () => {
    setUser({
      ...store.user,
      published_at: new Date().toISOString(),
    });
    await publishProfile();
    setPublished(true);
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

  createEffect(() => {
    if (isCurrentUserProfile) {
      if (published()) {
        mainButton.offClick(publish);
        mainButton.offClick(navigateToEdit);
        mainButton.setParams({
          text: 'Just open the app',
          isVisible: true,
          isEnabled: true,
        });
        mainButton.onClick(navigateBack);
        return;
      } else if (!store.user.published_at) {
        mainButton.setParams({
          text: 'Publish',
          isVisible: true,
          isEnabled: true,
        });
        mainButton.offClick(navigateToEdit);
        mainButton.onClick(publish);
      } else {
        mainButton.setParams({
          text: 'Edit',
          isVisible: true,
          isEnabled: true,
        });
        mainButton.offClick(publish);
        mainButton.onClick(navigateToEdit);
      }
    } else {
      mainButton.setParams({
        text: 'Collaborate',
        isVisible: true,
        isEnabled: true,
      });
      mainButton.offClick(navigateToEdit);
      mainButton.onClick(navigateToCollaborate);
    }

    onCleanup(() => {
      mainButton.offClick(navigateToCollaborate);
      mainButton.offClick(publish);
      mainButton.offClick(navigateBack);
      mainButton.offClick(navigateToEdit);
    });
  });

  onCleanup(async () => {
    mainButton.hide();
  });

  return (
    <div>
      <Suspense fallback={<div>Loading...</div>}>
        <Switch>
          <Match when={published() && isCurrentUserProfile}>
            <ActionDonePopup
              action="Profile published"
              description="Now you can find people, create and join collaborations. Have fun!"
              callToAction="There are 12 people you might be interested to collaborate with"
            />
          </Match>
          <Match when={query.data}>
            <div class="h-fit min-h-screen bg-secondary">
              <Switch>
                <Match when={isCurrentUserProfile && !store.user.published_at}>
                  <ActionButton text="Edit" onClick={navigateToEdit} />
                </Match>
                <Match
                  when={
                    isCurrentUserProfile &&
                    store.user.hidden_at &&
                    store.user.published_at
                  }
                >
                  <ActionButton text="Show" onClick={show} />
                </Match>
                <Match
                  when={
                    isCurrentUserProfile &&
                    !store.user.hidden_at &&
                    store.user.published_at
                  }
                >
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
                <p class="text-3xl text-main">{query.data.title}</p>
                <p class="mt-1 text-lg font-normal text-secondary">
                  {' '}
                  {query.data.description}
                </p>
                <div class="mt-5 flex flex-row flex-wrap items-center justify-start gap-1">
                  <For each={query.data.badges}>
                    {badge => (
                      <div
                        class="flex h-10 flex-row items-center justify-center gap-[5px] rounded-2xl border border-main px-2.5"
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
                        class="flex h-[60px] w-full flex-row items-center justify-start gap-2.5 rounded-2xl border border-main px-2.5"
                        style={{
                          'background-color': `#${op.color}`,
                        }}
                      >
                        <div class="flex size-10 items-center justify-center rounded-full bg-secondary">
                          <span class="material-symbols-rounded text-main">
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
          </Match>
        </Switch>
      </Suspense>
    </div>
  );
}
// background: ;

const ActionButton = (props: { text: string; onClick: () => void }) => {
  return (
    <button
      class="absolute left-4 top-4 z-10 h-8 w-20 rounded-lg bg-button px-2.5 text-button"
      onClick={props.onClick}
    >
      {props.text}
    </button>
  );
};
