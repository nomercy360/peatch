import { FormLayout } from '~/components/edit/layout';
import { useMainButton } from '~/hooks/useMainButton';
import { useNavigate } from '@solidjs/router';
import { createEffect, onCleanup } from 'solid-js';
import { editCollaboration, setEditCollaboration } from '~/store';
import SelectLocation from '~/components/edit/selectLocation';
import { createCollaboration } from '~/api';

export default function SelectBadges() {
  const mainButton = useMainButton();

  const navigate = useNavigate();

  const createCollab = async () => {
    const created = await createCollaboration(editCollaboration);
    // reload the page to get the new collaboration
    location.reload();
    navigate('/collaborations/' + created.id);
  };

  mainButton.onClick( createCollab);

  createEffect(() => {
    if (editCollaboration.country && editCollaboration.country_code) {
      mainButton.setParams({
        text: 'Choose & Save',
        isVisible: true,
        isEnabled: true,
      });
    } else {
      mainButton.setParams({
        text: 'Choose & Save',
        isVisible: true,
        isEnabled: false,
      });
    }
  });

  onCleanup(() => {
    mainButton.offClick( createCollab);
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
