import api from '../../../shared/api/api'
import { Conversation } from '../models/conversation'
import { Message } from '../../message/models/Message'

export const conversationAPI = {
  getUserConversations: async () => {
    const response = await api.get('/conversations')

    if (response.status === 200) { 
      return response.data.map((json) => Conversation.fromJSON(json))
    }
},
  getConversation: async (conversation_id) => {
    const response = await api.get(`/conversations/${conversation_id}`)
    return response.data.map((json) => new Message(json))
  },
  createConversation: async (members) => {
    const response = await api.post('/conversations', {
      members: members,
    })
    return new Conversation(response.data)
  },
}
