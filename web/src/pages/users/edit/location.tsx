import { FormLayout } from '~/components/edit/layout';
import { useMainButton } from '~/hooks/useMainButton';
import { createEffect, onCleanup } from 'solid-js';
import { editUser, setEditUser } from '~/store';
import SelectLocation from '~/components/edit/selectLocation';
import { useNavigate } from '@solidjs/router';

export default function SelectBadges() {
  const mainButton = useMainButton();

  const navigate = useNavigate();

  const navigateToDescription = async () => {
    navigate('/users/edit/description', { state: { back: true } });
  };

  mainButton
    .onClick(navigateToDescription);

  createEffect(() => {
    if (editUser.country && editUser.country_code) {
      mainButton.enable('Next');
    } else {
      mainButton.disable('Next');
    }
  });

  onCleanup(() => {
    mainButton.offClick(navigateToDescription);
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
