import { describe, it, expect, beforeEach, vi, afterEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useJumpStore } from './jumps.js'

// Mock the api module so we can control responses without a real backend
vi.mock('../api.js', () => ({
  api: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}))

import { api } from '../api.js'

const mockJump = {
  id: 1,
  number: 1,
  date: '2025-01-01',
  dropzone: 'Skydive Empuriabrava',
  jumpType: 'FF',
  aircraft: 'Caravan',
}

const mockListResponse = {
  items: [mockJump],
  total: 1,
  totalPages: 1,
}

describe('useJumpStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('initializes with empty state', () => {
    const store = useJumpStore()
    expect(store.items).toEqual([])
    expect(store.total).toBe(0)
    expect(store.page).toBe(1)
    expect(store.perPage).toBe(25)
    expect(store.loading).toBe(false)
    expect(store.error).toBeNull()
  })

  it('fetchJumps populates items and total on success', async () => {
    api.get.mockResolvedValueOnce(mockListResponse)
    const store = useJumpStore()
    await store.fetchJumps()
    expect(store.items).toEqual([mockJump])
    expect(store.total).toBe(1)
    expect(store.totalPages).toBe(1)
    expect(store.loading).toBe(false)
    expect(store.error).toBeNull()
  })

  it('fetchJumps sets error and clears items on failure', async () => {
    api.get.mockRejectedValueOnce(new Error('Network error'))
    const store = useJumpStore()
    await store.fetchJumps()
    expect(store.items).toEqual([])
    expect(store.error).toBe('Network error')
    expect(store.loading).toBe(false)
  })

  it('setSort toggles order when same field is selected again', () => {
    const store = useJumpStore()
    store.sortBy = 'number'
    store.order = 'asc'
    store.setSort('number')
    expect(store.order).toBe('desc')
    store.setSort('number')
    expect(store.order).toBe('asc')
  })

  it('setSort changes field and resets order', () => {
    const store = useJumpStore()
    store.setSort('date')
    expect(store.sortBy).toBe('date')
    expect(store.order).toBe('asc')
  })

  it('setPage respects boundaries', () => {
    const store = useJumpStore()
    store.totalPages = 5
    store.setPage(3)
    expect(store.page).toBe(3)
    store.setPage(0) // below 1 — should be rejected
    expect(store.page).toBe(3)
    store.setPage(6) // above totalPages — should be rejected
    expect(store.page).toBe(3)
  })

  it('setPerPage resets page to 1', () => {
    const store = useJumpStore()
    store.page = 3
    store.setPerPage(50)
    expect(store.perPage).toBe(50)
    expect(store.page).toBe(1)
  })

  it('resetFilters clears all filter fields', () => {
    const store = useJumpStore()
    store.filters.q = 'empuriabrava'
    store.filters.jumpType = 'WS'
    store.page = 4
    store.resetFilters()
    expect(store.filters.q).toBe('')
    expect(store.filters.jumpType).toBe('')
    expect(store.page).toBe(1)
  })

  it('hasActiveFilters returns true when any filter is set', () => {
    const store = useJumpStore()
    expect(store.hasActiveFilters()).toBe(false)
    store.filters.q = 'test'
    expect(store.hasActiveFilters()).toBe(true)
    store.resetFilters()
    store.filters.cutaway = false
    expect(store.hasActiveFilters()).toBe(true)
  })

  it('initFromQuery populates state from URL params', () => {
    const store = useJumpStore()
    store.initFromQuery({
      page: '2',
      per_page: '50',
      sort: 'date',
      order: 'asc',
      q: 'empuriabrava',
      jump_type: 'WS',
    })
    expect(store.page).toBe(2)
    expect(store.perPage).toBe(50)
    expect(store.sortBy).toBe('date')
    expect(store.order).toBe('asc')
    expect(store.filters.q).toBe('empuriabrava')
    expect(store.filters.jumpType).toBe('WS')
  })

  it('toQuery omits default values', () => {
    const store = useJumpStore()
    // defaults: page=1, perPage=25, sort=number, order=desc
    const q = store.toQuery()
    expect(q).not.toHaveProperty('page')
    expect(q).not.toHaveProperty('per_page')
    expect(q).not.toHaveProperty('sort')
    expect(q).not.toHaveProperty('order')
  })

  it('toQuery includes non-default values', () => {
    const store = useJumpStore()
    store.page = 2
    store.filters.q = 'tunnel'
    const q = store.toQuery()
    expect(q.page).toBe('2')
    expect(q.q).toBe('tunnel')
  })

  it('createJump calls api.post then refreshes list', async () => {
    api.post.mockResolvedValueOnce(mockJump)
    api.get.mockResolvedValue(mockListResponse)
    const store = useJumpStore()
    const result = await store.createJump({ date: '2025-01-01', jumpType: 'FF' })
    expect(api.post).toHaveBeenCalledWith('/jumps', { date: '2025-01-01', jumpType: 'FF' })
    expect(result).toEqual(mockJump)
    expect(api.get).toHaveBeenCalled() // triggered fetchJumps
  })

  it('deleteJump calls api.delete then refreshes list', async () => {
    api.delete.mockResolvedValueOnce(undefined)
    api.get.mockResolvedValue(mockListResponse)
    const store = useJumpStore()
    await store.deleteJump(1)
    expect(api.delete).toHaveBeenCalledWith('/jumps/1')
    expect(api.get).toHaveBeenCalled()
  })
})
