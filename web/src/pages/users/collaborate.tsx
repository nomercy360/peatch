import { createQuery } from '@tanstack/solid-query';
import { CDN_URL, createUserCollaboration, fetchProfile } from '../../api';
import { useNavigate, useParams } from '@solidjs/router';
import { store } from '~/store';
import {
  createEffect,
  createSignal,
  onCleanup,
  Show,
  Suspense,
} from 'solid-js';

import { createStore } from 'solid-js/store';
import { CreateUserCollaboration } from '../../../gen';
import TextArea from '~/components/TextArea';
import { usePopup } from '~/hooks/usePopup';
import { useMainButton } from '@tma.js/sdk-solid';

export default function Collaborate() {
  const params = useParams();
  const userId = params.id;

  const [created, setCreated] = createSignal(false);
  const mainButton = useMainButton();
  const { showAlert } = usePopup();
  const navigate = useNavigate();

  const navigateBack = () => {
    navigate(-1 );
  }

  const [collaborationRequest, setCollaborationRequest] =
    createStore<CreateUserCollaboration>({
      user_id: Number(userId),
      message: '',
      requester_id: store.user.id!,
    });


  const query = createQuery(() => ({
    queryKey: ['profiles', userId],
    queryFn: () => fetchProfile(Number(userId)),
  }));

  const postCollaboration = async () => {
    if (!store.user.published_at) {
      showAlert('You must publish your profile first');
      return;
    }
    try {
      await createUserCollaboration(collaborationRequest);
      setCreated(true);
    } catch (e) {
      console.error(e);
    }
  };

  createEffect(() => {
    if (created()) {
      mainButton().off('click', postCollaboration);
      mainButton().on('click', navigateBack);
      mainButton().setParams({
        text: 'Back to ' + query.data.first_name + '\'s profile',
        isEnabled: true,
        isVisible: true,
      });

    } else {
      mainButton().setParams({
        text: 'Send message',
        isEnabled: collaborationRequest.message !== '',
        isVisible: true,
      });
      mainButton().on('click', postCollaboration);
    }
  });

  onCleanup(() => {
    mainButton().off('click', postCollaboration);
    mainButton().off('click', navigateBack);
  });

  return (
    <Suspense fallback={<div>Loading...</div>}>
      <Show when={query.data} fallback={<div>Loading...</div>}>
        <div class="flex flex-col items-center justify-center p-4">
          <p class="text-3xl">Collaborate with {query.data.first_name}</p>
          <div class="mt-5 flex w-full flex-row items-center justify-center">
            <img
              class="z-10 size-24 rounded-3xl border-2 border-white object-cover object-center"
              src={CDN_URL + '/' + store.user.avatar_url}
              alt="User Avatar"
            />
            <img
              class="-ml-4 size-24 rounded-3xl border-2 border-white object-cover object-center"
              src={CDN_URL + '/' + query.data.avatar_url}
              alt="User Avatar"
            />
          </div>
          <p class="text-3xl">{query.data.title}</p>
          <p class="text-sm text-gray">{query.data.description}</p>
          <TextArea
            value={collaborationRequest.message!}
            setValue={(value: string) =>
              setCollaborationRequest('message', value)
            }
            placeholder="Write a message to start collaboration"
          />
        </div>
      </Show>
    </Suspense>
  );
}
