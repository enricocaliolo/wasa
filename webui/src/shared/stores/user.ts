import { ref } from 'vue'
import { defineStore } from 'pinia'
import { User } from '@/modules/auth/models/user'

export const useUserStore = defineStore('user', () => {
  const user = ref<User | undefined>()
  function setUser(_user: User) {
    user.value = _user
  }

  return { user, setUser }
})
