<script setup>
const rows = 5
const cols = [
  { width: '50px' },
  { width: '100px' },
  { width: '140px' },
  { width: '100px' },
  { width: '60px' },
  { width: '80px' },
  { width: '60px' },
  { width: '70px' },
  { width: '50px' },
]
</script>

<template>
  <div class="skeleton-wrapper">
    <div class="skeleton-table">
      <!-- Header shimmer -->
      <div class="skeleton-header">
        <div v-for="(col, i) in cols" :key="'h'+i" class="skeleton-cell" :style="{ width: col.width }">
          <div class="shimmer shimmer-header"></div>
        </div>
      </div>
      <!-- Row shimmers -->
      <div v-for="r in rows" :key="r" class="skeleton-row" :style="{ animationDelay: `${r * 80}ms` }">
        <div v-for="(col, i) in cols" :key="'c'+i" class="skeleton-cell" :style="{ width: col.width }">
          <div class="shimmer" :style="{ width: '80%' }"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.skeleton-wrapper {
  border: 1px solid var(--color-surface-600);
  border-radius: 12px;
  background-color: var(--color-surface-800);
  overflow: hidden;
}

.skeleton-header {
  display: flex;
  gap: 1rem;
  padding: 0.75rem 1rem;
  border-bottom: 1px solid var(--color-surface-600);
}

.skeleton-row {
  display: flex;
  gap: 1rem;
  padding: 0.75rem 1rem;
  border-bottom: 1px solid var(--color-surface-700);
  animation: fadeIn 400ms ease forwards;
  opacity: 0;
}

@keyframes fadeIn {
  to { opacity: 1; }
}

.skeleton-cell {
  flex-shrink: 0;
}

.shimmer {
  height: 14px;
  border-radius: 4px;
  background: linear-gradient(
    90deg,
    var(--color-surface-700) 25%,
    var(--color-surface-600) 50%,
    var(--color-surface-700) 75%
  );
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite ease-in-out;
}

.shimmer-header {
  height: 10px;
  width: 70%;
}

@keyframes shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}
</style>
