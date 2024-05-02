import { createResource, For } from 'solid-js';
import { A } from '@solidjs/router';
import { CDN_URL, fetchPreview } from '~/api';

export default function ProfilePublished() {
  const [previewImages, _] = createResource(async () => {
    const res = await fetchPreview();
    return res.map((image: any) => CDN_URL + '/' + image.avatar_url);
  });

  return (
    <div class="flex h-screen w-full flex-col items-center justify-between bg-secondary p-5 text-center">
      <img
        src="/confetti.png"
        alt="Confetti"
        class="absolute inset-x-0 top-0 mx-auto w-full"
      />
      <div class="flex flex-col items-center justify-start">
        <span class="material-symbols-rounded text-peatch-green text-[60px] text-main">
          check_circle
        </span>
        <p class="text-3xl text-main">Profile published</p>
        <p class="mt-2 text-2xl text-secondary">
          Now you can find people, create and join collaborations. Have fun!
        </p>
      </div>
      <div class="flex flex-col items-center justify-center">
        <div class="flex w-full flex-row items-center justify-center">
          <For each={previewImages()}>
            {(image, idx) => (
              <img
                src={image}
                alt="User Avatar"
                class="-ml-1 size-11 rounded-lg border-2 border-main object-cover object-center"
                classList={{
                  'ml-0': idx() === 0,
                  'z-20': idx() === 0,
                  'z-10': idx() === 1,
                }}
              />
            )}
          </For>
        </div>
        <p class="mt-4 max-w-xs text-lg text-secondary">
          There are 12 people you might be interested to collaborate with
        </p>
        <A
          class="mt-2 flex h-12 w-full items-center justify-center text-sm font-medium text-link"
          href="/users"
        >
          Show them
        </A>
      </div>
    </div>
  );
}
