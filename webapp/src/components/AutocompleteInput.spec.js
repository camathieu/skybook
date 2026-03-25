import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import AutocompleteInput from './AutocompleteInput.vue'

// Mock the api module
vi.mock('../api.js', () => ({
  api: {
    get: vi.fn(),
  },
}))
import { api } from '../api.js'

describe('AutocompleteInput', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('renders the input element', () => {
    const wrapper = mount(AutocompleteInput, {
      props: { field: 'dropzone', modelValue: '' },
    })
    expect(wrapper.find('input').exists()).toBe(true)
  })

  it('emits update:modelValue on input', async () => {
    const wrapper = mount(AutocompleteInput, {
      props: { field: 'dropzone', modelValue: '' },
    })
    const input = wrapper.find('input')
    await input.setValue('Emp')
    expect(wrapper.emitted('update:modelValue')).toBeTruthy()
    expect(wrapper.emitted('update:modelValue')[0]).toEqual(['Emp'])
  })

  it('does not show suggestions list initially', () => {
    const wrapper = mount(AutocompleteInput, {
      props: { field: 'dropzone', modelValue: '' },
    })
    expect(wrapper.find('ul.suggestions').exists()).toBe(false)
  })

  it('shows suggestions after debounced API call', async () => {
    api.get.mockResolvedValueOnce(['Empuriabrava', 'Elsinore'])
    const wrapper = mount(AutocompleteInput, {
      props: { field: 'dropzone', modelValue: '' },
    })
    const input = wrapper.find('input')
    await input.setValue('Emp')
    // advance past the 200ms debounce
    vi.advanceTimersByTime(250)
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()
    expect(api.get).toHaveBeenCalledWith('/jumps/autocomplete/dropzone?q=Emp')
  })

  it('re-fetches all suggestions when input is cleared', async () => {
    api.get.mockResolvedValueOnce(['Empuriabrava', 'Perris'])
    const wrapper = mount(AutocompleteInput, {
      props: { field: 'dropzone', modelValue: 'Emp' },
    })
    const input = wrapper.find('input')
    await input.setValue('')
    vi.advanceTimersByTime(250)
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()
    // Clearing the field fetches all recent values (no q param)
    expect(api.get).toHaveBeenCalledWith('/jumps/autocomplete/dropzone')
  })

  it('shows suggestions on focus without typing (on-focus recent values)', async () => {
    api.get.mockResolvedValueOnce(['Empuriabrava', 'Perris', 'DeLand'])
    const wrapper = mount(AutocompleteInput, {
      props: { field: 'dropzone', modelValue: '' },
    })
    const input = wrapper.find('input')
    await input.trigger('focus')
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()
    // Should call API with no prefix (empty field)
    expect(api.get).toHaveBeenCalledWith('/jumps/autocomplete/dropzone')
  })

  it('closes dropdown on Escape key', async () => {
    api.get.mockResolvedValueOnce(['Empuriabrava'])
    const wrapper = mount(AutocompleteInput, {
      props: { field: 'dropzone', modelValue: '' },
    })

    // Manually set open state to true to simulate open dropdown
    wrapper.vm.open = true
    wrapper.vm.suggestions = ['Empuriabrava']
    await wrapper.vm.$nextTick()

    await wrapper.find('input').trigger('keydown', { key: 'Escape' })
    expect(wrapper.vm.open).toBe(false)
  })

  // ---- Static options mode ----

  it('shows all static options on focus without making an API call', async () => {
    const wrapper = mount(AutocompleteInput, {
      props: { field: 'altitude', modelValue: '', options: ['5000', '10000', '15000'] },
    })
    await wrapper.find('input').trigger('focus')
    await wrapper.vm.$nextTick()
    await wrapper.vm.$nextTick()
    expect(api.get).not.toHaveBeenCalled()
    expect(wrapper.vm.suggestions).toEqual(['5000', '10000', '15000'])
    expect(wrapper.vm.open).toBe(true)
  })

  it('filters static options by typed prefix (client-side, no API call)', async () => {
    const wrapper = mount(AutocompleteInput, {
      props: { field: 'altitude', modelValue: '', options: ['5000', '10000', '12000', '13000', '15000', '20000'] },
    })
    await wrapper.find('input').setValue('1')
    await wrapper.vm.$nextTick()
    expect(api.get).not.toHaveBeenCalled()
    // Only options starting with "1" should be shown
    expect(wrapper.vm.suggestions).toEqual(['10000', '12000', '13000', '15000'])
  })

  it('allows emitting a custom value not in static options', async () => {
    const wrapper = mount(AutocompleteInput, {
      props: { field: 'altitude', modelValue: '', options: ['5000', '10000', '15000'] },
    })
    const input = wrapper.find('input')
    await input.setValue('12500')
    expect(wrapper.emitted('update:modelValue')).toBeTruthy()
    expect(wrapper.emitted('update:modelValue').at(-1)).toEqual(['12500'])
  })
})
