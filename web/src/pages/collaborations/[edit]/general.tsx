import { FormLayout } from '~/components/edit/layout';
import { editCollaboration, editUser, setEditCollaboration } from '~/store';
import { useMainButton } from '~/hooks/useMainButton';
import { createEffect, onCleanup } from 'solid-js';
import { useNavigate } from '@solidjs/router';
import TextArea from '~/components/TextArea';

export default function GeneralInfo() {
  const mainButton = useMainButton();

  const navigate = useNavigate();

  const navigateNext = () => {
    navigate('/collaborations/edit/badges', { state: { back: true } });
  };

  mainButton.onClick(navigateNext);

  createEffect(() => {
    if (editCollaboration.title && editCollaboration.description) {
      mainButton.setParams({
        isEnabled: true,
        isVisible: true,
        text: 'Next',
      });
    } else {
      mainButton.setParams({
        isEnabled: false,
        isVisible: true,
        text: 'Next',
      });
    }
  });

  onCleanup(() => {
    mainButton.offClick(navigateNext);
  });

  return (
    <FormLayout
      title="Describe collaboration"
      description="This will help people to understand it clearly"
      screen={1}
      totalScreens={4}
    >
      <div class="mt-5 flex w-full flex-col items-center justify-start gap-3">
        <input
          class="h-10 w-full rounded-lg bg-peatch-bg px-2.5 text-main placeholder:text-hint"
          placeholder="Title"
          value={editCollaboration.title}
          onInput={e => setEditCollaboration('title', e.currentTarget.value)}
        />
        <button
          class="flex h-10 w-full items-center justify-between"
          onClick={() =>
            setEditCollaboration('is_payable', !editCollaboration.is_payable)
          }
        >
          <p class="text-sm text-main">Is it this opportunity payable?</p>
          <span
            class="size-6 rounded-lg border border-peatch-stroke"
            classList={{
              'bg-peatch-blue': !editCollaboration.is_payable,
              'bg-peatch-bg': editCollaboration.is_payable,
            }}
          ></span>
        </button>
        <TextArea
          value={editCollaboration.description}
          setValue={d => setEditCollaboration('description', d)}
          placeholder="For example: 32 y.o. serial entrepreneur & product director with architecture, product design, marketing & tech development background."
        ></TextArea>
      </div>
    </FormLayout>
  );
}
