import { writable, derived } from 'svelte/store';

export type PageName = 'home' | 'editor' | 'glossary' | 'projects' | 'settings';

export interface ToastMessage {
  id: string;
  type: 'success' | 'error' | 'warning' | 'info';
  message: string;
  duration: number;
}

export interface ModalState {
  isOpen: boolean;
  component: string | null;
  props: Record<string, unknown>;
}

export interface UIState {
  sidebarOpen: boolean;
  sidebarCollapsed: boolean;
  currentPage: PageName;
  toasts: ToastMessage[];
  modal: ModalState;
  isWindowMaximized: boolean;
}

const initialState: UIState = {
  sidebarOpen: true,
  sidebarCollapsed: false,
  currentPage: 'home',
  toasts: [],
  modal: {
    isOpen: false,
    component: null,
    props: {}
  },
  isWindowMaximized: false
};

function createUIStore() {
  const { subscribe, set, update } = writable<UIState>(initialState);

  let toastIdCounter = 0;

  return {
    subscribe,

    toggleSidebar() {
      update(s => ({ ...s, sidebarOpen: !s.sidebarOpen }));
    },

    setSidebarOpen(open: boolean) {
      update(s => ({ ...s, sidebarOpen: open }));
    },

    toggleSidebarCollapsed() {
      update(s => ({ ...s, sidebarCollapsed: !s.sidebarCollapsed }));
    },

    setPage(page: PageName) {
      update(s => ({ ...s, currentPage: page }));
    },

    showToast(type: ToastMessage['type'], message: string, duration = 5000) {
      const id = `toast-${++toastIdCounter}`;
      const toast: ToastMessage = { id, type, message, duration };
      
      update(s => ({ ...s, toasts: [...s.toasts, toast] }));

      if (duration > 0) {
        setTimeout(() => {
          this.dismissToast(id);
        }, duration);
      }

      return id;
    },

    dismissToast(id: string) {
      update(s => ({
        ...s,
        toasts: s.toasts.filter(t => t.id !== id)
      }));
    },

    clearToasts() {
      update(s => ({ ...s, toasts: [] }));
    },

    openModal(component: string, props: Record<string, unknown> = {}) {
      update(s => ({
        ...s,
        modal: { isOpen: true, component, props }
      }));
    },

    closeModal() {
      update(s => ({
        ...s,
        modal: { isOpen: false, component: null, props: {} }
      }));
    },

    setWindowMaximized(maximized: boolean) {
      update(s => ({ ...s, isWindowMaximized: maximized }));
    },

    reset() {
      set(initialState);
    }
  };
}

export const ui = createUIStore();

export const isSidebarVisible = derived(ui, $ui => $ui.sidebarOpen);
export const hasToasts = derived(ui, $ui => $ui.toasts.length > 0);
export const isModalOpen = derived(ui, $ui => $ui.modal.isOpen);
