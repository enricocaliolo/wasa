import type { NavigationGuardNext, RouteLocationNormalized } from 'vue-router'

export function authMiddleware(
  to: RouteLocationNormalized,
  from: RouteLocationNormalized,
  next: NavigationGuardNext,
) {
  const username = localStorage.getItem('username')
  if (!username && to.path !== '/login') {
    next('/login')
  } else {
    next()
  }
}
