import { Telegram } from './telegram'

export {}

declare global {
	interface Window {
		Telegram: Telegram
	}
}

export interface UserResponse {
	id: string
	first_name: string
	last_name: string
	date_of_birth: string
	username: string
	telegram_id: number
	verification_status: string
	photos: Photo[]
	location?: {
		name: string
		country_name: string
	}
}

export interface Photo {
	id: string
	photo_url: string
	order_index: number
	verification_status: 'pending' | 'approved' | 'denied'
}
