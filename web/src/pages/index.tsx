import {
  createEffect,
  createResource,
  createSignal,
  For,
  Match,
  Show,
  Suspense,
  Switch,
} from 'solid-js';
import { store } from '~/store';
import { CDN_URL, fetchPreview } from '~/api';
import FillProfilePopup from '~/components/FillProfilePopup';
import { Link } from '~/components/Link';
import { useMainButton } from '~/hooks/useMainButton';

export default function Index() {
  const [profilePopup, setProfilePopup] = createSignal(false);

  const [previewImages, _] = createResource(async () => {
    const res = await fetchPreview();
    return res.map((image: any) => CDN_URL + '/' + image.avatar_url);
  });

  const getUserLink = () => {
    if (store.user.first_name && store.user.last_name) {
      return '/users/' + store.user?.id;
    } else {
      return '/users/edit';
    }
  };

  const mainButton = useMainButton();

  mainButton.hide();

  const closePopup = () => {
    setProfilePopup(false);
    window.Telegram.WebApp.CloudStorage.setItem('profilePopup', 'closed');
  };

  const updateProfilePopup = (err: any, value: any) => {
    setProfilePopup(value !== 'closed');
  };

  createEffect(() => {
    window.Telegram.WebApp.CloudStorage.getItem(
      'profilePopup',
      updateProfilePopup,
    );
    // window.Telegram.WebApp.CloudStorage.removeItem('profilePopup');
  });

  return (
    <div class="flex min-h-screen flex-col bg-secondary px-4">
      <Show when={!store.user.published_at && profilePopup()}>
        <FillProfilePopup onClose={closePopup} />
      </Show>
      <Link
        class="flex flex-row items-center justify-between py-4"
        href={getUserLink()}
      >
        <p class="text-3xl text-main">
          Bonsoir, {store.user?.first_name || store.user?.username}!
        </p>
        <Switch>
          <Match when={store.user.avatar_url}>
            <img
              class="size-11 rounded-xl border border-main object-cover object-center"
              src={CDN_URL + '/' + store.user.avatar_url}
              alt="User Avatar"
            />
          </Match>
          <Match when={!store.user.avatar_url}>
            <div class="flex size-11 items-center justify-center rounded-xl border-2 border-main bg-main">
              <span class="material-symbols-rounded text-peatch-main">
                account_circle
              </span>
            </div>
          </Match>
        </Switch>
      </Link>
      <div class="h-px w-full bg-main" />
      <Link class="flex flex-col items-start justify-start py-4" href="/users">
        <div class="flex w-full flex-row items-center justify-start">
          <Suspense fallback={<ImagesLoader />}>
            <For each={previewImages()}>
              {(image, idx) => (
                <img
                  src={image}
                  alt="User Avatar"
                  class="-ml-1 size-11 rounded-xl border border-main object-cover object-center"
                  classList={{
                    'ml-0': idx() === 0,
                    'z-20': idx() === 0,
                    'z-10': idx() === 1,
                  }}
                />
              )}
            </For>
          </Suspense>
        </div>
        <div class="flex flex-row items-center justify-between">
          <p class="mt-2 text-3xl text-main">
            <span class="text-accent">Explore people</span> you may like to
            collaborate
          </p>
          <span class="material-symbols-rounded text-[48px] text-pink">
            maps_ugc
          </span>
        </div>
        <p class="mt-1.5 text-sm text-hint">
          Figma Wizards, Consultants, Founders, and more
        </p>
      </Link>
      <div class="h-px w-full bg-main" />
      <Link
        class="flex flex-col items-start justify-start py-4"
        href="/collaborations"
      >
        <div class="flex w-full flex-row items-center justify-start">
          <div class="z-20 flex size-11 flex-col items-center justify-center rounded-xl border border-main bg-orange">
            <span class="material-symbols-rounded text-white">
              self_improvement
            </span>
          </div>
          <div class="z-10 -ml-1 flex size-11 flex-col items-center justify-center rounded-xl border border-main bg-red">
            <span class="material-symbols-rounded text-white">wine_bar</span>
          </div>
          <div class="-ml-1 flex size-11 flex-col items-center justify-center rounded-xl border border-main bg-blue">
            <span class="material-symbols-rounded text-white">
              directions_run
            </span>
          </div>
        </div>
        <div class="flex flex-row items-start justify-between">
          <p class="mt-2 text-3xl text-main">
            <span class="text-pink">Find collaborations</span> that you may be
            interested to join
          </p>
          <span class="material-symbols-rounded text-[48px] text-red">
            arrow_circle_right
          </span>
        </div>
        <p class="mt-1.5 text-sm text-hint">
          Yoga practice, Running, Grabbing a coffee, and more
        </p>
      </Link>
      <div class="h-px w-full bg-main" />
      <button
        class="flex flex-col items-start justify-start py-4"
        onClick={() =>
          window.Telegram.WebApp.openTelegramLink('https://t.me/telegram')
        }
      >
        <div class="flex flex-row items-start justify-between text-start">
          <p class="mt-2 text-3xl text-main">
            <span class="text-green">Join community</span> to talk with founders
            and users. Discuss and solve problems together
          </p>
          <span class="material-symbols-rounded text-[48px] text-green">
            forum
          </span>
        </div>
      </button>
    </div>
  );
}

const ImagesLoader = () => {
  return (
    <div class="flex w-full flex-row items-center justify-start">
      <For each={[1, 2, 3] as any}>
        {(image, idx) => (
          <div
            class="-ml-1 size-11 rounded-xl border border-main bg-hint"
            classList={{
              'ml-0': idx() === 0,
              'z-20': idx() === 0,
              'z-10': idx() === 1,
            }}
          />
        )}
      </For>
    </div>
  );
};
