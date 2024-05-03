import {
  createEffect,
  createSignal,
  For,
  Match,
  Show,
  Suspense,
  Switch,
} from 'solid-js';
import { Badge } from '../../../gen';

export function SelectBadge(props: {
  selected: number[];
  setSelected: (selected: number[]) => void;
  search: string;
  setSearch: (search: string) => void;
  badges: Badge[];
  onCreateBadgeButtonClick: () => void;
}) {
  const [filteredBadges, setFilteredBadges] = createSignal(props.badges);

  const onBadgeClick = (badgeId: number) => {
    if (props.selected.includes(badgeId!)) {
      props.setSelected(props.selected.filter(b => b !== badgeId));
    } else if (props.selected.length < 10) {
      props.setSelected([...props.selected, badgeId!]);
    }
  };

  createEffect(() => {
    if (props.badges && props.badges.length > 0) {
      setFilteredBadges(
        props.badges.filter(badge =>
          badge.text?.toLowerCase().includes(props.search.toLowerCase()),
        ),
      );
    }
  });

  return (
    <>
      <div class="mt-5 flex h-10 w-full flex-row items-center justify-between rounded-lg bg-main px-2.5">
        <input
          autofocus
          autocomplete="off"
          autocorrect="off"
          autocapitalize="off"
          spellcheck={false}
          class="h-10 w-full bg-transparent text-main placeholder:text-hint focus:outline-none"
          placeholder="Search for a badge"
          type="text"
          onInput={e => props.setSearch(e.currentTarget.value)}
          value={props.search}
        />
        <Show when={props.search}>
          <button
            class="flex h-10 items-center justify-center px-2.5 text-sm text-hint"
            onClick={() => props.setSearch('')}
          >
            Clear
          </button>
        </Show>
      </div>
      <div class="flex h-11 w-full flex-row items-center justify-between">
        <Switch>
          <Match when={filteredBadges()?.length > 0}>
            <div></div>
            <div class="flex h-11 items-center justify-center text-sm text-hint">
              {props.selected.length} / 10
            </div>
          </Match>
          <Match when={filteredBadges()?.length === 0}>
            <button
              class="size-full text-start text-sm text-hint"
              onClick={() =>
                props.selected.length < 10 && props.onCreateBadgeButtonClick()
              }
            >
              Can’t find such thing. <span class="text-accent">Create it</span>
            </button>
            <p class="text-nowrap text-sm text-hint">
              {props.selected.length} of 10
            </p>
          </Match>
        </Switch>
      </div>
      <div class="flex min-h-fit w-full flex-row flex-wrap items-center justify-start gap-1">
        <Suspense fallback={<div>Loading...</div>}>
          <For each={filteredBadges()}>
            {badge => (
              <button
                onClick={() => onBadgeClick(badge.id!)}
                class="flex h-10 flex-row items-center justify-center gap-[5px] rounded-2xl border px-2.5"
                style={{
                  'background-color': `${props.selected.includes(badge.id!) ? `#${badge.color}` : 'var(--tg-theme-secondary-bg-color)'}`,
                  'border-color': `${props.selected.includes(badge.id!) ? `#${badge.color}` : 'var(--tg-theme-secondary-bg-color)'}`,
                }}
              >
                <span
                  class="material-symbols-rounded"
                  style={{
                    color: `${props.selected.includes(badge.id!) ? 'white' : `#${badge.color}`}`,
                  }}
                >
                  {String.fromCodePoint(parseInt(badge.icon!, 16))}
                </span>
                <p class="text-sm font-semibold"
                   classList={{
                     'text-white': props.selected.includes(badge.id!),
                     'text-main': !props.selected.includes(badge.id!),
                   }}
                >{badge.text}</p>
              </button>
            )}
          </For>
        </Suspense>
      </div>
    </>
  );
}
