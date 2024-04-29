import { FormLayout } from '../layout';

export default function FillDescription(props: {
  description: string;
  setDescription: (description: string) => void;
}) {
  const resizer = (e: any) => {
    e.target.style.height = 'auto';
    e.target.style.height = e.target.scrollHeight + 'px';

    props.setDescription(e.target.value);
  };

  return (
    <FormLayout
      title="Introduce yourself"
      description="Tell others about your backround, achievments and goals"
    >
      <div class="relative mt-5 h-fit min-h-56 w-full rounded-lg bg-peatch-bg">
        <textarea
          class="size-full bg-transparent p-2.5 text-black placeholder:text-gray"
          placeholder="For example: 32 y.o. serial entrepreneur & product director with architecture, product design, marketing & tech development background. "
          value={props.description}
          onInput={e => resizer(e)}
        ></textarea>
        <span class="absolute bottom-2 right-2 text-sm text-gray">0/500</span>
      </div>
    </FormLayout>
  );
}
