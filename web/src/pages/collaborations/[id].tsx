import { createEffect, createResource, createSignal, For, Match, onCleanup, Suspense, Switch } from 'solid-js';
import { useNavigate, useParams, useSearchParams } from '@solidjs/router';
import { CDN_URL, fetchCollaboration, hideCollaboration, publishCollaboration, showCollaboration } from '~/lib/api';
import { store } from '~/store';
import ActionDonePopup from '~/components/ActionDonePopup';
import { useMainButton } from '~/lib/useMainButton';
import { usePopup } from '~/lib/usePopup';

export default function Collaboration() {
  const mainButton = useMainButton();

  const [wasPublished, setWasPublished] = createSignal(false);
  const [isCurrentUserCollab, setIsCurrentUserCollab] = createSignal(false);

  const navigate = useNavigate();

  const { showAlert } = usePopup();

  const params = useParams();
  const [searchParams] = useSearchParams();

  const collabId = params.id;

  const [collaboration, { mutate, refetch }] = createResource(() =>
    fetchCollaboration(Number(collabId)),
  );

  createEffect(async () => {
    if (searchParams.refetch) await refetch();
  });

  createEffect(() => {
    if (collaboration().id) {
      setIsCurrentUserCollab(store.user.id === collaboration().user_id);
    }
  });

  const closePopup = () => {
    setWasPublished(false);
  };

  const publish = async () => {
    await publishCollaboration(Number(collabId));
    window.Telegram.WebApp.HapticFeedback.impactOccurred('light');
    mutate({ ...collaboration(), published_at: new Date().toISOString() });
    setWasPublished(true);
  };

  const hide = async () => {
    await hideCollaboration(Number(collabId));
    window.Telegram.WebApp.HapticFeedback.impactOccurred('light');
    mutate({ ...collaboration(), hidden_at: new Date().toISOString() });
  };

  const show = async () => {
    await showCollaboration(Number(collabId));
    window.Telegram.WebApp.HapticFeedback.impactOccurred('light');
    mutate({ ...collaboration(), hidden_at: null });
  };

  const navigateToCollaborate = async () => {
    if (!store.user.published_at) {
      showAlert('Publish your profile first, so collaborators will see you');
    } else if (store.user.hidden_at) {
      showAlert('Unhide your profile first, so collaborators will see you');
    } else {
      navigate(`/collaborations/${collabId}/collaborate`, {
        state: { back: true },
      });
    }
  };

  const navigateToEdit = () => {
    navigate('/collaborations/edit/' + collabId, {
      state: { from: '/collaborations/' + collabId },
    });
  };

  createEffect(() => {
    if (isCurrentUserCollab()) {
      if (!collaboration().published_at) {
        mainButton.enable('Publish');
        mainButton.onClick(publish);
      } else {
        if (wasPublished()) {
          mainButton.enable('Back to collaboration');
          mainButton.onClick(closePopup);
        } else {
          mainButton.enable('Edit');
          mainButton.onClick(navigateToEdit);
        }
      }
    } else {
      mainButton.onClick(navigateToCollaborate);
      mainButton.enable('Collaborate');
    }

    onCleanup(() => {
      mainButton.offClick(closePopup);
      mainButton.offClick(publish);
      mainButton.offClick(navigateToEdit);
      mainButton.offClick(navigateToCollaborate);
    });
  })

  onCleanup(async () => {
    mainButton.hide();
  });

  return (
    <Suspense fallback={<Loader />}>
      <Switch>
        <Match when={wasPublished() && isCurrentUserCollab()}>
          <ActionDonePopup
            action="Collaboration published!"
            description="We have shared your collaboration with the community"
            callToAction="There are 12 people you might be interested to collaborate with"
          />
        </Match>
        <Match when={!collaboration.loading}>
          <div class="h-fit min-h-screen bg-secondary">
            <Switch>
              <Match
                when={isCurrentUserCollab() && !collaboration().published_at}
              >
                <ActionButton text="Edit" onClick={navigateToEdit} />
              </Match>
              <Match
                when={
                  isCurrentUserCollab() &&
                  collaboration().hidden_at &&
                  collaboration().published_at
                }
              >
                <ActionButton text="Show" onClick={show} />
              </Match>
              <Match
                when={
                  isCurrentUserCollab() &&
                  !collaboration().hidden_at &&
                  collaboration().published_at
                }
              >
                <ActionButton text="Hide" onClick={hide} />
              </Match>
            </Switch>
            <div
              class="flex w-full flex-col items-start justify-start px-4 pb-5 pt-4"
              style={{
                'background-color': `#${collaboration().opportunity.color}`,
              }}
            >
							<span class="material-symbols-rounded text-[48px] text-white">
								{String.fromCodePoint(
                  parseInt(collaboration().opportunity.icon!, 16),
                )}
							</span>
              <p class="text-3xl text-white">
                {collaboration().opportunity.text}:
              </p>
              <p class="text-3xl text-white">{collaboration().title}:</p>
              <div class="mt-4 flex w-full flex-row items-center justify-start gap-2">
                <img
                  class="size-11 rounded-xl object-cover"
                  src={CDN_URL + '/' + collaboration().user?.avatar_url}
                  alt="User Avatar"
                />
                <div>
                  <p class="text-sm font-bold text-white">
                    {collaboration().user?.first_name}{' '}
                    {collaboration().user?.last_name}:
                  </p>
                  <p class="text-sm text-white">
                    {collaboration().user?.title}
                  </p>
                </div>
              </div>
            </div>
            <div class="px-4 py-2.5">
              <p class="text-lg font-normal text-secondary">
                {collaboration().description}
              </p>
              <div class="mt-5 flex flex-row flex-wrap items-center justify-start gap-1">
                <For each={collaboration().badges}>
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
}

const Loader = () => {
  return (
    <div class="flex h-screen flex-col items-start justify-start bg-secondary">
      <div class="h-[260px] w-full bg-main" />
      <div class="flex flex-col items-start justify-start p-4">
        <div class="h-36 w-full rounded bg-main" />
        <div class="mt-4 flex w-full flex-row flex-wrap items-center justify-start gap-2">
          <div class="h-10 w-40 rounded-2xl bg-main" />
          <div class="h-10 w-32 rounded-2xl bg-main" />
          <div class="h-10 w-36 rounded-2xl bg-main" />
          <div class="h-10 w-24 rounded-2xl bg-main" />
          <div class="h-10 w-40 rounded-2xl bg-main" />
          <div class="h-10 w-28 rounded-2xl bg-main" />
          <div class="h-10 w-32 rounded-2xl bg-main" />
        </div>
      </div>
    </div>
  );
}
