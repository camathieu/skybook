// Global Vitest test setup
// This runs before each test file.

// Make ResizeObserver available in jsdom (used by some Vue components)
global.ResizeObserver = class ResizeObserver {
  observe() {}
  unobserve() {}
  disconnect() {}
}
