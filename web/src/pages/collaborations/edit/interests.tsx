import { FormLayout } from '~/components/edit/layout';
import { useMainButton } from '~/lib/useMainButton';
import { useNavigate, useParams } from '@solidjs/router';
import { createEffect, createResource, onCleanup } from 'solid-js';
import { editCollaboration, setEditCollaboration } from '~/store';
import { fetchOpportunities } from '~/lib/api';
import { SelectOpportunity } from '~/components/edit/selectOpp';

export default function SelectOpportunities() {
  const mainButton = useMainButton();

  const navigate = useNavigate();
  const params = useParams();
  const idPath = params.id ? '/' + params.id : '';

  const navigateNext = () => {
    navigate(`/collaborations/edit${idPath}/location`, {
      state: { back: true },
    });
  };

  const [opportunities] = createResource(() => fetchOpportunities());

  mainButton.onClick(navigateNext);

  createEffect(() => {
    if (
      editCollaboration.opportunity_id &&
      editCollaboration.opportunity_id > 0
    ) {
      mainButton.enable('Next');
    } else {
      mainButton.disable('Next');
    }
  });

  onCleanup(() => {
    mainButton.offClick(navigateNext);
  });

  return (
    <FormLayout
      title="Select a theme"
      description="This will help us to recommend it to other people"
      screen={3}
      totalScreens={6}
    >
      <SelectOpportunity
        selected={editCollaboration.opportunity_id}
        setSelected={b => setEditCollaboration('opportunity_id', b as any)}
        opportunities={opportunities()}
      />
    </FormLayout>
  );
}
