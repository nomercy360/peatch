import {
  createEffect,
  createSignal,
  For,
  Match,
  onCleanup,
  Suspense,
  Switch,
} from 'solid-js';
import { useNavigate, useParams, useSearchParams } from '@solidjs/router';
import {
  CDN_URL,
  fetchCollaboration,
  hideCollaboration,
  publishCollaboration,
  showCollaboration,
} from '~/api';
import { createQuery } from '@tanstack/solid-query';
import { setUser, store } from '~/store';
import ActionDonePopup from '~/components/ActionDonePopup';
import { useMainButton } from '~/hooks/useMainButton';
import { useNavigation } from '~/hooks/useNavigation';
import { usePopup } from '~/hooks/usePopup';

export default function Collaboration() {
  const mainButton = useMainButton();

  const [published, setPublished] = createSignal(false);

  const navigate = useNavigate();

  const { showConfirm } = usePopup();

  const { navigateBack } = useNavigation();

  const params = useParams();
  const [searchParams, _] = useSearchParams();

  const collabId = params.id;

  const query = createQuery(() => ({
    queryKey: ['collaborations', collabId],
    queryFn: () => fetchCollaboration(Number(collabId)),
  }));

  createEffect(async () => {
    if (searchParams.refetch) {
      await query.refetch();
      if (query.data.id === store.user.id) {
        setUser(query.data);
      }
    }
  });

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
    await query.refetch();
  };

  const show = async () => {
    await showCollaboration(Number(collabId));
    await query.refetch();
  };

  const navigateToCollaborate = async () => {
    if (store.user.published_at && !store.user.hidden_at) {
      navigate(`/collaborations/${collabId}/collaborate`);
    } else {
      showConfirm(
        'You must publish your profile first',
        (ok: boolean) => ok && navigate('/users/edit'),
      );
    }
  };

  const navigateToEdit = () => {
    navigate('/collaborations/edit/' + collabId, {
      state: { from: '/collaborations/' + collabId },
    });
  };

  createEffect(() => {
    if (isCurrentUserCollab()) {
      if (published()) {
        mainButton.offClick(publish);
        mainButton.offClick(navigateToEdit);
        mainButton.enable('Get back');
        mainButton.onClick(navigateBack);
        return;
      } else if (!query.data.published_at) {
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
      <Suspense fallback={<div>Loading...</div>}>
        <Switch>
          <Match when={published() && isCurrentUserCollab()}>
            <ActionDonePopup
              action="Collaboration published!"
              description="We have shared your collaboration with the community"
              callToAction="There are 12 people you might be interested to collaborate with"
            />
          </Match>
          <Match when={query.data}>
            <div class="h-fit min-h-screen bg-secondary">
              <Switch>
                <Match when={isCurrentUserCollab() && !query.data.published_at}>
                  <ActionButton text="Edit" onClick={navigateToEdit} />
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
                <p class="text-lg font-normal text-secondary">
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
      class="absolute right-4 top-4 z-10 h-9 w-[90px] rounded-xl bg-black/80 px-2.5 text-button"
      onClick={() => props.onClick()}
    >
      {props.text}
    </button>
  );
};
