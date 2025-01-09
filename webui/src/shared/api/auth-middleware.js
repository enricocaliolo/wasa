import { useUserStore } from "../stores/user_store"

export function authMiddleware(to, from, next) {
  const userStore = useUserStore()

  if (to.name === 'login') {
    next()
    return
  }
  if (userStore.user.userId === undefined) {
    next('/login')
  } else {
    next()
  }
}