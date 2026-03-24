import { test, expect } from '@playwright/test'

// Helper: fill and submit the jump creation form
async function createJump(page, { date, type, dropzone, aircraft }) {
  await page.click('[data-testid="new-jump-btn"]')
  await page.waitForSelector('[data-testid="jump-modal"]')

  await page.fill('#f-date', date)
  await page.selectOption('#f-jump-type', type)
  // AutocompleteInput renders a plain <input> inside, target via label id
  await page.fill('#f-dropzone input', dropzone)
  await page.fill('#f-aircraft input', aircraft)

  await page.click('[data-testid="jump-form-submit"]')
  await expect(page.locator('[data-testid="jump-modal"]')).not.toBeVisible({ timeout: 10_000 })
}

test.describe('Jump CRUD flows', () => {

  test.beforeEach(async ({ page }) => {
    await page.goto('/')
    await page.waitForSelector('[data-testid="jump-list"]', { timeout: 10_000 })
  })

  test('create a jump and verify it appears in the list', async ({ page }) => {
    await createJump(page, {
      date: '2025-06-15T10:00',
      type: 'FF',
      dropzone: 'Skydive Empuriabrava',
      aircraft: 'Caravan',
    })

    await expect(page.locator('text=Skydive Empuriabrava').first()).toBeVisible()
    // Toast should appear confirming creation
    await expect(page.locator('.toast--success')).toContainText('Jump created')
  })

  test('edit a jump and verify changes persist', async ({ page }) => {
    await createJump(page, {
      date: '2025-06-15T10:00',
      type: 'FF',
      dropzone: 'Skydive Elsinore',
      aircraft: 'Twin Otter',
    })

    // Click the jump row/card to open the edit modal
    const jumpRow = page.locator('[data-testid="jump-row"], [data-testid="jump-card"]').first()
    await jumpRow.click()
    await page.waitForSelector('[data-testid="jump-modal"]')

    // Change the dropzone
    const dropzoneInput = page.locator('#f-dropzone input')
    await dropzoneInput.clear()
    await dropzoneInput.fill('Skydive Perris')
    await page.click('[data-testid="jump-form-submit"]')

    await expect(page.locator('text=Skydive Perris').first()).toBeVisible()
    await expect(page.locator('text=Skydive Elsinore')).not.toBeVisible()
    // Toast should appear confirming edit
    await expect(page.locator('.toast--success')).toContainText('Jump updated')
  })

  test('delete a jump and verify it is removed', async ({ page }) => {
    await createJump(page, {
      date: '2025-07-04T09:00',
      type: 'WS',
      dropzone: 'Skydive Chicago',
      aircraft: 'Otter',
    })

    // Open edit modal then delete
    const jumpRow = page.locator('[data-testid="jump-row"], [data-testid="jump-card"]').first()
    await jumpRow.click()
    await page.waitForSelector('[data-testid="jump-modal"]')

    await page.click('[data-testid="jump-delete-btn"]')
    await page.waitForSelector('[data-testid="confirm-modal"]')
    await page.click('[data-testid="confirm-delete-btn"]')

    await expect(page.locator('text=Skydive Chicago')).not.toBeVisible()
    // Toast should appear confirming deletion
    await expect(page.locator('.toast--success')).toContainText('Jump deleted')
  })

  test('insert at position and verify row count', async ({ page }) => {
    await createJump(page, { date: '2025-01-01T10:00', type: 'FF', dropzone: 'DZ Alpha', aircraft: 'Caravan' })
    await createJump(page, { date: '2025-01-02T10:00', type: 'FF', dropzone: 'DZ Beta', aircraft: 'Caravan' })

    // Insert at position 1
    await page.click('[data-testid="new-jump-btn"]')
    await page.waitForSelector('[data-testid="jump-modal"]')
    await page.fill('#f-number-insert', '1')
    await page.fill('#f-date', '2025-01-01T09:00')
    await page.selectOption('#f-jump-type', 'FF')
    await page.fill('#f-dropzone input', 'DZ Inserted')
    await page.fill('#f-aircraft input', 'Caravan')
    await page.click('[data-testid="jump-form-submit"]')
    await expect(page.locator('[data-testid="jump-modal"]')).not.toBeVisible()

    const rows = page.locator('[data-testid="jump-row"], [data-testid="jump-card"]')
    await expect(rows).toHaveCount(3)
  })

  test('search jumps filters the list', async ({ page }) => {
    await createJump(page, {
      date: '2025-09-01T11:00',
      type: 'HOP',
      dropzone: 'UniqueSearchTestDZ',
      aircraft: 'Cessna',
    })

    const searchInput = page.locator('[data-testid="search-input"]')
    await searchInput.fill('UniqueSearchTestDZ')
    await page.waitForTimeout(400)

    await expect(page.locator('text=UniqueSearchTestDZ').first()).toBeVisible()
  })
})
