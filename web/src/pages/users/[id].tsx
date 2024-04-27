import { useButtons } from '../../hooks/useBackButton';
import { createEffect, For, onCleanup, Show, Suspense } from 'solid-js';
import { useNavigate, useParams } from '@solidjs/router';
import { CDN_URL, fetchProfile, followUser, hideProfile, showProfile, unfollowUser } from '../../api';
import { createQuery } from '@tanstack/solid-query';
import { User } from '../../../gen';
import { setUser, store } from '../../store';

export default function UserProfile() {
  const { mainButton, backButton } = useButtons();

  const navigate = useNavigate();
  const params = useParams();

  const userId = params.id;

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
    await showProfile();
  };

  const hide = async () => {
    setUser({
      ...store.user,
      published_at: undefined,
    });
    await hideProfile();
  };

  const follow = async () => {
    store.following.push(Number(userId));
    await followUser(Number(userId));
  };

  const unfollow = async () => {
    store.following = store.following.filter(f => f !== Number(userId));
    await unfollowUser(Number(userId));
  };

  const collaborate = async () => {
    navigate(`/collaborate/${userId}`);
  };

  const pushToEdit = () => {
    navigate(`/users/edit`);
  };

  createEffect(() => {
    if (isCurrentUserProfile) {
      if (!store.user.published_at) {
        mainButton.offClick(pushToEdit);
        mainButton.setVisible(true);
        mainButton.setText('Publish');
        mainButton.setColor('#FF5A5F');
        mainButton.onClick(publish);
      } else {
        mainButton.offClick(publish);
        mainButton.setVisible(true);
        mainButton.setText('Edit');
        mainButton.setColor('#D9D9D9');
        mainButton.setTextColor('#000000');
        mainButton.onClick(pushToEdit);
      }
    } else {
      mainButton.offClick(pushToEdit);
      mainButton.setVisible(true);
      mainButton.setText('Collaborate');
      mainButton.setColor('#F3A333');
      mainButton.onClick(collaborate);
    }
  });

  return (
    <div>
      <Suspense fallback={<div>Loading...</div>}>
        <Show when={query.data}>
          <ProfileCard
            user={query.data}
            buttonTitle={
              isCurrentUserProfile
                ? store.user.published_at
                  ? 'Hide'
                  : 'Edit'
                : store.following.includes(Number(userId))
                  ? 'Unfollow'
                  : 'Follow'
            }
            onAction={
              isCurrentUserProfile
                ? store.user.published_at
                  ? hide
                  : pushToEdit
                : store.following.includes(Number(userId))
                  ? unfollow
                  : follow
            }
          />
        </Show>
      </Suspense>
    </div>
  );
}
// background: ;

const ProfileCard = (props: {
  user: User;
  buttonTitle: string;
  onAction: () => void;
}) => {
  return (
    <div class="min-h-screen">
      <ActionButton text={props.buttonTitle} onClick={props.onAction} />
      <div class="image-container">
        <img
          src={CDN_URL + '/' + props.user.avatar_url}
          alt="avatar"
          class="aspect-square size-full object-cover"
        />
      </div>
      <div class="px-4 py-2.5">
        <p class="text-3xl text-pink">
          {props.user.first_name} {props.user.last_name}:
        </p>
        <p class="text-3xl text-black">{props.user.title}</p>
        <p class="text-lg font-normal"> {props.user.description}</p>
        <div class="mt-5 flex flex-row flex-wrap items-center justify-start gap-1">
          <For each={props.user.badges}>
            {badge => (
              <div
                class="flex h-10 flex-row items-center justify-center gap-[5px] rounded-2xl border border-peatch-stroke px-2.5"
                style={{
                  'background-color': `${badge.color}`,
                }}
              >
                <span class="material-symbols-rounded text-white">
                  {String.fromCodePoint(parseInt(badge.icon!, 16))}
                </span>
                <p class="text-sm font-semibold text-white">{badge.text}</p>
              </div>
            )}
          </For>
        </div>
        <div class="mt-5 flex w-full flex-col items-center justify-start gap-1">
          <For each={props.user.opportunities}>
            {op => (
              <div
                class="flex h-[60px] w-full flex-row items-center justify-start gap-2.5 rounded-2xl border border-peatch-stroke px-2.5"
                style={{
                  'background-color': op.color,
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
  );
};

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
