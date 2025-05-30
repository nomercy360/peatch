import { For, Show } from 'solid-js'
import icons from '../../assets/icons.json'
import { BadgeResponse } from '~/gen'

export default function CreateBadge(props: {
  createBadge: BadgeResponse
  setCreateBadge: any
}) {
  const colors = [
    'EF5DA8',
    'F9A826',
    '2D9CDB',
    '27AE60',
    '6D214F',
    'F2C94C',
    'F2994A',
  ]

  return (
    <>
      <div class="mt-5 flex w-full flex-row items-center justify-between gap-2.5">
        <For each={colors}>
          {color => (
            <button
              class="flex size-10 items-center justify-center rounded-full"
              style={{ 'background-color': `#${color}` }}
              onClick={() =>
                props.setCreateBadge({ ...props.createBadge, color: color })
              }
            >
              <Show when={props.createBadge.color === color}>
                <span class="material-symbols-rounded text-[20px] text-white">
                  check
                </span>
              </Show>
            </button>
          )}
        </For>
      </div>
      <div class="mt-5 grid w-full grid-cols-7 gap-2.5">
        <For each={icons}>
          {icon => (
            <button
              class="flex aspect-square items-center justify-center rounded-lg"
              style={{
                'background-color':
                  icon === props.createBadge.icon
                    ? `#${props.createBadge.color}`
                    : 'var(--tg-theme-secondary-bg-color)',
              }}
              onClick={() =>
                props.setCreateBadge({ ...props.createBadge, icon: icon })
              }
            >
              <span
                class="material-symbols-rounded text-[24px]"
                style={{
                  color: icon === props.createBadge.icon ? 'white' : '#B6B6B6',
                }}
              >
                {String.fromCharCode(parseInt(icon, 16))}
              </span>
            </button>
          )}
        </For>
      </div>
    </>
  )
}
