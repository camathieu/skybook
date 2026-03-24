<script setup>
const props = defineProps({
  jump: { type: Object, required: true },
})
const emit = defineEmits(['edit'])

function formatDate(dateStr) {
  if (!dateStr) return '—'
  const d = new Date(dateStr)
  return d.toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric' })
}

function formatAltitude(alt) {
  if (alt == null) return null
  return alt.toLocaleString() + ' ft'
}
</script>

<template>
  <div class="card jump-card" data-testid="jump-card" tabindex="0" style="cursor:pointer" @click="emit('edit', jump)" @keydown.enter="emit('edit', jump)">
    <div class="card-header">
      <span class="jump-number">#{{ jump.number }}</span>
      <span class="jump-date">{{ formatDate(jump.date) }}</span>
    </div>

    <div class="card-body">
      <div class="card-field">
        <span class="field-label">Dropzone</span>
        <span class="field-value">{{ jump.dropzone || '—' }}</span>
      </div>
      <div class="card-field">
        <span class="field-label">Type</span>
        <span class="field-value"><span class="badge">{{ jump.jumpType }}</span></span>
      </div>
      <div v-if="formatAltitude(jump.altitude)" class="card-field">
        <span class="field-label">Altitude</span>
        <span class="field-value">{{ formatAltitude(jump.altitude) }}</span>
      </div>
      <div v-if="jump.aircraft" class="card-field">
        <span class="field-label">Aircraft</span>
        <span class="field-value">{{ jump.aircraft }}</span>
      </div>
    </div>

    <div v-if="jump.nightJump || jump.oxygenJump || jump.cutaway" class="card-flags">
      <span v-if="jump.nightJump" class="flag-badge flag-night">🌙 Night</span>
      <span v-if="jump.oxygenJump" class="flag-badge flag-o2">O₂</span>
      <span v-if="jump.cutaway" class="flag-badge flag-cutaway">✂ Cutaway</span>
    </div>
  </div>
</template>

<style scoped>
.jump-card {
  cursor: pointer;
  transition: border-color var(--transition-fast), transform var(--transition-fast);
  -webkit-tap-highlight-color: transparent;
}

.jump-card:active {
  transform: scale(0.98);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.75rem;
  padding-bottom: 0.625rem;
  border-bottom: 1px solid var(--color-surface-600);
}

.jump-number {
  font-family: var(--font-mono);
  font-weight: 700;
  font-size: 1.125rem;
  color: var(--color-accent-teal);
}

.jump-date {
  font-size: 0.8125rem;
  color: var(--color-text-muted);
}

.card-body {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.5rem;
}

.card-field {
  display: flex;
  flex-direction: column;
  gap: 0.125rem;
}

.field-label {
  font-size: 0.6875rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--color-text-muted);
}

.field-value {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
}

.card-flags {
  display: flex;
  gap: 0.5rem;
  margin-top: 0.75rem;
  padding-top: 0.625rem;
  border-top: 1px solid var(--color-surface-700);
  flex-wrap: wrap;
}

.flag-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.25rem 0.5rem;
  font-size: 0.75rem;
  font-weight: 500;
  border-radius: 6px;
  background-color: var(--color-surface-700);
  color: var(--color-text-secondary);
}

.flag-night {
  color: var(--color-warning);
}

.flag-o2 {
  color: var(--color-accent-teal);
  background: rgba(20, 184, 166, 0.15);
}

.flag-cutaway {
  color: var(--color-danger);
  background: rgba(239, 68, 68, 0.15);
}
</style>
