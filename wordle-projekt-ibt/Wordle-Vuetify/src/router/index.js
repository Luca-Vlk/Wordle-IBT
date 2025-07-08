import { createRouter, createWebHistory } from 'vue-router'

import LoginPage from '../pages/LoginPage.vue'
import MainPage from '../pages/MainPage.vue'
import StatsPage from '../pages/StatsPage.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'login',
      component: LoginPage
    },
    {
      path: '/home/:user',
      name: 'home',
      component: MainPage
    },
    {
      path: '/stats/:user',
      name: 'stats',
      component: StatsPage
    }
  ],
})

export default router
