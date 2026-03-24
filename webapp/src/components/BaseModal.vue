<script setup>
import { onMounted, onBeforeUnmount } from 'vue'

const props = defineProps({
  zIndex: {
    type: Number,
    default: 1500,
  },
  maxWidth: {
    type: String,
    default: '700px',
  },
  dataTestid: {
    type: String,
    default: '',
  },
  alignTop: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['close'])

function onOverlayClick(e) {
  // Only close if clicking the overlay itself, not the dialog
  if (e.target === e.currentTarget) {
    emit('close')
  }
}

function onKeydown(e) {
  if (e.key === 'Escape') {
    emit('close')
  }
}

onMounted(() => document.addEventListener('keydown', onKeydown))
onBeforeUnmount(() => document.removeEventListener('keydown', onKeydown))
</script>

<template>
  <Teleport to="body">
    <div
      class="base-overlay"
      :class="{ 'base-overlay--top': alignTop }"
      :style="{ zIndex }"
      :data-testid="dataTestid"
      @click="onOverlayClick"
    >
      <div class="base-dialog" :style="{ maxWidth }">
        <slot />
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
.base-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.65);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem;
  overflow-y: auto;
}

.base-overlay--top {
  align-items: flex-start;
  padding: 2rem 1rem;
}

.base-dialog {
  width: 100%;
  animation: dialog-in 0.18s ease;
}

@keyframes dialog-in {
  from { opacity: 0; transform: translateY(-8px) scale(0.98); }
  to   { opacity: 1; transform: translateY(0) scale(1); }
}

/* Mobile: dialog fills bottom as a sheet */
@media (max-width: 639px) {
  .base-overlay--top {
    align-items: flex-end;
    padding: 0;
  }
}
</style>
