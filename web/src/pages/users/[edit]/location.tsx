import { FormLayout } from '~/components/edit/layout';
import { useMainButton } from '@tma.js/sdk-solid';
import { useNavigate } from '@solidjs/router';
import { createEffect, onCleanup } from 'solid-js';
import { editUser, setEditUser } from '~/store';
import SelectLocation from '~/components/edit/selectLocation';

export default function SelectBadges() {
  const mainButton = useMainButton();

  const navigate = useNavigate();

  const navigateToDescription = async () => {
    navigate('/users/edit/description');
  };

  mainButton()
    .setParams({ text: 'Next', isVisible: true, isEnabled: false })
    .on('click', navigateToDescription);

  createEffect(() => {
    if (editUser.country && editUser.country_code) {
      mainButton().enable();
    }
  });

  onCleanup(() => {
    mainButton().off('click', navigateToDescription);
  });

  return (
    <FormLayout
      title="Where do you live?"
      description="It will appears in your profile card, everyone will see it"
      screen={4}
      totalScreens={6}
    >
      <SelectLocation
        city={editUser.city}
        setCity={c => setEditUser('city', c)}
        country={editUser.country}
        setCountry={c => setEditUser('country', c)}
        countryCode={editUser.country_code}
        setCountryCode={c => setEditUser('country_code', c)}
      />
    </FormLayout>
  );
}
