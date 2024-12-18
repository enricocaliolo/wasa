import api from '../../../shared/api/api'

export const conversationAPI = {
  getUserConversation: async () => {
    const response = await api.get('/conversations', { headers: { Authorization: `Bearer 1` } })

    if (response.status == 200) {
      alert('mensagens')
      console.log(response.data)
    }
  },
}
