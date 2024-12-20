import { ref } from 'vue'
import { defineStore } from 'pinia'
import { Conversation } from '@/modules/conversation/models/conversation'
import { messagesAPI } from '@/modules/message/api/message_api'
import { Message } from '@/modules/message/models/message'
import { useUserStore } from './user'

export const useCurrentConversationStore = defineStore('currentConversation', () => {
  const currentConversation = ref<Conversation>()
  const userStore = useUserStore()

  function setCurrentConversation(conversation: Conversation) {
    currentConversation.value = conversation
  }

  async function sendMessage(new_message: string) {
    const data = await messagesAPI.sendMessage(
      currentConversation.value?.conversationId as number,
      new_message,
    )
    const message = Message.fromJSON(data)
    message.sender = userStore.getUser()

    currentConversation.value?.messages.push(message)

    return message
  }

  return { currentConversation, setCurrentConversation, sendMessage }
})
