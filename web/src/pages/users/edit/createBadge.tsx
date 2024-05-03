import { FormLayout } from '~/components/edit/layout';
import { useMainButton } from '~/hooks/useMainButton';
import { useNavigate, useSearchParams } from '@solidjs/router';
import { createEffect, onCleanup } from 'solid-js';
import { editUser, setEditUser } from '~/store';
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

      setEditUser('badge_ids', [...editUser.badge_ids, id]);
    }
  };

  const navigate = useNavigate();

  const onCreateBadgeButtonClick = async () => {
    await publishBadge();
    navigate('/users/edit/badges', {
      state: { from: '/users/edit' },
    });
  };

  mainButton
    .onClick(onCreateBadgeButtonClick);

  createEffect(() => {
    if (createBadge.icon && createBadge.color && createBadge.text) {
      mainButton.enable('Next');
    } else {
      mainButton.disable('Next');
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
