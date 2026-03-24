<script setup>
import { computed } from 'vue'
import { useJumpStore } from '../stores/jumps.js'

const store = useJumpStore()

const perPageOptions = [25, 50, 100]

const showingFrom = computed(() => {
  if (store.total === 0) return 0
  return (store.page - 1) * store.perPage + 1
})

const showingTo = computed(() => {
  return Math.min(store.page * store.perPage, store.total)
})

const canPrev = computed(() => store.page > 1)
const canNext = computed(() => store.page < store.totalPages)
</script>

<template>
  <div v-if="store.total > 0" class="pagination">
    <span class="pagination-info">
      Showing <strong>{{ showingFrom }}–{{ showingTo }}</strong> of <strong>{{ store.total }}</strong> jumps
    </span>

    <div class="pagination-controls">
      <button
        class="pagination-btn"
        :disabled="!canPrev"
        @click="store.setPage(store.page - 1)"
        aria-label="Previous page"
      >
        ← Prev
      </button>

      <span class="pagination-page">
        Page {{ store.page }} / {{ store.totalPages }}
      </span>

      <button
        class="pagination-btn"
        :disabled="!canNext"
        @click="store.setPage(store.page + 1)"
        aria-label="Next page"
      >
        Next →
      </button>
    </div>

    <div class="per-page">
      <span class="per-page-label">Per page:</span>
      <div class="per-page-group">
        <button
          v-for="pp in perPageOptions"
          :key="pp"
          class="per-page-btn"
          :class="{ active: store.perPage === pp }"
          :data-testid="`per-page-${pp}`"
          @click="store.setPerPage(pp)"
        >
          {{ pp }}
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.pagination {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding: 1rem 0;
  flex-wrap: wrap;
}

.pagination-info {
  font-size: 0.8125rem;
  color: var(--color-text-muted);
}

.pagination-info strong {
  color: var(--color-text-secondary);
  font-weight: 600;
}

.pagination-controls {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.pagination-btn {
  display: inline-flex;
  align-items: center;
  min-height: 36px;
  padding: 0.375rem 0.875rem;
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  background-color: var(--color-surface-700);
  border: 1px solid var(--color-surface-600);
  border-radius: 8px;
  cursor: pointer;
  transition: all var(--transition-fast);
  user-select: none;
}

.pagination-btn:hover:not(:disabled) {
  background-color: var(--color-surface-600);
  color: var(--color-text-primary);
}

.pagination-btn:disabled {
  opacity: 0.3;
  cursor: default;
}

.pagination-page {
  font-size: 0.8125rem;
  font-family: var(--font-mono);
  color: var(--color-text-muted);
}

/* Per-page */
.per-page {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.per-page-label {
  font-size: 0.75rem;
  color: var(--color-text-muted);
}

.per-page-group {
  display: flex;
  border: 1px solid var(--color-surface-600);
  border-radius: 8px;
  overflow: hidden;
}

.per-page-btn {
  padding: 0.375rem 0.625rem;
  font-size: 0.75rem;
  font-family: var(--font-mono);
  color: var(--color-text-muted);
  background: var(--color-surface-700);
  border: none;
  border-right: 1px solid var(--color-surface-600);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.per-page-btn:last-child {
  border-right: none;
}

.per-page-btn:hover {
  background-color: var(--color-surface-600);
  color: var(--color-text-secondary);
}

.per-page-btn.active {
  background: rgba(20, 184, 166, 0.15);
  color: var(--color-accent-teal);
  font-weight: 600;
}

/* Mobile */
@media (max-width: 767px) {
  .pagination {
    flex-direction: column;
    align-items: center;
    text-align: center;
  }

  .pagination-controls {
    width: 100%;
    justify-content: center;
  }

  .pagination-btn {
    min-height: 44px;
    padding: 0.5rem 1.25rem;
  }

  .per-page-btn {
    min-height: 44px;
    padding: 0.5rem 0.75rem;
  }
}
</style>
