import { createEffect } from 'solid-js';

// Unified hook to interact with both BackButton and MainButton
export function useButtons() {
  // BackButton functions
  const setBackVisibility = (visible: boolean) => {
    if (visible) {
      window.Telegram.WebApp.BackButton.show();
    } else {
      window.Telegram.WebApp.BackButton.isVisible = false;
    }
  };

  const onBackClick = (callback: () => void) => {
    window.Telegram.WebApp.BackButton.onClick(callback);
  };

  const offBackClick = (callback: () => void) => {
    window.Telegram.WebApp.BackButton.offClick(callback);
  };

  // MainButton functions
  const setMainText = (text: string) => {
    window.Telegram.WebApp.MainButton.setText(text);
  };

  const setMainColor = (color: string) => {
    window.Telegram.WebApp.MainButton.setParams({ color });
  };

  const setMainTextColor = (textColor: string) => {
    window.Telegram.WebApp.MainButton.setParams({ text_color: textColor });
  };

  const setMainVisibility = (visible: boolean) => {
    if (visible) {
      window.Telegram.WebApp.MainButton.show();
    } else {
      window.Telegram.WebApp.MainButton.isVisible = false;
    }
  };

  const setMainActive = (active: boolean) => {
    if (active) {
      window.Telegram.WebApp.MainButton.enable();
    } else {
      window.Telegram.WebApp.MainButton.disable();
    }
  };

  const onMainClick = (callback: () => void) => {
    window.Telegram.WebApp.MainButton.onClick(callback);
  };

  const offMainClick = (callback: () => void) => {
    window.Telegram.WebApp.MainButton.offClick(callback);
  };

  const showMainProgress = (leaveActive = false) => {
    window.Telegram.WebApp.MainButton.showProgress(leaveActive);
  };

  const hideMainProgress = () => {
    window.Telegram.WebApp.MainButton.hideProgress();
  };

  return {
    // BackButton controls
    setBackVisibility,
    onBackClick,
    offBackClick,
    // MainButton controls
    setMainText,
    setMainColor,
    setMainTextColor,
    setMainVisibility,
    setMainActive,
    onMainClick,
    offMainClick,
    showMainProgress,
    hideMainProgress,
  };
}
