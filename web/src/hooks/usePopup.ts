export function usePopup() {
  const showAlert = (text: string, callback?: () => any) => {
    window.Telegram.WebApp.showAlert(text, callback);

    return {
      showAlert,
    }
  }

  return {
    showAlert,
  }
}