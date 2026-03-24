<script setup>
import { ref } from 'vue'

const mobileMenuOpen = ref(false)

function toggleMobileMenu() {
  mobileMenuOpen.value = !mobileMenuOpen.value
}

function closeMobileMenu() {
  mobileMenuOpen.value = false
}
</script>

<template>
  <header class="header">
    <div class="header-inner">
      <!-- Left: hamburger (mobile) + logo -->
      <div class="header-left">
        <button
          class="hamburger"
          @click="toggleMobileMenu"
          aria-label="Toggle navigation"
        >
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
            <template v-if="!mobileMenuOpen">
              <line x1="3" y1="6" x2="21" y2="6" />
              <line x1="3" y1="12" x2="21" y2="12" />
              <line x1="3" y1="18" x2="21" y2="18" />
            </template>
            <template v-else>
              <line x1="6" y1="6" x2="18" y2="18" />
              <line x1="6" y1="18" x2="18" y2="6" />
            </template>
          </svg>
        </button>

        <router-link to="/" class="logo" @click="closeMobileMenu">
          <span class="logo-icon">✦</span>
          <span class="logo-text">SkyBook</span>
        </router-link>
      </div>

      <!-- Center: nav tabs (desktop) -->
      <nav class="nav-tabs" aria-label="Main navigation">
        <router-link to="/" class="nav-tab" exact-active-class="active">
          Jumps
        </router-link>
        <span class="nav-tab disabled" title="Coming in v9">BASE</span>
        <span class="nav-tab disabled" title="Coming in v10">Tunnel</span>
      </nav>

      <!-- Right: CTA -->
      <div class="header-right">
        <button class="btn-primary new-jump-btn" aria-label="New jump">
          <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
            <path d="M8 2a.75.75 0 0 1 .75.75v4.5h4.5a.75.75 0 0 1 0 1.5h-4.5v4.5a.75.75 0 0 1-1.5 0v-4.5h-4.5a.75.75 0 0 1 0-1.5h4.5v-4.5A.75.75 0 0 1 8 2Z"/>
          </svg>
          <span class="new-jump-label">New Jump</span>
          <span class="kbd desktop-only">N</span>
        </button>
      </div>
    </div>

    <!-- Mobile nav drawer -->
    <transition name="slide">
      <nav v-if="mobileMenuOpen" class="mobile-nav" @click="closeMobileMenu">
        <router-link to="/" class="mobile-nav-item active">
          Jumps
        </router-link>
        <span class="mobile-nav-item disabled">BASE — coming soon</span>
        <span class="mobile-nav-item disabled">Tunnel — coming soon</span>
      </nav>
    </transition>
  </header>

  <!-- Backdrop -->
  <transition name="fade">
    <div v-if="mobileMenuOpen" class="backdrop" @click="closeMobileMenu" />
  </transition>
</template>

<style scoped>
.header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  background-color: var(--color-surface-800);
  border-bottom: 1px solid var(--color-surface-600);
}

.header-inner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: var(--header-height);
  padding: 0 1.5rem;
  max-width: 1400px;
  margin: 0 auto;
}

/* Left section */
.header-left {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.hamburger {
  display: none;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  background: none;
  border: none;
  color: var(--color-text-secondary);
  cursor: pointer;
  border-radius: 8px;
  transition: background-color var(--transition-fast);
  -webkit-tap-highlight-color: transparent;
}
.hamburger:hover {
  background-color: var(--color-surface-700);
}

.logo {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  text-decoration: none;
  color: var(--color-text-primary);
}
.logo-icon {
  font-size: 1.25rem;
  background: linear-gradient(135deg, var(--color-accent-orange), var(--color-accent-teal));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}
.logo-text {
  font-weight: 700;
  font-size: 1.125rem;
  letter-spacing: -0.025em;
}

/* Nav tabs (desktop) */
.nav-tabs {
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.nav-tab {
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--color-text-muted);
  text-decoration: none;
  border-radius: 8px;
  transition: color var(--transition-fast), background-color var(--transition-fast);
  user-select: none;
}
.nav-tab:hover:not(.disabled) {
  color: var(--color-text-secondary);
  background-color: var(--color-surface-700);
}
.nav-tab.active {
  color: var(--color-text-primary);
  background-color: var(--color-surface-700);
}
.nav-tab.disabled {
  opacity: 0.4;
  cursor: default;
}

/* Right section */
.header-right {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.new-jump-btn {
  gap: 0.375rem;
  padding: 0.5rem 1rem;
  font-size: 0.8125rem;
}

.desktop-only {
  display: inline-flex;
}

/* Mobile nav drawer */
.mobile-nav {
  display: none;
  flex-direction: column;
  padding: 0.5rem 1rem 1rem;
  background-color: var(--color-surface-800);
  border-bottom: 1px solid var(--color-surface-600);
}

.mobile-nav-item {
  display: flex;
  align-items: center;
  padding: 0.875rem 1rem;
  min-height: 44px;
  font-size: 1rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  text-decoration: none;
  border-radius: 8px;
  transition: background-color var(--transition-fast);
  -webkit-tap-highlight-color: transparent;
}
.mobile-nav-item:hover:not(.disabled) {
  background-color: var(--color-surface-700);
}
.mobile-nav-item.active {
  color: var(--color-text-primary);
}
.mobile-nav-item.disabled {
  opacity: 0.4;
}

.backdrop {
  position: fixed;
  inset: 0;
  z-index: 90;
  background-color: rgba(0, 0, 0, 0.5);
}

/* Slide transition */
.slide-enter-active,
.slide-leave-active {
  transition: transform var(--transition-normal), opacity var(--transition-normal);
}
.slide-enter-from,
.slide-leave-to {
  transform: translateY(-8px);
  opacity: 0;
}

/* =========== Responsive =========== */

@media (max-width: 767px) {
  .header-inner {
    height: var(--header-height-mobile);
    padding: 0 1rem;
  }

  .hamburger {
    display: flex;
  }

  .nav-tabs {
    display: none;
  }

  .mobile-nav {
    display: flex;
  }

  .new-jump-label {
    display: none;
  }

  .new-jump-btn {
    padding: 0.5rem;
    border-radius: 50%;
    width: 44px;
    height: 44px;
  }

  .desktop-only {
    display: none;
  }
}
</style>
