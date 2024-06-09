export function usePopup() {
	const showAlert = (text: string, callback?: () => any) => {
		window.Telegram.WebApp.showAlert(text, callback)

		return {
			showAlert,
		}
	}

	const showConfirm = (text: string, callback: (ok: boolean) => any) => {
		window.Telegram.WebApp.showConfirm(text, callback)

		return {
			showConfirm,
		}
	}

	return {
		showAlert,
		showConfirm,
	}
}
