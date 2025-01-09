import { useUserStore } from '@/shared/stores/user_store'
import api from '../../../shared/api/api'
import { User } from '../models/user'

export const userAPI = {
  login: async (_username) => {
    try {
      const userStore = useUserStore()
      const response = await api.put('/session', {
        username: _username,
      })

      if (response.data.id) {
        const user = {
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
  findUser: async (_username) => {
    try {
      const response = await api.get('/users/search', {
        params: {
          username: _username,
        },
      })

      const user = User.fromJSON(response.data)
      user.username = _username

      return user
    } catch (error) {
      throw error
    }
  },
}
