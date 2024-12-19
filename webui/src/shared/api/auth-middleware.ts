import type { NavigationGuardNext, RouteLocationNormalized } from 'vue-router'
import { useUserStore } from '../stores/user'

export function authMiddleware(
  to: RouteLocationNormalized,
  from: RouteLocationNormalized,
  next: NavigationGuardNext,
) {
  const userStore = useUserStore()

  // Skip login page
  if (to.name === 'login') {
    next()
    return
  }

  // Check user ID in store
  if (userStore.user === -1) {
    next('/login')
  } else {
    next()
  }
}
