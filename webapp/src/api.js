/**
 * SkyBook API client.
 *
 * Centralized HTTP client using fetch().
 * Base URL: /api/v1
 * In dev mode, Vite proxies /api/* to localhost:8080.
 */

const BASE_URL = '/api/v1'

/**
 * Make an API request and return parsed JSON.
 * @param {string} path - API path (e.g. '/jumps')
 * @param {object} options - fetch options
 * @returns {Promise<any>}
 */
async function request(path, options = {}) {
  const url = `${BASE_URL}${path}`

  const headers = {
    ...options.headers,
  }

  if (options.body) {
    headers['Content-Type'] = 'application/json'
  }

  const response = await fetch(url, { ...options, headers })

  if (!response.ok) {
    const body = await response.json().catch(() => ({}))
    const error = new Error(body.error || `Request failed: ${response.status}`)
    error.status = response.status
    error.body = body
    throw error
  }

  // 204 No Content
  if (response.status === 204) {
    return null
  }

  return response.json()
}

export const api = {
  get: (path) => request(path),
  post: (path, data) => request(path, { method: 'POST', body: JSON.stringify(data) }),
  put: (path, data) => request(path, { method: 'PUT', body: JSON.stringify(data) }),
  delete: (path) => request(path, { method: 'DELETE' }),
}
