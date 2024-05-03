import { SelectBadge } from '~/components/edit/selectBadge';
import { FormLayout } from '~/components/edit/layout';
import { useMainButton } from '~/hooks/useMainButton';
import { useNavigate } from '@solidjs/router';
import { createEffect, createSignal, onCleanup } from 'solid-js';
import { editCollaboration, setEditCollaboration } from '~/store';
import { fetchBadges } from '~/api';
import { createQuery } from '@tanstack/solid-query';

export default function SelectBadges() {
  const mainButton = useMainButton();

  const [badgeSearch, setBadgeSearch] = createSignal('');

  const navigate = useNavigate();

  const navigateNext = () => {
    navigate('/collaborations/edit/interests', { state: { back: true } });
  };

  const navigateCreateBadge = () => {
    navigate('/collaborations/edit/create-badge?badge_name=' + badgeSearch(), {
      state: { back: true },
    });
  };

  const fetchBadgeQuery = createQuery(() => ({
    queryKey: ['badges'],
    queryFn: () => fetchBadges(),
  }));

  mainButton.onClick(navigateNext);

  createEffect(() => {
    if (editCollaboration.badge_ids.length > 0) {
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
      title="Who are you looking for?"
      description="This will help us to recommend it to other people"
      screen={2}
      totalScreens={6}
    >
      <SelectBadge
        selected={editCollaboration.badge_ids}
        setSelected={b => setEditCollaboration('badge_ids', b)}
        onCreateBadgeButtonClick={navigateCreateBadge}
        search={badgeSearch()}
        setSearch={setBadgeSearch}
        badges={fetchBadgeQuery.data}
      />
    </FormLayout>
  );
}
