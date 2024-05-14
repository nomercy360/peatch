import { FormLayout } from '~/components/edit/layout';
import { useMainButton } from '~/lib/useMainButton';
import { useNavigate } from '@solidjs/router';
import { createEffect, createResource, onCleanup } from 'solid-js';
import { editUser, setEditUser } from '~/store';
import { fetchOpportunities } from '~/lib/api';
import { SelectOpportunity } from '~/components/edit/selectOpp';

export default function SelectOpportunities() {
  const mainButton = useMainButton();

  const navigate = useNavigate();

  const navigateNext = () => {
    navigate('/users/edit/location', { state: { back: true } });
  };

  const [opportunities] = createResource(() => fetchOpportunities());

  mainButton.onClick(navigateNext);

  createEffect(() => {
    if (editUser.opportunity_ids.length) {
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
      title="What are you open for?"
      description="This will help us to recommend you to other people"
      screen={3}
      totalScreens={6}
    >
      <SelectOpportunity
        selected={editUser.opportunity_ids}
        setSelected={b => setEditUser('opportunity_ids', b as any)}
        opportunities={opportunities()!}
      />
    </FormLayout>
  );
}
