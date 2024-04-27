import { createEffect, createSignal, For, Show, Suspense } from 'solid-js';
import { FormLayout } from '../../pages/users/edit';
import { Opportunity } from '../../../gen';

export function SelectOpportunity(props: {
  selected: number[];
  setSelected: (selected: number[]) => void;
  opportunities: Opportunity[];
}) {
  const [filtered, setFiltered] = createSignal(props.opportunities);
  const [search, setSearch] = createSignal('');

  const onBadgeClick = (badgeId: number) => {
    if (props.selected.includes(badgeId!)) {
      props.setSelected(props.selected.filter(b => b !== badgeId));
    } else {
      props.setSelected([...props.selected, badgeId!]);
    }
  };

  createEffect(() => {
    if (props.opportunities && props.opportunities.length > 0) {
      setFiltered(
        props.opportunities.filter(
          op =>
            op.text?.toLowerCase().includes(search().toLowerCase()) ||
            op.description?.toLowerCase().includes(search().toLowerCase()),
        ),
      );
    }
  });

  return (
    <FormLayout
      title="What are you open for?"
      description="This will help us to recommend you to other people"
    >
      <div class="mt-5 flex h-10 w-full flex-row items-center justify-between rounded-lg bg-peatch-bg px-2.5">
        <input
          class="w-full bg-transparent text-black placeholder:text-gray focus:outline-none"
          placeholder="Search collaboration opportunities"
          type="text"
          onInput={e => setSearch(e.currentTarget.value)}
          value={search()}
        />
        <Show when={search()}>
          <button
            class="flex h-full items-center justify-center px-2.5 text-sm text-peatch-dark-gray"
            onClick={() => setSearch('')}
          >
            Clear
          </button>
        </Show>
      </div>
      <div class="flex h-11 w-full flex-row items-center justify-between">
        <div></div>
        <p class="text-sm text-gray">{props.selected.length} / 10</p>
      </div>
      <div class="flex w-full flex-row flex-wrap items-center justify-start gap-1">
        <Suspense fallback={<div>Loading...</div>}>
          <For each={filtered()}>
            {op => (
              <button
                onClick={() => onBadgeClick(op.id!)}
                class="flex h-[60px] w-full flex-row items-center justify-start gap-2.5 rounded-2xl border border-peatch-stroke px-2.5"
                style={{
                  'background-color': `${props.selected.includes(op.id!) ? op.color : '#F8F8F8'}`,
                }}
              >
                <div class="flex size-10 items-center justify-center rounded-full bg-white">
                  <span class="material-symbols-rounded text-black">
                    {String.fromCodePoint(parseInt(op.icon!, 16))}
                  </span>
                </div>

                <div
                  class="text-start"
                  classList={{
                    'text-white': props.selected.includes(op.id!),
                    'text-peatch-gray': !props.selected.includes(op.id!),
                  }}
                >
                  <p class="text-sm font-semibold">{op.text}</p>
                  <p class="text-xs">{op.description}</p>
                </div>
              </button>
            )}
          </For>
        </Suspense>
      </div>
    </FormLayout>
  );
}
