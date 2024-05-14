import { store } from '~/store';
import { CreateCollaboration } from '../../gen';

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

export const fetchUsers = async (search: string) => {
	return await apiFetch({
		endpoint: '/users?search=' + search + '&page=1&limit=120',
	})
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

export const updateUser = async (user: any) => {
	return await apiFetch({
		endpoint: '/users',
		method: 'PUT',
		body: user,
		showProgress: false,
	})
}

export const uploadToS3 = (
	url: string,
	file: File,
	onProgress: (e: ProgressEvent) => void,
	onFinished: () => void,
): Promise<void> => {
	return new Promise<void>((resolve, reject) => {
		const req = new XMLHttpRequest()
		req.onreadystatechange = () => {
			if (req.readyState === 4) {
				if (req.status === 200) {
					onFinished()
					resolve()
				} else {
					reject(new Error('Failed to upload file'))
				}
			}
		}
		req.upload.addEventListener('progress', onProgress)
		req.open('PUT', url)
		req.send(file)
	})
}

export const fetchPresignedUrl = async (file: string) => {
	const { path, url } = await apiFetch({
		endpoint: `/presigned-url?filename=${file}`,
		showProgress: false,
	})

	return { path, url }
}

export const fetchProfile = async (username: string) => {
	return await apiFetch({ endpoint: `/users/${username}` })
}

export const followUser = async (userID: number) => {
	return await apiFetch({
		endpoint: `/users/${userID}/follow`,
		method: 'POST',
		showProgress: false,
	})
}

export const unfollowUser = async (userID: number) => {
	return await apiFetch({
		endpoint: `/users/${userID}/follow`,
		method: 'DELETE',
		showProgress: false,
	})
}

export const hideProfile = async () => {
	return await apiFetch({
		endpoint: '/users/hide',
		method: 'POST',
		showProgress: false,
	})
}

export const showProfile = async () => {
	return await apiFetch({
		endpoint: '/users/show',
		method: 'POST',
		showProgress: false,
	})
}

export const publishProfile = async () => {
	return await apiFetch({
		endpoint: '/users/publish',
		method: 'POST',
	})
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

export const createUserCollaboration = async (
	receiverID?: number,
	message?: string,
) => {
	return await apiFetch({
		endpoint: `/users/${receiverID}/collaborations/requests`,
		method: 'POST',
		body: { message },
	})
}

export const fetchPreview = async () => {
	return await apiFetch({ endpoint: '/user-preview' }).then((res: any) =>
		res.map(
			(image: { avatar_url: string }) => CDN_URL + '/' + image.avatar_url,
		),
	)
}

export const publishCollaboration = async (collaborationID: number) => {
	return await apiFetch({
		endpoint: `/collaborations/${collaborationID}/publish`,
		method: 'POST',
	})
}

export const hideCollaboration = async (collaborationID: number) => {
	return await apiFetch({
		endpoint: `/collaborations/${collaborationID}/hide`,
		method: 'POST',
	})
}

export const showCollaboration = async (collaborationID: number) => {
	return await apiFetch({
		endpoint: `/collaborations/${collaborationID}/show`,
		method: 'POST',
	})
}

export const fetchCollaboration = async (collaborationID: number) => {
	return await apiFetch({ endpoint: `/collaborations/${collaborationID}` })
}

export const findUserCollaborationRequest = async (handle: string) => {
	return await apiFetch({
		endpoint: `/users/${handle}/collaborations/requests`,
	})
}

export const findCollaborationRequest = async (collaborationID: number) => {
	return await apiFetch({
		endpoint: `/collaborations/${collaborationID}/requests`,
	})
}

export const createCollaborationRequest = async (
	collaborationID: number,
	message: string,
) => {
	return await apiFetch({
		endpoint: `/collaborations/${collaborationID}/requests`,
		method: 'POST',
		body: { message },
	})
}

export const searchLocations = async (search: string) => {
	return await apiFetch({ endpoint: `/locations?search=${search}` })
}

export const fetchFollowing = async (id: number) => {
  return await apiFetch({ endpoint: `/users/${id}/following` });
};

export const fetchFollowers = async (id: number) => {
  return await apiFetch({ endpoint: `/users/${id}/followers` });
};
