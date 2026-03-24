<script setup>
import { ref, watch, onMounted, onUnmounted } from 'vue'
import { useJumpStore } from '../stores/jumps.js'

const store = useJumpStore()
const searchInput = ref(null)
const localQuery = ref(store.filters.q)

let debounceTimer = null

watch(localQuery, (val) => {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    store.filters.q = val
  }, 300)
})

// Sync from store to local (e.g. on reset)
watch(() => store.filters.q, (val) => {
  if (val !== localQuery.value) {
    localQuery.value = val
  }
})

function clear() {
  localQuery.value = ''
  store.filters.q = ''
}

function handleKeydown(e) {
  if (e.key === '/' && document.activeElement?.tagName !== 'INPUT') {
    e.preventDefault()
    searchInput.value?.focus()
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
  clearTimeout(debounceTimer)
})
</script>

<template>
  <div class="search-bar">
    <svg class="search-icon" width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
      <path d="M11.742 10.344a6.5 6.5 0 1 0-1.397 1.398h-.001q.044.06.098.115l3.85 3.85a1 1 0 0 0 1.415-1.414l-3.85-3.85a1 1 0 0 0-.115-.1zM12 6.5a5.5 5.5 0 1 1-11 0 5.5 5.5 0 0 1 11 0"/>
    </svg>
    <input
      ref="searchInput"
      v-model="localQuery"
      type="text"
      class="search-input"
      data-testid="search-input"
      placeholder="Search jumps..."
      aria-label="Search jumps"
    />
    <button v-if="localQuery" class="search-clear" @click="clear" aria-label="Clear search">
      ×
    </button>
    <span v-else class="kbd desktop-only">/</span>
  </div>
</template>

<style scoped>
.search-bar {
  position: relative;
  display: flex;
  align-items: center;
  flex: 1;
  max-width: 400px;
}

.search-icon {
  position: absolute;
  left: 0.875rem;
  color: var(--color-text-muted);
  pointer-events: none;
}

.search-input {
  width: 100%;
  padding: 0.625rem 2.5rem 0.625rem 2.5rem;
  min-height: 44px;
  font-size: 0.875rem;
  color: var(--color-text-primary);
  background-color: var(--color-surface-700);
  border: 1px solid var(--color-surface-600);
  border-radius: 10px;
  transition: border-color var(--transition-fast);
}

.search-input::placeholder {
  color: var(--color-text-muted);
}

.search-input:focus {
  outline: none;
  border-color: var(--color-accent-teal);
  box-shadow: 0 0 0 3px rgba(20, 184, 166, 0.15);
}

.search-clear {
  position: absolute;
  right: 0.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border: none;
  background: none;
  color: var(--color-text-muted);
  font-size: 1.25rem;
  cursor: pointer;
  border-radius: 6px;
  transition: color var(--transition-fast), background-color var(--transition-fast);
}

.search-clear:hover {
  color: var(--color-text-primary);
  background-color: var(--color-surface-600);
}

.desktop-only {
  position: absolute;
  right: 0.75rem;
}

@media (max-width: 767px) {
  .search-bar {
    max-width: none;
  }

  .desktop-only {
    display: none;
  }
}
</style>
