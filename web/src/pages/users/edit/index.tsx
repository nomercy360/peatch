import { setEditUser, store } from '~/store';
import { RouteSectionProps, useNavigate } from '@solidjs/router';

export default function EditUser(props: RouteSectionProps) {
  const navigate = useNavigate();

  setEditUser({
    first_name: store.user.first_name || '',
    last_name: store.user.last_name || '',
    title: store.user.title || '',
    description: store.user.description || '',
    avatar_url: store.user.avatar_url || '',
    city: store.user.city || '',
    country: store.user.country || '',
    country_code: store.user.country_code || '',
    badge_ids: store.user.badges?.map(b => b.id) || ([] as any),
    opportunity_ids: store.user.opportunities?.map(o => o.id) || ([] as any),
  });

  // if first name or last name or title is not set, redirect to the first step
  if (!store.user.first_name || !store.user.last_name || !store.user.title) {
    navigate('/users/edit');
  }

  return <div>{props.children}</div>;
}
