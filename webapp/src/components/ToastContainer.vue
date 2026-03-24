<script setup>
import { useToastStore } from '../stores/toast.js'

const store = useToastStore()
</script>

<template>
  <div class="toast-container" aria-live="polite" aria-atomic="false">
    <TransitionGroup name="toast">
      <div
        v-for="toast in store.toasts"
        :key="toast.id"
        class="toast"
        :class="`toast--${toast.type}`"
        role="status"
      >
        <!-- Icon -->
        <svg v-if="toast.type === 'success'" class="toast__icon" width="18" height="18" viewBox="0 0 18 18" fill="none">
          <circle cx="9" cy="9" r="9" fill="currentColor" opacity="0.15"/>
          <path d="M5 9.5L7.5 12L13 6" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        <svg v-else-if="toast.type === 'error'" class="toast__icon" width="18" height="18" viewBox="0 0 18 18" fill="none">
          <circle cx="9" cy="9" r="9" fill="currentColor" opacity="0.15"/>
          <path d="M6 6L12 12M12 6L6 12" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
        </svg>
        <svg v-else class="toast__icon" width="18" height="18" viewBox="0 0 18 18" fill="none">
          <circle cx="9" cy="9" r="9" fill="currentColor" opacity="0.15"/>
          <circle cx="9" cy="5.5" r="1" fill="currentColor"/>
          <path d="M9 8V13" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
        </svg>

        <span class="toast__message">{{ toast.message }}</span>

        <button
          class="toast__close"
          @click="store.removeToast(toast.id)"
          aria-label="Dismiss notification"
        >
          ×
        </button>
      </div>
    </TransitionGroup>
  </div>
</template>

<style scoped>
.toast-container {
  position: fixed;
  bottom: 1.5rem;
  right: 1.5rem;
  z-index: 9999;
  display: flex;
  flex-direction: column-reverse;
  gap: 0.5rem;
  pointer-events: none;
  max-width: 380px;
}

.toast {
  display: flex;
  align-items: center;
  gap: 0.625rem;
  padding: 0.75rem 1rem;
  border-radius: 10px;
  font-size: 0.8125rem;
  font-weight: 500;
  line-height: 1.4;
  pointer-events: auto;
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.35), 0 1px 4px rgba(0, 0, 0, 0.25);
}

/* Success */
.toast--success {
  background-color: rgba(20, 184, 166, 0.12);
  border: 1px solid rgba(20, 184, 166, 0.3);
  color: var(--color-accent-teal);
}

/* Error */
.toast--error {
  background-color: rgba(239, 68, 68, 0.12);
  border: 1px solid rgba(239, 68, 68, 0.3);
  color: var(--color-danger);
}

/* Info */
.toast--info {
  background-color: rgba(96, 165, 250, 0.12);
  border: 1px solid rgba(96, 165, 250, 0.3);
  color: var(--color-info);
}

.toast__icon {
  flex-shrink: 0;
}

.toast__message {
  flex: 1;
  color: var(--color-text-primary);
}

.toast__close {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border: none;
  background: none;
  color: var(--color-text-muted);
  font-size: 1.125rem;
  cursor: pointer;
  border-radius: 6px;
  transition: background-color 0.15s, color 0.15s;
}

.toast__close:hover {
  background-color: rgba(255, 255, 255, 0.08);
  color: var(--color-text-primary);
}

/* Animations */
.toast-enter-active {
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

.toast-leave-active {
  transition: all 0.2s ease-in;
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(40px) scale(0.95);
}

.toast-leave-to {
  opacity: 0;
  transform: translateX(40px) scale(0.95);
}

.toast-move {
  transition: transform 0.25s ease;
}

/* Mobile */
@media (max-width: 767px) {
  .toast-container {
    bottom: 1rem;
    right: 1rem;
    left: 1rem;
    max-width: none;
  }

  .toast__close {
    width: 44px;
    height: 44px;
    font-size: 1.25rem;
  }
}
</style>
