import { User } from '../gen';
import { createStore } from 'solid-js/store';

export const [store, setStore] = createStore<{
  user: User;
  token: string;
  following: number[];
}>({
  user: null as any,
  token: null as any,
  following: [],
});

export const setUser = (user: User) => setStore('user', user);

export const setToken = (token: string) => setStore('token', token);

export const setFollowing = (following: number[]) =>
  setStore('following', following);
