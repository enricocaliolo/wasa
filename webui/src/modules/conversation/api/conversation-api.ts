import api from '../../../shared/api/api'
import { Conversation, type ConversationDTO } from '../models/conversation'

export const conversationAPI = {
  getUserConversation: async () => {
    const response = await api.get('/conversations', { headers: { Authorization: `Bearer 1` } })

    if (response.status == 200) {
      return response.data.map((json: ConversationDTO) => new Conversation(json))
    }
  },
}
