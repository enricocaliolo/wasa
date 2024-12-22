<script setup lang="ts">
import { computed } from 'vue'
import type { Conversation } from '../models/conversation'
import { useConversationStore } from '@/shared/stores/conversation_store'
import api from '@/shared/api/api'
import { conversationAPI } from '../api/conversation-api'

const props = defineProps<{
  conversation: Conversation
}>()

// const getLastMessage = computed(() => {
//   const lastMessage = props.conversation.messages[props.conversation.messages.length - 1]
//   return lastMessage ? lastMessage.content : ''
// })

// const getLastMessageSender = computed(() => {
//   const lastMessage = props.conversation.messages[props.conversation.messages.length - 1]
//   return lastMessage ? lastMessage.sender.username : ''
// })

const currentConversationStore = useConversationStore()

async function getConversation(conversation: Conversation) {
  const messages = await conversationAPI.getConversation(conversation.conversationId)
  console.log(messages)
  conversation.messages = messages

  currentConversationStore.setCurrentConversation(conversation)
}
</script>

<template>
  <div class="conversation-preview" @click="getConversation(conversation)">
    <span class="name">
      {{ conversation.name }}
    </span>
    <!-- <p>
      <span>{{ getLastMessageSender }}: </span>{{ getLastMessage }}
    </p> -->
  </div>
</template>

<style scoped>
.conversation-preview {
  height: 72px;
  width: 100%;
  background-color: green;
  padding: 1rem;
  margin-top: 1rem;
  border: 1px solid gold;
}

.name {
  font-size: 1.5rem;
  font-weight: bold;
}
</style>
