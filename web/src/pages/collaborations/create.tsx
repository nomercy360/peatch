import { createStore } from 'solid-js/store';
import { createEffect, createSignal, Match, onCleanup, Switch } from 'solid-js';
import { useButtons } from '../../hooks/useBackButton';
import { Badge, CreateCollaboration } from '../../../gen';
import { useNavigate } from '@solidjs/router';
import { SelectBadge } from '../../components/edit/selectBadge';
import { createQuery } from '@tanstack/solid-query';
import { createCollaboration, fetchBadges, fetchOpportunities, postBadge } from '../../api';
import CreateBadge from '../../components/edit/createBadge';
import { SelectOpportunity } from '../../components/edit/selectOpp';
import SelectLocation from '../../components/edit/selectLocation';
import { FormLayout } from '../../components/edit/layout';
import DescribeCollaboration from '../../components/edit/collaboration/descriptiom';
import { ProgressBar } from '../../components/edit/progress';

const totalScreens = 4;

export default function Create() {
  const [collab, setCollab] = createStore<CreateCollaboration>({
    title: '',
    description: '',
    city: undefined,
    country: '',
    country_code: '',
    badge_ids: [],
    opportunity_id: 0,
    is_payable: false,
  });

  const [screen, setScreen] = createSignal(1);
  const [createBadgeOpen, setCreateBadgeOpen] = createSignal(false);
  const [badgeSearch, setBadgeSearch] = createSignal('');

  const { mainButton, backButton } = useButtons();

  const navigate = useNavigate();

  const nextScreen = () => {
    if (screen() === totalScreens) {
      return;
    }
    setScreen(screen() + 1);
  };

  const prevScreen = () => {
    if (screen() === 1) {
      return;
    }
    setScreen(screen() - 1);
  };

  const goBack = () => {
    navigate('/');
  };

  const fetchBadgeQuery = createQuery(() => ({
    queryKey: ['badges'],
    queryFn: () => fetchBadges(),
  }));

  const fetchOpportunityQuery = createQuery(() => ({
    queryKey: ['opportunities'],
    queryFn: () => fetchOpportunities(),
  }));

  const createCollab = async () => {
    const created = await createCollaboration(collab);
    navigate('/collaborations/' + created.id);
  };

  const publishBadge = async () => {
    if (createBadge.text && createBadge.color && createBadge.icon) {
      const { id } = await postBadge(
        createBadge.text,
        createBadge.color,
        createBadge.icon,
      );

      setCollab('badge_ids', [...collab.badge_ids, id]);

      setCreateBadgeOpen(false);
    }

    await fetchBadgeQuery.refetch();
  };

  createEffect(() => {
    switch (screen()) {
      case 1:
        backButton.offClick(prevScreen);
        backButton.onClick(goBack);
        backButton.setVisible();
        mainButton.onClick(nextScreen);
        mainButton.setVisible('Next');
        if (collab.title && collab.description) {
          mainButton.setActive(true);
        } else {
          mainButton.setActive(false);
        }
        break;
      case 2:
        backButton.offClick(goBack);

        if (createBadgeOpen()) {
          backButton.onClick(() => setCreateBadgeOpen(false));

          mainButton.setActive(false);
          mainButton.setText('Create ' + badgeSearch());
          mainButton.offClick(nextScreen);
          mainButton.onClick(publishBadge);

          if (createBadge.text && createBadge.color && createBadge.icon) {
            mainButton.setActive(true);
          }
        } else {
          backButton.onClick(prevScreen);

          mainButton.offClick(publishBadge);
          mainButton.setText('Next');
          mainButton.onClick(nextScreen);
          mainButton.setActive(false);

          if (collab.badge_ids.length) {
            mainButton.setActive(true);
          }
        }
        break;
      case 3:
        if (collab.opportunity_id) {
          mainButton.setActive(true);
        } else {
          mainButton.setActive(false);
        }
        mainButton.onClick(nextScreen);
        mainButton.offClick(publishBadge);
        mainButton.setText('Next');
        break;
      case 4:
        if (collab.country && collab.country_code) {
          mainButton.setActive(true);
        } else {
          mainButton.setActive(false);
        }
        mainButton.setText('Save');
        mainButton.offClick(nextScreen);
        mainButton.onClick(createCollab);
        break;
    }
  });

  const [createBadge, setCreateBadge] = createStore<Badge>({
    text: '',
    color: '',
    icon: '',
  });

  const onBadgeModalOpen = () => {
    setCreateBadgeOpen(true);
    setCreateBadge({
      text: badgeSearch(),
      color: '#EF5DA8',
      icon: '',
    });
  };

  onCleanup(() => {
    mainButton.hide();
    mainButton.offClick(nextScreen);
    backButton.hide();
    backButton.offClick(prevScreen);
    backButton.offClick(goBack);
  });

  return (
    <div class="flex h-screen flex-col items-center justify-start p-3.5">
      <ProgressBar screen={screen()} totalScreens={totalScreens} />
      <Switch>
        <Match when={screen() === 1}>
          <DescribeCollaboration
            title={collab.title}
            setTitle={t => setCollab('title', t)}
            description={collab.description}
            setDescription={d => setCollab('description', d)}
            isPayable={collab.is_payable}
            setIsPayable={p => setCollab('is_payable', p)}
          />
        </Match>
        <Match when={screen() === 2 && !createBadgeOpen()}>
          <FormLayout
            title="Who are you looking for?"
            description="This will help us to recommend it to other people"
          >
            <SelectBadge
              selected={collab.badge_ids}
              setSelected={b => setCollab('badge_ids', b)}
              setCreateBadgeModal={onBadgeModalOpen}
              search={badgeSearch()}
              setSearch={setBadgeSearch}
              badges={fetchBadgeQuery.data}
            />
          </FormLayout>
        </Match>
        <Match when={screen() === 2 && createBadgeOpen()}>
          <CreateBadge
            createBadge={createBadge}
            setCreateBadge={setCreateBadge}
          />
        </Match>
        <Match when={screen() === 3}>
          <FormLayout
            title="Select a theme"
            description="This will help us to recommend it to other people"
          >
            <SelectOpportunity
              selected={collab.opportunity_id}
              setSelected={b => setCollab('opportunity_id', b as number)}
              opportunities={fetchOpportunityQuery.data}
            />
          </FormLayout>
        </Match>
        <Match when={screen() === 4}>
          <FormLayout
            title="Any special location?"
            description="This will help us to recommend it to other people"
          >
            <SelectLocation
              city={collab.city}
              setCity={c => setCollab('city', c)}
              country={collab.country}
              setCountry={c => setCollab('country', c)}
              countryCode={collab.country_code}
              setCountryCode={c => setCollab('country_code', c)}
            />
          </FormLayout>
        </Match>
      </Switch>
    </div>
  );
}
