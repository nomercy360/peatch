import { FormLayout } from '../../pages/users/edit';

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
      <textarea
        class="mt-5 h-fit min-h-40 w-full rounded-lg bg-peatch-bg p-2.5 text-black placeholder:text-gray"
        placeholder="For example: 32 y.o. serial entrepreneur & product director with architecture, product design, marketing & tech development background. "
        value={props.description}
        onInput={e => resizer(e)}
      ></textarea>
    </FormLayout>
  );
}
