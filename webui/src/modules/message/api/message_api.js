import api from '../../../shared/api/api'

export const messagesAPI = {
  sendMessage: async (conversation_id, message) => {
    const response = await api.post(`/conversations/${conversation_id}`, {
      content: message,
      content_type: 'text',
    })
    if (response.status === 201) {
      return response.data
    }
  },
  sendRepliedMessage: async(conversation_id, message, replied_to_message) => {
    const response = await api.post(`/conversations/${conversation_id}/reply`, {
      content: message,
      content_type: 'text',
      replied_to: replied_to_message.messageId
    })
    if(response.status === 201) {
      return response.data
    }
  }
}
