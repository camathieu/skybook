import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import CustomSelect from './CustomSelect.vue'

// Stub document event listeners used for outside-click
const addSpy = vi.spyOn(document, 'addEventListener')
const removeSpy = vi.spyOn(document, 'removeEventListener')

const OPTIONS = ['Dropzone A', 'Dropzone B', 'Dropzone C']

function mountSelect(props = {}) {
  return mount(CustomSelect, {
    props: {
      modelValue: '',
      options: OPTIONS,
      placeholder: 'All dropzones',
      ariaLabel: 'Filter by dropzone',
      ...props,
    },
    attachTo: document.body,
  })
}

describe('CustomSelect', () => {
  it('renders the trigger button with placeholder when no value selected', () => {
    const wrapper = mountSelect()
    expect(wrapper.find('[data-testid="custom-select-trigger"]').text()).toContain('All dropzones')
  })

  it('renders selected label when modelValue is set', () => {
    const wrapper = mountSelect({ modelValue: 'Dropzone B' })
    expect(wrapper.find('[data-testid="custom-select-trigger"]').text()).toContain('Dropzone B')
  })

  it('dropdown is closed initially', () => {
    const wrapper = mountSelect()
    expect(wrapper.find('[data-testid="custom-select-menu"]').exists()).toBe(false)
  })

  it('opens dropdown on trigger click', async () => {
    const wrapper = mountSelect()
    await wrapper.find('[data-testid="custom-select-trigger"]').trigger('click')
    expect(wrapper.find('[data-testid="custom-select-menu"]').exists()).toBe(true)
  })

  it('closes dropdown on second trigger click', async () => {
    const wrapper = mountSelect()
    const trigger = wrapper.find('[data-testid="custom-select-trigger"]')
    await trigger.trigger('click')
    await trigger.trigger('click')
    expect(wrapper.find('[data-testid="custom-select-menu"]').exists()).toBe(false)
  })

  it('renders all options when open', async () => {
    const wrapper = mountSelect()
    await wrapper.find('[data-testid="custom-select-trigger"]').trigger('click')
    for (const opt of OPTIONS) {
      expect(wrapper.find(`[data-testid="custom-select-option-${opt}"]`).exists()).toBe(true)
    }
  })

  it('emits update:modelValue with option value when option clicked', async () => {
    const wrapper = mountSelect()
    await wrapper.find('[data-testid="custom-select-trigger"]').trigger('click')
    await wrapper.find('[data-testid="custom-select-option-Dropzone B"]').trigger('click')
    expect(wrapper.emitted('update:modelValue')).toBeTruthy()
    expect(wrapper.emitted('update:modelValue')[0]).toEqual(['Dropzone B'])
  })

  it('closes dropdown after selecting an option', async () => {
    const wrapper = mountSelect()
    await wrapper.find('[data-testid="custom-select-trigger"]').trigger('click')
    await wrapper.find('[data-testid="custom-select-option-Dropzone A"]').trigger('click')
    expect(wrapper.find('[data-testid="custom-select-menu"]').exists()).toBe(false)
  })

  it('emits empty string when clear (All) option is clicked', async () => {
    const wrapper = mountSelect({ modelValue: 'Dropzone A' })
    await wrapper.find('[data-testid="custom-select-trigger"]').trigger('click')
    await wrapper.find('[data-testid="custom-select-clear"]').trigger('click')
    expect(wrapper.emitted('update:modelValue')[0]).toEqual([''])
  })

  it('shows checkmark on currently selected option', async () => {
    const wrapper = mountSelect({ modelValue: 'Dropzone C' })
    await wrapper.find('[data-testid="custom-select-trigger"]').trigger('click')
    const selectedBtn = wrapper.find('[data-testid="custom-select-option-Dropzone C"]')
    expect(selectedBtn.classes()).toContain('cselect__option--selected')
  })

  it('registers and unregisters outside-click listener', () => {
    const wrapper = mountSelect()
    expect(addSpy).toHaveBeenCalledWith('mousedown', expect.any(Function))
    wrapper.unmount()
    expect(removeSpy).toHaveBeenCalledWith('mousedown', expect.any(Function))
  })
})
