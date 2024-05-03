import { createEffect, createSignal, For, Match, onCleanup, Suspense, Switch } from 'solid-js';
import { useNavigate, useParams, useSearchParams } from '@solidjs/router';
import { CDN_URL, fetchProfile, followUser, hideProfile, publishProfile, showProfile, unfollowUser } from '~/api';
import { createQuery } from '@tanstack/solid-query';
import { setFollowing, setUser, store } from '~/store';
import ActionDonePopup from '../../components/ActionDonePopup';
import { useMainButton } from '~/hooks/useMainButton';
import { useNavigation } from '~/hooks/useNavigation';
import { usePopup } from '~/hooks/usePopup';

export default function UserProfile() {
  const mainButton = useMainButton();
  const [published, setPublished] = createSignal(false);

  const navigate = useNavigate();

  const { navigateBack } = useNavigation();

  const params = useParams();
  const [searchParams, _] = useSearchParams();

  const { showConfirm } = usePopup();

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
    } else {
      showConfirm(
        'You must publish your profile first',
        (ok: boolean) =>
          ok && navigate('/users/edit', { state: { back: true } }),
      );
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
        mainButton.enable('Get back');
        mainButton.onClick(navigateBack);
        return;
      } else if (!store.user.published_at) {
        mainButton.enable('Publish');
        mainButton.offClick(navigateToEdit);
        mainButton.onClick(publish);
      } else {
        mainButton.enable('Edit');
        mainButton.offClick(publish);
        mainButton.onClick(navigateToEdit);
      }
    } else {
      mainButton.enable('Collaborate');
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
      <Suspense fallback={<Loader />}>
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
              <div class="p-2">
                <img
                  src={CDN_URL + '/' + query.data.avatar_url}
                  alt="avatar"
                  class="aspect-square size-full rounded-xl object-cover"
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
                        <div class="flex size-10 shrink-0 items-center justify-center rounded-full bg-secondary">
                          <span class="material-symbols-rounded shrink-0 text-main">
                            {String.fromCodePoint(parseInt(op.icon!, 16))}
                          </span>
                        </div>
                        <div class="text-start text-white">
                          <p class="text-sm font-semibold">{op.text}</p>
                          <p class="text-xs leading-tight text-white/60">
                            {op.description}
                          </p>
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
      class="absolute right-4 top-4 z-10 h-9 w-[90px] rounded-xl bg-black/80 px-2.5 text-button"
      onClick={props.onClick}
    >
      {props.text}
    </button>
  );
};

const Loader = () => {
  return (
    <div class="flex h-screen flex-col items-start justify-start bg-secondary p-2">
      <div class="aspect-square w-full rounded-xl bg-main" />
      <div class="flex flex-col items-start justify-start p-2">
        <div class="mt-2 h-6 w-1/2 rounded bg-main" />
        <div class="mt-2 h-6 w-1/2 rounded bg-main" />
        <div class="mt-2 h-20 w-full rounded bg-main" />
        <div class="mt-4 flex w-full flex-row flex-wrap items-center justify-start gap-2">
          <div class="h-10 w-40 rounded-2xl bg-main" />
          <div class="h-10 w-32 rounded-2xl bg-main" />
          <div class="h-10 w-40 rounded-2xl bg-main" />
          <div class="h-10 w-28 rounded-2xl bg-main" />
          <div class="h-10 w-32 rounded-2xl bg-main" />
        </div>
      </div>
    </div>
  );
};
