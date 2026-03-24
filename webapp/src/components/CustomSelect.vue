<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'

const props = defineProps({
  modelValue: {
    type: String,
    default: '',
  },
  options: {
    type: Array,
    default: () => [],
  },
  placeholder: {
    type: String,
    default: 'All',
  },
  ariaLabel: {
    type: String,
    default: 'Select option',
  },
})

const emit = defineEmits(['update:modelValue'])

const open = ref(false)
const triggerRef = ref(null)
const menuRef = ref(null)

const selectedLabel = computed(() =>
  props.options.find(o => o === props.modelValue) || props.modelValue || props.placeholder
)

function select(value) {
  emit('update:modelValue', value)
  open.value = false
}

function clear() {
  emit('update:modelValue', '')
  open.value = false
}

function handleOutsideClick(e) {
  if (
    triggerRef.value && !triggerRef.value.contains(e.target) &&
    menuRef.value && !menuRef.value.contains(e.target)
  ) {
    open.value = false
  }
}

onMounted(() => document.addEventListener('mousedown', handleOutsideClick))
onBeforeUnmount(() => document.removeEventListener('mousedown', handleOutsideClick))
</script>

<template>
  <div class="cselect" :class="{ 'cselect--active': modelValue }" role="combobox" :aria-expanded="open" :aria-label="ariaLabel">
    <button
      ref="triggerRef"
      class="cselect__trigger"
      type="button"
      :aria-label="ariaLabel"
      aria-haspopup="listbox"
      data-testid="custom-select-trigger"
      @click="open = !open"
    >
      <span class="cselect__label" :class="{ 'cselect__label--selected': modelValue }">
        {{ selectedLabel }}
      </span>
      <svg class="cselect__chevron" :class="{ 'cselect__chevron--open': open }" width="12" height="12" viewBox="0 0 12 12" fill="none">
        <path d="M2 4L6 8L10 4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
    </button>

    <Transition name="dropdown">
      <div v-if="open" ref="menuRef" class="cselect__menu" role="listbox" data-testid="custom-select-menu">
        <!-- Clear / All option -->
        <button
          class="cselect__option"
          :class="{ 'cselect__option--selected': !modelValue }"
          type="button"
          role="option"
          :aria-selected="!modelValue"
          data-testid="custom-select-clear"
          @click="clear"
        >
          <span class="cselect__option-label">{{ placeholder }}</span>
          <svg v-if="!modelValue" class="cselect__check" width="12" height="12" viewBox="0 0 12 12" fill="none">
            <path d="M2 6L5 9L10 3" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </button>

        <div v-if="options.length" class="cselect__divider" />

        <div class="cselect__scroll">
          <button
            v-for="opt in options"
            :key="opt"
            class="cselect__option"
            :class="{ 'cselect__option--selected': modelValue === opt }"
            type="button"
            role="option"
            :aria-selected="modelValue === opt"
            :data-testid="`custom-select-option-${opt}`"
            @click="select(opt)"
          >
            <span class="cselect__option-label">{{ opt }}</span>
            <svg v-if="modelValue === opt" class="cselect__check" width="12" height="12" viewBox="0 0 12 12" fill="none">
              <path d="M2 6L5 9L10 3" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </button>

          <div v-if="!options.length" class="cselect__empty">No options</div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.cselect {
  position: relative;
  display: inline-block;
}

/* Trigger button */
.cselect__trigger {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  min-height: 40px;
  padding: 0.5rem 0.75rem;
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--color-text-muted);
  background-color: var(--color-surface-700);
  border: 1px solid var(--color-surface-600);
  border-radius: 8px;
  cursor: pointer;
  white-space: nowrap;
  transition: border-color 0.15s, color 0.15s, background-color 0.15s;
  user-select: none;
  -webkit-tap-highlight-color: transparent;
}

.cselect__trigger:hover {
  border-color: var(--color-surface-500);
  color: var(--color-text-secondary);
}

.cselect--active .cselect__trigger {
  border-color: var(--color-accent-teal);
  color: var(--color-accent-teal);
  background-color: rgba(20, 184, 166, 0.08);
}

.cselect__label {
  color: inherit;
}

.cselect__label--selected {
  font-weight: 600;
}

/* Chevron */
.cselect__chevron {
  color: var(--color-text-muted);
  flex-shrink: 0;
  transition: transform 0.2s;
}

.cselect__chevron--open {
  transform: rotate(180deg);
}

/* Dropdown menu */
.cselect__menu {
  position: absolute;
  top: calc(100% + 6px);
  left: 0;
  z-index: 100;
  min-width: 180px;
  padding-top: 0.25rem;
  background-color: var(--color-surface-800, #1e2330);
  border: 1px solid var(--color-surface-600, #3a4055);
  border-radius: 10px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4), 0 2px 8px rgba(0, 0, 0, 0.3);
  overflow: hidden;
}

/* Scrollable list capped at ~10 items (40px each) */
.cselect__scroll {
  max-height: 400px;
  overflow-y: auto;
  /* padding-bottom ensures the last item has spacing even when scrolled to bottom;
     this must be on the scroll child (not the overflow:hidden parent) to avoid clipping */
  padding-bottom: 0.25rem;
  scrollbar-width: thin;
  scrollbar-color: var(--color-surface-600) transparent;
}

.cselect__scroll::-webkit-scrollbar {
  width: 4px;
}

.cselect__scroll::-webkit-scrollbar-track {
  background: transparent;
}

.cselect__scroll::-webkit-scrollbar-thumb {
  background-color: var(--color-surface-600);
  border-radius: 4px;
}

.cselect__divider {
  height: 1px;
  background-color: var(--color-surface-700, #252b3b);
  margin: 0.25rem 0;
}

/* Option */
.cselect__option {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: 0.5625rem 0.875rem;
  font-size: 0.8125rem;
  color: var(--color-text-secondary);
  background: none;
  border: none;
  cursor: pointer;
  text-align: left;
  transition: background-color 0.1s, color 0.1s;
  gap: 0.5rem;
}

.cselect__option:hover {
  background-color: var(--color-surface-700, #252b3b);
  color: var(--color-text-primary);
}

.cselect__option--selected {
  color: var(--color-accent-teal);
  font-weight: 600;
}

.cselect__option--selected:hover {
  color: var(--color-accent-teal);
}

.cselect__option-label {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.cselect__check {
  flex-shrink: 0;
  color: var(--color-accent-teal);
}

.cselect__empty {
  padding: 0.5rem 0.875rem;
  font-size: 0.8125rem;
  color: var(--color-text-muted);
  font-style: italic;
}

/* Dropdown animation */
.dropdown-enter-active,
.dropdown-leave-active {
  transition: opacity 0.15s, transform 0.15s;
  transform-origin: top left;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: scaleY(0.9) translateY(-4px);
}

/* Mobile */
@media (max-width: 767px) {
  .cselect__trigger {
    min-height: 44px;
    width: 100%;
    justify-content: space-between;
  }

  .cselect {
    width: 100%;
  }

  .cselect__menu {
    width: 100%;
    max-height: 260px;
    overflow-y: auto;
  }

  .cselect__scroll {
    max-height: 200px;
  }
}
</style>
