import {
  createEffect,
  createSignal,
  For,
  onMount,
  Show,
  Suspense,
} from 'solid-js'
import { Link } from '~/components/link'
import BadgeList from '~/components/badge-list'
import useDebounce from '~/lib/useDebounce'
import { store } from '~/store'
import FillProfilePopup from '~/components/fill-profile-popup'
import { LocationBadge } from '~/components/location-badge'
import { useTranslations } from '~/lib/locale-context'
import { useInfiniteQuery } from '@tanstack/solid-query'
import { verificationStatus, UserProfileResponse } from '~/gen'
import { fetchUsers } from '~/lib/api'
import { useNavigation } from '~/lib/useNavigation'
import { Motion, Presence } from 'solid-motionone'

export const [search, setSearch] = createSignal('')

export default function FeedPage() {
  const { t } = useTranslations()
  const navigation = useNavigation()

  const updateSearch = useDebounce(setSearch, 350)

  const query = useInfiniteQuery(() => ({
    queryKey: ['users', search()],
    queryFn: fetchUsers,
    getNextPageParam: lastPage => lastPage.nextPage,
    initialPageParam: 1,
  }))

  const [scroll, setScroll] = createSignal(0)
  const [profilePopup, setProfilePopup] = createSignal(false)
  const [communityPopup, setCommunityPopup] = createSignal(false)
  const [isLoadingMore, setIsLoadingMore] = createSignal(false)

  const loadMoreUsers = () => {
    if (query.hasNextPage && !query.isFetchingNextPage && !isLoadingMore()) {
      setIsLoadingMore(true)
      query.fetchNextPage().finally(() => setIsLoadingMore(false))
    }
  }

  createEffect(() => {
    const onScroll = () => {
      setScroll(window.scrollY)

      const feedElement = document.getElementById('feed')
      if (feedElement) {
        const { scrollTop, scrollHeight, clientHeight } = feedElement
        if (scrollHeight - scrollTop - clientHeight < 300) {
          loadMoreUsers()
        }
      }
    }

    const feedElement = document.getElementById('feed')
    if (feedElement) {
      feedElement.addEventListener('scroll', onScroll)
      return () => feedElement.removeEventListener('scroll', onScroll)
    }
  })

  onMount(() => {
    window.Telegram.WebApp.CloudStorage.getItem(
      'profilePopup',
      updateProfilePopup,
    )
    window.Telegram.WebApp.CloudStorage.getItem(
      'communityPopup',
      updateCommunityPopup,
    )

    window.Telegram.WebApp.disableClosingConfirmation()
    // window.Telegram.WebApp.CloudStorage.removeItem('profilePopup')
    // window.Telegram.WebApp.CloudStorage.removeItem('communityPopup')
  })

  const closePopup = (name: string) => {
    switch (name) {
      case 'profilePopup':
        setProfilePopup(false)
        break
      case 'communityPopup':
        setCommunityPopup(false)
        break
    }
    window.Telegram.WebApp.CloudStorage.setItem(name, 'closed')
  }

  const updateProfilePopup = (err: unknown, value: unknown) => {
    setProfilePopup(value !== 'closed')
  }

  const updateCommunityPopup = (err: unknown, value: unknown) => {
    setCommunityPopup(value !== 'closed')
  }

  const allUsers = () => {
    if (!query.data) return []
    return query.data.pages.flatMap(page => page.data)
  }

  return (
    <div class="flex h-screen flex-col overflow-hidden">
      <Motion.div
        class="flex w-full flex-shrink-0 flex-col items-center justify-between space-y-4 border-b p-4"
        initial={{ y: -10, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        transition={{ duration: 0.3 }}
      >
        <Presence exitBeforeEnter>
          <Show
            when={
              store.user.verification_status ==
                verificationStatus.VerificationStatusUnverified &&
              profilePopup()
            }
          >
            <Motion.div
              initial={{ scale: 0.9, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.9, opacity: 0 }}
              transition={{ duration: 0.3 }}
            >
              <FillProfilePopup onClose={() => closePopup('profilePopup')} />
            </Motion.div>
          </Show>
        </Presence>
        <Presence exitBeforeEnter>
          <Show
            when={
              communityPopup() &&
              store.user.verification_status ==
                verificationStatus.VerificationStatusVerified
            }
          >
            <Motion.div
              initial={{ scale: 0.9, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.9, opacity: 0 }}
              transition={{ duration: 0.3 }}
            >
              <OpenCommunityPopup
                onClose={() => closePopup('communityPopup')}
              />
            </Motion.div>
          </Show>
        </Presence>
        <div class="relative flex h-10 w-full flex-row items-center justify-center rounded-lg bg-secondary">
          <input
            class="h-full w-full bg-transparent px-2.5 placeholder:text-secondary-foreground"
            placeholder={t('common.search.people')}
            type="text"
            value={search()}
            onInput={e => updateSearch(e.currentTarget.value)}
          />
          <Presence exitBeforeEnter>
            <Show when={search()}>
              <Motion.button
                class="absolute right-2.5 flex size-5 shrink-0 items-center justify-center rounded-full bg-secondary"
                onClick={() => setSearch('')}
                initial={{ scale: 0, opacity: 0 }}
                animate={{ scale: 1, opacity: 1 }}
                exit={{ scale: 0, opacity: 0 }}
                transition={{ duration: 0.2 }}
              >
                <span class="material-symbols-rounded text-[20px] text-secondary">
                  close
                </span>
              </Motion.button>
            </Show>
          </Presence>
        </div>
      </Motion.div>
      <div
        class="flex h-full w-full flex-shrink-0 flex-col overflow-y-auto pb-20"
        id="feed"
      >
        <Suspense>
          <For each={allUsers()}>
            {(user, index) => (
              <Motion.div
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.3, delay: index() * 0.05 }}
              >
                <UserCard user={user} scroll={scroll()} index={index()} />
                <Motion.div
                  class="h-px w-full bg-border"
                  initial={{ scaleX: 0 }}
                  animate={{ scaleX: 1 }}
                  transition={{ duration: 0.3, delay: index() * 0.05 + 0.2 }}
                  style={{ 'transform-origin': 'left center' }}
                />
              </Motion.div>
            )}
          </For>

          <Show when={query.isFetchingNextPage}>
            <Motion.div
              class="flex justify-center p-4"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              transition={{ duration: 0.3 }}
            >
              <Motion.div
                class="h-10 w-10 rounded-full border-4 border-primary border-t-transparent"
                animate={{ rotate: 360 }}
                transition={{ duration: 1, repeat: Infinity, easing: 'linear' }}
              />
            </Motion.div>
          </Show>

          <Show when={!query.hasNextPage && allUsers().length > 0}>
            <Motion.div
              class="p-4 text-center text-secondary-foreground"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              transition={{ duration: 0.3 }}
            >
              {t('common.search.noMoreResults')}
            </Motion.div>
          </Show>

          <Show when={allUsers().length === 0 && !query.isLoading}>
            <Motion.div
              class="p-4 text-center text-secondary-foreground"
              initial={{ opacity: 0, scale: 0.9 }}
              animate={{ opacity: 1, scale: 1 }}
              transition={{ duration: 0.3 }}
            >
              {t('common.search.noResults')}
            </Motion.div>
          </Show>
        </Suspense>
      </div>
    </div>
  )
}

type UserCardProps = {
  user: UserProfileResponse
  scroll: number
  index: number
}

const UserCard = (props: UserCardProps) => {
  const shortenDescription = (description: string) => {
    if (description.length <= 120) return description
    return description.slice(0, 120) + '...'
  }
  return (
    <Motion.div transition={{ duration: 0.2 }}>
      <Link
        class="flex flex-col items-start px-4 pb-5 pt-4 text-start"
        href={`/users/${props.user.id}`}
        state={{ from: '/' }}
      >
        <Motion.img
          class="size-10 rounded-xl object-cover"
          src={`https://assets.peatch.io/cdn-cgi/image/width=100/${props.user.avatar_url}`}
          loading="lazy"
          alt="User Avatar"
          initial={{ scale: 0 }}
          animate={{ scale: 1 }}
          transition={{ duration: 0.3, delay: props.index * 0.05 }}
        />
        <Motion.p
          class="mt-3 text-3xl font-semibold capitalize text-primary"
          initial={{ opacity: 0, x: -10 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.3, delay: props.index * 0.05 + 0.1 }}
        >
          {props.user.name?.trimEnd()}:
        </Motion.p>
        <Motion.p
          class="text-3xl capitalize"
          initial={{ opacity: 0, x: -10 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.3, delay: props.index * 0.05 + 0.15 }}
        >
          {props.user.title}
        </Motion.p>
        <Motion.p
          class="mt-2 text-sm text-secondary-foreground"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.3, delay: props.index * 0.05 + 0.2 }}
        >
          {shortenDescription(props.user.description!)}
        </Motion.p>
        <Motion.div
          initial={{ opacity: 0, y: 5 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.3, delay: props.index * 0.05 + 0.25 }}
        >
          <LocationBadge
            country={props.user.location?.country_name}
            city={props.user.location?.name}
            countryCode={props.user.location?.country_code}
          />
        </Motion.div>
        <Show when={props.user.badges && props.user.badges.length > 0}>
          <Motion.div
            initial={{ opacity: 0, y: 5 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.3, delay: props.index * 0.05 + 0.3 }}
          >
            <BadgeList badges={props.user.badges || []} position="start" />
          </Motion.div>
        </Show>
      </Link>
    </Motion.div>
  )
}

const OpenCommunityPopup = (props: { onClose: () => void }) => {
  return (
    <Motion.div
      class="relative w-full rounded-xl bg-secondary p-3 text-center"
      initial={{ scale: 0.9, opacity: 0 }}
      animate={{ scale: 1, opacity: 1 }}
      transition={{ duration: 0.3 }}
    >
      <Motion.button
        class="absolute right-4 top-4 flex size-6 items-center justify-center rounded-full bg-background"
        onClick={() => props.onClose()}
        press={{ scale: 0.98 }}
      >
        <span class="material-symbols-rounded text-[20px] text-secondary-foreground">
          close
        </span>
      </Motion.button>
      <Motion.div
        class="text-green flex items-center justify-center gap-1 text-2xl font-extrabold"
        initial={{ y: -10, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        transition={{ duration: 0.3, delay: 0.1 }}
      >
        <Motion.span
          class="material-symbols-rounded text-[36px] text-green-400"
          animate={{ rotate: [0, 5, -5, 0] }}
          transition={{ duration: 2, repeat: Infinity }}
        >
          maps_ugc
        </Motion.span>
        Join community
      </Motion.div>
      <Motion.p
        class="mt-2 text-base font-normal text-secondary-foreground"
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ duration: 0.3, delay: 0.2 }}
      >
        To talk with founders and users. Discuss and solve problems together
      </Motion.p>
      <Motion.button
        class="mt-4 flex h-10 w-full items-center justify-center rounded-xl bg-primary text-sm font-semibold"
        onClick={() =>
          window.Telegram.WebApp.openTelegramLink(
            'https://t.me/peatch_community',
          )
        }
        press={{ scale: 0.98 }}
        initial={{ y: 10, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        transition={{ duration: 0.3, delay: 0.3 }}
      >
        Open Peatch Community
      </Motion.button>
    </Motion.div>
  )
}
