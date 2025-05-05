import { createStore } from 'solid-js/store'


export const [store, setStore] = createStore<{
	token: string
}>({
	token: null as any,
})


export const setToken = (token: string) => setStore('token', token)


