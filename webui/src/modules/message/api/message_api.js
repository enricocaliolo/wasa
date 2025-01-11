import api from '../../../shared/api/api'

export const messagesAPI = {
  sendMessage: async (conversation_id, message) => {
    try{const response = await api.post(`/conversations/${conversation_id}`, {
      content: message,
      content_type: 'text',
    })
    if (response.status === 201) {
      return response.data
    }}
    catch(e) {
      console.log(e)
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
  },
  sendForwardedMessage: async(source_conversation_id, destination_conversation_id, message) => {
    const response = await api.post(`/conversations/${source_conversation_id}/forward`, {
      content: message,
      content_type: 'text',
      destination_conversation_id: destination_conversation_id
    })
    if(response.status === 201) {
      return response.data
    }
  },
  commentMessage: async (conversation_id, message_id, comment) => {
    const response = await api.put(`/conversations/${conversation_id}/messages/${message_id}`, {
      reaction: comment
    })
    if(response.status === 202) {
      return response.data
    }
  },
  uncommentMessage: async (conversation_id, message_id, reaction_id) => {
    const response = await api.delete(`/conversations/${conversation_id}/messages/${message_id}/reactions/${reaction_id}`)
    if(response.status === 202) {
      return true
    }
  },
}
