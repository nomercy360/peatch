import { createQuery } from '@tanstack/solid-query';
import { CDN_URL, createUserCollaboration, fetchProfile } from '../../api';
import { useNavigate, useParams } from '@solidjs/router';
import { store } from '../../store';
import { createEffect, createSignal, onCleanup, Show, Suspense } from 'solid-js';
import { useButtons } from '../../hooks/useBackButton';
import { createStore } from 'solid-js/store';
import { CreateUserCollaboration } from '../../../gen';
import TextArea from '../../components/textArea';
import { usePopup } from '../../hooks/usePopup';

export default function Collaborate() {
  const params = useParams();
  const userId = params.id;

  const { mainButton, backButton } = useButtons();
  const [created, setCreated] = createSignal(false);

  const navigate = useNavigate();
  const { showAlert } = usePopup();

  const [collaborationRequest, setCollaborationRequest] = createStore<CreateUserCollaboration>({
    user_id: Number(userId),
    message: '',
    requester_id: store.user.id!,
  });

  const backToProfile = () => {
    navigate('/users/' + userId);
  };

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
      const resp = await  createUserCollaboration(collaborationRequest);
      setCreated(true);
    } catch (e) {
      console.error(e);
    }
  };

  createEffect(() => {
    backButton.setVisible();
    backButton.onClick(backToProfile);
    if (created()) {
      mainButton.offClick(postCollaboration);
      mainButton.onClick(backToProfile);
      mainButton.setText('Back to ' + query.data.first_name + '\'s profile');
    } else {
      mainButton.setVisible('Send message')
      mainButton.onClick(postCollaboration);
      if (collaborationRequest.message === '') {
        mainButton.setActive(false);
      } else {
        mainButton.setActive(true);
      }
    }
  });

  onCleanup(() => {
    backButton.offClick(backToProfile);
    backButton.offClick(backToProfile);
    mainButton.offClick(postCollaboration);
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
            setValue={(value: string) => setCollaborationRequest('message', value)}
            placeholder="Write a message to start collaboration"
          />
        </div>
      </Show>
    </Suspense>
  );
}
