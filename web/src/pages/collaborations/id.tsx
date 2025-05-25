import {
  createEffect,
  createSignal,
  For,
  Match,
  onCleanup,
  Show,
  Suspense,
  Switch,
} from 'solid-js'
import { useNavigate, useParams, useSearchParams } from '@solidjs/router'
import { fetchCollaboration } from '~/lib/api'
import { store } from '~/store'
import { useMainButton } from '~/lib/useMainButton'
import { useQuery } from '@tanstack/solid-query'
import { Motion } from 'solid-motionone'

export default function Collaboration() {
  const mainButton = useMainButton()

  const [isCurrentUserCollab, setIsCurrentUserCollab] = createSignal(false)

  const navigate = useNavigate()
  const params = useParams()
  const [searchParams] = useSearchParams()

  const collabId = params.id

  const query = useQuery(() => ({
    queryKey: ['collaborations', collabId],
    queryFn: () => fetchCollaboration(collabId),
  }))

  createEffect(async () => {
    if (searchParams.refetch) {
      await query.refetch()
    }
  })

  createEffect(() => {
    if (query.data?.id) {
      setIsCurrentUserCollab(store.user.id === query.data.user.id)
    }
  })

  const navigateToEdit = () => {
    navigate('/collaborations/edit/' + collabId, {
      state: { from: '/collaborations/' + collabId },
    })
  }

  const openUserTelegram = () => {
    window.Telegram.WebApp.openTelegramLink(
      `https://t.me/${query.data.user.username}`,
    )
  }

  createEffect(() => {
    if (isCurrentUserCollab()) {
      mainButton.enable('Edit')
      mainButton.onClick(navigateToEdit)
    } else {
      mainButton.enable('Contact in Telegram')
      mainButton.onClick(openUserTelegram)
    }
  })

  onCleanup(async () => {
    mainButton.hide()
    mainButton.offClick(navigateToEdit)
    mainButton.offClick(openUserTelegram)
  })

  return (
    <Suspense fallback={<AnimatedLoader />}>
      <Switch>
        <Match when={!query.isLoading}>
          <Motion.div
            class="h-fit min-h-screen bg-secondary"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ duration: 0.3 }}
          >
            <Switch>
              <Match when={isCurrentUserCollab() && !query.data.published_at}>
                <Motion.div
                  initial={{ opacity: 0, scale: 0.9 }}
                  animate={{ opacity: 1, scale: 1 }}
                  transition={{ duration: 0.3, delay: 0.2 }}
                >
                  <ActionButton text="Edit" onClick={navigateToEdit} />
                </Motion.div>
              </Match>
            </Switch>
            <Motion.div
              class="flex w-full flex-col items-start justify-start px-4 pb-5 pt-4"
              style={{
                'background-color': `#${query.data.opportunity.color}`,
              }}
              initial={{ y: -20, opacity: 0 }}
              animate={{ y: 0, opacity: 1 }}
              transition={{ duration: 0.4, easing: 'ease-out' }}
            >
              <Motion.span
                class="material-symbols-rounded text-[48px] text-white"
                initial={{ scale: 0, rotate: -180 }}
                animate={{ scale: 1, rotate: 0 }}
                transition={{ duration: 0.5, easing: 'ease-out' }}
              >
                {String.fromCodePoint(
                  parseInt(query.data.opportunity.icon!, 16),
                )}
              </Motion.span>
              <Motion.p
                class="text-3xl text-white"
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ duration: 0.3, delay: 0.2 }}
              >
                {query.data.opportunity.text}:
              </Motion.p>
              <Motion.p
                class="text-3xl text-white"
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ duration: 0.3, delay: 0.25 }}
              >
                {query.data.title}:
              </Motion.p>
              <Motion.div
                class="mt-4 flex w-full flex-row items-center justify-start gap-2"
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.3, delay: 0.3 }}
              >
                <Motion.img
                  class="size-11 rounded-xl object-cover"
                  src={`https://assets.peatch.io/cdn-cgi/image/width=100/${query.data.user?.avatar_url}`}
                  alt="User Avatar"
                  initial={{ scale: 0.5, opacity: 0 }}
                  animate={{ scale: 1, opacity: 1 }}
                  transition={{ duration: 0.3, delay: 0.35 }}
                />
                <Motion.div
                  initial={{ opacity: 0, x: -10 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ duration: 0.3, delay: 0.4 }}
                >
                  <p class="text-sm font-bold text-white">
                    {query.data.user?.name}
                  </p>
                  <Show when={query.data.user?.title}>
                    <p class="text-sm text-white">{query.data.user?.title}</p>
                  </Show>
                </Motion.div>
              </Motion.div>
            </Motion.div>
            <Motion.div
              class="px-4 py-2.5"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              transition={{ duration: 0.3, delay: 0.4 }}
            >
              <Motion.p
                class="mt-1 text-start text-sm font-normal text-secondary-foreground"
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.3, delay: 0.45 }}
              >
                {query.data.description}
              </Motion.p>
              <Motion.div
                class="mt-5 flex flex-row flex-wrap items-center justify-start gap-1"
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ duration: 0.3, delay: 0.5 }}
              >
                <For each={query.data.badges}>
                  {(badge, index) => (
                    <Motion.div
                      class="flex h-10 flex-row items-center justify-center gap-[5px] rounded-2xl border px-2.5"
                      style={{
                        'background-color': `#${badge.color}`,
                        'border-color': `#${badge.color}`,
                      }}
                      initial={{ scale: 0, opacity: 0 }}
                      animate={{ scale: 1, opacity: 1 }}
                      transition={{
                        duration: 0.2,
                        delay: 0.55 + index() * 0.05,
                        easing: 'ease-out',
                      }}
                    >
                      <span class="material-symbols-rounded text-white">
                        {String.fromCodePoint(parseInt(badge.icon!, 16))}
                      </span>
                      <p class="text-sm font-semibold text-white">
                        {badge.text}
                      </p>
                    </Motion.div>
                  )}
                </For>
              </Motion.div>
            </Motion.div>
          </Motion.div>
        </Match>
      </Switch>
    </Suspense>
  )
}

const ActionButton = (props: { text: string; onClick: () => void }) => {
  return (
    <Motion.button
      class="absolute right-4 top-4 z-10 h-9 w-[90px] rounded-xl bg-black/80 px-2.5 text-white"
      onClick={() => props.onClick()}
      press={{ scale: 0.95 }}
    >
      {props.text}
    </Motion.button>
  )
}

const AnimatedLoader = () => {
  return (
    <div class="flex h-screen flex-col items-start justify-start bg-secondary">
      <Motion.div
        class="bg-main h-[260px] w-full"
        initial={{ opacity: 0 }}
        animate={{ opacity: [0.3, 0.6, 0.3] }}
        transition={{ duration: 1.5, repeat: Infinity }}
      />
      <div class="flex flex-col items-start justify-start p-4">
        <Motion.div
          class="bg-main h-36 w-full rounded"
          initial={{ opacity: 0 }}
          animate={{ opacity: [0.3, 0.6, 0.3] }}
          transition={{ duration: 1.5, repeat: Infinity, delay: 0.1 }}
        />
        <div class="mt-4 flex w-full flex-row flex-wrap items-center justify-start gap-2">
          <For each={[40, 32, 36, 24, 40, 28, 32]}>
            {(width, index) => (
              <Motion.div
                class="bg-main h-10 rounded-2xl"
                style={{ width: `${width * 4}px` }}
                initial={{ opacity: 0 }}
                animate={{ opacity: [0.3, 0.6, 0.3] }}
                transition={{
                  duration: 1.5,
                  repeat: Infinity,
                  delay: 0.2 + index() * 0.05,
                }}
              />
            )}
          </For>
        </div>
      </div>
    </div>
  )
}
