import { FormLayout } from '~/components/edit/layout';
import { useMainButton } from '~/hooks/useMainButton';
import { useNavigate } from '@solidjs/router';
import { createEffect, onCleanup } from 'solid-js';
import { editCollaboration, setEditCollaboration } from '~/store';
import { fetchOpportunities } from '~/api';
import { createQuery } from '@tanstack/solid-query';
import { SelectOpportunity } from '~/components/edit/selectOpp';

export default function SelectOpportunities() {
  const mainButton = useMainButton();

  const navigate = useNavigate();

  const navigateNext = () => {
    navigate('/collaborations/edit/location', { state: { back: true } });
  };

  const fetchOpportunityQuery = createQuery(() => ({
    queryKey: ['opportunities'],
    queryFn: () => fetchOpportunities(),
  }));

  mainButton
    .setParams({ text: 'Next', isVisible: true, isEnabled: false })
    .onClick(navigateNext);

  createEffect(() => {
    if (editCollaboration.opportunity_id) {
      mainButton.enable();
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
        opportunities={fetchOpportunityQuery.data}
      />
    </FormLayout>
  );
}
