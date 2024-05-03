import { SelectBadge } from '~/components/edit/selectBadge';
import { FormLayout } from '~/components/edit/layout';
import { useMainButton } from '~/hooks/useMainButton';
import { useNavigate, useParams, useSearchParams } from '@solidjs/router';
import { createEffect, createSignal, onCleanup } from 'solid-js';
import { editCollaboration, setEditCollaboration } from '~/store';
import { fetchBadges } from '~/api';
import { createQuery } from '@tanstack/solid-query';
import { Badge } from '../../../../gen';

export default function SelectBadges() {
  const mainButton = useMainButton();

  const [badgeSearch, setBadgeSearch] = createSignal('');

  const navigate = useNavigate();
  const idPath = useParams().id ? '/' + useParams().id : '';

  const [searchParams, _] = useSearchParams();

  const navigateNext = () => {
    console.log('NAVIGATE: ', `/collaborations/edit${idPath}/interests`);
    navigate(`/collaborations/edit${idPath}/interests`, {
      state: { back: true },
    });
  };

  const navigateCreateBadge = () => {
    navigate(
      `/collaborations/edit${idPath}/create-badge?badge_name=` + badgeSearch(),
      {
        state: { back: true },
      },
    );
  };

  const fetchBadgeQuery = createQuery(() => ({
    queryKey: ['badges'],
    // then push selected to the top
    queryFn: () =>
      fetchBadges().then(badges => {
        const selected = editCollaboration.badge_ids;
        return [
          ...selected.map(id => badges.find((b: Badge) => b.id === id)),
          ...badges.filter((b: Badge) => !selected.includes(b.id!)),
        ];
      }),
  }));

  createEffect(async () => {
    if (searchParams.refetch) {
      await fetchBadgeQuery.refetch();
    }
  });

  mainButton.onClick(navigateNext);

  createEffect(() => {
    if (editCollaboration.badge_ids && editCollaboration.badge_ids.length > 0) {
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
        badges={fetchBadgeQuery.data!}
      />
    </FormLayout>
  );
}
