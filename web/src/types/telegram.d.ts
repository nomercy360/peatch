// CloudStorage Interface Definition
interface CloudStorage {
  setItem(
    key: string,
    value: string,
    callback?: (error: Error | null, success: boolean) => void,
  ): CloudStorage;

  getItem(
    key: string,
    callback: (error: Error | null, value: string) => void,
  ): string;

  getItems(
    keys: string[],
    callback: (error: Error | null, values: string[]) => void,
  ): void;

  removeItem(
    key: string,
    callback?: (error: Error | null, success: boolean) => void,
  ): CloudStorage;

  removeItems(
    keys: string[],
    callback?: (error: Error | null, success: boolean) => void,
  ): CloudStorage;

  getKeys(callback: (error: Error | null, keys: string[]) => void): void;
}

// Telegram Interface Definitions
interface Telegram {
  WebView: WebView;
  Utils: Utils;
  WebApp: WebApp;
}

interface Utils {
  notificationOccurred(warning: string): void;
}

interface WebApp {
  initData: string;
  initDataUnsafe: InitDataUnsafe;
  version: string;
  colorScheme: string;
  themeParams: ThemeParams;
  isExpanded: boolean;
  viewportHeight: number;
  viewportStableHeight: number;
  isClosingConfirmationEnabled: boolean;
  headerColor: string;
  backgroundColor: string;
  BackButton: BackButton;
  MainButton: MainButton;
  HapticFeedback: Utils;
  CloudStorage: CloudStorage;

  openTelegramLink(url: string): void;

  showAlert(message: string, callback?: () => void): void;

  showConfirm(message: string, callback: (ok: boolean) => void): void;

  expand(): void;

  ready(): void;

  onEvent(event: string, callback: () => void): void;

  offEvent(event: string, callback: () => void): void;
}

interface BackButton {
  isVisible: boolean;

  onClick(callback: () => void): void;

  offClick(callback: () => void): void;

  show(): void;

  hide(): void;

  setParams(params: { text_color?: string }): void;
}

export interface MainButton {
  onClick: any;
  text: string;
  color: string;
  offClick: any;
  textColor: string;
  isVisible: boolean;
  isProgressVisible: boolean;
  isActive: boolean;

  setParams(params: {
    text_color?: string;
    color?: string;
    text?: string;
    is_active?: boolean;
    is_visible?: boolean;
  }): MainButton;

  showProgress(leaveActive: boolean): void;

  hideProgress(): void;

  disable(): void;

  setText(nextText: string): void;

  show(): void;

  enable(): void;
}

interface InitDataUnsafe {
  query_id: string;
  user: User;
  auth_date: string;
  hash: string;
}

interface User {
  id: number;
  first_name: string;
  last_name: string;
  username: string;
  language_code: string;
}

interface ThemeParams {
  bg_color: string;
  text_color: string;
  hint_color: string;
  link_color: string;
  button_color: string;
  button_text_color: string;
  secondary_bg_color: string;
  header_bg_color: string;
  accent_text_color: string;
  section_bg_color: string;
  section_header_text_color: string;
  subtitle_text_color: string;
  destructive_text_color: string;
}

interface WebView {
  initParams: InitParams;
  isIframe: boolean;
}

interface InitParams {
  tgWebAppData: string;
  tgWebAppVersion: string;
  tgWebAppThemeParams: string;
}
