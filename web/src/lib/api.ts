import { store } from '~/store'
import { CreateCollaboration, UpdateUserRequest } from '~/gen'

export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL as string
export const CDN_URL = 'https://assets.peatch.io'

export const apiFetch = async ({
																 endpoint,
																 method = 'GET',
																 body = null,
																 showProgress = true,
																 responseContentType = 'json' as 'json' | 'blob',
															 }: {
	endpoint: string
	method?: 'GET' | 'POST' | 'PUT' | 'DELETE'
	body?: any
	showProgress?: boolean
	responseContentType?: string
}) => {
	const headers: { [key: string]: string } = {
		'Content-Type': 'application/json',
		Authorization: `Bearer ${store.token}`,
	}

	try {
		showProgress && window.Telegram.WebApp.MainButton.showProgress(false)

		const response = await fetch(`${API_BASE_URL}/api${endpoint}`, {
			method,
			headers,
			body: body ? JSON.stringify(body) : undefined,
		})

		if (!response.ok) {
			const errorResponse = await response.json()
			throw { code: response.status, message: errorResponse.message }
		}

		switch (response.status) {
			case 204:
				return true
			default:
				return response[responseContentType as 'json' | 'blob']()
		}
	} finally {
		showProgress && window.Telegram.WebApp.MainButton.hideProgress()
	}
}

export const fetchUsers = async ({ pageParam = 1, queryKey }: any) => {
	const [_, search] = queryKey
	const response = await apiFetch({
		endpoint: `/users?search=${search}&page=${pageParam}&limit=20`,
	})

	return {
		data: response,
		nextPage: response.length === 20 ? pageParam + 1 : undefined,
	}
}

export const fetchBadges = async () => {
	return await apiFetch({
		endpoint: '/badges',
		showProgress: false,
	})
}

export const postBadge = async (text: string, color: string, icon: string) => {
	return await apiFetch({
		endpoint: '/badges',
		method: 'POST',
		body: { text, color, icon },
	})
}

export const fetchOpportunities = async () => {
	return await apiFetch({ endpoint: '/opportunities', showProgress: false })
}

export const updateUser = async (user: UpdateUserRequest) => {
	return await apiFetch({
		endpoint: '/users',
		method: 'PUT',
		body: user,
		showProgress: false,
	})
}

export const fetchProfile = async (id: string) => {
	const endpoint = `/users/${id}`

	return await apiFetch({ endpoint })
}


export const followUser = async (userID: string) => {
	const resp = await apiFetch({
		endpoint: `/users/${userID}/follow`,
		method: 'POST',
		showProgress: false,
	})

	if (resp?.status && resp.status.includes('bot_blocked')) {
		throw {
			botBlocked: true,
			username: resp.username,
		}
	}

	return resp
}

export const createCollaboration = async (collaboration: any) => {
	return await apiFetch({
		endpoint: '/collaborations',
		method: 'POST',
		body: collaboration,
	})
}

export const updateCollaboration = async (
	id: number,
	collaboration: CreateCollaboration,
) => {
	return await apiFetch({
		endpoint: '/collaborations/' + id,
		method: 'PUT',
		body: collaboration,
	})
}

export const fetchCollaborations = async (search: any) => {
	return await apiFetch({ endpoint: '/collaborations?search=' + search })
}

export const fetchCollaboration = async (collaborationID: string) => {
	return await apiFetch({ endpoint: `/collaborations/${collaborationID}` })
}

export const searchLocations = async (search: string) => {
	return await apiFetch({ endpoint: `/locations?search=${search}` })
}

export const uploadUserAvatar = async (file: File) => {
	const formData = new FormData()
	formData.append('photo', file)

	const response = await fetch(`${API_BASE_URL}/api/users/avatar`, {
		method: 'POST',
		headers: {
			Authorization: `Bearer ${store.token}`,
		},
		body: formData,
	})

	if (!response.ok) {
		const errorResponse = await response.json()
		throw { code: response.status, message: errorResponse.message }
	}

	return await response.json()
}
