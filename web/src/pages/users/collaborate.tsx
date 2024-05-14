import { CDN_URL, createUserCollaboration, fetchProfile, findUserCollaborationRequest } from '~/lib/api';
import { useNavigate, useParams } from '@solidjs/router';
import { store } from '~/store';
import { createEffect, createResource, createSignal, Match, onCleanup, Show, Switch } from 'solid-js';
import TextArea from '~/components/TextArea';
import { usePopup } from '~/lib/usePopup';
import { useMainButton } from '~/lib/useMainButton';
import BadgeList from '~/components/BadgeList';
import ActionDonePopup from '~/components/ActionDonePopup';

export default function Collaborate() {
  const params = useParams();
  const username = params.handle;

  const [created, setCreated] = createSignal(false);
  const mainButton = useMainButton();
  const { showConfirm } = usePopup();
  const navigate = useNavigate();

  const backToProfile = () => {
    navigate(`/users/${username}`, { state: { from: '/users' } });
  };

  const [message, setMessage] = createSignal('');

  const [profile] = createResource(() => fetchProfile(username));

  const [existedRequest] = createResource(async () => {
    try {
      return await findUserCollaborationRequest(username);
    } catch (e: unknown) {
      if ((e as { status: number }).status === 404) {
        return null;
      }
    }
  });

  const postCollaboration = async () => {
    if (!store.user.published_at) {
      showConfirm(
        'You must publish your profile first',
        (ok: boolean) =>
          ok && navigate('/users/edit', { state: { back: true } }),
      );
      return;
    }
    try {
      await createUserCollaboration(profile()?.id, message());
      setCreated(true);
    } catch (e) {
      console.error(e);
    }
  };

  createEffect(() => {
    if (created() || existedRequest()) {
      mainButton.offClick(postCollaboration);
      mainButton.onClick(backToProfile);
      mainButton.enable(`Back to ${profile()?.first_name}'s profile`);
    } else if (!existedRequest.loading && !existedRequest()) {
      mainButton.onClick(postCollaboration);
      if (message()) {
        mainButton.enable('Send message');
      } else {
        mainButton.disable('Send message');
      }
    }

    onCleanup(() => {
      mainButton.offClick(postCollaboration);
      mainButton.offClick(backToProfile);
    });
  })

  return (
    <Switch>
      <Match when={created()}>
        <ActionDonePopup
          action="Message sent"
          description={`Once ${profile()?.first_name} accepts your invitation, we'll share your contacts`}
          callToAction={`There are 12 people with a similar profiles like ${profile()?.first_name}`}
        />
      </Match>
      <Match when={existedRequest.loading && !profile()}>
        <div />
      </Match>
      <Match when={!existedRequest.loading && profile()}>
        <Show when={existedRequest()}>
          <ActionDonePopup
            action="Message sent"
            description={`Once ${profile()?.first_name} accepts your invitation, we'll share your contacts`}
            callToAction={`There are 12 people with a similar profiles like ${profile()?.first_name}`}
          />
        </Show>
        <Show when={!existedRequest()}>
          <div class="flex flex-col items-center justify-center bg-secondary p-4">
            <div class="mb-4 mt-1 flex flex-col items-center justify-center text-center">
              <p class="max-w-[220px] text-3xl text-main">
                Collaborate with {profile()?.first_name}
              </p>
              <div class="my-5 flex w-full flex-row items-center justify-center">
                <img
                  class="z-10 size-24 rounded-3xl border-2 border-secondary object-cover object-center"
                  src={CDN_URL + '/' + store.user.avatar_url}
                  alt="User Avatar"
                />
                <img
                  class="-ml-4 size-24 rounded-3xl border-2 border-secondary object-cover object-center"
                  src={CDN_URL + '/' + profile()?.avatar_url}
                  alt="User Avatar"
                />
              </div>
              <Show when={profile()?.badges && profile().badges.length > 0}>
                <BadgeList
                  badges={profile()?.badges}
                  position="center"
                  city={profile()?.city}
                  country={profile()?.country}
                  countryCode={profile()?.country_code}
                />
              </Show>
            </div>
            <TextArea
              value={message()}
              setValue={(value: string) => setMessage(value)}
              placeholder="Write a message to start collaboration"
            />
          </div>
        </Show>
      </Match>
    </Switch>
  );
}
