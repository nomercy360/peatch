import { createEffect, createResource, createSignal, For, Match, onCleanup, Suspense, Switch } from 'solid-js';
import { useNavigate, useParams, useSearchParams } from '@solidjs/router';
import { CDN_URL, fetchProfile, followUser, hideProfile, publishProfile, showProfile, unfollowUser } from '~/lib/api';
import { setUser, store } from '~/store';
import ActionDonePopup from '~/components/ActionDonePopup';
import { useMainButton } from '~/lib/useMainButton';
import { usePopup } from '~/lib/usePopup';

export default function UserProfile() {
  const mainButton = useMainButton();
  const [published, setPublished] = createSignal(false);

  const navigate = useNavigate();

  const params = useParams();
  const [searchParams] = useSearchParams();

  const username = params.handle;

  const [profile, { mutate, refetch }] = createResource(() =>
    fetchProfile(username),
  );

  const { showAlert } = usePopup();

  const isCurrentUserProfile = store.user.username === username;

  const navigateToEdit = () => {
    navigate('/users/edit', { state: { back: true } });
  };

  createEffect(async () => {
    if (searchParams.refetch) {
      await refetch();
      if (profile()?.id === store.user.id) setUser(profile()!);
    }
  });

  const navigateToCollaborate = async () => {
    if (!store.user.published_at) {
      showAlert(
        `Publish your profile first, so ${profile()?.first_name} will see it`,
      );
    } else if (store.user.hidden_at) {
      showAlert(
        `Unhide your profile first, so ${profile()?.first_name} will see it`,
      );
    } else {
      navigate(`/users/${username}/collaborate`, { state: { back: true } });
    }
  };

  const closePopup = () => {
    setPublished(false);
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
    if (!profile()) return;
    mutate({
      ...profile(),
      is_following: true,
      followers_count: profile()?.followers_count + 1,
    });
    window.Telegram.WebApp.HapticFeedback.impactOccurred('light');
    await followUser(Number(profile()?.id));
  };

  const unfollow = async () => {
    mutate({
      ...profile(),
      is_following: false,
      followers_count: profile()?.followers_count - 1,
    });
    window.Telegram.WebApp.HapticFeedback.impactOccurred('light');
    await unfollowUser(Number(profile()?.id));
  };

  createEffect(() => {
    if (isCurrentUserProfile) {
      if (!store.user.published_at) {
        mainButton.enable('Publish');
        mainButton.onClick(publish);
      } else {
        if (published()) {
          mainButton.onClick(closePopup);
          mainButton.enable('Back to profile');
        } else {
          mainButton.enable('Edit');
          mainButton.onClick(navigateToEdit);
        }
      }
    } else {
      mainButton.enable('Collaborate');
      mainButton.onClick(navigateToCollaborate);
    }

    onCleanup(() => {
      mainButton.offClick(navigateToCollaborate);
      mainButton.offClick(publish);
      mainButton.offClick(closePopup);
      mainButton.offClick(navigateToEdit);
    });
  });

  onCleanup(async () => {
    mainButton.hide();
  });

  const [contentCopied, setContentCopied] = createSignal(false);

  async function copyToClipboard() {
    try {
      await navigator.clipboard.writeText(window.location.href);
      setContentCopied(true);
      window.Telegram.WebApp.HapticFeedback.impactOccurred('light');
      setTimeout(() => setContentCopied(false), 2000);
    } catch (err) {
      console.error('Failed to copy: ', err);
      window.Telegram.WebApp.sendData(window.location.href);
    }
  }

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
          <Match when={profile()}>
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
                <Match when={!isCurrentUserProfile && !profile()?.is_following}>
                  <ActionButton text="Follow" onClick={follow} />
                </Match>
                <Match when={!isCurrentUserProfile && profile()?.is_following}>
                  <ActionButton text="Unfollow" onClick={unfollow} />
                </Match>
              </Switch>
              <div class="p-2">
                <img
                  src={CDN_URL + '/' + profile()?.avatar_url}
                  alt="avatar"
                  class="aspect-square size-full rounded-xl object-cover"
                />
              </div>
              <div class="px-4 py-2.5">
                <div class="flex flex-row items-center justify-between pb-4">
                  <div class="flex h-8 flex-row items-center space-x-2 text-sm font-semibold">
										<span class="flex flex-row items-center text-main">
											{profile().followers_count}
										</span>
                    <span class="text-secondary">following</span>
                    <span class="text-secondary">·</span>
                    <span class="flex flex-row items-center text-main">
											{profile().following_count}
										</span>
                    <span class="text-secondary">followers</span>
                  </div>
                  <button
                    class="flex h-8 flex-row items-center space-x-2 bg-transparent px-2.5"
                    classList={{
                      'text-main': !contentCopied(),
                      'text-green': contentCopied(),
                    }}
                    onClick={copyToClipboard}
                  >
                    <span class="text-sm font-semibold">share app profile</span>
                    <span class="material-symbols-rounded text-[14px]">
											{contentCopied() ? 'check_circle' : 'content_copy'}
										</span>
                  </button>
                </div>
                <p class="text-3xl text-pink">
                  {profile()?.first_name} {profile()?.last_name}:
                </p>
                <p class="text-3xl text-main">{profile()?.title}</p>
                <p class="mt-1 text-lg font-normal text-secondary">
                  {profile()?.description}
                </p>
                <div class="mt-5 flex flex-row flex-wrap items-center justify-start gap-1">
                  <For each={profile()?.badges}>
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
                  <For each={profile()?.opportunities}>
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
      class="absolute right-4 top-4 z-10 h-9 w-[90px] rounded-xl bg-black/80 px-2.5 text-sm font-semibold text-button"
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
