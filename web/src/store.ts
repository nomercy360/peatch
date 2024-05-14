import { CreateCollaboration, UpdateUserRequest, User } from '../gen';
import { createStore } from 'solid-js/store';
import { createSignal } from 'solid-js';

export const [store, setStore] = createStore<{
  user: User
  token: string
  following: number[]
}>({
  user: null as any,
  token: null as any,
  following: [],
})

export const setUser = (user: User) => setStore('user', user);

export const setToken = (token: string) => setStore('token', token);

export const setFollowing = (following: number[]) =>
  setStore('following', following)

export const [editUser, setEditUser] = createStore<UpdateUserRequest>({
  first_name: '',
  last_name: '',
  title: '',
  description: '',
  avatar_url: '',
  city: '',
  country: '',
  country_code: '',
  badge_ids: [],
  opportunity_ids: [],
})

export const [editCollaboration, setEditCollaboration] =
  createStore<CreateCollaboration>({
    badge_ids: [],
    city: '',
    country: '',
    country_code: '',
    description: '',
    is_payable: false,
    opportunity_id: 0,
    title: '',
  })

export const [editCollaborationId, setEditCollaborationId] =
  createSignal<number>(0)
