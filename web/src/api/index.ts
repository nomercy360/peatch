import { store } from '../store';

export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL as string;
export const CDN_URL = 'https://assets.peatch.io';

export const apiFetch = async ({
                                 endpoint,
                                 method = 'GET',
                                 body = null,
                                 responseType = 'json',
                                 catchError = true,
                                 showProgress = true,
                               }: {
  endpoint: string;
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE';
  body?: any;
  responseType?: 'json' | 'blob';
  catchError?: boolean;
  showProgress?: boolean;
}) => {
  const headers: { [key: string]: string } = {
    'Content-Type': 'application/json',
    Authorization: `Bearer ${store.token}`,
  };

  try {
    if (showProgress) {
      window.Telegram.WebApp.MainButton.showProgress(false);
    }
    const response = await fetch(`${API_BASE_URL}/api${endpoint}`, {
      method,
      headers,
      body: body ? JSON.stringify(body) : undefined,
    });

    if (!response.ok) {
      if (catchError) {
        const { message } = await response.json();
        console.error('Error:', message || 'Error! Please try again');
        window.Telegram.WebApp.BackButton.isVisible = false;
      }
      return false;
    }

    if (response.status === 204) {
      return true;
    } else {
      return responseType === 'json' ? response.json() : response.blob();
    }
  } catch (error) {
    console.error('Error:', error);
  } finally {
    if (showProgress) {
      window.Telegram.WebApp.MainButton.hideProgress();
    }
  }
};

export const fetchUsers = async (search: any) => {
  return await apiFetch({ endpoint: '/users?search=' + search });
};

export const fetchBadges = async () => {
  return await apiFetch({ endpoint: '/badges' });
};

export const postBadge = async (text: string, color: string, icon: string) => {
  return await apiFetch({
    endpoint: '/badges',
    method: 'POST',
    body: { text, color, icon },
  });
};

export const fetchOpportunities = async () => {
  return await apiFetch({ endpoint: '/opportunities' });
};

export const updateUser = async (user: any) => {
  return await apiFetch({
    endpoint: `/users`,
    method: 'PUT',
    body: user,
  });
};

export const uploadToS3 = (
  url: string,
  file: File,
  onProgress: (e: ProgressEvent) => void,
  onFinished: () => void,
): Promise<void> => {
  return new Promise<void>((resolve, reject) => {
    const req = new XMLHttpRequest();
    req.onreadystatechange = () => {
      if (req.readyState === 4) {
        if (req.status === 200) {
          onFinished();
          resolve();
        } else {
          reject(new Error('Failed to upload file'));
        }
      }
    };
    req.upload.addEventListener('progress', onProgress);
    req.open('PUT', url);
    req.send(file);
  });
};

export const fetchPresignedUrl = async (file: string) => {
  const { path, url } = await apiFetch({
    endpoint: `/presigned-url?filename=${file}`,
  });

  return { path, url };
};

export const fetchProfile = async (userID: number) => {
  return await apiFetch({ endpoint: `/users/${userID}` });
};

export const followUser = async (userID: number) => {
  return await apiFetch({
    endpoint: `/users/${userID}/follow`,
    method: 'POST',
  });
};

export const unfollowUser = async (userID: number) => {
  return await apiFetch({
    endpoint: `/users/${userID}/follow`,
    method: 'DELETE',
  });
};

export const hideProfile = async () => {
  return await apiFetch({
    endpoint: '/users/hide',
    method: 'POST',
  });
};

export const showProfile = async () => {
  return await apiFetch({
    endpoint: '/users/show',
    method: 'POST',
  });
};

export const publishProfile = async () => {
  return await apiFetch({
    endpoint: '/users/publish',
    method: 'POST',
  });
};

export const collaborateUser = async (userID: number) => {
  return await apiFetch({
    endpoint: `/users/${userID}/collaborate`,
    method: 'POST',
  });
};

export const createCollaboration = async (collaboration: any) => {
  return await apiFetch({
    endpoint: '/collaborations',
    method: 'POST',
    body: collaboration,
  });
};

export const updateCollaboration = async (collaboration: any) => {
  return await apiFetch({
    endpoint: '/collaborations',
    method: 'PUT',
    body: collaboration,
  });
};

export const fetchCollaborations = async () => {
  return await apiFetch({ endpoint: '/collaborations' });
};

export const createUserCollaboration = async (collaboration: any) => {
  return await apiFetch({
    endpoint: '/collaborations',
    method: 'POST',
    body: collaboration,
  });
};
