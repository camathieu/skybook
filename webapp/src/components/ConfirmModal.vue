<script setup>
import BaseModal from './BaseModal.vue'

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
  <BaseModal :z-index="2000" data-testid="confirm-modal" @close="$emit('cancel')">
    <div class="dialog" role="alertdialog" aria-modal="true">
        <h3 class="dialog-title">{{ title }}</h3>
        <p v-if="message" class="dialog-message">{{ message }}</p>
        <div class="dialog-actions">
          <button class="btn-secondary" :disabled="loading" @click="$emit('cancel')">
            Cancel
          </button>
          <button
            class="btn-primary"
            :class="{ 'btn-primary--danger': danger }"
            data-testid="confirm-delete-btn"
            :disabled="loading"
            @click="$emit('confirm')"
          >
            <span v-if="loading" class="spinner" />
            <span>{{ confirmText }}</span>
          </button>
        </div>
    </div>
  </BaseModal>
</template>

<style scoped>
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
</style>

