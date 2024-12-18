import api from '../../../shared/api/api'

export const conversationAPI = {
  getUserConversation: async () => {
    localStorage.getItem('username')
    const response = await api.get('/conversations')

    if (response.status == 200) {
      alert('mensagens')
      console.log(response.data)
    }
  },
}
