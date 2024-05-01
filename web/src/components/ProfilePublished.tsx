import { For } from 'solid-js';

export default function ProfilePublished() {
  const images = ['/thumb.png', '/thumb.png', '/thumb.png'];

  return (
    <div class="h-screen p-5 bg-peatch-bg w-full text-center flex flex-col items-center justify-between">
      <img src="/confetti.png" alt="Confetti" class="w-full mx-auto absolute top-0 left-0 right-0" />
      <div class="flex flex-col items-center justify-start">
        <span class="material-symbols-rounded text-[60px] text-peatch-green">check_circle</span>
        <p class="text-3xl">Profile published</p>
        <p class="text-2xl mt-2">Now you can find people, create and join collaborations. Have fun!</p>
      </div>
      <div class="flex flex-col items-center justify-center">
        <div class="flex w-full flex-row items-center justify-center">
          <For each={images}>
            {(image, idx) => (
              <img
                src={image}
                alt="User Avatar"
                class="-ml-1 size-11 rounded-xl border-2 border-white object-cover object-center"
                classList={{
                  'ml-0': idx() === 0,
                  'z-20': idx() === 0,
                  'z-10': idx() === 1,
                }}
              />
            )}
          </For>
        </div>
        <p class="mt-4 text-lg max-w-xs">
          There are 12 people you might be interested to collaborate with
        </p>
        <a class="mt-2 text-sm text-peatch-blue h-12 w-full flex items-center justify-center"
           href="/users">
          Show them</a>
      </div>
    </div>
  );
}