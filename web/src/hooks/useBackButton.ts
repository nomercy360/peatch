// Unified hook to interact with both BackButton and MainButton
export function useButtons() {
  const backButton = {
    setVisible: (visible: boolean) => {
      if (visible) {
        window.Telegram.WebApp.BackButton.show();
      } else {
        window.Telegram.WebApp.BackButton.isVisible = false;
      }
    },
    onClick: (callback: () => void) => {
      window.Telegram.WebApp.BackButton.onClick(callback);
    },
    offClick: (callback: () => void) => {
      window.Telegram.WebApp.BackButton.offClick(callback);
    },
  };

  const mainButton = {
    setText: (text: string) => {
      window.Telegram.WebApp.MainButton.setText(text);
    },
    setColor: (color: string) => {
      window.Telegram.WebApp.MainButton.setParams({ color });
    },
    setTextColor: (textColor: string) => {
      window.Telegram.WebApp.MainButton.setParams({ text_color: textColor });
    },
    setVisible: (visible: boolean) => {
      if (visible) {
        window.Telegram.WebApp.MainButton.show();
      } else {
        window.Telegram.WebApp.MainButton.isVisible = false;
      }
    },
    setActive: (active: boolean) => {
      if (active) {
        window.Telegram.WebApp.MainButton.enable();
      } else {
        window.Telegram.WebApp.MainButton.disable();
      }
    },
    onClick: (callback: () => void) => {
      window.Telegram.WebApp.MainButton.onClick(callback);
    },
    offClick: (callback: () => void) => {
      window.Telegram.WebApp.MainButton.offClick(callback);
    },
    showProgress: (leaveActive = false) => {
      window.Telegram.WebApp.MainButton.showProgress(leaveActive);
    },
    hideProgress: () => {
      window.Telegram.WebApp.MainButton.hideProgress();
    },
  };

  return { backButton, mainButton };
}
