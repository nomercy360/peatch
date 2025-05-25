import { SelectBadge } from '~/components/edit/selectBadge'
import { FormLayout } from '~/components/edit/layout'
import { useMainButton } from '~/lib/useMainButton'
import { useNavigate, useSearchParams } from '@solidjs/router'
import { createEffect, createSignal, onCleanup } from 'solid-js'
import { editUser, setEditUser } from '~/store'
import { fetchBadges } from '~/lib/api'
import { useTranslations } from '~/lib/locale-context'
import { useQuery } from '@tanstack/solid-query'
import { BadgeResponse } from '~/gen'

export default function SelectBadges() {
  const mainButton = useMainButton()
  const { t } = useTranslations()

  const [badgeSearch, setBadgeSearch] = createSignal('')

  const navigate = useNavigate()
  const [searchParams] = useSearchParams()

  const navigateNext = () => {
    navigate('/users/edit/interests', { state: { back: true } })
  }

  const navigateCreateBadge = () => {
    navigate('/users/edit/create-badge?badge_name=' + badgeSearch(), {
      state: { back: true },
    })
  }

  const fetchBadgeQuery = useQuery(() => ({
    queryKey: ['badges'],
    // then push selected to the top
    queryFn: () =>
      fetchBadges().then(badges => {
        const selected = editUser.badge_ids
        return [
          ...selected.map((id: string) =>
            badges.find((b: BadgeResponse) => b.id === id),
          ),
          ...badges.filter((b: BadgeResponse) => !selected.includes(b.id!)),
        ]
      }),
  }))

  createEffect(async () => {
    if (searchParams.refetch) {
      await fetchBadgeQuery.refetch()
    }
  })

  mainButton.onClick(navigateNext)

  createEffect(() => {
    if (editUser.badge_ids.length) {
      mainButton.enable(t('common.buttons.next'))
    } else {
      mainButton.disable(t('common.buttons.next'))
    }
  })

  onCleanup(() => {
    mainButton.offClick(navigateNext)
  })

  return (
    <FormLayout
      title={t('pages.users.edit.badges.title')}
      description={t('pages.users.edit.badges.description')}
      screen={2}
      totalScreens={6}
    >
      <SelectBadge
        selected={editUser.badge_ids}
        setSelected={b => setEditUser('badge_ids', b)}
        onCreateBadgeButtonClick={navigateCreateBadge}
        search={badgeSearch()}
        setSearch={setBadgeSearch}
        badges={fetchBadgeQuery.data!}
      />
    </FormLayout>
  )
}
