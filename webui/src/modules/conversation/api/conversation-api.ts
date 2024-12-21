import { Message, type MessageDTO } from '@/modules/message/models/message'
import api from '../../../shared/api/api'
import { Conversation, type ConversationDTO } from '../models/conversation'

export const conversationAPI = {
  getUserConversations: async () => {
    const response = await api.get('/conversations', { headers: { Authorization: `Bearer 1` } })

    if (response.status == 200) {
      return response.data.map((json: ConversationDTO) => new Conversation(json))
    }
  },
  getConversation: async (conversation_id: number) => {
    const response = await api.get(`/conversations/${conversation_id}`)
    return response.data.map((json: MessageDTO) => new Message(json))
  },
  createConversation: async () => {
    const response = await api.post('/conversations')
  },
}
