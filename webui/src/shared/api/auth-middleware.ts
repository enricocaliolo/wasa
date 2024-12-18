import type { NavigationGuardNext, RouteLocationNormalized } from 'vue-router'

export function authMiddleware(
  to: RouteLocationNormalized,
  from: RouteLocationNormalized,
  next: NavigationGuardNext,
) {
  if (to.path !== '/login' && !localStorage.getItem('username')) {
    next('/login')
  } else {
    next()
  }
}
