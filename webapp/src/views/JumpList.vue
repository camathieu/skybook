<script setup>
import { watch, onMounted, computed, ref, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useJumpStore } from '../stores/jumps.js'
import SearchBar from '../components/SearchBar.vue'
import FilterBar from '../components/FilterBar.vue'
import JumpTable from '../components/JumpTable.vue'
import JumpCard from '../components/JumpCard.vue'
import JumpSkeleton from '../components/JumpSkeleton.vue'
import Pagination from '../components/Pagination.vue'

const route = useRoute()
const router = useRouter()
const store = useJumpStore()

const ready = ref(false)

const showEmpty = computed(() => !store.loading && !store.error && store.total === 0 && !store.hasActiveFilters())
const showNoResults = computed(() => !store.loading && !store.error && store.total === 0 && store.hasActiveFilters())
const showData = computed(() => !store.loading && !store.error && store.total > 0)

// Initialize: read URL → fetch → then start watching for changes
onMounted(async () => {
  store.initFromQuery(route.query)
  await store.fetchJumps()
  await nextTick()
  ready.value = true
})

// After initialization, watch store state changes → update URL + refetch
watch(
  () => store.toQuery(),
  (newQuery) => {
    if (!ready.value) return
    router.replace({ query: newQuery })
    store.fetchJumps()
  },
)

// Watch filter changes → reset to page 1
watch(
  () => ({ ...store.filters }),
  () => {
    if (!ready.value) return
    if (store.page !== 1) {
      store.page = 1
    }
  },
)

function retry() {
  store.fetchJumps()
}
</script>

<template>
  <div class="jump-list">
    <!-- Toolbar: search + filters -->
    <div class="toolbar">
      <SearchBar />
      <FilterBar />
    </div>

    <!-- Loading skeleton -->
    <JumpSkeleton v-if="store.loading" />

    <!-- Error state -->
    <div v-else-if="store.error" class="state-card">
      <div class="state-icon">⚠</div>
      <h2 class="state-title">Something went wrong</h2>
      <p class="state-subtitle">{{ store.error }}</p>
      <button class="btn-primary" @click="retry">Try again</button>
    </div>

    <!-- Empty state (no jumps at all) -->
    <div v-else-if="showEmpty" class="state-card">
      <div class="state-icon pulse">✦</div>
      <h1 class="state-title">Your logbook is empty</h1>
      <p class="state-subtitle">Log your first jump to get started</p>
      <button class="btn-primary state-cta">
        <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
          <path d="M8 2a.75.75 0 0 1 .75.75v4.5h4.5a.75.75 0 0 1 0 1.5h-4.5v4.5a.75.75 0 0 1-1.5 0v-4.5h-4.5a.75.75 0 0 1 0-1.5h4.5v-4.5A.75.75 0 0 1 8 2Z"/>
        </svg>
        Log First Jump
      </button>
    </div>

    <!-- No results from filters -->
    <div v-else-if="showNoResults" class="state-card">
      <div class="state-icon">∅</div>
      <h2 class="state-title">No jumps found</h2>
      <p class="state-subtitle">Try adjusting your search or filters</p>
      <button class="btn-secondary" @click="store.resetFilters()">Clear all filters</button>
    </div>

    <!-- Data: table + cards -->
    <template v-else-if="showData">
      <!-- Desktop table -->
      <div class="desktop-only">
        <JumpTable />
      </div>

      <!-- Mobile cards -->
      <div class="mobile-only">
        <div class="card-list">
          <JumpCard v-for="jump in store.items" :key="jump.id" :jump="jump" />
        </div>
      </div>

      <Pagination />
    </template>
  </div>
</template>

<style scoped>
.jump-list {
  max-width: 1200px;
  margin: 0 auto;
  padding: 1.5rem;
}

/* Toolbar */
.toolbar {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  margin-bottom: 1.5rem;
}

/* State cards (empty, error, no results) */
.state-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 4rem 1rem;
  min-height: 40vh;
}

.state-icon {
  font-size: 3rem;
  margin-bottom: 1.5rem;
  background: linear-gradient(135deg, var(--color-accent-orange), var(--color-accent-teal));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.state-icon.pulse {
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.7; transform: scale(1.05); }
}

.state-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0 0 0.5rem;
}

.state-subtitle {
  font-size: 1rem;
  color: var(--color-text-muted);
  margin: 0 0 2rem;
}

.state-cta {
  padding: 0.75rem 2rem;
  font-size: 1rem;
}

/* Responsive toggle */
.desktop-only {
  display: block;
}

.mobile-only {
  display: none;
}

.card-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

@media (max-width: 767px) {
  .jump-list {
    padding: 1rem;
  }

  .desktop-only {
    display: none;
  }

  .mobile-only {
    display: block;
  }

  .state-card {
    padding: 3rem 1rem;
    min-height: 30vh;
  }
}
</style>
