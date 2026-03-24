<script setup>
</script>

<template>
  <div class="jump-list">
    <div class="empty-state">
      <div class="empty-icon">✦</div>
      <h1 class="empty-title">Your logbook is empty</h1>
      <p class="empty-subtitle">Log your first jump to get started</p>
      <button class="btn-primary empty-cta">
        <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
          <path d="M8 2a.75.75 0 0 1 .75.75v4.5h4.5a.75.75 0 0 1 0 1.5h-4.5v4.5a.75.75 0 0 1-1.5 0v-4.5h-4.5a.75.75 0 0 1 0-1.5h4.5v-4.5A.75.75 0 0 1 8 2Z"/>
        </svg>
        Log First Jump
      </button>
    </div>
  </div>
</template>

<style scoped>
.jump-list {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem 1.5rem;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 4rem 1rem;
  min-height: 50vh;
}

.empty-icon {
  font-size: 3rem;
  margin-bottom: 1.5rem;
  background: linear-gradient(135deg, var(--color-accent-orange), var(--color-accent-teal));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.7; transform: scale(1.05); }
}

.empty-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0 0 0.5rem;
}

.empty-subtitle {
  font-size: 1rem;
  color: var(--color-text-muted);
  margin: 0 0 2rem;
}

.empty-cta {
  padding: 0.75rem 2rem;
  font-size: 1rem;
}

@media (max-width: 767px) {
  .jump-list {
    padding: 1.5rem 1rem;
  }

  .empty-state {
    padding: 3rem 1rem;
    min-height: 40vh;
  }
}
</style>
