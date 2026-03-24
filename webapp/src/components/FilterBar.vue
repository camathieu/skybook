<script setup>
import { ref, computed, onMounted } from 'vue'
import { useJumpStore } from '../stores/jumps.js'
import { api } from '../api.js'

const store = useJumpStore()
const showMobileFilters = ref(false)
const dropzoneOptions = ref([])

const jumpTypes = [
  'FF', 'WS', 'FS', 'CRW', 'HOP', 'CF', 'AFF',
  'TANDEM', 'DEMO', 'XRW', 'ANGLE', 'TRACKING', 'CP', 'WINGSUIT', 'OTHER'
]

const activeFilters = computed(() => {
  const chips = []
  if (store.filters.dropzone) chips.push({ key: 'dropzone', label: `DZ: ${store.filters.dropzone}` })
  if (store.filters.jumpType) chips.push({ key: 'jumpType', label: `Type: ${store.filters.jumpType}` })
  if (store.filters.dateFrom) chips.push({ key: 'dateFrom', label: `From: ${store.filters.dateFrom}` })
  if (store.filters.dateTo) chips.push({ key: 'dateTo', label: `To: ${store.filters.dateTo}` })
  if (store.filters.cutaway !== null) chips.push({ key: 'cutaway', label: store.filters.cutaway ? 'Cutaway' : 'No cutaway' })
  if (store.filters.night !== null) chips.push({ key: 'night', label: store.filters.night ? 'Night' : 'Day only' })
  return chips
})

function removeFilter(key) {
  if (key === 'cutaway' || key === 'night') {
    store.filters[key] = null
  } else {
    store.filters[key] = ''
  }
}

async function loadDropzones() {
  try {
    const results = await api.get('/jumps/autocomplete/dropzone')
    dropzoneOptions.value = results.map(r => r.value)
  } catch {
    dropzoneOptions.value = []
  }
}

function toggleBool(key) {
  if (store.filters[key] === null) {
    store.filters[key] = true
  } else if (store.filters[key] === true) {
    store.filters[key] = false
  } else {
    store.filters[key] = null
  }
}

onMounted(() => {
  loadDropzones()
})
</script>

<template>
  <!-- Mobile trigger -->
  <button class="filter-toggle mobile-only" @click="showMobileFilters = !showMobileFilters">
    <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
      <path d="M1.5 1.5A.5.5 0 0 1 2 1h12a.5.5 0 0 1 .39.812l-4.89 6.13V13.5a.5.5 0 0 1-.724.447l-2.5-1.25A.5.5 0 0 1 6 12.25V7.942L1.11 1.812A.5.5 0 0 1 1.5 1.5z"/>
    </svg>
    Filters
    <span v-if="activeFilters.length" class="filter-count">{{ activeFilters.length }}</span>
  </button>

  <!-- Filter content -->
  <div class="filter-bar" :class="{ 'mobile-open': showMobileFilters }">
    <div class="filter-row">
      <!-- Jump Type -->
      <select
        v-model="store.filters.jumpType"
        class="filter-select"
        aria-label="Filter by jump type"
      >
        <option value="">All types</option>
        <option v-for="t in jumpTypes" :key="t" :value="t">{{ t }}</option>
      </select>

      <!-- Dropzone -->
      <select
        v-model="store.filters.dropzone"
        class="filter-select"
        aria-label="Filter by dropzone"
      >
        <option value="">All dropzones</option>
        <option v-for="dz in dropzoneOptions" :key="dz" :value="dz">{{ dz }}</option>
      </select>

      <!-- Date range -->
      <input
        v-model="store.filters.dateFrom"
        type="date"
        class="filter-date"
        aria-label="Date from"
      />
      <span class="date-sep">–</span>
      <input
        v-model="store.filters.dateTo"
        type="date"
        class="filter-date"
        aria-label="Date to"
      />

      <!-- Boolean toggles -->
      <button
        class="filter-bool"
        :class="{ active: store.filters.cutaway === true, inactive: store.filters.cutaway === false }"
        @click="toggleBool('cutaway')"
        title="Toggle cutaway filter"
      >
        ✂ Cutaway
      </button>
      <button
        class="filter-bool"
        :class="{ active: store.filters.night === true, inactive: store.filters.night === false }"
        @click="toggleBool('night')"
        title="Toggle night filter"
      >
        🌙 Night
      </button>
    </div>

    <!-- Active filter chips -->
    <div v-if="activeFilters.length" class="filter-chips">
      <span v-for="chip in activeFilters" :key="chip.key" class="chip">
        {{ chip.label }}
        <button class="chip-remove" @click="removeFilter(chip.key)" aria-label="Remove filter">×</button>
      </span>
      <button class="chip chip-clear" @click="store.resetFilters()">Clear all</button>
    </div>
  </div>
</template>

<style scoped>
/* Mobile toggle */
.filter-toggle {
  display: none;
}

.mobile-only {
  display: none;
}

/* Filter bar */
.filter-bar {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.filter-row {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.filter-select,
.filter-date {
  min-height: 40px;
  padding: 0.5rem 0.75rem;
  font-size: 0.8125rem;
  color: var(--color-text-primary);
  background-color: var(--color-surface-700);
  border: 1px solid var(--color-surface-600);
  border-radius: 8px;
  transition: border-color var(--transition-fast);
  cursor: pointer;
}

.filter-select:focus,
.filter-date:focus {
  outline: none;
  border-color: var(--color-accent-teal);
}

.filter-select {
  min-width: 120px;
}

.filter-date {
  width: 140px;
}

.date-sep {
  color: var(--color-text-muted);
  font-size: 0.875rem;
}

/* Boolean toggles */
.filter-bool {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  min-height: 40px;
  padding: 0.5rem 0.75rem;
  font-size: 0.8125rem;
  color: var(--color-text-muted);
  background-color: var(--color-surface-700);
  border: 1px solid var(--color-surface-600);
  border-radius: 8px;
  cursor: pointer;
  transition: all var(--transition-fast);
  user-select: none;
  -webkit-tap-highlight-color: transparent;
}

.filter-bool:hover {
  border-color: var(--color-surface-500);
  color: var(--color-text-secondary);
}

.filter-bool.active {
  background: rgba(20, 184, 166, 0.15);
  border-color: rgba(20, 184, 166, 0.3);
  color: var(--color-accent-teal);
}

.filter-bool.inactive {
  background: rgba(239, 68, 68, 0.1);
  border-color: rgba(239, 68, 68, 0.2);
  color: var(--color-text-muted);
}

/* Filter chips */
.filter-chips {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  flex-wrap: wrap;
}

.chip {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.25rem 0.625rem;
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  background-color: var(--color-surface-700);
  border: 1px solid var(--color-surface-600);
  border-radius: 9999px;
}

.chip-remove {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  border: none;
  background: none;
  color: var(--color-text-muted);
  font-size: 0.875rem;
  cursor: pointer;
  border-radius: 50%;
  transition: color var(--transition-fast);
}

.chip-remove:hover {
  color: var(--color-danger);
}

.chip-clear {
  cursor: pointer;
  color: var(--color-text-muted);
  border-style: dashed;
  transition: color var(--transition-fast);
}

.chip-clear:hover {
  color: var(--color-text-primary);
}

/* Mobile */
@media (max-width: 767px) {
  .mobile-only {
    display: inline-flex;
  }

  .filter-toggle {
    display: inline-flex;
    align-items: center;
    gap: 0.375rem;
    min-height: 44px;
    padding: 0.5rem 1rem;
    font-size: 0.875rem;
    font-weight: 500;
    color: var(--color-text-secondary);
    background-color: var(--color-surface-700);
    border: 1px solid var(--color-surface-600);
    border-radius: 10px;
    cursor: pointer;
    -webkit-tap-highlight-color: transparent;
  }

  .filter-count {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-width: 20px;
    height: 20px;
    font-size: 0.6875rem;
    font-weight: 700;
    color: white;
    background: linear-gradient(135deg, var(--color-accent-orange), var(--color-accent-teal));
    border-radius: 9999px;
    padding: 0 0.25rem;
  }

  .filter-bar {
    display: none;
    padding: 1rem;
    background-color: var(--color-surface-800);
    border: 1px solid var(--color-surface-600);
    border-radius: 12px;
  }

  .filter-bar.mobile-open {
    display: flex;
  }

  .filter-row {
    flex-direction: column;
    align-items: stretch;
  }

  .filter-select,
  .filter-date {
    width: 100%;
    min-height: 44px;
  }

  .date-sep {
    display: none;
  }
}
</style>
