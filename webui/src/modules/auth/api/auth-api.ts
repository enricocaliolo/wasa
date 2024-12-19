import { useUserStore } from '@/shared/stores/user'
import api from '../../../shared/api/api'

export const authAPI = {
  login: async (_username: string): Promise<string> => {
    try {
      const userStore = useUserStore()
      const response = await api.put('/session', {
        username: _username,
      })

      if (response.data.id) {
        userStore.setUserId(response.data.id)
      }

      return response.data
    } catch (error) {
      throw error
    }
  },
}
