import { SelectBadge } from '~/components/edit/selectBadge';
import { FormLayout } from '~/components/edit/layout';
import { useMainButton } from '~/hooks/useMainButton';
import { useNavigate } from '@solidjs/router';
import { createEffect, createSignal, onCleanup } from 'solid-js';
import { editUser, setEditUser } from '~/store';
import { fetchBadges } from '~/api';
import { createQuery } from '@tanstack/solid-query';

export default function SelectBadges() {
  const mainButton = useMainButton();

  const [badgeSearch, setBadgeSearch] = createSignal('');

  const navigate = useNavigate();

  const navigateNext = () => {
    navigate('/users/edit/interests', { state: { back: true } });
  };

  const navigateCreateBadge = () => {
    navigate('/users/edit/create-badge?badge_name=' + badgeSearch(), {
      state: { back: true },
    });
  };

  const fetchBadgeQuery = createQuery(() => ({
    queryKey: ['badges'],
    queryFn: () => fetchBadges(),
  }));

  mainButton
    .onClick(navigateNext);

  createEffect(() => {
    if (editUser.badge_ids.length) {
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
      title="What describes you?"
      description="This will help us to recommend you to other people"
      screen={2}
      totalScreens={6}
    >
      <SelectBadge
        selected={editUser.badge_ids}
        setSelected={b => setEditUser('badge_ids', b)}
        onCreateBadgeButtonClick={navigateCreateBadge}
        search={badgeSearch()}
        setSearch={setBadgeSearch}
        badges={fetchBadgeQuery.data}
      />
    </FormLayout>
  );
}
