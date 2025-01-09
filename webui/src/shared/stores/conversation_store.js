import { ref } from 'vue'
import { defineStore } from 'pinia'
import { messagesAPI } from '@/modules/message/api/message_api'
import { useUserStore } from './user_store'
import { Message } from '../../modules/message/models/Message'

export const useConversationStore = defineStore('conversationStore', () => {
  const conversations = ref([])
  const currentConversation = ref()
  const userStore = useUserStore()

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

  async function addConversation(conv) {
    conversations.value.push(conv)
  }

  return {
    conversations,
    currentConversation,
    setCurrentConversation,
    sendMessage,
    addConversation,
  }
})
