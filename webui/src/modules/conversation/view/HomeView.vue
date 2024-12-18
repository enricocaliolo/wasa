<script setup lang="ts">
// import { conversationAPI } from '@/modules/conversation/api/conversation-api'
import LoadingSpinner from '@/shared/components/LoadingSpinner.vue'
import ConversationComponent from '../components/ConversationComponent.vue'
import { computed, onMounted, ref } from 'vue'
import { conversationAPI } from '../api/conversation-api'
import type { Conversation } from '../models/conversation'

const isLoading = ref(true)
const conversations = ref<Conversation[]>([])

onMounted(async () => {
  try {
    conversations.value = await conversationAPI.getUserConversation()
  } catch (error) {
    console.error('Failed to fetch conversations:', error)
  }
})
</script>

<template>
  <main v-if="isLoading">
    <div class="conversations-box">
      <header>
        <h1>Search</h1>
      </header>
      <ConversationComponent
        v-for="conversation in conversations"
        :key="conversation.conversationId"
        :conversation="conversation"
      >
      </ConversationComponent>
    </div>
    <div class="current-conversation"></div>
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
