<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import type { Conversation } from '../models/conversation'
import { conversationAPI } from '../api/conversation-api'
import { useCurrentConversationStore } from '@/shared/stores/current_conversation_store'
import { ConversationListItem } from './index'

const currentConversationStore = useCurrentConversationStore()
const conversations = ref<Conversation[]>([])
const searchInput = ref('')

onMounted(async () => {
  try {
    conversations.value = await conversationAPI.getUserConversation()
  } catch (error) {
    console.error('Failed to fetch conversations:', error)
  }
})

function changeCurrentConversation(conversation: Conversation) {
  currentConversationStore.setCurrentConversation(conversation)
}

const filteredConversations = computed(() => {
  return conversations.value.filter((conv) =>
    conv.name.toLowerCase().includes(searchInput.value.toLowerCase()),
  )
})
</script>

<template>
  <div class="conversations-box">
    <header>
      <input type="text" placeholder="Type a message..." v-model="searchInput" />
      <button>+</button>
    </header>
    <ConversationListItem
      v-for="conversation in filteredConversations"
      :key="conversation.conversationId"
      :conversation="conversation"
      @click="changeCurrentConversation(conversation)"
    >
    </ConversationListItem>
  </div>
</template>

<style scoped>
header {
  display: flex;
  padding: 0.5em;
  gap: 5px;
}

.conversations-box {
  background-color: blue;
}
</style>
