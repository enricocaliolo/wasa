import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '@/modules/auth/view/LoginView.vue'
import HomeView from '@/modules/conversation/view/HomeView.vue'
import { authMiddleware } from '../api/auth-middleware'
// import { authMiddleware } from '../api/auth-middleware'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', name: 'home', component: HomeView },
    {
      path: '/login',
      name: 'login',
      component: LoginView,
    },
  ],
})

router.beforeEach(authMiddleware)

export default router
