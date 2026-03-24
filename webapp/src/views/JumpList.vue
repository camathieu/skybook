<script setup>
import { watch, onMounted, onUnmounted, computed, ref, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useJumpStore } from '../stores/jumps.js'
import SearchBar from '../components/SearchBar.vue'
import FilterBar from '../components/FilterBar.vue'
import JumpTable from '../components/JumpTable.vue'
import JumpCard from '../components/JumpCard.vue'
import JumpSkeleton from '../components/JumpSkeleton.vue'
import Pagination from '../components/Pagination.vue'
import JumpModal from '../components/JumpModal.vue'

const route = useRoute()
const router = useRouter()
const store = useJumpStore()

const ready = ref(false)

// --- Modal state ---
const showModal = ref(false)
const editingJump = ref(null) // null = create, object = edit

function openCreate() {
  editingJump.value = null
  showModal.value = true
}

function openEdit(jump) {
  editingJump.value = jump
  showModal.value = true
}

function onModalClose() {
  showModal.value = false
  editingJump.value = null
}

// --- Keyboard shortcut: N = new jump ---
function onGlobalKey(e) {
  // Ignore when typing in inputs/textareas
  const tag = document.activeElement?.tagName
  if (tag === 'INPUT' || tag === 'TEXTAREA' || tag === 'SELECT') return
  if (e.key === 'n' || e.key === 'N') openCreate()
}

onMounted(() => window.addEventListener('keydown', onGlobalKey))
onUnmounted(() => window.removeEventListener('keydown', onGlobalKey))

// --- Data states ---
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
  <div class="jump-list" data-testid="jump-list">
    <!-- Page header with New Jump button -->
    <div class="page-header">
      <h1 class="page-title">Jumps</h1>
      <button class="btn-new-jump" data-testid="new-jump-btn" @click="openCreate" title="New Jump (N)">
        <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
          <path d="M8 2a.75.75 0 0 1 .75.75v4.5h4.5a.75.75 0 0 1 0 1.5h-4.5v4.5a.75.75 0 0 1-1.5 0v-4.5h-4.5a.75.75 0 0 1 0-1.5h4.5v-4.5A.75.75 0 0 1 8 2Z"/>
        </svg>
        New Jump
      </button>
    </div>

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
      <button class="btn-primary state-cta" @click="openCreate">
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
        <JumpTable @edit="openEdit" />
      </div>

      <!-- Mobile cards -->
      <div class="mobile-only">
        <div class="card-list">
          <JumpCard v-for="jump in store.items" :key="jump.id" :jump="jump" @edit="openEdit" />
        </div>
      </div>

      <Pagination />
    </template>
  </div>

  <!-- Jump Modal (create / edit) -->
  <JumpModal v-if="showModal" data-testid="jump-modal" :jump="editingJump" @close="onModalClose" />
</template>

<style scoped>
.jump-list {
  max-width: 1200px;
  margin: 0 auto;
  padding: 1.5rem;
}

/* Page header */
.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1.25rem;
}

.page-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--color-text-primary);
  margin: 0;
}

.btn-new-jump {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  background: linear-gradient(135deg, var(--color-accent-orange), var(--color-accent-teal));
  color: #0f1923;
  border: none;
  border-radius: 8px;
  padding: 0.625rem 1.25rem;
  font-size: 0.875rem;
  font-weight: 700;
  cursor: pointer;
  min-height: 44px;
  transition: opacity 0.15s, transform 0.1s;
}

.btn-new-jump:hover {
  opacity: 0.9;
  transform: translateY(-1px);
}

.btn-new-jump:active {
  transform: translateY(0);
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
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 2rem;
  font-size: 1rem;
}

/* Responsive toggle */
.desktop-only { display: block; }
.mobile-only  { display: none; }

.card-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

@media (max-width: 767px) {
  .jump-list {
    padding: 1rem;
  }

  .desktop-only { display: none; }
  .mobile-only  { display: block; }

  .state-card {
    padding: 3rem 1rem;
    min-height: 30vh;
  }
}
</style>
