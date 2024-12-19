<script setup lang="ts">
// import { conversationAPI } from '@/modules/conversation/api/conversation-api'
import LoadingSpinner from '@/shared/components/LoadingSpinner.vue'
import { computed, onMounted, ref } from 'vue'
import { conversationAPI } from '../api/conversation-api'
import { Conversation } from '../models/conversation'
import { ConversationListItem, ConversationView } from '../components/index.ts'

const isLoading = ref(true)
const conversations = ref<Conversation[]>([])
const currentConversation = ref<Conversation>()

onMounted(async () => {
  try {
    conversations.value = await conversationAPI.getUserConversation()
  } catch (error) {
    console.error('Failed to fetch conversations:', error)
  }
})

function changeCurrentConversation(conversation: Conversation) {
  currentConversation.value = conversation
  debugger
}
</script>

<template>
  <main v-if="isLoading">
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
    <ConversationView v-if="currentConversation" :conversation="currentConversation" />
    <div v-else class="current-conversation"></div>
  </main>
  <main v-else>
    <LoadingSpinner />
  </main>
</template>

<style scoped>
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
