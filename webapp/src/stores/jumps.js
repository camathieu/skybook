import { defineStore } from 'pinia'
import { ref, reactive, watch } from 'vue'
import { api } from '../api.js'

export const useJumpStore = defineStore('jumps', () => {
  // --- State ---
  const items = ref([])
  const total = ref(0)
  const page = ref(1)
  const perPage = ref(25)
  const totalPages = ref(0)
  const loading = ref(false)
  const error = ref(null)

  const sortBy = ref('number')
  const order = ref('desc')

  const filters = reactive({
    q: '',
    dropzone: '',
    aircraft: '',
    jumpType: '',
    dateFrom: '',
    dateTo: '',
    cutaway: null,
    night: null,
  })

  // --- Actions ---

  async function fetchJumps() {
    loading.value = true
    error.value = null

    try {
      const params = new URLSearchParams()

      // Pagination
      params.set('page', String(page.value))
      params.set('per_page', String(perPage.value))

      // Sorting
      if (sortBy.value) params.set('sort', sortBy.value)
      if (order.value) params.set('order', order.value)

      // Filters
      if (filters.q) params.set('q', filters.q)
      if (filters.dropzone) params.set('dropzone', filters.dropzone)
      if (filters.aircraft) params.set('aircraft', filters.aircraft)
      if (filters.jumpType) params.set('jump_type', filters.jumpType)
      if (filters.dateFrom) params.set('date_from', filters.dateFrom)
      if (filters.dateTo) params.set('date_to', filters.dateTo)
      if (filters.cutaway !== null) params.set('cutaway', String(filters.cutaway))
      if (filters.night !== null) params.set('night', String(filters.night))

      const data = await api.get(`/jumps?${params.toString()}`)

      items.value = data.items
      total.value = data.total
      totalPages.value = data.totalPages

    } catch (err) {
      error.value = err.message || 'Failed to load jumps'
      items.value = []
      total.value = 0
      totalPages.value = 0
    } finally {
      loading.value = false
    }
  }

  async function createJump(data) {
    const jump = await api.post('/jumps', data)
    await fetchJumps() // refresh list to maintain exact sorting/numbering
    return jump
  }

  async function updateJump(id, data) {
    const jump = await api.put(`/jumps/${id}`, data)
    await fetchJumps()
    return jump
  }

  async function deleteJump(id) {
    await api.delete(`/jumps/${id}`)
    await fetchJumps()
  }

  function setSort(field) {
    if (sortBy.value === field) {
      order.value = order.value === 'asc' ? 'desc' : 'asc'
    } else {
      sortBy.value = field
      order.value = field === 'number' ? 'desc' : 'asc'
    }
  }

  function setPage(p) {
    if (p >= 1 && p <= totalPages.value) {
      page.value = p
    }
  }

  function setPerPage(pp) {
    perPage.value = pp
    page.value = 1
  }

  function resetFilters() {
    filters.q = ''
    filters.dropzone = ''
    filters.aircraft = ''
    filters.jumpType = ''
    filters.dateFrom = ''
    filters.dateTo = ''
    filters.cutaway = null
    filters.night = null
    page.value = 1
  }

  function hasActiveFilters() {
    return !!(
      filters.q ||
      filters.dropzone ||
      filters.aircraft ||
      filters.jumpType ||
      filters.dateFrom ||
      filters.dateTo ||
      filters.cutaway !== null ||
      filters.night !== null
    )
  }

  /**
   * Initialize store state from URL query params.
   */
  function initFromQuery(query) {
    if (query.page) page.value = parseInt(query.page, 10) || 1
    if (query.per_page) perPage.value = parseInt(query.per_page, 10) || 25
    const allowedSort = ['number', 'date', 'dropzone', 'altitude']
    if (query.sort && allowedSort.includes(query.sort)) sortBy.value = query.sort
    if (query.order && (query.order === 'asc' || query.order === 'desc')) order.value = query.order
    if (query.q) filters.q = query.q
    if (query.dropzone) filters.dropzone = query.dropzone
    if (query.jump_type) filters.jumpType = query.jump_type
    if (query.date_from) filters.dateFrom = query.date_from
    if (query.date_to) filters.dateTo = query.date_to
    if (query.cutaway) filters.cutaway = query.cutaway === 'true'
    if (query.night) filters.night = query.night === 'true'
  }

  /**
   * Build URL query params from current store state.
   */
  function toQuery() {
    const q = {}
    if (page.value > 1) q.page = String(page.value)
    if (perPage.value !== 25) q.per_page = String(perPage.value)
    if (sortBy.value !== 'number') q.sort = sortBy.value
    if (order.value !== 'desc') q.order = order.value
    if (filters.q) q.q = filters.q
    if (filters.dropzone) q.dropzone = filters.dropzone
    if (filters.jumpType) q.jump_type = filters.jumpType
    if (filters.dateFrom) q.date_from = filters.dateFrom
    if (filters.dateTo) q.date_to = filters.dateTo
    if (filters.cutaway !== null) q.cutaway = String(filters.cutaway)
    if (filters.night !== null) q.night = String(filters.night)
    return q
  }

  return {
    // State
    items,
    total,
    page,
    perPage,
    totalPages,
    loading,
    error,
    sortBy,
    order,
    filters,

    // Actions
    fetchJumps,
    createJump,
    updateJump,
    deleteJump,
    setSort,
    setPage,
    setPerPage,
    resetFilters,
    hasActiveFilters,
    initFromQuery,
    toQuery,
  }
})
