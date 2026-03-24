<script setup>
import { ref, reactive, computed, watch } from 'vue'
import { useJumpStore } from '../stores/jumps.js'
import AutocompleteInput from './AutocompleteInput.vue'
import ConfirmModal from './ConfirmModal.vue'

const props = defineProps({
  jump: { type: Object, default: null }, // null = create mode
})

const emit = defineEmits(['close'])

const store = useJumpStore()

// ----- Helpers -----
const JUMP_TYPES = [
  'FF', 'WS', 'FS', 'CRW', 'HOP', 'CF', 'AFF',
  'TANDEM', 'DEMO', 'XRW', 'ANGLE', 'TRACKING', 'CP', 'WINGSUIT', 'OTHER',
]

const LANDING_OPTIONS = ['Stand-up', 'Sliding', 'PLF', 'Off-DZ', 'Water']

const isEdit = computed(() => !!props.jump)
const modalTitle = computed(() =>
  isEdit.value ? `Edit Jump #${props.jump.number}` : 'New Jump',
)

// ----- Form state -----
function toLocalDatetime(isoDate) {
  if (!isoDate) return ''
  // Format as YYYY-MM-DDTHH:mm for datetime-local input
  const d = new Date(isoDate)
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`
}

function todayDatetime() {
  return toLocalDatetime(new Date().toISOString())
}

const form = reactive({
  number: props.jump?.number ?? '',
  date: toLocalDatetime(props.jump?.date) || todayDatetime(),
  dropzone: props.jump?.dropzone ?? '',
  jumpType: props.jump?.jumpType ?? 'FF',
  aircraft: props.jump?.aircraft ?? '',
  altitude: props.jump?.altitude ?? '',
  freefallTime: props.jump?.freefallTime ?? '',
  canopySize: props.jump?.canopySize ?? '',
  lo: props.jump?.lo ?? '',
  event: props.jump?.event ?? '',
  landing: props.jump?.landing ?? '',
  nightJump: props.jump?.nightJump ?? false,
  oxygenJump: props.jump?.oxygenJump ?? false,
  cutaway: props.jump?.cutaway ?? false,
  description: props.jump?.description ?? '',
})

// ----- State -----
const saving = ref(false)
const saveError = ref('')
const showDeleteConfirm = ref(false)
const deleting = ref(false)
const isDirty = ref(false)

watch(form, () => { isDirty.value = true }, { deep: true })

// ----- Validation -----
const errors = reactive({ date: '', dropzone: '', jumpType: '' })

function validate() {
  errors.date = form.date ? '' : 'Date is required'
  errors.dropzone = form.dropzone.trim() ? '' : 'Dropzone is required'
  errors.jumpType = form.jumpType ? '' : 'Jump type is required'
  return !errors.date && !errors.dropzone && !errors.jumpType
}

// ----- Build payload -----
function buildPayload() {
  const payload = {
    date: new Date(form.date).toISOString(),
    dropzone: form.dropzone.trim(),
    jumpType: form.jumpType,
    aircraft: form.aircraft.trim() || undefined,
    altitude: form.altitude !== '' ? Number(form.altitude) : undefined,
    freefallTime: form.freefallTime !== '' ? Number(form.freefallTime) : undefined,
    canopySize: form.canopySize !== '' ? Number(form.canopySize) : undefined,
    lo: form.lo.trim() || undefined,
    event: form.event.trim() || undefined,
    landing: form.landing || undefined,
    nightJump: form.nightJump,
    oxygenJump: form.oxygenJump,
    cutaway: form.cutaway,
    description: form.description.trim() || undefined,
  }

  // Number field: only set if user specified one (create with insert, or edit with reposition)
  if (form.number !== '' && Number(form.number) > 0) {
    payload.number = Number(form.number)
  }

  return payload
}

// ----- Submit -----
async function submit() {
  if (!validate()) return
  saving.value = true
  saveError.value = ''
  try {
    if (isEdit.value) {
      await store.updateJump(props.jump.id, buildPayload())
    } else {
      await store.createJump(buildPayload())
    }
    emit('close', { success: true })
  } catch (err) {
    saveError.value = err.message || 'Something went wrong'
  } finally {
    saving.value = false
  }
}

// ----- Delete -----
async function confirmDelete() {
  deleting.value = true
  try {
    await store.deleteJump(props.jump.id)
    emit('close', { success: true, deleted: true })
  } catch (err) {
    saveError.value = err.message || 'Delete failed'
    showDeleteConfirm.value = false
  } finally {
    deleting.value = false
  }
}

// ----- Close / escape -----
function requestClose() {
  if (isDirty.value) {
    if (!confirm('You have unsaved changes. Close anyway?')) return
  }
  emit('close')
}

function onKey(e) {
  if (e.key === 'Escape') requestClose()
}
</script>

<template>
  <Teleport to="body">
    <div class="modal-overlay" @click.self="requestClose" @keydown="onKey">
      <div class="modal" role="dialog" aria-modal="true" :aria-labelledby="'modal-title'">
        <!-- Header -->
        <div class="modal-header">
          <h2 id="modal-title" class="modal-title">{{ modalTitle }}</h2>
          <button class="close-btn" aria-label="Close" @click="requestClose">✕</button>
        </div>

        <!-- Form -->
        <form class="modal-body" @submit.prevent="submit" novalidate>

          <!-- 1. Core section -->
          <fieldset class="section">
            <legend class="section-title">Core</legend>
            <div class="grid">
              <!-- Jump Number (edit/insert) -->
              <div class="field" v-if="isEdit">
                <label for="f-number" class="label">Jump #</label>
                <input id="f-number" type="number" class="form-input" v-model="form.number" min="1" />
              </div>
              <div class="field" v-else>
                <label for="f-number-insert" class="label">Insert at # <span class="hint">(optional)</span></label>
                <input id="f-number-insert" type="number" class="form-input" v-model="form.number" min="1" placeholder="Leave blank to append" />
              </div>

              <!-- Date -->
              <div class="field">
                <label for="f-date" class="label required">Date & Time</label>
                <input id="f-date" type="datetime-local" class="form-input" v-model="form.date" required />
                <span v-if="errors.date" class="field-error">{{ errors.date }}</span>
              </div>

              <!-- Dropzone -->
              <div class="field">
                <label for="f-dropzone" class="label required">Dropzone</label>
                <AutocompleteInput
                  id="f-dropzone"
                  field="dropzone"
                  v-model="form.dropzone"
                  placeholder="e.g. Skydive City"
                  required
                />
                <span v-if="errors.dropzone" class="field-error">{{ errors.dropzone }}</span>
              </div>

              <!-- Jump Type -->
              <div class="field">
                <label for="f-jump-type" class="label required">Jump Type</label>
                <select id="f-jump-type" class="form-input form-select" v-model="form.jumpType" required>
                  <option v-for="t in JUMP_TYPES" :key="t" :value="t">{{ t }}</option>
                </select>
                <span v-if="errors.jumpType" class="field-error">{{ errors.jumpType }}</span>
              </div>
            </div>
          </fieldset>

          <!-- 2. Details section -->
          <fieldset class="section">
            <legend class="section-title">Details</legend>
            <div class="grid">
              <div class="field">
                <label for="f-aircraft" class="label">Aircraft</label>
                <AutocompleteInput id="f-aircraft" field="aircraft" v-model="form.aircraft" placeholder="e.g. Cessna 182" />
              </div>
              <div class="field">
                <label for="f-altitude" class="label">Exit Altitude (ft)</label>
                <input id="f-altitude" type="number" class="form-input" v-model="form.altitude" min="0" placeholder="e.g. 14000" />
              </div>
              <div class="field">
                <label for="f-freefall" class="label">Freefall (s)</label>
                <input id="f-freefall" type="number" class="form-input" v-model="form.freefallTime" min="0" placeholder="e.g. 60" />
              </div>
              <div class="field">
                <label for="f-canopy" class="label">Canopy Size (sqft)</label>
                <input id="f-canopy" type="number" class="form-input" v-model="form.canopySize" min="0" placeholder="e.g. 190" />
              </div>
              <div class="field">
                <label for="f-landing" class="label">Landing</label>
                <select id="f-landing" class="form-input form-select" v-model="form.landing">
                  <option value="">— select —</option>
                  <option v-for="opt in LANDING_OPTIONS" :key="opt" :value="opt">{{ opt }}</option>
                </select>
              </div>
            </div>
          </fieldset>

          <!-- 3. People & Events -->
          <fieldset class="section">
            <legend class="section-title">People & Events</legend>
            <div class="grid">
              <div class="field">
                <label for="f-lo" class="label">Load Organizer / Coach</label>
                <AutocompleteInput id="f-lo" field="lo" v-model="form.lo" placeholder="Name" />
              </div>
              <div class="field">
                <label for="f-event" class="label">Event</label>
                <AutocompleteInput id="f-event" field="event" v-model="form.event" placeholder="e.g. Summerfest 2025" />
              </div>
            </div>
          </fieldset>

          <!-- 4. Flags -->
          <fieldset class="section">
            <legend class="section-title">Flags</legend>
            <div class="flags-grid">
              <label class="toggle">
                <input type="checkbox" v-model="form.nightJump" />
                <span class="toggle-label">🌙 Night</span>
              </label>
              <label class="toggle">
                <input type="checkbox" v-model="form.oxygenJump" />
                <span class="toggle-label">O₂ Oxygen</span>
              </label>
              <label class="toggle">
                <input type="checkbox" v-model="form.cutaway" />
                <span class="toggle-label">✂ Cutaway</span>
              </label>
            </div>
          </fieldset>

          <!-- 5. Notes -->
          <fieldset class="section">
            <legend class="section-title">Notes</legend>
            <textarea
              class="form-input notes-input"
              v-model="form.description"
              placeholder="Freeform debrief / notes…"
              rows="3"
            />
          </fieldset>

          <!-- Save error -->
          <p v-if="saveError" class="save-error">{{ saveError }}</p>

          <!-- Footer -->
          <div class="modal-footer">
            <button
              v-if="isEdit"
              type="button"
              class="btn-danger"
              :disabled="saving || deleting"
              @click="showDeleteConfirm = true"
            >
              Delete
            </button>
            <div class="footer-right">
              <button type="button" class="btn-secondary" :disabled="saving" @click="requestClose">Cancel</button>
              <button type="submit" class="btn-primary" :disabled="saving">
                <span v-if="saving" class="spinner" />
                {{ isEdit ? 'Save' : 'Create Jump' }}
              </button>
            </div>
          </div>
        </form>
      </div>
    </div>

    <!-- Delete Confirmation -->
    <ConfirmModal
      v-if="showDeleteConfirm"
      :title="`Delete Jump #${props.jump?.number}?`"
      message="This will renumber all subsequent jumps. This action cannot be undone."
      confirm-text="Delete"
      :danger="true"
      :loading="deleting"
      @confirm="confirmDelete"
      @cancel="showDeleteConfirm = false"
    />
  </Teleport>
</template>

<style scoped>
/* Overlay */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.7);
  z-index: 1500;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding: 2rem 1rem;
  overflow-y: auto;
}

/* Modal */
.modal {
  background: var(--color-surface-800);
  border: 1px solid var(--color-surface-600);
  border-radius: 12px;
  width: 100%;
  max-width: 700px;
  box-shadow: 0 24px 80px rgba(0, 0, 0, 0.6);
  animation: modal-in 0.18s ease;
  overflow: hidden;
}

@keyframes modal-in {
  from { opacity: 0; transform: translateY(-12px); }
  to   { opacity: 1; transform: translateY(0); }
}

/* Header */
.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.25rem 1.5rem;
  border-bottom: 1px solid var(--color-surface-600);
}

.modal-title {
  font-size: 1.125rem;
  font-weight: 700;
  color: var(--color-text-primary);
  margin: 0;
}

.close-btn {
  background: none;
  border: none;
  color: var(--color-text-muted);
  cursor: pointer;
  font-size: 1rem;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: background 0.15s, color 0.15s;
}

.close-btn:hover {
  background: var(--color-surface-700);
  color: var(--color-text-primary);
}

/* Body */
.modal-body {
  padding: 1.25rem 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

/* Sections */
.section {
  border: none;
  padding: 0;
  margin: 0;
}

.section-title {
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--color-text-muted);
  margin-bottom: 0.75rem;
}

/* Grid */
.grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.75rem 1rem;
}

/* Field */
.field {
  display: flex;
  flex-direction: column;
  gap: 0.375rem;
}

.label {
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--color-text-secondary);
}

.label.required::after {
  content: ' *';
  color: var(--color-accent-orange);
}

.hint {
  font-weight: 400;
  color: var(--color-text-muted);
}

/* Inputs */
.form-input {
  background: var(--color-surface-700);
  border: 1px solid var(--color-surface-600);
  border-radius: 8px;
  color: var(--color-text-primary);
  padding: 0.625rem 0.75rem;
  font-size: 0.875rem;
  min-height: 44px;
  transition: border-color 0.15s;
  width: 100%;
  box-sizing: border-box;
}

.form-input:focus {
  outline: none;
  border-color: var(--color-accent-teal);
  box-shadow: 0 0 0 2px rgba(94, 234, 212, 0.15);
}

.form-select {
  cursor: pointer;
  appearance: auto;
  color-scheme: dark;
}

.notes-input {
  resize: vertical;
  min-height: 80px;
}

.field-error {
  font-size: 0.75rem;
  color: #f87171;
}

/* Flags */
.flags-grid {
  display: flex;
  gap: 1.25rem;
  flex-wrap: wrap;
}

.toggle {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  min-height: 44px;
}

.toggle input[type="checkbox"] {
  width: 18px;
  height: 18px;
  accent-color: var(--color-accent-teal);
  cursor: pointer;
}

.toggle-label {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  user-select: none;
}

/* Error */
.save-error {
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 8px;
  padding: 0.625rem 0.875rem;
  font-size: 0.875rem;
  color: #f87171;
  margin: 0;
}

/* Footer */
.modal-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 0.5rem;
  border-top: 1px solid var(--color-surface-600);
  margin-top: 0.25rem;
}

.footer-right {
  display: flex;
  gap: 0.75rem;
  margin-left: auto;
}

.btn-danger {
  background: rgba(239, 68, 68, 0.15);
  color: #f87171;
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 8px;
  padding: 0.625rem 1rem;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  min-height: 44px;
  transition: background 0.15s, color 0.15s;
}

.btn-danger:hover {
  background: rgba(239, 68, 68, 0.25);
  color: #fca5a5;
}

.spinner {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255,255,255,0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
  display: inline-block;
  vertical-align: middle;
  margin-right: 4px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Mobile: full screen sheet */
@media (max-width: 639px) {
  .modal-overlay {
    align-items: flex-end;
    padding: 0;
  }

  .modal {
    border-radius: 12px 12px 0 0;
    max-height: 92dvh;
    overflow-y: auto;
  }

  .grid {
    grid-template-columns: 1fr;
  }
}
</style>
