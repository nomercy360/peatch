import { createSignal, createEffect } from 'solid-js'

type ToastType = 'success' | 'info' | 'warning' | 'error'

interface ToastAction {
  text: string
  onClick: () => void
}

interface Toast {
  id: number
  message: string
  type: ToastType
  action?: ToastAction
}

export const [toasts, setToasts] = createSignal<Toast[]>([])

export const addToast = (
  message: string,
  type: ToastType = 'info',
  action?: ToastAction,
) => {
  const id = Date.now()
  setToasts([...toasts(), { id, message, type, action }])

  setTimeout(() => {
    setToasts(toasts().filter(toast => toast.id !== id))
  }, 5000)
}

const Toast = () => {
  createEffect(() => {
    const currentToasts = toasts()
    if (currentToasts.length > 5) {
      const newToasts = currentToasts.slice(1)
      setToasts(newToasts)
    }
  })

  const getIcon = (type: ToastType): string => {
    switch (type) {
      case 'success':
        return 'check_circle'
      case 'error':
        return 'error'
      case 'warning':
        return 'warning'
      case 'info':
      default:
        return 'info'
    }
  }

  const getBackgroundColor = (type: ToastType): string => {
    switch (type) {
      case 'success':
        return 'bg-success'
      case 'error':
        return 'bg-destructive'
      case 'warning':
        return 'bg-warning'
      case 'info':
      default:
        return 'bg-accent'
    }
  }

  return (
    <div class="fixed bottom-4 left-1/2 z-50 -translate-x-1/2 transform space-y-2">
      {toasts().map(toast => (
        <div
          class={`flex w-[calc(100vw-2rem)] items-center justify-between rounded-lg ${getBackgroundColor(
            toast.type,
          )} border border-white border-opacity-20 px-4 py-3 text-sm font-medium shadow-lg animate-in fade-in`}
        >
          <div class="flex items-center">
            <span class="material-symbols-rounded mr-2 text-[20px]">
              {getIcon(toast.type)}
            </span>
            <span class="flex-1">{toast.message}</span>
          </div>
          {toast.action && (
            <button
              onClick={toast.action.onClick}
              class="ml-4 rounded-md bg-white bg-opacity-20 px-3 py-1 text-xs transition-colors hover:bg-opacity-30"
            >
              {toast.action.text}
            </button>
          )}
        </div>
      ))}
    </div>
  )
}

export default Toast
