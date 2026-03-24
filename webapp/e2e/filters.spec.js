import { test, expect } from '@playwright/test'

// Helper to create a jump quickly
async function seedJump(page, dropzone, aircraft, type = 'FF') {
  await page.click('[data-testid="new-jump-btn"]')
  await page.waitForSelector('[data-testid="jump-modal"]')
  await page.fill('#f-date', '2025-06-15T10:00')
  await page.fill('#f-dropzone input', dropzone)
  await page.fill('#f-aircraft input', aircraft)
  await page.selectOption('#f-jump-type', type)
  await page.click('[data-testid="jump-form-submit"]')
  await expect(page.locator('[data-testid="jump-modal"]')).not.toBeVisible({ timeout: 10_000 })
}

// Helper to click a CustomSelect and pick an option by text
async function selectCustomOption(page, ariaLabel, optionText) {
  const trigger = page.locator(`[aria-label="${ariaLabel}"] [data-testid="custom-select-trigger"]`)
  await trigger.click()
  const menu = page.locator(`[aria-label="${ariaLabel}"] [data-testid="custom-select-menu"]`)
  await expect(menu).toBeVisible()
  await menu.locator(`text="${optionText}"`).click()
  await expect(menu).not.toBeVisible()
}

test.describe('Filter bar (CustomSelect)', () => {

  test.beforeEach(async ({ page }) => {
    await page.goto('/')
    await page.waitForSelector('[data-testid="jump-list"]', { timeout: 10_000 })
    // Seed two jumps with distinct dropzones and aircraft
    await seedJump(page, 'FilterDZ-Alpha', 'Cessna182', 'FF')
    await seedJump(page, 'FilterDZ-Beta', 'TwinOtter', 'WS')
  })

  test('filter by dropzone via CustomSelect narrows the list', async ({ page }) => {
    await selectCustomOption(page, 'Filter by dropzone', 'FilterDZ-Alpha')
    await expect(page.locator('text=FilterDZ-Alpha').first()).toBeVisible()
    await expect(page.locator('text=FilterDZ-Beta')).not.toBeVisible()
  })

  test('filter by aircraft via CustomSelect narrows the list', async ({ page }) => {
    await selectCustomOption(page, 'Filter by aircraft', 'TwinOtter')
    await expect(page.locator('text=FilterDZ-Beta').first()).toBeVisible()
    await expect(page.locator('text=FilterDZ-Alpha')).not.toBeVisible()
  })

  test('filter by jump type via CustomSelect narrows the list', async ({ page }) => {
    await selectCustomOption(page, 'Filter by jump type', 'WS')
    await expect(page.locator('text=FilterDZ-Beta').first()).toBeVisible()
    await expect(page.locator('text=FilterDZ-Alpha')).not.toBeVisible()
  })

  test('clearing a filter via All option restores the full list', async ({ page }) => {
    // Apply dropzone filter first
    await selectCustomOption(page, 'Filter by dropzone', 'FilterDZ-Alpha')
    await expect(page.locator('text=FilterDZ-Beta')).not.toBeVisible()

    // Clear by choosing "All dropzones"
    await selectCustomOption(page, 'Filter by dropzone', 'All dropzones')
    await expect(page.locator('text=FilterDZ-Beta').first()).toBeVisible()
    await expect(page.locator('text=FilterDZ-Alpha').first()).toBeVisible()
  })

  test('CustomSelect dropdown opens and closes', async ({ page }) => {
    const trigger = page.locator('[aria-label="Filter by dropzone"] [data-testid="custom-select-trigger"]')
    const menu = page.locator('[aria-label="Filter by dropzone"] [data-testid="custom-select-menu"]')

    await expect(menu).not.toBeVisible()
    await trigger.click()
    await expect(menu).toBeVisible()
    await trigger.click()
    await expect(menu).not.toBeVisible()
  })
})
