import {
  createEffect,
  createSignal,
  For,
  Match,
  onCleanup,
  Suspense,
  Switch,
} from 'solid-js';
import { useNavigate, useParams } from '@solidjs/router';
import {
  CDN_URL,
  fetchCollaboration,
  hideCollaboration,
  publishCollaboration,
  showCollaboration,
} from '~/api';
import { createQuery } from '@tanstack/solid-query';
import { store } from '~/store';
import { usePopup } from '~/hooks/usePopup';
import ProfilePublished from '~/components/ProfilePublished';
import { useMainButton } from '~/hooks/useMainButton';
import { useNavigation } from '~/hooks/useNavigation';

export default function Collaboration() {
  const mainButton = useMainButton();

  const [published, setPublished] = createSignal(false);

  const navigate = useNavigate();
  const { navigateBack } = useNavigation();
  const params = useParams();
  const { showAlert } = usePopup();

  const collabId = params.id;

  const query = createQuery(() => ({
    queryKey: ['collaborations', collabId],
    queryFn: () => fetchCollaboration(Number(collabId)),
  }));

  const [isCurrentUserCollab, setIsCurrentUserCollab] = createSignal(false);

  createEffect(() => {
    setIsCurrentUserCollab(query.data?.user_id === store.user.id);
  });

  const publish = async () => {
    await publishCollaboration(Number(collabId));
    setPublished(true);
  };

  const hide = async () => {
    await hideCollaboration(Number(collabId));
  };

  const show = async () => {
    await showCollaboration(Number(collabId));
  };

  const collaborate = async () => {
    if (store.user.published_at && !store.user.hidden_at) {
      navigate(`/collaborations/${collabId}/collaborate`);
    } else {
      showAlert('Fill and publish your profile first');
    }
  };

  const pushToEdit = () => {
    navigate(`/collaborations/${collabId}/edit`);
  };

  createEffect(() => {
    if (isCurrentUserCollab()) {
      if (published()) {
        mainButton.offClick(publish);
        mainButton.offClick(pushToEdit);
        mainButton.setParams({
          text: 'Just open the app',
          isVisible: true,
          isEnabled: true,
        });
        mainButton.onClick(navigateBack);
        return;
      }
      if (!store.user.published_at) {
        // mainButton.offClick(pushToEdit);
        // mainButton.setVisible('Publish');
        // mainButton.onClick(publish);
        mainButton.setParams({
          text: 'Publish',
          isVisible: true,
          isEnabled: true,
        });
        mainButton.onClick(publish);
      } else {
        // mainButton.offClick(publish);
        // mainButton.setVisible('Edit');
        // mainButton.onClick(pushToEdit);
      }
    } else {
      // mainButton.offClick(pushToEdit);
      // mainButton.setVisible('Collaborate');
      // mainButton.onClick(collaborate);
    }
  });

  onCleanup(() => {
    // mainButton.hide();
    // mainButton.offClick(collaborate);
    // mainButton.offClick(publish);
    // mainButton.offClick(back);
    // mainButton.offClick(pushToEdit);
  });

  return (
    <div>
      <Suspense fallback={<div>Loading...</div>}>
        <Switch>
          <Match when={published() && isCurrentUserCollab()}>
            <ProfilePublished />
          </Match>
          <Match when={query.data}>
            <div class="min-h-screen">
              <Switch>
                <Match when={isCurrentUserCollab() && !query.data.published_at}>
                  <ActionButton text="Edit" onClick={pushToEdit} />
                </Match>
                <Match
                  when={
                    isCurrentUserCollab() &&
                    query.data.hidden_at &&
                    query.data.published_at
                  }
                >
                  <ActionButton text="Show" onClick={show} />
                </Match>
                <Match
                  when={
                    isCurrentUserCollab() &&
                    !query.data.hidden_at &&
                    query.data.published_at
                  }
                >
                  <ActionButton text="Hide" onClick={hide} />
                </Match>
              </Switch>
              <div class="flex w-full flex-col items-start justify-start bg-pink px-4 pb-5 pt-4">
                <span class="material-symbols-rounded text-[48px] text-white">
                  {String.fromCodePoint(
                    parseInt(query.data.opportunity.icon!, 16),
                  )}
                </span>
                <p class="text-3xl text-white">
                  {query.data.opportunity.text}:
                </p>
                <p class="text-3xl text-white">{query.data.title}:</p>
                <div class="mt-4 flex w-full flex-row items-center justify-start gap-2">
                  <img
                    class="size-10 rounded-2xl object-cover"
                    src={CDN_URL + '/' + query.data.user?.avatar_url}
                    alt="User Avatar"
                  />
                  <div>
                    <p class="text-sm font-bold text-white">
                      {query.data.user?.first_name} {query.data.user?.last_name}
                      :
                    </p>
                    <p class="text-sm text-white">{query.data.user?.title}</p>
                  </div>
                </div>
              </div>
              <div class="px-4 py-2.5">
                <p class="text-lg font-normal text-main">
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
      class="absolute left-4 top-4 z-10 h-8 w-20 rounded-lg px-2.5 text-white"
      onClick={props.onClick}
      style={{ background: 'rgba(255, 255, 255, 0.20)' }}
    >
      {props.text}
    </button>
  );
};
