<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { conversationAPI } from '../api/conversation-api'
import { useConversationStore } from '@/shared/stores/conversation_store'
import  ConversationListItem  from './ConversationListItem.vue'
import ModalComponent from './ModalComponent.vue'

const conversationStore = useConversationStore()
const searchInput = ref('')
const showModal = ref(false)

onMounted(async () => {
  try {
    conversationStore.conversations = await conversationAPI.getUserConversations()
  } catch (error) {
    console.error('Failed to fetch conversations:', error)
  }
})

const filteredConversations = computed(() => {
  return conversationStore.conversations.filter((conv) =>
    conv.name.toLowerCase().includes(searchInput.value.toLowerCase()),
  )
})

watch(showModal, async (isOpen) => {
  if (isOpen) {
  }
})
</script>

<template>
  <div class="conversations-box">
    <header>
      <input type="text" placeholder="Type a message..." v-model="searchInput" />
      <button @click="showModal = true">+</button>
      <ModalComponent :show="showModal" @close="showModal = false" />
    </header>
    <ConversationListItem
      v-for="conversation in filteredConversations"
      :key="conversation.conversationId"
      :conversation="conversation"
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
