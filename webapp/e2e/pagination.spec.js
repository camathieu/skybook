import { test, expect } from '@playwright/test'

test.describe('Pagination and Per Page', () => {

  test.beforeEach(async ({ page }) => {
    await page.goto('/')
    await page.waitForSelector('[data-testid="jump-list"]', { timeout: 10_000 })
  })

  test('per-page button sends correct per_page param and updates active state', async ({ page }) => {
    // Wait for the initial jump list API call to complete
    await page.waitForResponse(resp =>
      resp.url().includes('/api/v1/jumps') && !resp.url().includes('autocomplete') && resp.status() === 200
    )

    // Click the 50 per-page button and verify the API is called with per_page=50
    const [response] = await Promise.all([
      page.waitForResponse(resp =>
        resp.url().includes('/api/v1/jumps') &&
        resp.url().includes('per_page=50') &&
        !resp.url().includes('autocomplete') &&
        resp.status() === 200
      ),
      page.click('[data-testid="per-page-50"]'),
    ])

    // Verify the response was successful and contains the right perPage value
    const body = await response.json()
    expect(body.perPage).toBe(50)

    // Verify the 50 button is now active
    const btn50 = page.locator('[data-testid="per-page-50"]')
    await expect(btn50).toHaveClass(/active/)

    // Verify the 25 button is no longer active
    const btn25 = page.locator('[data-testid="per-page-25"]')
    await expect(btn25).not.toHaveClass(/active/)
  })

  test('switching per-page resets to page 1', async ({ page }) => {
    // Wait for initial load
    await page.waitForResponse(resp =>
      resp.url().includes('/api/v1/jumps') && !resp.url().includes('autocomplete') && resp.status() === 200
    )

    // Click 100 per-page and verify the request includes page=1
    const [response] = await Promise.all([
      page.waitForResponse(resp =>
        resp.url().includes('/api/v1/jumps') &&
        resp.url().includes('per_page=100') &&
        !resp.url().includes('autocomplete') &&
        resp.status() === 200
      ),
      page.click('[data-testid="per-page-100"]'),
    ])

    const body = await response.json()
    expect(body.perPage).toBe(100)
    expect(body.page).toBe(1)
  })
})
