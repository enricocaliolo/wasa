import api from '../../../shared/api/api'

export const messagesAPI = {
  sendMessage: async (conversation_id, message,) => {
    try{
      const response = await api.post(`/conversations/${conversation_id}`, {
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
  sendImage: async (conversation_id, image) => {
    try{
      const response = await api.post(
        `/conversations/${conversation_id}`, 
        image, // Send the raw file directly
        {
          headers: {
            'Content-Type': 'image/jpeg',
        }}
      )
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

  // sendMessage:async ({ 
  //   conversation_id, 
  //   content, 
  //   content_type = 'text',
  //   replied_to_message = null,
  //   destination_conversation_id = null
  // }) => {
  //   try {
  //     let endpoint = `/conversations/${conversation_id}`
  //     let payload = {}
  //     let config = {}
      
  //     // Set up config for image content type
  //     if (content_type === 'image/jpeg') {
  //       config = {
  //         headers: {
  //           'Content-Type': 'image/jpeg'
  //         }
  //       }
  //       payload = content // Raw image file
  //     } else {
  //       payload = {
  //         content,
  //         content_type
  //       }
  //     }
      
  //     // Handle different message scenarios
  //     if (destination_conversation_id) {
  //       endpoint += '/forward'
  //       payload = content_type === 'image/jpeg' ? content : {
  //         content,
  //         content_type,
  //         destination_conversation_id
  //       }
  //     } else if (replied_to_message) {
  //       endpoint += '/reply'
  //       payload = content_type === 'image/jpeg' ? content : {
  //         content,
  //         content_type,
  //         replied_to: replied_to_message.messageId
  //       }
  //     }
  
  //     const response = await api.post(endpoint, payload, config)
  
  //     if (response.status === 201) {
  //       return response.data
  //     }
  //   } catch (error) {
  //     console.log(error)
  //   }
  // }

}
