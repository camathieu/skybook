import { test, expect } from '@playwright/test'

/**
 * Mobile layout E2E tests.
 * These run at 375px viewport width (iPhone SE profile from playwright.config.js).
 */
test.describe('Mobile responsive layout', () => {

  test('jump list renders as cards on mobile viewport', async ({ page }) => {
    await page.goto('/')
    await page.waitForSelector('[data-testid="jump-list"]', { timeout: 10_000 })

    // On mobile we expect jump-card elements instead of table rows
    // The app switches from table to cards at small breakpoints
    const viewport = page.viewportSize()
    if (viewport && viewport.width <= 767) {
      // Create one jump to verify the card layout
      await page.click('[data-testid="new-jump-btn"]')
      await page.waitForSelector('[data-testid="jump-modal"]')
      await page.fill('#jump-date', '2025-06-01')
      await page.selectOption('#jump-type', 'FF')
      await page.fill('#jump-dropzone', 'Mobile DZ')
      await page.fill('#jump-aircraft', 'Caravan')
      await page.click('[data-testid="jump-form-submit"]')
      await expect(page.locator('[data-testid="jump-modal"]')).not.toBeVisible()

      // Mobile shows cards, not table rows
      await expect(page.locator('[data-testid="jump-card"]').first()).toBeVisible()
    }
  })

  test('navigation is accessible on mobile', async ({ page }) => {
    await page.goto('/')
    // The nav should be visible or there should be a hamburger menu on mobile
    const nav = page.locator('nav, [data-testid="mobile-nav-btn"]')
    await expect(nav.first()).toBeVisible()
  })

  test('new jump button is reachable on mobile', async ({ page }) => {
    await page.goto('/')
    await page.waitForSelector('[data-testid="jump-list"]', { timeout: 10_000 })

    // The "New Jump" button must be visible and tappable
    const btn = page.locator('[data-testid="new-jump-btn"]')
    await expect(btn).toBeVisible()
    // Verify touch target is at least 44px tall (WCAG / mobile guidelines)
    const box = await btn.boundingBox()
    expect(box?.height).toBeGreaterThanOrEqual(44)
  })

  test('jump modal form is usable on mobile', async ({ page }) => {
    await page.goto('/')
    await page.waitForSelector('[data-testid="jump-list"]', { timeout: 10_000 })

    await page.click('[data-testid="new-jump-btn"]')
    await page.waitForSelector('[data-testid="jump-modal"]', { timeout: 5_000 })

    // Modal should fill the screen or be scrollable on mobile
    const modal = page.locator('[data-testid="jump-modal"]')
    await expect(modal).toBeVisible()

    // Form fields should be usable (44px min touch targets)
    const dateInput = page.locator('#jump-date')
    await expect(dateInput).toBeVisible()
    const box = await dateInput.boundingBox()
    expect(box?.height).toBeGreaterThanOrEqual(44)
  })
})
