import { ref } from 'vue'
import { defineStore } from 'pinia'
import { messagesAPI } from '@/modules/message/api/message_api'
import { useUserStore } from './user_store'
import { Message } from '../../modules/message/models/message'
import { Reaction } from '../../modules/message/models/reaction'
import { conversationAPI } from '../../modules/conversation/api/conversation-api'

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

  async function sendImage(image) {
    const data = await messagesAPI.sendMessage(
      currentConversation.value.conversationId,
      image,
      content_type = 'image'
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
    console.log(message)
    replyMessage.value = message
  }

  async function addReaction(message_id, _reaction) {
    const data = await messagesAPI.commentMessage(
      currentConversation.value.conversationId,
      message_id,
      _reaction
    )
    const reaction = Reaction.fromJSON(data)
    
    const messageToUpdate = currentConversation.value.messages.find(
      message => message.messageId === message_id
    )
    
    if (messageToUpdate) {
      if (!messageToUpdate.reactions) {
        messageToUpdate.reactions = []
      }
      messageToUpdate.reactions.push(reaction)
    }

    return reaction
  }

  async function deleteReaction(message_id, reaction_id) {
    await messagesAPI.uncommentMessage(
      currentConversation.value.conversationId,
      message_id,
      reaction_id
    )

    for (const message of currentConversation.value.messages) {
      if (message.messageId === message_id) {
        message.reactions = message.reactions.filter(
          reaction => reaction.reactionId !== reaction_id
        )
      }
    }

    return
  }

  async function updateGroupName(name) {
    await conversationAPI.updateGroupName(currentConversation.value.conversationId, name)
    currentConversation.value.name = name
    return
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
    setReplyMessage,
    addReaction,
    deleteReaction,
    updateGroupName,
    sendImage
  }
})