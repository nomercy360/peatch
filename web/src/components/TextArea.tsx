import { createSignal, Match, Switch } from 'solid-js';

export default function TextArea(props: {
  value: string;
  setValue: (value: string) => void;
  placeholder: string;
}) {
  const [count, setCount] = createSignal(0);

  const resizer = (e: any) => {
    // e.target.style.height = 'auto';
    // e.target.style.height = e.target.scrollHeight + 'px';

    props.setValue(e.target.value);

    const count = e.target.value.length;
    setCount(count);
  };

  return (
    <div class="relative mt-5 h-80 w-full rounded-lg bg-main">
      <textarea
        class="size-full resize-none bg-transparent p-2.5 text-main placeholder:text-hint"
        placeholder={props.placeholder}
        value={props.value}
        onInput={e => resizer(e)}
        autocomplete="off"
        autocapitalize="off"
        spellcheck={false}
      ></textarea>
      <Switch>
        <Match when={count() > 0}>
          <div class="absolute bottom-2 left-2 text-sm text-hint">
            {count()} / 500
          </div>
        </Match>
        <Match when={count() === 0}>
          <div class="absolute bottom-2 left-2 text-sm text-hint">
            max 500 characters
          </div>
        </Match>
      </Switch>
    </div>
  );
}
