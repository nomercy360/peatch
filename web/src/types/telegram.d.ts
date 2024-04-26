export interface Telegram {
  WebView: WebView
  Utils: Utils
  WebApp: WebApp
}

export interface Utils {
  result?: unknown

  notificationOccurred(warning: string): void
}

export interface WebApp {
  initData: string
  initDataUnsafe: InitDataUnsafe
  version: string
  colorScheme: string
  themeParams: ThemeParams
  isExpanded: boolean
  viewportHeight: number
  viewportStableHeight: number
  isClosingConfirmationEnabled: boolean
  headerColor: string
  backgroundColor: string
  BackButton: BackButton
  MainButton: MainButton
  HapticFeedback: Utils

  openTelegramLink(url: string): void

  showAlert(message: string, callback: () => void): void

  showConfirm(message: string, callback: () => void): void

  expand(): void

  ready(): void

  onEvent(event: string, callback: () => void): void

  offEvent(event: string, callback: () => void): void
}

export interface BackButton {
  isVisible: boolean

  onClick(callback: () => void): void

  offClick(callback: () => void): void

  show(): void
  hide(): void

  setParams(param: { text_color?: string }): any
}

export interface MainButton {
  onClick: any
  text: string
  color: string
  offClick: any
  textColor: string
  isVisible: boolean
  isProgressVisible: boolean
  isActive: boolean

  setParams(param: { text_color?: string; color?: string; text?: string }): any

  showProgress(leaveActive: boolean): void

  hideProgress(): void

  disable(): void

  setText(next: string): void;

  show(): void;

  enable(): void;
}

export interface InitDataUnsafe {
  query_id: string
  user: User
  auth_date: string
  hash: string
}

export interface User {
  id: number
  first_name: string
  last_name: string
  username: string
  language_code: string
}

export interface ThemeParams {
  bg_color: string
  text_color: string
  hint_color: string
  link_color: string
  button_color: string
  button_text_color: string
  secondary_bg_color: string
}

export interface WebView {
  initParams: InitParams
  isIframe: boolean
}

export interface InitParams {
  tgWebAppData: string
  tgWebAppVersion: string
  tgWebAppThemeParams: string
}
