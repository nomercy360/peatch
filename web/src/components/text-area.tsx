import { createEffect, createSignal, Match, Switch } from 'solid-js'
import { useTranslations } from '~/lib/locale-context'

export default function TextArea(props: {
  value: string
  setValue: (value: string) => void
  placeholder: string
}) {
  const [count, setCount] = createSignal(0)
  const maxLength = 500

  const { t } = useTranslations()

  createEffect(() => {
    setCount(props.value.length)
  })

  return (
    <div class="relative mt-5 h-80 w-full rounded-lg bg-secondary pb-6">
      <textarea
        class="text-main placeholder:text-hint size-full h-full w-full resize-none bg-transparent p-2.5 focus:outline-none"
        placeholder={props.placeholder}
        value={props.value}
        onInput={e => props.setValue((e.target as HTMLTextAreaElement).value)}
        autocomplete="off"
        autocapitalize="off"
        spellcheck={false}
        maxLength={maxLength}
      />
      <Switch>
        <Match when={count() > 0}>
          <div class="text-hint absolute bottom-2 left-2 text-sm">
            {count()} / {maxLength}
          </div>
        </Match>
        <Match when={count() === 0}>
          <div class="text-hint absolute bottom-2 left-2 text-sm">
            {t('common.textarea.maxLength', { maxLength })}
          </div>
        </Match>
      </Switch>
    </div>
  )
}
