import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useToastStore } from './toast.js'

describe('useToastStore', () => {
  beforeEach(() => {
    vi.useFakeTimers()
    setActivePinia(createPinia())
  })

  it('addToast adds a toast with correct fields', () => {
    const store = useToastStore()
    store.addToast('Hello', 'success', 3000)
    expect(store.toasts).toHaveLength(1)
    expect(store.toasts[0]).toMatchObject({ message: 'Hello', type: 'success' })
    expect(store.toasts[0].id).toBeGreaterThan(0)
  })

  it('auto-removes toast after duration', () => {
    const store = useToastStore()
    store.addToast('Gone soon', 'info', 2000)
    expect(store.toasts).toHaveLength(1)
    vi.advanceTimersByTime(2000)
    expect(store.toasts).toHaveLength(0)
  })

  it('removeToast removes a specific toast by id', () => {
    const store = useToastStore()
    store.addToast('A', 'success', 0)
    store.addToast('B', 'error', 0)
    const idA = store.toasts[0].id
    store.removeToast(idA)
    expect(store.toasts).toHaveLength(1)
    expect(store.toasts[0].message).toBe('B')
  })

  it('enforces max 5 toasts by removing oldest', () => {
    const store = useToastStore()
    for (let i = 0; i < 7; i++) {
      store.addToast(`Toast ${i}`, 'info', 0)
    }
    expect(store.toasts).toHaveLength(5)
    // The two oldest (0, 1) should have been shifted out
    expect(store.toasts[0].message).toBe('Toast 2')
    expect(store.toasts[4].message).toBe('Toast 6')
  })

  it('defaults to success type and 4000ms duration', () => {
    const store = useToastStore()
    store.addToast('Default')
    expect(store.toasts[0].type).toBe('success')
    // Should still be visible before 4s
    vi.advanceTimersByTime(3999)
    expect(store.toasts).toHaveLength(1)
    // Should be gone after 4s
    vi.advanceTimersByTime(1)
    expect(store.toasts).toHaveLength(0)
  })
})
