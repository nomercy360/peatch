import { createSignal, For, Show } from 'solid-js';

type ToastType = 'success' | 'info' | 'warning' | 'error'

interface ToastAction {
  text: string;
  onClick: () => void;
}

interface Toast {
  id: number;
  message: string;
  type: ToastType;
  action?: ToastAction;
}

export const [toasts, setToasts] = createSignal<Toast[]>([]);
let toastCounter = 0;

export const addToast = (
  message: string,
  type: ToastType = 'info',
  action?: ToastAction,
) => {
  const id = ++toastCounter;
  setToasts(prev => [...prev, { id, message, type, action }]);

  setTimeout(() => {
    setToasts(current => current.filter(toast => toast.id !== id));
  }, 3000);
};

const Toast = () => {
  // createEffect(() => {
  //   const currentToasts = toasts();
  //   if (currentToasts.length > 5) {
  //     const newToasts = currentToasts.slice(-5);
  //     setToasts(newToasts);
  //   }
  // });

  const getIcon = (type: ToastType): string => {
    switch (type) {
      case 'success':
        return 'check_circle';
      case 'error':
        return 'error';
      case 'warning':
        return 'warning';
      case 'info':
      default:
        return 'info';
    }
  };

  const getIconColor = (type: ToastType): string => {
    switch (type) {
      case 'success':
        return 'text-green-600';
      case 'error':
        return 'text-red-600';
      case 'warning':
        return 'text-yellow-600';
      case 'info':
      default:
        return 'text-blue-600';
    }
  };

  return (
    <div class="fixed top-4 left-0 right-0 z-50 pointer-events-none mx-auto flex flex-col items-center gap-2 px-4">
      <For each={toasts()}>
        {(toast) => (
          <div
            class="flex w-full max-w-sm items-center gap-3 rounded-2xl bg-background px-4 py-3 shadow-lg border border-secondary pointer-events-auto"
          >
              <span
                class={`material-symbols-rounded text-[24px] ${getIconColor(toast.type)}`}
              >
                {getIcon(toast.type)}
              </span>
            <span
              class="flex-1 text-sm font-medium text-foreground"
            >
                {toast.message}
              </span>
            <Show when={toast.action}>
              <button
                onClick={toast.action!.onClick}
                class="rounded-xl bg-secondary px-3 py-1.5 text-xs font-medium text-secondary-foreground transition-all hover:bg-secondary/80"
              >
                {toast.action!.text}
              </button>
            </Show>
          </div>
        )}
      </For>
    </div>
  );
};

export default Toast;