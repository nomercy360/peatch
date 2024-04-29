import { FormLayout } from '../layout';

export default function DescribeCollaboration(props: {
  description: string;
  setDescription: (description: string) => void;
  isPayable: boolean;
  setIsPayable: (isPayable: boolean) => void;
  title: string;
  setTitle: (title: string) => void;
}) {
  const resizer = (e: any) => {
    e.target.style.height = 'auto';
    e.target.style.height = e.target.scrollHeight + 'px';

    props.setDescription(e.target.value);
  };

  return (
    <FormLayout
      title="Describe collaboration"
      description="This will help people to understand it clearly"
    >
      <div class="mt-5 flex w-full flex-col items-center justify-start gap-3">
        <input
          class="h-10 w-full rounded-lg bg-peatch-bg px-2.5 text-black placeholder:text-gray"
          placeholder="Title"
          value={props.title}
          onInput={e => props.setTitle(e.currentTarget.value)}
        />
        <button
          class="flex h-10 w-full items-center justify-between"
          onClick={() => props.setIsPayable(!props.isPayable)}
        >
          <p class="text-sm text-black">Is it this opportunity payable?</p>
          <span
            class="size-6 rounded-lg border border-peatch-stroke"
            classList={{
              'bg-peatch-blue': props.isPayable,
              'bg-peatch-bg': !props.isPayable,
            }}
          ></span>
        </button>
        <div class="relative h-fit min-h-56 w-full rounded-lg bg-peatch-bg">
          <textarea
            class="size-full bg-transparent p-2.5 text-black placeholder:text-gray"
            placeholder="For example: 32 y.o. serial entrepreneur & product director with architecture, product design, marketing & tech development background. "
            value={props.description}
            onInput={e => resizer(e)}
          ></textarea>
          <span class="absolute bottom-2 right-2 text-sm text-gray">0/500</span>
        </div>
      </div>
    </FormLayout>
  );
}
