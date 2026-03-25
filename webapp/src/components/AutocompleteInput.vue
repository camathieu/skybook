<script setup>
import { ref, watch, computed, onBeforeUnmount } from 'vue'
import { api } from '../api.js'

const props = defineProps({
  field: { type: String, required: true },
  modelValue: { type: String, default: '' },
  placeholder: { type: String, default: '' },
  required: { type: Boolean, default: false },
  id: { type: String, default: '' },
  // Static options array — when provided, use client-side filtering instead of API calls
  options: { type: Array, default: null },
  // Pass-through to <input inputmode> for mobile keyboard control (e.g. 'numeric')
  inputmode: { type: String, default: '' },
})

const emit = defineEmits(['update:modelValue'])

const inputRef = ref(null)
const suggestions = ref([])
const open = ref(false)
const activeIndex = ref(-1)

let debounceTimer = null

// Sync local value with modelValue
const localValue = ref(props.modelValue)
watch(() => props.modelValue, (v) => { localValue.value = v })

// Whether we're in static-options mode (no API call)
const isStaticMode = computed(() => Array.isArray(props.options))

function filterStaticOptions(prefix) {
  if (!prefix) return props.options
  const lower = prefix.toLowerCase()
  return props.options.filter(o => o.toLowerCase().startsWith(lower))
}

async function fetchSuggestions(prefix) {
  if (isStaticMode.value) {
    suggestions.value = filterStaticOptions(prefix)
    open.value = suggestions.value.length > 0
    return
  }
  try {
    const url = prefix
      ? `/jumps/autocomplete/${props.field}?q=${encodeURIComponent(prefix)}`
      : `/jumps/autocomplete/${props.field}`
    const data = await api.get(url)
    suggestions.value = data || []
    open.value = suggestions.value.length > 0
  } catch {
    suggestions.value = []
    open.value = false
  }
}

async function onFocus() {
  // Show recent values immediately on focus (even if the field is empty)
  activeIndex.value = -1
  clearTimeout(debounceTimer)
  await fetchSuggestions(localValue.value.trim())
}

async function onInput(e) {
  const val = e.target.value
  localValue.value = val
  emit('update:modelValue', val)
  activeIndex.value = -1

  clearTimeout(debounceTimer)
  if (isStaticMode.value) {
    // No debounce needed for client-side filtering
    suggestions.value = filterStaticOptions(val.trim())
    open.value = suggestions.value.length > 0
  } else {
    // Intentional: clearing the field re-fetches all recent values (empty prefix → show all).
    // This makes the dropdown act as a "recent values" list when the field is blank.
    debounceTimer = setTimeout(() => fetchSuggestions(val.trim()), 200)
  }
}

function select(val) {
  localValue.value = val
  emit('update:modelValue', val)
  suggestions.value = []
  open.value = false
  activeIndex.value = -1
}

function onKeydown(e) {
  if (!open.value) return
  if (e.key === 'ArrowDown') {
    e.preventDefault()
    activeIndex.value = Math.min(activeIndex.value + 1, suggestions.value.length - 1)
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    activeIndex.value = Math.max(activeIndex.value - 1, 0)
  } else if (e.key === 'Enter') {
    e.preventDefault()
    if (activeIndex.value >= 0) select(suggestions.value[activeIndex.value])
  } else if (e.key === 'Escape') {
    open.value = false
    activeIndex.value = -1
  }
}

function onBlur() {
  // Delay so click on suggestion registers first
  setTimeout(() => {
    open.value = false
    activeIndex.value = -1
  }, 150)
}

onBeforeUnmount(() => clearTimeout(debounceTimer))
</script>

<template>
  <div class="autocomplete-wrap">
    <input
      :id="id"
      ref="inputRef"
      type="text"
      class="form-input"
      :value="localValue"
      :placeholder="placeholder"
      :required="required"
      :inputmode="inputmode || undefined"
      autocomplete="off"
      @focus="onFocus"
      @input="onInput"
      @keydown="onKeydown"
      @blur="onBlur"
    />
    <ul v-if="open" class="suggestions" role="listbox">
      <li
        v-for="(item, idx) in suggestions"
        :key="item"
        class="suggestion-item"
        :class="{ active: idx === activeIndex }"
        role="option"
        @mousedown.prevent="select(item)"
      >
        {{ item }}
      </li>
    </ul>
  </div>
</template>

<style scoped>
.autocomplete-wrap {
  position: relative;
  width: 100%;
}

.form-input {
  width: 100%;
  background: var(--color-surface-700);
  border: 1px solid var(--color-surface-600);
  border-radius: 8px;
  color: var(--color-text-primary);
  padding: 0.625rem 0.75rem;
  font-size: 0.875rem;
  min-height: 44px;
  transition: border-color 0.15s;
  box-sizing: border-box;
}

.form-input:focus {
  outline: none;
  border-color: var(--color-accent-teal);
  box-shadow: 0 0 0 2px rgba(94, 234, 212, 0.15);
}

.suggestions {
  position: absolute;
  top: calc(100% + 4px);
  left: 0;
  right: 0;
  background: var(--color-surface-800);
  border: 1px solid var(--color-surface-600);
  border-radius: 8px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.3);
  z-index: 1000;
  max-height: 200px;
  overflow-y: auto;
  padding: 0;
  margin: 0;
  list-style: none;
}

.suggestion-item {
  padding: 0.625rem 0.75rem;
  font-size: 0.875rem;
  color: var(--color-text-primary);
  cursor: pointer;
  transition: background 0.1s;
  min-height: 44px;
  display: flex;
  align-items: center;
}

.suggestion-item:hover,
.suggestion-item.active {
  background: var(--color-surface-700);
  color: var(--color-accent-teal);
}
</style>
