<script setup lang="ts">
import LoadingSpinner from '@/shared/components/LoadingSpinner.vue'
import { onMounted, ref } from 'vue'
import { conversationAPI } from '../api/conversation-api'
import { Conversation } from '../models/conversation'
import { ConversationListItem, ConversationView } from '../components/index.ts'
import { useCurrentConversationStore } from '@/shared/stores/current_conversation_store.ts'

const currentConversationStore = useCurrentConversationStore()
const conversations = ref<Conversation[]>([])

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
</script>

<template>
  <main>
    <div class="conversations-box">
      <header>
        <h1>Search</h1>
      </header>
      <ConversationListItem
        v-for="conversation in conversations"
        :key="conversation.conversationId"
        :conversation="conversation"
        @click="changeCurrentConversation(conversation)"
      >
      </ConversationListItem>
    </div>
    <ConversationView
      v-if="currentConversationStore.currentConversation"
      :conversation="currentConversationStore.currentConversation"
    />
    <div v-else class="current-conversation"></div>
  </main>
</template>

<style scoped>
* {
  border: 1px solid red;
}

main {
  height: 100vh;
  width: 90vw;
  margin: 0 auto;
  max-width: 1500px;
  display: grid;
  grid-template-columns: 1fr 2fr;
  padding: 5vh 0;
}

.conversations-box {
  background-color: blue;
}

.current-conversation {
  background-color: red;
}
</style>
