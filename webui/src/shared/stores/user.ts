import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useUserStore = defineStore('user', () => {
  const user = ref(-1)
  function setUserId(id: number) {
    user.value = id
  }

  return { user, setUserId }
})
