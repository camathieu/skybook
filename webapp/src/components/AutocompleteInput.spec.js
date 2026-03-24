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

  it('clears suggestions when input is emptied', async () => {
    const wrapper = mount(AutocompleteInput, {
      props: { field: 'dropzone', modelValue: 'Emp' },
    })
    const input = wrapper.find('input')
    await input.setValue('')
    vi.advanceTimersByTime(250)
    await wrapper.vm.$nextTick()
    expect(wrapper.find('ul.suggestions').exists()).toBe(false)
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
})
