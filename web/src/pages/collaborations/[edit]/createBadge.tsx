import { FormLayout } from '~/components/edit/layout';
import { useMainButton } from '~/hooks/useMainButton';
import { useNavigate, useSearchParams } from '@solidjs/router';
import { createEffect, onCleanup } from 'solid-js';
import { editCollaboration, setEditCollaboration } from '~/store';
import { postBadge } from '~/api';
import { createStore } from 'solid-js/store';
import CreateBadge from '~/components/edit/createBadge';

export default function SelectBadges() {
  const mainButton = useMainButton();

  const [searchParams, _] = useSearchParams();

  const [createBadge, setCreateBadge] = createStore({
    text: searchParams.badge_name,
    color: 'EF5DA8',
    icon: '',
  });

  const publishBadge = async () => {
    if (createBadge.text && createBadge.color && createBadge.icon) {
      const { id } = await postBadge(
        createBadge.text,
        createBadge.color,
        createBadge.icon,
      );

      setEditCollaboration('badge_ids', [...editCollaboration.badge_ids, id]);
    }
  };

  const navigate = useNavigate();

  const onCreateBadgeButtonClick = async () => {
    await publishBadge();
    navigate('/collaboration/edit/badges', { state: { back: true } });
  };

  mainButton
    .setParams({ text: 'Next', isVisible: true, isEnabled: false })
    .onClick(onCreateBadgeButtonClick);

  createEffect(() => {
    if (createBadge.icon && createBadge.color && createBadge.text) {
      mainButton.enable();
    }
  });

  onCleanup(() => {
    mainButton.offClick(onCreateBadgeButtonClick);
  });

  return (
    <FormLayout
      title={`Creating ${createBadge.text}`}
      description="This will help us to recommend you to other people"
      screen={2}
      totalScreens={6}
    >
      <CreateBadge createBadge={createBadge} setCreateBadge={setCreateBadge} />
    </FormLayout>
  );
}
