import { createStore } from 'solid-js/store'
import { createSignal } from 'solid-js'
import { CityResponse } from '~/gen'
import { UserResponse } from '~/gen'

export const [store, setStore] = createStore<{
	user: UserResponse
	token: string
	following: number[]
}>({
	user: null as any,
	token: null as any,
	following: [],
})

export const setUser = (user: UserResponse) => setStore('user', user)

export const setToken = (token: string) => setStore('token', token)

export const [editUser, setEditUser] = createStore<{
	first_name: string,
	last_name: string,
	title: string,
	description: string,
	location: CityResponse
	badge_ids: string[]
	opportunity_ids: string[]
}>({
	first_name: '',
	last_name: '',
	title: '',
	description: '',
	location: {},
	badge_ids: [],
	opportunity_ids: [],
})

export const [editCollaboration, setEditCollaboration] =
	createStore<{
		badge_ids: string[]
		location: CityResponse
		description: string
		is_payable: boolean
		opportunity_id: string
		title: string
	}>({
		badge_ids: [],
		location: {},
		description: '',
		is_payable: false,
		opportunity_id: '',
		title: '',
	})


export const [editCollaborationId, setEditCollaborationId] =
	createSignal<number>(0)

