import api from '../../../shared/api/api'

export const authAPI = {
  login: async (_username: string): Promise<string> => {
    try {
      const response = await api.put('/session', {
        username: _username,
      })

      if (response.data.username) {
        localStorage.setItem('username', response.data.username)
      }

      return response.data
    } catch (error) {
      throw error
    }
  },
}
