import { FormLayout } from '~/components/edit/layout';
import { useMainButton } from '~/hooks/useMainButton';
import { useNavigate } from '@solidjs/router';
import { createEffect, onCleanup } from 'solid-js';
import { editCollaboration, editCollaborationId, setEditCollaboration } from '~/store';
import SelectLocation from '~/components/edit/selectLocation';
import { createCollaboration, updateCollaboration } from '~/api';

export default function SelectBadges() {
  const mainButton = useMainButton();

  const navigate = useNavigate();

  const createCollab = async () => {
    const created = await createCollaboration(editCollaboration);
    navigate('/collaborations/' + created.id);
  };

  const editCollab = async () => {
    await updateCollaboration(editCollaborationId(), editCollaboration);
    navigate('/collaborations/' + editCollaborationId() + '?refetch=true');
  };

  const createOrEditCollab = async () => {
    if (editCollaborationId()) {
      await editCollab();
    } else {
      await createCollab();
    }
  };

  mainButton.onClick(createOrEditCollab);

  createEffect(() => {
    if (editCollaboration.country && editCollaboration.country_code) {
      mainButton.enable('Choose & Save');
    } else {
      mainButton.disable('Choose & Save');
    }
  });

  onCleanup(() => {
    mainButton.offClick(createOrEditCollab);
  });

  return (
    <FormLayout
      title="Any special location?"
      description="This will help us to recommend it to other people"
      screen={4}
      totalScreens={6}
    >
      <SelectLocation
        city={editCollaboration.city}
        setCity={c => setEditCollaboration('city', c)}
        country={editCollaboration.country}
        setCountry={c => setEditCollaboration('country', c)}
        countryCode={editCollaboration.country_code}
        setCountryCode={c => setEditCollaboration('country_code', c)}
      />
    </FormLayout>
  );
}
