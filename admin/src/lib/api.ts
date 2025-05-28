import { UpdateUserRequest, VerificationStatus } from '~/gen'

export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL as string

export function getToken(): string | null {
	return localStorage.getItem('admin-token')
}

export function setToken(token: string) {
	localStorage.setItem('admin-token', token)
}

export function clearToken() {
	localStorage.removeItem('admin-token')
}

export async function apiRequest(endpoint: string, options: RequestInit = {}) {
	const token = getToken()

	try {
		const response = await fetch(`${API_BASE_URL}${endpoint}`, {
			...options,
			headers: {
				'Content-Type': 'application/json',
				...(token ? { 'x-api-token': token } : {}),
				...(options.headers || {}),
			},
		})

		if (response.status === 401) {
			clearToken()
			// Only redirect if we're not already on the login page
			if (window.location.pathname !== '/login') {
				window.location.href = '/login'
			}
			throw new Error('Unauthorized')
		}

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

export async function checkAuth(): Promise<boolean> {
	try {
		await apiRequest('/admin/me')
		return true
	} catch {
		return false
	}
}

export async function login(token: string): Promise<boolean> {
	setToken(token)
	const isValid = await checkAuth()
	if (!isValid) {
		clearToken()
	}
	return isValid
}

export const fetchUsers = async ({ page = 0, limit = 10, status = 'verified' }: {
	page?: number,
	limit?: number,
	status?: string
}) => {
	const params = new URLSearchParams()
	params.append('page', (page + 1).toString()) // API uses 1-based pagination
	params.append('per_page', limit.toString())

	// Only add status param if it's provided and not empty
	if (status) {
		params.append('status', status)
	}

	const users = await apiRequest(`/admin/users?${params.toString()}`)

	// The API returns an array of users, not a paginated response
	return {
		users,
		total: users.length === limit ? (page + 1) * limit + 1 : page * limit + users.length,
	}
}

export const fetchCollaborations = async ({ page = 0, limit = 10, status = '' }: {
	page?: number,
	limit?: number,
	status?: string
}) => {
	const params = new URLSearchParams()
	params.append('page', (page + 1).toString()) // API uses 1-based pagination
	params.append('per_page', limit.toString())
	
	// Only add status param if it's provided and not empty
	if (status) {
		params.append('status', status)
	}
	
	const collaborations = await apiRequest(`/admin/collaborations?${params.toString()}`)

	// The API returns an array of collaborations, similar to users endpoint
	return {
		collaborations,
		total: collaborations.length === limit ? (page + 1) * limit + 1 : page * limit + collaborations.length,
	}
}

export const updateUser = async (userId: any, data: UpdateUserRequest) => {
	return await apiRequest(`/api/users`, {
		method: 'PUT',
		body: JSON.stringify({ ...data, id: userId }),
	})
}

export const updateUserStatus = async (userId: any, status: VerificationStatus) => {
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

export const deleteUser = async (userId: string) => {
	return await apiRequest(`/admin/users/${userId}`, {
		method: 'DELETE',
	})
}

export const deleteUsers = async (userIds: string[]) => {
	// For now, delete users one by one
	// TODO: Add batch delete endpoint when available
	const promises = userIds.map(id => deleteUser(id))
	return Promise.all(promises)
}

export const deleteCollaboration = async (collaborationId: string) => {
	return await apiRequest(`/admin/collaborations/${collaborationId}`, {
		method: 'DELETE',
	})
}

export const deleteCollaborations = async (collaborationIds: string[]) => {
	// For now, delete collaborations one by one
	// TODO: Add batch delete endpoint when available
	const promises = collaborationIds.map(id => deleteCollaboration(id))
	return Promise.all(promises)
}
