import { useMainButton } from '~/hooks/useMainButton';
import { useNavigate } from '@solidjs/router';
import { createEffect, createSignal, onCleanup } from 'solid-js';
import { editUser, setEditUser } from '~/store';
import TextArea from '~/components/TextArea';
import { FormLayout } from '~/components/edit/layout';

export default function Description() {
  const mainButton = useMainButton();

  const navigate = useNavigate();

  const navigateToImageUpload = async () => {
    navigate('/users/edit/image', { state: { back: true } });
  };

  mainButton.onClick(navigateToImageUpload);

  createEffect(() => {
    if (editUser.description) {
      mainButton.enable('Next');
    } else {
      mainButton.disable('Next');
    }
  });

  onCleanup(() => {
    mainButton.offClick(navigateToImageUpload);
  });

  return (
    <FormLayout
      title="Introduce yourself"
      description="Tell others about your backround, achievments and goals"
      screen={5}
      totalScreens={6}
    >
      <TextArea
        value={editUser.description}
        setValue={d => setEditUser('description', d)}
        placeholder="For example: 32 y.o. serial entrepreneur & product director with architecture, product design, marketing & tech development background. "
      />
    </FormLayout>
  );
}
