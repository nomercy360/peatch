// Unified hook to interact with both BackButton and MainButton
export function useButtons() {
  const backButton = {
    setVisible: () => {
      window.Telegram.WebApp.BackButton.show();
    },
    hide() {
      window.Telegram.WebApp.BackButton.isVisible = false;
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
    setVisible: (text: string) => {
      window.Telegram.WebApp.MainButton.setParams({
        is_visible: true,
        text_color: '#FFFFFF',
        color: '#3F8AF7',
        text,
      });
    },
    hide: () => {
      window.Telegram.WebApp.MainButton.isVisible = false;
    },
    setActive: (active: boolean) => {
      if (active) {
        //window.Telegram.WebApp.MainButton.enable();
        window.Telegram.WebApp.MainButton.setParams({
          color: '#3F8AF7',
          is_active: true,
          text_color: '#FFFFFF',
        });
      } else {
        //window.Telegram.WebApp.MainButton.disable();
        window.Telegram.WebApp.MainButton.setParams({
          color: '#BEDDFC',
          is_active: false,
          text_color: '#FFFFFF',
        });
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
