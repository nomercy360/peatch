import { createEffect, createSignal, For, Show, Suspense } from 'solid-js';
import { Opportunity } from '../../../gen';

export function SelectOpportunity(props: {
  selected: number[] | number;
  setSelected: (selected: number[] | number) => void;
  opportunities: Opportunity[];
}) {
  const [filtered, setFiltered] = createSignal(props.opportunities);
  const [search, setSearch] = createSignal('');

  const onClick = (oppId: number) => {
    if (Array.isArray(props.selected)) {
      if (props.selected.includes(oppId)) {
        props.setSelected(props.selected.filter(b => b !== oppId));
      } else if (props.selected.length < 10) {
        props.setSelected([...props.selected, oppId]);
      }
    } else {
      props.setSelected(oppId);
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

  const includes = (oppId: number) => {
    if (Array.isArray(props.selected)) {
      return props.selected.includes(oppId);
    }
    return props.selected === oppId;
  };

  return (
    <>
      <div class="mt-5 flex h-10 w-full flex-row items-center justify-between rounded-lg bg-peatch-main px-2.5">
        <input
          class="h-10 w-full bg-transparent text-main placeholder:text-hint focus:outline-none"
          placeholder="Search collaboration opportunities"
          type="text"
          onInput={e => setSearch(e.currentTarget.value)}
          value={search()}
        />
        <Show when={search()}>
          <button
            class="flex h-full items-center justify-center px-2.5 text-sm text-hint"
            onClick={() => setSearch('')}
          >
            Clear
          </button>
        </Show>
      </div>
      <div class="flex h-11 w-full flex-row items-center justify-between">
        <div></div>
        <div class="flex h-11 items-center justify-center text-sm text-hint">
          {Array.isArray(props.selected)
            ? `${props.selected.length} / 10`
            : 'choose one'}
        </div>
      </div>
      <div class="flex w-full flex-row flex-wrap items-center justify-start gap-1">
        <Suspense fallback={<div>Loading...</div>}>
          <For each={filtered()}>
            {op => (
              <button
                onClick={() => onClick(op.id!)}
                class="flex h-[60px] w-full flex-row items-center justify-start gap-2.5 rounded-2xl border border-peatch-stroke px-2.5"
                style={{
                  'background-color': `${includes(op.id!) ? `#${op.color}` : 'var(--tg-theme-section-bg-color)'}`,
                  'border-color': `${includes(op.id!) ? `#${op.color}` : 'var(--tg-theme-section-bg-color)'}`,
                }}
              >
                <div class="flex size-10 items-center justify-center rounded-full bg-peatch-hint  ">
                  <span class="material-symbols-rounded text-main">
                    {String.fromCodePoint(parseInt(op.icon!, 16))}
                  </span>
                </div>

                <div
                  class="text-start text-main"
                >
                  <p class="text-sm font-semibold">{op.text}</p>
                  <p class="text-xs">{op.description}</p>
                </div>
              </button>
            )}
          </For>
        </Suspense>
      </div>
    </>
  );
}
