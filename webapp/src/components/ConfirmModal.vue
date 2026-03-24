<script setup>
defineProps({
  title: { type: String, default: 'Are you sure?' },
  message: { type: String, default: '' },
  confirmText: { type: String, default: 'Confirm' },
  danger: { type: Boolean, default: false },
  loading: { type: Boolean, default: false },
})
defineEmits(['confirm', 'cancel'])
</script>

<template>
  <Teleport to="body">
    <div class="overlay" @click.self="$emit('cancel')">
      <div class="dialog" role="alertdialog" aria-modal="true">
        <h3 class="dialog-title">{{ title }}</h3>
        <p v-if="message" class="dialog-message">{{ message }}</p>
        <div class="dialog-actions">
          <button class="btn-secondary" :disabled="loading" @click="$emit('cancel')">
            Cancel
          </button>
          <button
            class="btn-confirm"
            :class="{ danger }"
            :disabled="loading"
            @click="$emit('confirm')"
          >
            <span v-if="loading" class="spinner" />
            <span>{{ confirmText }}</span>
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
.overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  z-index: 2000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem;
}

.dialog {
  background: var(--color-surface-800);
  border: 1px solid var(--color-surface-600);
  border-radius: 12px;
  padding: 1.5rem;
  width: 100%;
  max-width: 400px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
  animation: dialog-in 0.15s ease;
}

@keyframes dialog-in {
  from { opacity: 0; transform: scale(0.95); }
  to   { opacity: 1; transform: scale(1); }
}

.dialog-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0 0 0.5rem;
}

.dialog-message {
  font-size: 0.875rem;
  color: var(--color-text-muted);
  margin: 0 0 1.5rem;
  line-height: 1.5;
}

.dialog-actions {
  display: flex;
  gap: 0.75rem;
  justify-content: flex-end;
}

.btn-confirm {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.625rem 1.25rem;
  border-radius: 8px;
  border: none;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  background: var(--color-accent-teal);
  color: #0f1923;
  transition: opacity 0.15s;
  min-height: 44px;
}

.btn-confirm.danger {
  background: #ef4444;
  color: #fff;
}

.btn-confirm:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.spinner {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255,255,255,0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
  display: inline-block;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
