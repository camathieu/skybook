import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'jumps',
    component: () => import('./views/JumpList.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
