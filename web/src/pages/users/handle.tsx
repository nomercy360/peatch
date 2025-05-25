import {
  createEffect,
  createSignal,
  For,
  Match,
  onCleanup,
  Show,
  Switch,
} from 'solid-js'
import { Navigate, useNavigate, useParams } from '@solidjs/router'
import { fetchProfile, followUser, publishUserProfile } from '~/lib/api'
import { addToast } from '~/components/toast'
import { setUser, store } from '~/store'
import { useMainButton } from '~/lib/useMainButton'
import { queryClient } from '~/App'
import { UserProfileResponse, verificationStatus } from '~/gen'
import { useTranslations } from '~/lib/locale-context'
import { useMutation, useQuery } from '@tanstack/solid-query'
import { Motion } from 'solid-motionone'
import { useSecondaryButton } from '~/lib/useSecondaryButton'
import LinkEditor from '~/components/link-editor'

export default function UserProfilePage() {
  const mainButton = useMainButton()
  const secondaryButton = useSecondaryButton()
  const [badgesExpanded, setBadgesExpanded] = createSignal(false)
  const [linkDrawerOpen, setLinkDrawerOpen] = createSignal(false)

  const [opportunitiesExpanded, setOpportunitiesExpanded] = createSignal(false)

  const navigate = useNavigate()

  const params = useParams()
  const id = params.handle

  const { t } = useTranslations()

  const query = useQuery(() => ({
    queryKey: ['profiles', id],
    queryFn: () => fetchProfile(id),
    retry: 1,
  }))

  const followMutate = useMutation(() => ({
    mutationFn: (id: string) => followUser(id),
    retry: 0,
    onMutate: async (id: string) => {
      await queryClient.cancelQueries({ queryKey: ['profiles', id] })

      const previousProfile = queryClient.getQueryData([
        'profiles',
        id,
      ]) as UserProfileResponse

      queryClient.setQueryData(['profiles', id], (old: UserProfileResponse) => {
        if (old) {
          return {
            ...old,
            is_following: true,
          }
        }
        return old
      })

      return { previousProfile }
    },
    onSuccess: () => {
      addToast(t('pages.users.followSuccess'), 'success')
    },
    onError: (
      error: any,
      _id: string,
      context?: { previousProfile?: UserProfileResponse },
    ) => {
      if (context?.previousProfile) {
        queryClient.setQueryData(
          ['profiles', context.previousProfile.id],
          context.previousProfile,
        )
      }

      if (error.botBlocked) {
        const username = error.username
        if (username) {
          addToast(t('pages.users.botBlocked'), 'warning', {
            text: t('pages.users.messageUser'),
            onClick: () => {
              window.Telegram.WebApp.openTelegramLink(
                `https://t.me/${username}`,
              )
            },
          })
        } else {
          addToast(t('pages.users.botBlocked'), 'warning')
        }
      } else {
        addToast(t('pages.users.followError'), 'error')
      }
    },
  }))

  const publishMutate = useMutation(() => ({
    mutationFn: publishUserProfile,
    retry: 0,
    onSuccess: () => {
      addToast(t('pages.users.publishSuccess'), 'success')
      setUser({
        ...store.user,
        hidden_at: undefined,
      })
    },
    onError: (error: any) => {
      if (error.code === 400) {
        addToast(t('pages.users.profileIncomplete'), 'error')
      } else if (error.code === 403) {
        addToast(t('pages.users.profileBlocked'), 'error')
      } else {
        addToast(t('pages.users.publishError'), 'error')
      }
    },
  }))

  const isCurrentUserProfile = store.user.id === id
  // @ts-ignore - hidden_at might not be in types yet
  const isProfileHidden = () =>
    isCurrentUserProfile &&
    store.user.hidden_at !== null &&
    store.user.hidden_at !== undefined

  const navigateToEdit = () => {
    navigate('/users/edit', { state: { back: true } })
  }

  const publishProfile = async () => {
    window.Telegram.WebApp.HapticFeedback.impactOccurred('light')
    publishMutate.mutate()
  }

  const follow = async () => {
    if (!query.data) return
    followMutate.mutate(query.data.id)
    window.Telegram.WebApp.HapticFeedback.impactOccurred('light')
  }

  createEffect(() => {
    if (isCurrentUserProfile && !linkDrawerOpen()) {
      if (isProfileHidden()) {
        mainButton.offClick(navigateToEdit)
        mainButton.enable(t('common.buttons.publish'))
        mainButton.onClick(publishProfile)
        // Secondary button for edit profile
        secondaryButton.enable(t('common.buttons.edit'))
        secondaryButton.onClick(navigateToEdit)
      } else {
        secondaryButton.hide()
        secondaryButton.offClick(navigateToEdit)
        mainButton.offClick(publishProfile)
        mainButton.enable(t('common.buttons.edit'))
        mainButton.onClick(navigateToEdit)
      }
    } else if (linkDrawerOpen() && isCurrentUserProfile) {
      secondaryButton.hide()
      mainButton.offClick(navigateToEdit)
      secondaryButton.offClick(navigateToEdit)
    }
  })

  onCleanup(() => {
    mainButton.offClick(navigateToEdit)
    mainButton.offClick(publishProfile)
    secondaryButton.offClick(navigateToEdit)
    secondaryButton.hide()
    mainButton.hide()
  })

  function shareURL() {
    const url =
      'https://t.me/share/url?' +
      new URLSearchParams({
        url: 'https://t.me/peatch_bot/app?startapp=u' + id,
      }).toString() +
      '&text=' +
      t('pages.users.shareURLText', { name: query.data.name })

    window.Telegram.WebApp.openTelegramLink(url)
  }

  const showInfoPopup = () => {
    window.Telegram.WebApp.showAlert(t('pages.users.verificationStatusDenied'))
  }

  return (
    <div>
      <Switch>
        <Match when={query.isLoading}>
          <AnimatedLoader />
        </Match>
        <Match when={query.isError}>
          <Navigate href={'/404'} />
        </Match>
        <Match when={query.isSuccess}>
          <Motion.div
            class="flex h-fit min-h-screen flex-col items-center p-2 pb-8 text-center"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ duration: 0.3 }}
          >
            <Show
              when={
                isCurrentUserProfile &&
                store.user.verification_status ===
                  verificationStatus.VerificationStatusDenied
              }
            >
              <Motion.button
                onClick={showInfoPopup}
                class="absolute left-4 top-4 flex size-7 items-center justify-center rounded-lg bg-secondary"
                initial={{ scale: 0, opacity: 0 }}
                animate={{ scale: 1, opacity: 1 }}
                transition={{ duration: 0.3, delay: 0.5 }}
                press={{ scale: 0.95 }}
              >
                <span class="material-symbols-rounded text-[16px] text-secondary-foreground">
                  visibility_off
                </span>
              </Motion.button>
            </Show>
            <Motion.img
              alt="User Avatar"
              class="relative aspect-square w-32 rounded-xl bg-cover bg-center object-cover"
              src={`https://assets.peatch.io/cdn-cgi/image/width=400/${query.data.avatar_url}`}
              initial={{ scale: 0.5, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              transition={{ duration: 0.4, easing: 'ease-out' }}
            />
            <Show when={!isCurrentUserProfile}>
              <Motion.button
                onClick={shareURL}
                class="absolute right-3 top-3 flex size-8 flex-row items-center justify-center rounded-lg border bg-secondary px-3 text-accent-foreground transition-all duration-300"
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ duration: 0.3, delay: 0.3 }}
                press={{ scale: 0.95 }}
              >
                <span class="material-symbols-rounded text-[16px]">
                  ios_share
                </span>
              </Motion.button>
            </Show>
            <Show when={!isCurrentUserProfile}>
              <Motion.button
                class={`mt-4 flex h-8 flex-row items-center space-x-1 rounded-2xl border px-3 transition-all duration-300 ${
                  query.data.is_following
                    ? 'border-secondary bg-secondary text-secondary-foreground'
                    : 'border-primary bg-primary text-primary-foreground'
                }`}
                onClick={() => follow()}
                disabled={query.data.is_following}
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ duration: 0.3, delay: 0.4 }}
                press={{ scale: 0.95 }}
              >
                <Motion.span
                  class="material-symbols-rounded text-[16px]"
                  animate={{
                    rotate: query.data.is_following
                      ? 0
                      : [0, -20, 20, -10, 10, 0],
                  }}
                  transition={{ duration: 0.5 }}
                >
                  {query.data.is_following ? 'check' : 'waving_hand'}
                </Motion.span>
                <span class="text-sm">
                  {query.data.is_following
                    ? t('pages.users.saidHi')
                    : t('pages.users.sayHi')}
                </span>
              </Motion.button>
            </Show>
            <div class="w-full px-4 py-2.5">
              <Motion.p
                class="text-3xl font-semibold capitalize text-primary"
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ duration: 0.3, delay: 0.2 }}
              >
                {query.data.name}:
              </Motion.p>
              <Motion.p
                class="text-3xl capitalize"
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ duration: 0.3, delay: 0.25 }}
              >
                {query.data.title}
              </Motion.p>
              <Motion.p
                class="mt-1 text-start text-sm font-normal text-secondary-foreground"
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ duration: 0.3, delay: 0.3 }}
              >
                {query.data.description}
              </Motion.p>
              <Motion.div
                class="mt-3 flex flex-row flex-wrap items-center justify-start gap-1"
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ duration: 0.3, delay: 0.35 }}
              >
                <For
                  each={
                    badgesExpanded()
                      ? query.data.badges
                      : query.data.badges.slice(0, 3)
                  }
                >
                  {(badge, index) => (
                    <Motion.div
                      class="flex h-8 flex-row items-center justify-center gap-1 rounded-xl border px-2"
                      style={{
                        'background-color': `#${badge.color}`,
                        'border-color': `#${badge.color}`,
                      }}
                      initial={{ scale: 0, opacity: 0 }}
                      animate={{ scale: 1, opacity: 1 }}
                      transition={{
                        duration: 0.2,
                        delay: 0.1 + index() * 0.05,
                        easing: 'ease-out',
                      }}
                    >
                      <span class="material-symbols-rounded text-sm text-white">
                        {String.fromCodePoint(parseInt(badge.icon!, 16))}
                      </span>
                      <p class="text-xs font-semibold text-white">
                        {badge.text}
                      </p>
                    </Motion.div>
                  )}
                </For>
              </Motion.div>
              <Show when={query.data.badges.length > 3}>
                <Motion.div
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  transition={{ duration: 0.3, delay: 0.2 }}
                >
                  <ExpandButton
                    expanded={badgesExpanded()}
                    setExpanded={setBadgesExpanded}
                  />
                </Motion.div>
              </Show>
              <Motion.p
                class="pb-1 pt-3 text-start text-xl font-extrabold"
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.3, delay: 0.45 }}
              >
                {t('pages.users.availableFor')}
              </Motion.p>
              <Motion.div
                class="flex w-full flex-col items-center justify-start gap-1"
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ duration: 0.3, delay: 0.5 }}
              >
                <For
                  each={
                    opportunitiesExpanded()
                      ? query.data.opportunities
                      : query.data.opportunities.slice(0, 3)
                  }
                >
                  {(op, index) => (
                    <Motion.div
                      class="flex h-14 w-full flex-row items-center justify-start gap-2 rounded-xl bg-secondary px-2 text-secondary-foreground"
                      initial={{ opacity: 0, x: -20 }}
                      animate={{ opacity: 1, x: 0 }}
                      transition={{
                        duration: 0.3,
                        delay: 0.2 + index() * 0.05,
                      }}
                    >
                      <Motion.div
                        class="flex size-8 shrink-0 items-center justify-center rounded-full text-white"
                        style={{ 'background-color': `#${op.color}` }}
                        transition={{
                          duration: 0.5,
                          delay: 0.6 + index() * 0.05,
                        }}
                      >
                        <span class="material-symbols-rounded text-sm">
                          {String.fromCodePoint(parseInt(op.icon!, 16))}
                        </span>
                      </Motion.div>
                      <div class="text-start">
                        <p class="text-xs font-semibold text-foreground">
                          {op.text}
                        </p>
                        <p class="text-[10px] leading-tight">
                          {op.description}
                        </p>
                      </div>
                    </Motion.div>
                  )}
                </For>
                <Show when={query.data.opportunities.length > 3}>
                  <Motion.div
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    transition={{ duration: 0.3, delay: 0.3 }}
                  >
                    <ExpandButton
                      expanded={opportunitiesExpanded()}
                      setExpanded={setOpportunitiesExpanded}
                    />
                  </Motion.div>
                </Show>
              </Motion.div>
              <LinkEditor
                links={query.data.links || []}
                isCurrentUser={isCurrentUserProfile}
                onDrawerStateChange={setLinkDrawerOpen}
              />
            </div>
          </Motion.div>
        </Match>
      </Switch>
    </div>
  )
}

const ExpandButton = (props: {
  expanded: boolean
  setExpanded: (val: boolean) => void
}) => {
  return (
    <Motion.button
      class="flex h-8 w-full items-center justify-start rounded-xl bg-transparent text-xs font-medium text-secondary-foreground"
      onClick={() => props.setExpanded(!props.expanded)}
      transition={{ duration: 0.2 }}
    >
      <Motion.span
        class="material-symbols-rounded text-secondary-foreground"
        animate={{ rotate: props.expanded ? 180 : 0 }}
        transition={{ duration: 0.3 }}
      >
        expand_more
      </Motion.span>
      {props.expanded ? 'show less' : 'show more'}
    </Motion.button>
  )
}

const AnimatedLoader = () => {
  return (
    <div class="flex min-h-screen flex-col items-start justify-start bg-secondary p-2">
      <Motion.div
        class="aspect-square w-full rounded-xl bg-background"
        initial={{ opacity: 0 }}
        animate={{ opacity: [0.3, 0.6, 0.3] }}
        transition={{ duration: 1.5, repeat: Infinity }}
      />
      <div class="flex flex-col items-start justify-start p-2">
        <Motion.div
          class="mt-2 h-6 w-1/2 rounded bg-background"
          initial={{ opacity: 0 }}
          animate={{ opacity: [0.3, 0.6, 0.3] }}
          transition={{ duration: 1.5, repeat: Infinity, delay: 0.1 }}
        />
        <Motion.div
          class="mt-2 h-6 w-1/2 rounded bg-background"
          initial={{ opacity: 0 }}
          animate={{ opacity: [0.3, 0.6, 0.3] }}
          transition={{ duration: 1.5, repeat: Infinity, delay: 0.2 }}
        />
        <Motion.div
          class="mt-2 h-20 w-full rounded bg-background"
          initial={{ opacity: 0 }}
          animate={{ opacity: [0.3, 0.6, 0.3] }}
          transition={{ duration: 1.5, repeat: Infinity, delay: 0.3 }}
        />
        <div class="mt-4 flex w-full flex-row flex-wrap items-center justify-start gap-2">
          <For each={[40, 32, 40, 28, 32]}>
            {(width, index) => (
              <Motion.div
                class="h-10 rounded-2xl bg-background"
                style={{ width: `${width * 4}px` }}
                initial={{ opacity: 0 }}
                animate={{ opacity: [0.3, 0.6, 0.3] }}
                transition={{
                  duration: 1.5,
                  repeat: Infinity,
                  delay: 0.4 + index() * 0.1,
                }}
              />
            )}
          </For>
        </div>
      </div>
    </div>
  )
}
