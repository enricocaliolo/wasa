import { ref, type Ref } from 'vue'
import { defineStore } from 'pinia'
import { User } from '@/modules/auth/models/user'

export const useUserStore = defineStore('user', () => {
  const user = ref<User>() as Ref<User>
  function setUser(_user: User) {
    user.value = _user
  }
  function getUser() {
    return user.value
  }

  return { user, setUser, getUser }
})
