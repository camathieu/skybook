import { defineConfig, devices } from '@playwright/test'

/**
 * Playwright E2E test configuration.
 * Starts the Go backend + Vite dev server automatically before running tests.
 *
 * Set SKYBOOK_DATABASE_PATH to an in-memory path or temp file for isolation.
 */
export default defineConfig({
  testDir: './e2e',
  timeout: 30_000,
  retries: process.env.CI ? 2 : 0,
  reporter: [['list'], ['html', { open: 'never' }]],
  use: {
    baseURL: 'http://localhost:5173',
    trace: 'on-first-retry',
  },
  projects: [
    {
      name: 'Desktop Chrome',
      use: { ...devices['Desktop Chrome'] },
    },
    {
      name: 'Mobile Safari',
      use: { ...devices['iPhone SE'] },
    },
  ],
  webServer: [
    {
      // Go backend
      command: 'cd ../server && SKYBOOK_DATABASE_PATH=:memory: go run .',
      url: 'http://localhost:8080/health',
      timeout: 60_000,
      reuseExistingServer: !process.env.CI,
    },
    {
      // Vite dev server (proxies /api to :8080)
      command: 'npm run dev',
      url: 'http://localhost:5173',
      timeout: 30_000,
      reuseExistingServer: !process.env.CI,
    },
  ],
})
