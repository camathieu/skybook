<script setup>
import { useJumpStore } from '../stores/jumps.js'

const emit = defineEmits(['edit'])
const store = useJumpStore()

const columns = [
  { key: 'number', label: '#', sortable: true, class: 'col-number' },
  { key: 'date', label: 'Date', sortable: true, class: 'col-date' },
  { key: 'dropzone', label: 'Dropzone', sortable: true, class: 'col-dropzone' },
  { key: 'aircraft', label: 'Aircraft', sortable: false, class: 'col-aircraft' },
  { key: 'jumpType', label: 'Type', sortable: false, class: 'col-type' },
  { key: 'altitude', label: 'Altitude', sortable: true, class: 'col-altitude' },
  { key: 'freefallTime', label: 'Freefall', sortable: false, class: 'col-freefall' },
  { key: 'landing', label: 'Landing', sortable: false, class: 'col-landing' },
  { key: 'flags', label: '', sortable: false, class: 'col-flags' },
]

function formatDate(dateStr) {
  if (!dateStr) return '—'
  // Parse YYYY-MM-DD as local date to avoid timezone shift
  // (new Date("2025-03-08") is UTC, which shifts in -UTC timezones)
  const [y, m, d] = dateStr.split('-').map(Number)
  const date = new Date(y, m - 1, d)
  return date.toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric' })
}

function formatAltitude(alt) {
  if (alt == null) return '—'
  return alt.toLocaleString() + ' ft'
}

function formatFreefall(seconds) {
  if (seconds == null) return '—'
  return seconds + 's'
}

function handleSort(col) {
  if (!col.sortable) return
  store.setSort(col.key)
}

function sortIcon(col) {
  if (!col.sortable) return ''
  if (store.sortBy !== col.key) return '↕'
  return store.order === 'asc' ? '↑' : '↓'
}
</script>

<template>
  <div class="table-wrapper">
    <table class="jump-table">
      <thead>
        <tr>
          <th
            v-for="col in columns"
            :key="col.key"
            :class="[col.class, { sortable: col.sortable, active: store.sortBy === col.key }]"
            @click="handleSort(col)"
          >
            {{ col.label }}
            <span v-if="col.sortable" class="sort-icon">{{ sortIcon(col) }}</span>
          </th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="(jump, i) in store.items"
          :key="jump.id"
          class="jump-row"
          data-testid="jump-row"
          style="cursor:pointer"
          tabindex="0"
          @click="emit('edit', jump)"
          @keydown.enter="emit('edit', jump)"
        >
          <td class="col-number">
            <span class="jump-number">{{ jump.number }}</span>
          </td>
          <td class="col-date">{{ formatDate(jump.date) }}</td>
          <td class="col-dropzone">{{ jump.dropzone || '—' }}</td>
          <td class="col-aircraft">{{ jump.aircraft || '—' }}</td>
          <td class="col-type">
            <span class="badge">{{ jump.jumpType }}</span>
          </td>
          <td class="col-altitude">{{ formatAltitude(jump.altitude) }}</td>
          <td class="col-freefall">{{ formatFreefall(jump.freefallTime) }}</td>
          <td class="col-landing">{{ jump.landing || '—' }}</td>
          <td class="col-flags">
            <span v-if="jump.favorite" class="flag flag-fav" title="Favorite">★</span>
            <span v-if="jump.nightJump" class="flag flag-night" title="Night jump">🌙</span>
            <span v-if="jump.oxygenJump" class="flag flag-o2" title="Oxygen">O₂</span>
            <span v-if="jump.cutaway" class="flag flag-cutaway" title="Cutaway">✂</span>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.table-wrapper {
  overflow-x: auto;
  border: 1px solid var(--color-surface-600);
  border-radius: 12px;
  background-color: var(--color-surface-800);
}

.jump-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.875rem;
}

/* Header */
.jump-table thead {
  position: sticky;
  top: 0;
  z-index: 2;
}

.jump-table th {
  padding: 0.75rem 1rem;
  text-align: left;
  font-weight: 600;
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--color-text-muted);
  background-color: var(--color-surface-800);
  border-bottom: 1px solid var(--color-surface-600);
  white-space: nowrap;
  user-select: none;
}

.jump-table th.sortable {
  cursor: pointer;
  transition: color var(--transition-fast);
}

.jump-table th.sortable:hover {
  color: var(--color-text-secondary);
}

.jump-table th.active {
  color: var(--color-accent-teal);
}

.sort-icon {
  display: inline-block;
  margin-left: 0.25rem;
  font-size: 0.6875rem;
  opacity: 0.5;
}

.jump-table th.active .sort-icon {
  opacity: 1;
}

/* Rows */
.jump-table td {
  padding: 0.625rem 1rem;
  border-bottom: 1px solid var(--color-surface-700);
  color: var(--color-text-secondary);
  white-space: nowrap;
}

.jump-row {
  transition: background-color var(--transition-fast);
}

.jump-row:hover {
  background-color: var(--color-surface-700);
}

/* Column styling */
.jump-number {
  font-family: var(--font-mono);
  font-weight: 600;
  color: var(--color-accent-teal);
}

.col-number {
  width: 60px;
}

.col-flags {
  text-align: right;
}

.flag {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  margin-left: 0.25rem;
  font-size: 0.8125rem;
}

.flag-cutaway {
  color: var(--color-danger);
}

.flag-fav {
  color: #f5c542;
}

.flag-night {
  color: var(--color-warning);
}

.flag-o2 {
  font-family: var(--font-mono);
  font-size: 0.6875rem;
  font-weight: 600;
  color: var(--color-accent-teal);
  background: rgba(20, 184, 166, 0.15);
  padding: 0.125rem 0.375rem;
  border-radius: 4px;
}
</style>
