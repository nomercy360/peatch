import { MainButton } from '~/types/telegram';

export function useMainButton() {
  return {
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
    enable: () => {
      window.Telegram.WebApp.MainButton.enable();
    },
    disable: () => {
      window.Telegram.WebApp.MainButton.disable();
    },
    setParams: (params: {
      text?: string;
      isVisible?: boolean;
      color?: string;
      textColor?: string;
      isEnabled?: boolean;
    }) => {
      return window.Telegram.WebApp.MainButton.setParams({
        is_visible: params.isVisible,
        text: params.text,
        color: params.color,
        text_color: params.textColor,
        is_active: params.isEnabled,
      });
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
}
