import { describe, it, expect, vi, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import BaseModal from './BaseModal.vue'

describe('BaseModal', () => {
  let wrapper

  afterEach(() => {
    if (wrapper) wrapper.unmount()
    wrapper = null
  })

  function mountModal(props = {}, slots = {}) {
    wrapper = mount(BaseModal, {
      props: { dataTestid: 'test-modal', ...props },
      slots: {
        default: '<div class="dialog-content">Hello</div>',
        ...slots,
      },
      attachTo: document.body,
    })
    return wrapper
  }

  // Teleport sends content to <body>, so query from document.body
  function findInBody(selector) {
    return document.body.querySelector(selector)
  }

  it('renders slot content inside the dialog', () => {
    mountModal()
    expect(findInBody('.dialog-content').textContent).toBe('Hello')
  })

  it('emits close when clicking on the overlay (not the dialog)', async () => {
    mountModal()
    findInBody('[data-testid="test-modal"]').click()
    expect(wrapper.emitted('close')).toHaveLength(1)
  })

  it('does NOT emit close when clicking inside the dialog', async () => {
    mountModal()
    findInBody('.dialog-content').click()
    expect(wrapper.emitted('close')).toBeUndefined()
  })

  it('emits close on Escape key', () => {
    mountModal()
    document.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape' }))
    expect(wrapper.emitted('close')).toHaveLength(1)
  })

  it('applies custom z-index', () => {
    mountModal({ zIndex: 3000 })
    const overlay = findInBody('[data-testid="test-modal"]')
    expect(overlay.style.zIndex).toBe('3000')
  })

  it('applies custom max-width to dialog', () => {
    mountModal({ maxWidth: '400px' })
    const dialog = findInBody('.base-dialog')
    expect(dialog.style.maxWidth).toBe('400px')
  })

  it('applies data-testid attribute', () => {
    mountModal({ dataTestid: 'my-modal' })
    expect(findInBody('[data-testid="my-modal"]')).not.toBeNull()
  })

  it('registers and cleans up keydown listener', () => {
    const addSpy = vi.spyOn(document, 'addEventListener')
    const removeSpy = vi.spyOn(document, 'removeEventListener')
    mountModal()
    expect(addSpy).toHaveBeenCalledWith('keydown', expect.any(Function))
    wrapper.unmount()
    wrapper = null
    expect(removeSpy).toHaveBeenCalledWith('keydown', expect.any(Function))
    addSpy.mockRestore()
    removeSpy.mockRestore()
  })
})
