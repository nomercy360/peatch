import { createEffect, createSignal, For, Match, Show, Suspense, Switch } from 'solid-js';
import { FormLayout } from '../../pages/users/edit';
import { Badge } from '../../../gen';

export function SelectBadge(props: {
  selected: number[];
  setSelected: (selected: number[]) => void;
  setCreateBadgeModal: (open: boolean) => void;
  search: string;
  setSearch: (search: string) => void;
  badges: Badge[];
}) {
  const [filteredBadges, setFilteredBadges] = createSignal(props.badges);

  const onBadgeClick = (badgeId: number) => {
    if (props.selected.includes(badgeId!)) {
      props.setSelected(props.selected.filter(b => b !== badgeId));
    } else {
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
    <FormLayout
      title="What describes you?"
      description="This will help us to recommend you to other people"
    >
      <div class="mt-5 flex h-10 w-full flex-row items-center justify-between rounded-lg bg-peatch-bg px-2.5">
        <input
          class="w-full bg-transparent text-black placeholder:text-gray focus:outline-none"
          placeholder="Search for a badge"
          type="text"
          onInput={e => props.setSearch(e.currentTarget.value)}
          value={props.search}
        />
        <Show when={props.search}>
          <button
            class="flex h-full items-center justify-center px-2.5 text-sm text-peatch-dark-gray"
            onClick={() => props.setSearch('')}
          >
            Clear
          </button>
        </Show>
      </div>
      <div class="flex h-11 w-full flex-row items-center justify-between">
        <Switch>
          <Match when={filteredBadges().length > 0}>
            <div></div>
            <p class="text-sm text-gray">{props.selected.length} / 10</p>
          </Match>
          <Match when={filteredBadges().length === 0}>
            <button
              class="size-full border text-start text-sm"
              onClick={() => props.setCreateBadgeModal(true)}
            >
              Can’t find such thing.{' '}
              <span class="text-peatch-blue">Create it</span>
            </button>
            <p class="text-nowrap text-sm text-gray">
              {filteredBadges().length} of 10
            </p>
          </Match>
        </Switch>
      </div>
      <div class="flex w-full flex-row flex-wrap items-center justify-start gap-1">
        <Suspense fallback={<div>Loading...</div>}>
          <For each={filteredBadges()}>
            {badge => (
              <button
                onClick={() => onBadgeClick(badge.id!)}
                class="flex h-10 flex-row items-center justify-center gap-[5px] rounded-2xl border border-peatch-stroke px-2.5"
                style={{
                  'background-color': `${props.selected.includes(badge.id!) ? badge.color : 'white'}`,
                }}
              >
                <span
                  class="material-symbols-rounded"
                  style={{
                    color: `${props.selected.includes(badge.id!) ? 'white' : badge.color}`,
                  }}
                >
                  {String.fromCodePoint(parseInt(badge.icon!, 16))}
                </span>
                <p
                  class="text-sm font-semibold"
                  classList={{
                    'text-white': props.selected.includes(badge.id!),
                  }}
                >
                  {badge.text}
                </p>
              </button>
            )}
          </For>
        </Suspense>
      </div>
    </FormLayout>
  );
}
