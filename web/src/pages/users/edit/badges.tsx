import { SelectBadge } from '~/components/edit/selectBadge';
import { FormLayout } from '~/components/edit/layout';
import { useMainButton } from '~/lib/useMainButton';
import { useNavigate, useSearchParams } from '@solidjs/router';
import { createEffect, createResource, createSignal, onCleanup } from 'solid-js';
import { editUser, setEditUser } from '~/store';
import { fetchBadges } from '~/lib/api';
import { Badge } from '~/gen/types';

export default function SelectBadges() {
  const mainButton = useMainButton();

  const [badgeSearch, setBadgeSearch] = createSignal('');

  const navigate = useNavigate();
  const [searchParams] = useSearchParams();

  const navigateNext = () => {
    navigate('/users/edit/interests');
  };

  const navigateCreateBadge = () => {
    navigate('/users/edit/create-badge?badge_name=' + badgeSearch(), {
      state: { back: true },
    });
  };

  const [badges, { refetch }] = createResource(() =>
    fetchBadges().then(badges => {
      const selected = editUser.badge_ids;
      return [
        ...selected.map(id => badges.find((b: Badge) => b.id === id)),
        ...badges.filter((b: Badge) => !selected.includes(b.id!)),
      ];
    }),
  );

  createEffect(async () => {
    if (searchParams.refetch) await refetch();
  });

  mainButton.onClick(navigateNext);

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
        badges={badges()!}
      />
    </FormLayout>
  );
}
