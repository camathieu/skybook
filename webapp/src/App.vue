<script setup>
import AppHeader from './components/AppHeader.vue'
import ToastContainer from './components/ToastContainer.vue'
</script>

<template>
  <AppHeader />
  <main class="app-content">
    <router-view v-slot="{ Component }">
      <transition name="fade" mode="out-in">
        <component :is="Component" />
      </transition>
    </router-view>
  </main>
  <ToastContainer />
</template>

<style scoped>
.app-content {
  flex: 1;
  padding-top: var(--header-height);
}

@media (max-width: 767px) {
  .app-content {
    padding-top: var(--header-height-mobile);
  }
}
</style>
