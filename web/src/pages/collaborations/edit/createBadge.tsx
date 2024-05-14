import { FormLayout } from '~/components/edit/layout';
import { useMainButton } from '~/lib/useMainButton';
import { useNavigate, useParams, useSearchParams } from '@solidjs/router';
import { createEffect, onCleanup } from 'solid-js';
import { editCollaboration, setEditCollaboration } from '~/store';
import { postBadge } from '~/lib/api';
import { createStore } from 'solid-js/store';
import CreateBadge from '~/components/edit/createBadge';

export default function SelectBadges() {
  const mainButton = useMainButton();

  const [searchParams] = useSearchParams();

  const [createBadge, setCreateBadge] = createStore({
    text: searchParams.badge_name,
    color: 'EF5DA8',
    icon: '',
  });

  const idPath = useParams().id ? '/' + useParams().id : '';

  const publishBadge = async () => {
    if (createBadge.text && createBadge.color && createBadge.icon) {
      const { id } = await postBadge(
        createBadge.text,
        createBadge.color,
        createBadge.icon,
      );

      setEditCollaboration('badge_ids', [...editCollaboration.badge_ids, id]);
    }
  }

  const navigate = useNavigate();

  const onCreateBadgeButtonClick = async () => {
    await publishBadge();
    navigate(`/collaborations/edit${idPath}/badges?refetch=true`, {
      state: { from: '/collaborations/edit' + idPath },
    });
  };

  mainButton.onClick(onCreateBadgeButtonClick);

  createEffect(() => {
    if (createBadge.icon && createBadge.color && createBadge.text) {
      mainButton.enable('Create ' + createBadge.text);
    } else {
      mainButton.disable('Create ' + createBadge.text);
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
