import { ref } from 'vue'
import { defineStore } from 'pinia'
import { messagesAPI } from '@/modules/message/api/message_api'
import { useUserStore } from './user_store'
import { Message } from '../../modules/message/models/Message'

export const useConversationStore = defineStore('conversationStore', () => {
  const userStore = useUserStore()

  const conversations = ref([])
  const currentConversation = ref()
  const replyMessage= ref()

  function setCurrentConversation(conversation) {
    currentConversation.value = conversation
  }

  async function sendMessage(new_message) {
    const data = await messagesAPI.sendMessage(
      currentConversation.value.conversationId,
      new_message,
    )
    const message = Message.fromJSON(data)
    message.sender = userStore.getUser()

    currentConversation.value.messages.push(message)

    return message
  }

  async function sendRepliedMessage(new_message) {
    const data = await messagesAPI.sendRepliedMessage(
      currentConversation.value.conversationId,
      new_message,
      replyMessage.value
    )
    const message = Message.fromJSON(data)
    message.sender = userStore.getUser()

    currentConversation.value.messages.push(message)
    return message
  }

  async function sendForwardedMessage(source_conversation_id, destination_conversation_id, new_message) { 
    const data = await messagesAPI.sendForwardedMessage(
      source_conversation_id,
      destination_conversation_id,
      new_message
    )     
    const message = Message.fromJSON(data)
    message.sender = userStore.getUser()
    
    return message
  }

  async function addConversation(conv) {
    conversations.value.push(conv)
  }

  function setReplyMessage(message) {
    replyMessage.value = message
  }

  return {
    conversations,
    currentConversation,
    replyMessage,
    setCurrentConversation,
    sendMessage,
    sendRepliedMessage,
    sendForwardedMessage,
    addConversation,
    setReplyMessage
  }
})
