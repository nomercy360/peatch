import { createQuery } from '@tanstack/solid-query';
import { CDN_URL, fetchProfile } from '../../api';
import { useNavigate, useParams } from '@solidjs/router';
import { store } from '../../store';
import { createEffect, Show, Suspense } from 'solid-js';
import { useButtons } from '../../hooks/useBackButton';

export default function Collaborate() {
  const params = useParams();
  const userId = params.id;

  const { mainButton, backButton } = useButtons();

  const navigate = useNavigate();

  const back = () => {
    navigate('/');
  };

  const query = createQuery(() => ({
    queryKey: ['profiles', userId],
    queryFn: () => fetchProfile(Number(userId)),
  }));

  createEffect(() => {
    backButton.setVisible();
    backButton.onClick(back);
  });

  const postCollaboration = async () => {
    createUserCollaboration();
  };

  return (
    <Suspense fallback={<div>Loading...</div>}>
      <Show when={query.data}>
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
        </div>
      </Show>
    </Suspense>
  );
}
