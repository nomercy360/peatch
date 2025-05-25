import { setEditUser, store } from '~/store'
import { RouteSectionProps, useNavigate } from '@solidjs/router'

export default function EditUser(props: RouteSectionProps) {
  const navigate = useNavigate()

  setEditUser({
    name: store.user.name || '',
    title: store.user.title || '',
    description: store.user.description || '',
    avatar_url: store.user.avatar_url || '',
    location: store.user.location || {},
    // @ts-ignore
    badge_ids: store.user.badges?.map(b => b.id) || [],
    // @ts-ignore
    opportunity_ids: store.user.opportunities?.map(o => o.id) || [],
  })

  // if the first name or last name or title is not set, redirect to the first step
  if (!store.user.name || !store.user.title) {
    navigate('/users/edit')
  }

  return <div>{props.children}</div>
}
