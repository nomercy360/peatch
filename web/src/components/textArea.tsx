export default function TextArea(props: { value: string; setValue: (value: string) => void; placeholder: string }) {
  const resizer = (e: any) => {
    e.target.style.height = 'auto';
    e.target.style.height = e.target.scrollHeight + 'px';

    props.setValue(e.target.value);
  };

  return (
    <div class="relative mt-5 h-fit min-h-56 w-full rounded-lg bg-peatch-bg">
        <textarea
          class="size-full bg-transparent p-2.5 text-black placeholder:text-gray"
          placeholder={props.placeholder}
          value={props.value}
          onInput={e => resizer(e)}
        ></textarea>
      <span class="absolute bottom-2 right-2 text-sm text-gray">0/500</span>
    </div>
  );
}