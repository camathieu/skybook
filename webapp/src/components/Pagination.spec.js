import { describe, it, expect, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import Pagination from './Pagination.vue'
import { useJumpStore } from '../stores/jumps.js'

// Pagination.vue is tightly coupled to useJumpStore, so we test it through the store.
vi.mock('../api.js', () => ({ api: { get: vi.fn(), post: vi.fn(), put: vi.fn(), delete: vi.fn() } }))

describe('Pagination', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  function mountWithStore(overrides = {}) {
    const store = useJumpStore()
    Object.assign(store, {
      total: overrides.total ?? 100,
      page: overrides.page ?? 1,
      perPage: overrides.perPage ?? 25,
      totalPages: overrides.totalPages ?? 4,
    })
    return { wrapper: mount(Pagination), store }
  }

  it('renders nothing when total is 0', () => {
    const { wrapper } = mountWithStore({ total: 0 })
    expect(wrapper.find('.pagination').exists()).toBe(false)
  })

  it('renders pagination when total > 0', () => {
    const { wrapper } = mountWithStore({ total: 100 })
    expect(wrapper.find('.pagination').exists()).toBe(true)
  })

  it('shows correct "Showing X–Y of Z jumps" text', () => {
    const { wrapper } = mountWithStore({ total: 100, page: 2, perPage: 25 })
    const info = wrapper.find('.pagination-info').text()
    expect(info).toContain('26')  // showingFrom: (2-1)*25+1 = 26
    expect(info).toContain('50')  // showingTo: min(2*25, 100) = 50
    expect(info).toContain('100')
  })

  it('disables Prev button on first page', () => {
    const { wrapper } = mountWithStore({ page: 1, totalPages: 4 })
    const prevBtn = wrapper.find('[aria-label="Previous page"]')
    expect(prevBtn.attributes('disabled')).toBeDefined()
  })

  it('disables Next button on last page', () => {
    const { wrapper } = mountWithStore({ page: 4, totalPages: 4 })
    const nextBtn = wrapper.find('[aria-label="Next page"]')
    expect(nextBtn.attributes('disabled')).toBeDefined()
  })

  it('enables both buttons on a mid page', () => {
    const { wrapper } = mountWithStore({ page: 2, totalPages: 4 })
    const prevBtn = wrapper.find('[aria-label="Previous page"]')
    const nextBtn = wrapper.find('[aria-label="Next page"]')
    expect(prevBtn.attributes('disabled')).toBeUndefined()
    expect(nextBtn.attributes('disabled')).toBeUndefined()
  })

  it('Prev button calls store.setPage with page - 1', async () => {
    const { wrapper, store } = mountWithStore({ page: 3, totalPages: 5 })
    await wrapper.find('[aria-label="Previous page"]').trigger('click')
    expect(store.page).toBe(2)
  })

  it('Next button calls store.setPage with page + 1', async () => {
    const { wrapper, store } = mountWithStore({ page: 2, totalPages: 5 })
    await wrapper.find('[aria-label="Next page"]').trigger('click')
    expect(store.page).toBe(3)
  })

  it('per-page buttons render and mark active correctly', () => {
    const { wrapper } = mountWithStore({ perPage: 25 })
    const btns = wrapper.findAll('.per-page-btn')
    expect(btns).toHaveLength(3)  // 25, 50, 100
    const active = btns.filter(b => b.classes('active'))
    expect(active).toHaveLength(1)
    expect(active[0].text()).toBe('25')
  })

  it('clicking a per-page button updates store.perPage and resets to page 1', async () => {
    const { wrapper, store } = mountWithStore({ page: 3, perPage: 25, total: 100, totalPages: 4 })
    const btn50 = wrapper.findAll('.per-page-btn').find(b => b.text() === '50')
    await btn50.trigger('click')
    expect(store.perPage).toBe(50)
    expect(store.page).toBe(1)
  })
})
