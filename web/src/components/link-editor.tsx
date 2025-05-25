import { createEffect, createSignal, For, Show, onCleanup } from 'solid-js'
import { Motion, Presence } from 'solid-motionone'
import { Link } from '~/gen'
import { fetchLinkMetadata, updateUserLinks } from '~/lib/api'
import { detectLinkType, isValidUrl } from '~/lib/utils'
import { addToast } from '~/components/toast'
import { useTranslations } from '~/lib/locale-context'
import { useMutation, useQuery } from '@tanstack/solid-query'
import { useMainButton } from '~/lib/useMainButton'

interface LinkEditorProps {
  links: Link[]
  isCurrentUser: boolean
  onDrawerStateChange?: (isOpen: boolean) => void
}

export default function LinkEditor(props: LinkEditorProps) {
  const { t } = useTranslations()
  const mainButton = useMainButton()
  const [links, setLinks] = createSignal<Link[]>(props.links || [])
  const [editingLinks, setEditingLinks] = createSignal(false)
  const [newLinkUrl, setNewLinkUrl] = createSignal('')
  const [newLinkLabel, setNewLinkLabel] = createSignal('')
  const [editingLinkIndex, setEditingLinkIndex] = createSignal<number | null>(
    null,
  )
  const [drawerOpen, setDrawerOpen] = createSignal(false)
  const [linksExpanded, setLinksExpanded] = createSignal(false)

  createEffect(() => {
    setLinks(props.links || [])
  })

  const updateLinksMutation = useMutation(() => ({
    mutationFn: (links: Link[]) => updateUserLinks(links),
    onSuccess: () => {
      addToast(t('pages.users.linksUpdated'), 'success')
    },
    onError: () => {
      addToast(t('pages.users.linksUpdateError'), 'error')
    },
  }))

  const metadataQuery = useQuery(() => ({
    queryKey: ['linkMetadata', newLinkUrl()],
    queryFn: async () => {
      const url = newLinkUrl()
      if (!url || !isValidUrl(url)) return null
      return await fetchLinkMetadata(url)
    },
    enabled: !!newLinkUrl() && isValidUrl(newLinkUrl()),
  }))

  createEffect(() => {
    const metadata = metadataQuery.data
    if (metadata && metadata.title && !newLinkLabel()) {
      setNewLinkLabel(metadata.title)
    }
  })

  const handleSaveLink = () => {
    if (!isValidUrl(newLinkUrl())) {
      addToast(t('pages.users.invalidUrl'), 'error')
      return
    }

    const linkInfo = detectLinkType(newLinkUrl())

    if (editingLinkIndex() !== null) {
      // Update existing link
      const updatedLinks = [...links()]
      updatedLinks[editingLinkIndex()!] = {
        url: newLinkUrl(),
        label: newLinkLabel() || newLinkUrl(),
        type: linkInfo.type,
        order: editingLinkIndex()!,
      }
      setLinks(updatedLinks)

      // Save to server
      const orderedLinks = updatedLinks.map((link, index) => ({
        ...link,
        order: index,
      }))
      updateLinksMutation.mutate(orderedLinks)
    } else {
      // Add new link
      const newLink: Link = {
        url: newLinkUrl(),
        label: newLinkLabel() || newLinkUrl(),
        type: linkInfo.type,
        order: links().length,
      }
      const updatedLinks = [...links(), newLink]
      setLinks(updatedLinks)

      // Save to server
      const orderedLinks = updatedLinks.map((link, index) => ({
        ...link,
        order: index,
      }))
      updateLinksMutation.mutate(orderedLinks)
    }

    closeDrawer()
  }

  const openDrawerForEdit = (index: number) => {
    const link = links()[index]
    setNewLinkUrl(link.url || '')
    setNewLinkLabel(link.label || '')
    setEditingLinkIndex(index)
    setDrawerOpen(true)
  }

  const openDrawerForAdd = () => {
    setNewLinkUrl('')
    setNewLinkLabel('')
    setEditingLinkIndex(null)
    setDrawerOpen(true)
  }

  const closeDrawer = () => {
    setDrawerOpen(false)
    setNewLinkUrl('')
    setNewLinkLabel('')
    setEditingLinkIndex(null)
  }

  const handleRemoveLink = (index: number) => {
    const updatedLinks = links().filter((_, i) => i !== index)
    setLinks(updatedLinks)

    // Auto-save after removal
    const orderedLinks = updatedLinks.map((link, idx) => ({
      ...link,
      order: idx,
    }))
    updateLinksMutation.mutate(orderedLinks)
  }

  // Handle drawer state changes and main button
  createEffect(() => {
    if (drawerOpen()) {
      props.onDrawerStateChange?.(true)
      mainButton.enable(t('common.buttons.save'))
      mainButton.onClick(handleSaveLink)
    } else {
      props.onDrawerStateChange?.(false)
      mainButton.offClick(handleSaveLink)
    }
  })

  onCleanup(() => {
    mainButton.offClick(handleSaveLink)
  })

  return (
    <>
      <Show when={links().length > 0 || props.isCurrentUser}>
        <Motion.div
          class="mt-4 w-full"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.3, delay: 0.55 }}
        >
          <div class="mb-2 flex items-center justify-between">
            <Motion.p
              class="text-start text-xl font-extrabold"
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.3, delay: 0.6 }}
            >
              {t('pages.users.links')}
            </Motion.p>
            <Show when={props.isCurrentUser}>
              <Motion.button
                onClick={() => setEditingLinks(!editingLinks())}
                class="text-sm font-medium text-primary"
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ duration: 0.3, delay: 0.65 }}
                press={{ scale: 0.95 }}
              >
                {editingLinks()
                  ? t('common.buttons.cancel')
                  : t('common.buttons.edit')}
              </Motion.button>
            </Show>
          </div>
          <Motion.div
            class="flex w-full flex-col items-center justify-start gap-1"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ duration: 0.3, delay: 0.7 }}
          >
            <For each={linksExpanded() ? links() : links().slice(0, 3)}>
              {(link, index) => {
                const linkInfo = detectLinkType(link.url || '')
                return (
                  <Motion.div
                    class="flex h-12 w-full flex-row items-center justify-between rounded-xl bg-secondary px-3"
                    initial={{ opacity: 0, x: -20 }}
                    animate={{ opacity: 1, x: 0 }}
                    transition={{
                      duration: 0.3,
                      delay: 0.1 + index() * 0.05,
                    }}
                  >
                    <a
                      href={link.url}
                      target="_blank"
                      rel="noopener noreferrer"
                      class="flex flex-1 items-center gap-2 text-secondary-foreground"
                    >
                      <span class="material-symbols-rounded text-[20px]">
                        {linkInfo.icon}
                      </span>
                      <span class="truncate text-sm font-medium">
                        {link.label}
                      </span>
                    </a>
                    <Show when={editingLinks()}>
                      <div class="flex gap-1">
                        <Motion.button
                          onClick={() => openDrawerForEdit(index())}
                          class="flex size-7 items-center justify-center rounded-lg hover:bg-background"
                          press={{ scale: 0.9 }}
                        >
                          <span class="material-symbols-rounded text-[16px]">
                            edit
                          </span>
                        </Motion.button>
                        <Motion.button
                          onClick={() => handleRemoveLink(index())}
                          class="flex size-7 items-center justify-center rounded-lg hover:bg-background"
                          press={{ scale: 0.9 }}
                        >
                          <span class="material-symbols-rounded text-[16px]">
                            close
                          </span>
                        </Motion.button>
                      </div>
                    </Show>
                  </Motion.div>
                )
              }}
            </For>
            <Show when={links().length > 3}>
              <Motion.div
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ duration: 0.3, delay: 0.3 }}
              >
                <ExpandButton
                  expanded={linksExpanded()}
                  setExpanded={setLinksExpanded}
                />
              </Motion.div>
            </Show>
            <Show when={editingLinks()}>
              <Motion.button
                onClick={openDrawerForAdd}
                class="border-secondary-foreground/30 flex h-12 w-full items-center justify-center gap-2 rounded-xl border-2 border-dashed text-secondary-foreground"
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ duration: 0.3 }}
                press={{ scale: 0.98 }}
              >
                <span class="material-symbols-rounded text-[20px]">add</span>
                <span class="text-sm font-medium">
                  {t('pages.users.addLink')}
                </span>
              </Motion.button>
            </Show>
          </Motion.div>
        </Motion.div>
      </Show>

      {/* Bottom Drawer */}
      <Presence>
        <Show when={drawerOpen()}>
          <Motion.div
            class="fixed inset-0 z-50 flex items-end"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.2 }}
            onClick={closeDrawer}
          >
            <Motion.div
              class="absolute inset-0 bg-black/50"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
            />
            <Motion.div
              class="relative w-full rounded-t-3xl bg-background p-6"
              style={{ height: '50vh' }}
              initial={{ y: '100%' }}
              animate={{ y: 0 }}
              exit={{ y: '100%' }}
              transition={{ duration: 0.3, easing: 'ease-out' }}
              onClick={e => e.stopPropagation()}
            >
              <button
                onClick={closeDrawer}
                class="absolute right-4 top-4 flex size-8 items-center justify-center rounded-lg bg-secondary"
              >
                <span class="material-symbols-rounded text-[20px]">close</span>
              </button>

              <h2 class="mb-6 text-start text-xl font-bold">
                {editingLinkIndex() !== null
                  ? t('pages.users.editLink')
                  : t('pages.users.addLink')}
              </h2>

              <div class="flex w-full flex-col gap-4">
                <div>
                  <label class="mb-2 block text-start text-sm font-medium">
                    {t('pages.users.linkUrl')}
                  </label>
                  <input
                    type="url"
                    placeholder={t('pages.users.linkUrlPlaceholder')}
                    value={newLinkUrl()}
                    onInput={e => setNewLinkUrl(e.currentTarget.value)}
                    class="h-12 w-full rounded-xl border border-secondary bg-background px-4 text-sm outline-none focus:border-primary"
                  />
                </div>

                <div>
                  <label class="mb-2 block text-start text-sm font-medium">
                    {t('pages.users.linkLabel')}
                  </label>
                  <input
                    type="text"
                    placeholder={t('pages.users.linkLabelPlaceholder')}
                    value={newLinkLabel()}
                    onInput={e => setNewLinkLabel(e.currentTarget.value)}
                    class="h-12 w-full rounded-xl border border-secondary bg-background px-4 text-sm outline-none focus:border-primary"
                  />
                  <Show when={metadataQuery.isLoading}>
                    <p class="mt-2 text-xs text-secondary-foreground">
                      {t('pages.users.fetchingMetadata')}
                    </p>
                  </Show>
                </div>
              </div>
            </Motion.div>
          </Motion.div>
        </Show>
      </Presence>
    </>
  )
}

const ExpandButton = (props: {
  expanded: boolean
  setExpanded: (val: boolean) => void
}) => {
  const { t } = useTranslations()
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
      {props.expanded ? t('pages.users.showLess') : t('pages.users.showMore')}
    </Motion.button>
  )
}
