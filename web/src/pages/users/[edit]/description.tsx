import { useMainButton } from '@tma.js/sdk-solid';
import { useNavigate } from '@solidjs/router';
import { createEffect, onCleanup } from 'solid-js';
import { editUser, setEditUser } from '~/store';
import TextArea from '~/components/TextArea';
import { FormLayout } from '~/components/edit/layout';

export default function Description() {
  const mainButton = useMainButton();

  const navigate = useNavigate();

  const navigateToImageUpload = async () => {
    navigate('/users/edit/image');
  };

  mainButton()
    .setParams({ text: 'Next', isVisible: true, isEnabled: false })
    .on('click', navigateToImageUpload);

  createEffect(() => {
    if (editUser.description) {
      mainButton().enable();
    }
  });

  onCleanup(() => {
    mainButton().off('click', navigateToImageUpload);
  });

  return (
    <FormLayout
      title="Introduce yourself"
      description="Tell others about your backround, achievments and goals"
      screen={3}
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
