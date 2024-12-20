import { useUserStore } from '@/shared/stores/user'
import api from '../../../shared/api/api'
import { User, type UserDTO } from '../models/user'

export const authAPI = {
  login: async (_username: string): Promise<string> => {
    try {
      const userStore = useUserStore()
      const response = await api.put('/session', {
        username: _username,
      })

      if (response.data.id) {
        const user: UserDTO = {
          user_id: response.data.id,
          username: response.data.username,
          created_at: response.data.created_at,
        }
        userStore.setUser(User.fromJSON(user))
      }

      return response.data
    } catch (error) {
      throw error
    }
  },
}
