import { createStore } from 'solid-js/store';
import { store } from '../../store';
import { createEffect, createSignal, Match, onCleanup, Switch } from 'solid-js';
import { useButtons } from '../../hooks/useBackButton';
import { Badge, UpdateUserRequest, User } from '../../../gen';
import { useNavigate } from '@solidjs/router';
import { SelectBadge } from '../../components/edit/selectBadge';
import { createQuery } from '@tanstack/solid-query';
import {
  CDN_URL,
  fetchBadges,
  fetchOpportunities,
  fetchPresignedUrl,
  postBadge,
  updateUser,
  uploadToS3,
} from '../../api';
import CreateBadge from '../../components/edit/createBadge';
import { SelectOpportunity } from '../../components/edit/selectOpp';
import SelectLocation from '../../components/edit/selectLocation';
import FillDescription from '../../components/edit/fillDescription';
import ImageUpload from '../../components/edit/imageUpload';

const totalScreens = 6;

export default function EditUserProfile() {
  const [user, setUser] = createStore<UpdateUserRequest>({
    first_name: store.user.first_name || '',
    last_name: store.user.last_name || '',
    title: store.user.title || '',
    description: store.user.description || '',
    avatar_url: store.user.avatar_url || '',
    city: store.user.city || undefined,
    country: store.user.country || '',
    country_code: store.user.country_code || '',
    badge_ids: store.user.badges?.map(b => b.id) || ([] as any),
    opportunity_ids: store.user.opportunities?.map(o => o.id) || ([] as any),
  });

  const [screen, setScreen] = createSignal(1);
  const [createBadgeOpen, setCreateBadgeOpen] = createSignal(false);
  const [badgeSearch, setBadgeSearch] = createSignal('');
  const [_, setImgUploadProgress] = createSignal(0);

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

  createEffect(() => {
    mainButton.setVisible(true);
    mainButton.setText('Next');
    mainButton.onClick(nextScreen);
    backButton.setVisible(true);
  });

  onCleanup(() => {
    mainButton.setVisible(false);
    mainButton.offClick(nextScreen);
    backButton.setVisible(false);
    backButton.offClick(prevScreen);
    backButton.offClick(goBack);
  });

  const fetchBadgeQuery = createQuery(() => ({
    queryKey: ['badges'],
    queryFn: () => fetchBadges(),
  }));

  const fetchOpportunityQuery = createQuery(() => ({
    queryKey: ['opportunities'],
    queryFn: () => fetchOpportunities(),
  }));

  const saveUser = async () => {
    if (imgFile() && imgFile() !== null) {
      try {
        const { path, url } = await fetchPresignedUrl(imgFile()!.name);
        mainButton.showProgress(false);
        await uploadToS3(
          url,
          imgFile()!,
          e => {
            setImgUploadProgress(Math.round((e.loaded / e.total) * 100));
          },
          () => {
            setImgUploadProgress(0);
          },
        );

        setUser('avatar_url', path);
      } catch (e) {
        console.error(e);
        mainButton.hideProgress();
      }
    }
    store.user = await updateUser(user);
  };

  const publishBadge = async () => {
    if (createBadge.text && createBadge.color && createBadge.icon) {
      const { id } = await postBadge(
        createBadge.text,
        createBadge.color,
        createBadge.icon,
      );

      setUser('badge_ids', [...user.badge_ids, id]);

      setCreateBadgeOpen(false);
    }

    await fetchBadgeQuery.refetch();
  };

  createEffect(() => {
    switch (screen()) {
      case 1:
        backButton.offClick(prevScreen);
        backButton.onClick(goBack);
        if (user.first_name && user.last_name && user.title) {
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

          if (user.badge_ids.length) {
            mainButton.setActive(true);
          }
        }
        break;
      case 3:
        if (user.opportunity_ids.length) {
          mainButton.setActive(true);
        } else {
          mainButton.setActive(false);
        }
        break;
      case 4:
        if (user.country && user.country_code) {
          mainButton.setActive(true);
        } else {
          mainButton.setActive(false);
        }
        break;
      case 5:
        if (user.description) {
          mainButton.setActive(true);
        } else {
          mainButton.setActive(false);
        }
        mainButton.onClick(nextScreen);
        mainButton.offClick(saveUser);
        mainButton.setText('Next');
        break;
      case 6:
        if (user.avatar_url || imgFile()) {
          mainButton.setActive(true);
        } else {
          mainButton.setActive(false);
        }
        mainButton.setText('Save');
        mainButton.offClick(nextScreen);
        mainButton.onClick(saveUser);
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

  const [imgFile, setImgFile] = createSignal<File | null>(null);

  return (
    <div class="flex h-screen flex-col items-center justify-start p-3.5">
      <ProgressBar screen={screen()} />
      <Switch>
        <Match when={screen() === 1}>
          <GeneralInfo user={user} setUser={setUser} />
        </Match>
        <Match when={screen() === 2 && !createBadgeOpen()}>
          <SelectBadge
            selected={user.badge_ids}
            setSelected={b => setUser('badge_ids', b)}
            setCreateBadgeModal={onBadgeModalOpen}
            search={badgeSearch()}
            setSearch={setBadgeSearch}
            badges={fetchBadgeQuery.data}
          />
        </Match>
        <Match when={screen() === 2 && createBadgeOpen()}>
          <CreateBadge
            createBadge={createBadge}
            setCreateBadge={setCreateBadge}
          />
        </Match>
        <Match when={screen() === 3}>
          <SelectOpportunity
            selected={user.opportunity_ids}
            setSelected={b => setUser('opportunity_ids', b)}
            opportunities={fetchOpportunityQuery.data}
          />
        </Match>
        <Match when={screen() === 4}>
          <SelectLocation
            city={user.city}
            setCity={c => setUser('city', c)}
            country={user.country}
            setCountry={c => setUser('country', c)}
            countryCode={user.country_code}
            setCountryCode={c => setUser('country_code', c)}
          />
        </Match>
        <Match when={screen() === 5}>
          <FillDescription
            setDescription={d => setUser('description', d)}
            description={user.description}
          />
        </Match>
        <Match when={screen() === 6}>
          <ImageUpload
            imgURL={`${CDN_URL}/${user.avatar_url}`}
            imgFile={imgFile()}
            setImgFile={setImgFile}
          />
        </Match>
      </Switch>
    </div>
  );
}

export function ProgressBar(props: { screen: number }) {
  return (
    <div class="h-1.5 w-[160px] rounded-lg bg-peatch-bg">
      <div
        class="h-1.5 rounded-lg bg-peatch-blue"
        style={`width: ${(props.screen / totalScreens) * 100}%`}
      ></div>
    </div>
  );
}

export function FormLayout(props: {
  children: any;
  title: string;
  description: string;
}) {
  return (
    <>
      <p class="mt-2 text-3xl">{props.title}</p>
      <p class="mt-1 text-sm text-gray">{props.description}</p>
      {props.children}
    </>
  );
}

export function GeneralInfo(props: {
  user: User;
  setUser: (key: string, value: string) => void;
}) {
  return (
    <FormLayout
      title="Introduce yourself"
      description="It will appears in your profile card, everyone will see it"
    >
      <div class="mt-5 flex w-full flex-col items-center justify-start gap-3">
        <input
          class="h-10 w-full rounded-lg bg-peatch-bg px-2.5 text-black placeholder:text-gray"
          placeholder="First Name"
          value={props.user.first_name}
          onInput={e => props.setUser('first_name', e.currentTarget.value)}
        />
        <input
          class="h-10 w-full rounded-lg bg-peatch-bg px-2.5 text-black placeholder:text-gray"
          placeholder="Last Name"
          value={props.user.last_name}
          onInput={e => props.setUser('last_name', e.currentTarget.value)}
        />
        <input
          class="h-10 w-full rounded-lg bg-peatch-bg px-2.5 text-black placeholder:text-gray"
          placeholder="Title"
          value={props.user.title}
          onInput={e => props.setUser('title', e.currentTarget.value)}
        />
      </div>
    </FormLayout>
  );
}
