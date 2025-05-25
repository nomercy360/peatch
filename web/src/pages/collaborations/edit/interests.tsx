import { FormLayout } from '~/components/edit/layout'
import { useMainButton } from '~/lib/useMainButton'
import { useNavigate, useParams } from '@solidjs/router'
import { createEffect, onCleanup } from 'solid-js'
import { editCollaboration, setEditCollaboration } from '~/store'
import { fetchOpportunities } from '~/lib/api'
import { SelectOpportunity } from '~/components/edit/select-opp'
import { useQuery } from '@tanstack/solid-query'
import { useTranslations } from '~/lib/locale-context'

export default function SelectOpportunities() {
  const mainButton = useMainButton()

  const navigate = useNavigate()
  const params = useParams()
  const idPath = params.id ? '/' + params.id : ''

  const navigateNext = () => {
    navigate(`/collaborations/edit${idPath}/location`, {
      state: { back: true },
    })
  }

  const query = useQuery(() => ({
    queryKey: ['opportunities'],
    queryFn: () => fetchOpportunities(),
  }))

  mainButton.onClick(navigateNext)

  createEffect(() => {
    if (
      editCollaboration.opportunity_id &&
      editCollaboration.opportunity_id !== ''
    ) {
      mainButton.enable(t('common.buttons.next'))
    } else {
      mainButton.disable(t('common.buttons.next'))
    }
  })

  onCleanup(() => {
    mainButton.offClick(navigateNext)
  })

  const { t } = useTranslations()

  return (
    <FormLayout
      title={t('pages.collaborations.edit.interests.title')}
      description={t('pages.collaborations.edit.interests.description')}
      screen={3}
      totalScreens={6}
    >
      <SelectOpportunity
        selected={editCollaboration.opportunity_id || ''}
        setSelected={b => setEditCollaboration('opportunity_id', b as any)}
        opportunities={query.data}
        loading={query.isLoading}
      />
    </FormLayout>
  )
}
