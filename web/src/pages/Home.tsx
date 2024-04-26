import { For } from 'solid-js';

export default function Home() {
  const images = ['/thumb.png', '/thumb.png', '/thumb.png'];
  const name = 'John Doe';
  return (
    <div class="flex flex-col px-4">
      <div class="flex flex-row items-center justify-between py-4">
        <p class="text-3xl">Bonsoir, {name}!</p>
        <img src="/thumb.png" alt="User Avatar" class="size-10 rounded-xl" />
      </div>
      <div class="h-px w-full bg-peatch-stroke"></div>
      <a class="flex flex-col items-start justify-start py-4" href="/profiles">
        <div class="flex w-full flex-row items-center justify-start">
          <For each={images}>
            {(image, idx) => (
              <img
                src={image}
                alt="User Avatar"
                class="-ml-1 size-11 rounded-xl border-2 border-white"
                classList={{
                  'ml-0': idx() === 0,
                  'z-20': idx() === 0,
                  'z-10': idx() === 1,
                }}
              />
            )}
          </For>
        </div>
        <div class="flex flex-row items-center justify-between">
          <p class="mt-2 text-3xl">
            <span class="text-pink">Explore people</span> you may like to
            collaborate
          </p>
          <span class="material-symbols-rounded text-pink text-[48px]">
            maps_ugc
          </span>
        </div>
        <p class="mt-1.5 text-sm text-gray">
          Figma Wizards, Consultants, Founders, and more
        </p>
      </a>
      <div class="h-px w-full bg-peatch-stroke"></div>
      <div class="flex flex-col items-start justify-start py-4">
        <div class="flex w-full flex-row items-center justify-start">
          <div
            class="bg-orange z-20 flex size-11 flex-col items-center justify-center rounded-2xl border-2 border-white">
            <span class="material-symbols-rounded text-white">
              self_improvement
            </span>
          </div>
          <div
            class="bg-red z-10 -ml-1 flex size-11 flex-col items-center justify-center rounded-2xl border-2 border-white">
            <span class="material-symbols-rounded text-white">wine_bar</span>
          </div>
          <div
            class="bg-blue -ml-1 flex size-11 flex-col items-center justify-center rounded-2xl border-2 border-white">
            <span class="material-symbols-rounded text-white">
              directions_run
            </span>
          </div>
        </div>
        <div class="flex flex-row items-start justify-between">
          <p class="mt-2 text-3xl">
            <span class="text-pink">Find collaborations</span> that you may be
            interested to join
          </p>
          <span class="material-symbols-rounded text-red text-[48px]">
            arrow_circle_right
          </span>
        </div>
        <p class="mt-1.5 text-sm text-gray">
          Yoga practice, Running, Grabbing a coffee, and more
        </p>
      </div>
      <div class="h-px w-full bg-peatch-stroke"></div>
      <div class="flex flex-col items-start justify-start py-4">
        <div class="flex flex-row items-start justify-between">
          <p class="mt-2 text-3xl">
            <span class="text-green">Join community</span> to talk with founders
            and users. Discuss and solve problems together
          </p>
          <span class="material-symbols-rounded text-green text-[48px]">
            forum
          </span>
        </div>
      </div>
    </div>
  );
}
