import { defineStore } from 'pinia'
import { ref } from 'vue'

let nextId = 0

export const useToastStore = defineStore('toast', () => {
  const toasts = ref([])
  const MAX_TOASTS = 5

  /**
   * Show a toast notification.
   * @param {string} message - Text to display
   * @param {'success'|'error'|'info'} type - Visual style
   * @param {number} duration - Auto-dismiss in ms (0 = no auto-dismiss)
   */
  function addToast(message, type = 'success', duration = 4000) {
    const id = ++nextId
    toasts.value.push({ id, message, type })

    // Enforce max cap
    while (toasts.value.length > MAX_TOASTS) {
      toasts.value.shift()
    }

    if (duration > 0) {
      setTimeout(() => removeToast(id), duration)
    }
  }

  function removeToast(id) {
    const idx = toasts.value.findIndex(t => t.id === id)
    if (idx !== -1) {
      toasts.value.splice(idx, 1)
    }
  }

  return { toasts, addToast, removeToast }
})
