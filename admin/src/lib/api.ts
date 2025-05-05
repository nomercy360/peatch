import { store } from '~/store'
import { VerificationStatus } from '~/gen'

export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL as string

export async function apiRequest(endpoint: string, options: RequestInit = {}) {
	try {
		const response = await fetch(`${API_BASE_URL}${endpoint}`, {
			...options,
			headers: {
				'Content-Type': 'application/json',
				Authorization: `Bearer ${store.token}`,
				...(options.headers || {}),
			},
		})

		let data
		try {
			data = await response.json()
		} catch {
			throw new Error('Failed to get response from server')
		}

		if (!response.ok) {
			const errorMessage = Array.isArray(data?.error)
				? data.error.join('\n')
				: typeof data?.error === 'string'
					? data.error
					: 'An error occurred'

			throw new Error(errorMessage)
		}

		return data
	} catch (error) {
		const errorMessage = error instanceof Error ? error.message : 'An unexpected error occurred'
		throw new Error(errorMessage)
	}
}

export const fetchUsers = async ({ pageParam = 1, status = 'pending' }: { pageParam?: number, status?: string }) => {
	const response = await apiRequest(`/admin/users?page=${pageParam}&limit=20&status=${status}`)

	return {
		data: response,
		nextPage: response.length === 20 ? pageParam + 1 : undefined,
	}
}

export const fetchCollaborations = async ({ pageParam = 1, status = 'pending' }: { pageParam?: number, status?: string }) => {
	const response = await apiRequest(`/admin/collaborations?page=${pageParam}&limit=20&status=${status}`)

	return {
		data: response,
		nextPage: response.length === 20 ? pageParam + 1 : undefined,
	}
}

export const updateUserStatus = async (userId: string, status: VerificationStatus) => {
	return await apiRequest(`/admin/users/${userId}/verify`, {
		method: 'PUT',
		body: JSON.stringify({ status }),
	})
}

export const updateCollaborationStatus = async (userId: string, collaborationId: string, status: VerificationStatus) => {
	return await apiRequest(`/admin/users/${userId}/collaborations/${collaborationId}/verify`, {
		method: 'PUT',
		body: JSON.stringify({ status }),
	})
}
