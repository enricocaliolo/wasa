import { defineStore } from 'pinia'
import { User } from '../../modules/auth/models/user'
import { ref } from 'vue'

export const useUserStore = defineStore('userStore', () => {
  const user = ref(User.fromJSON({}))
  function setUser(_user) {
    user.value = _user
  }
  function getUser() {
    return user.value
  }

  return { user, setUser, getUser }
})
