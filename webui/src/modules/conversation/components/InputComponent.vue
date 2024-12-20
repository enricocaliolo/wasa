<script setup lang="ts">
import { useCurrentConversationStore } from '@/shared/stores/current_conversation_store'
import { ref } from 'vue'

const currentConversationStore = useCurrentConversationStore()
const messageInput = ref('')

const sendMessage = () => {
  currentConversationStore.sendMessage(messageInput.value)
}
</script>

<template>
  <footer class="input-wrapper">
    <input type="text" placeholder="Type a message..." v-model="messageInput" />
    <button @click="sendMessage">Send</button>
  </footer>
</template>

<style scoped>
.input-wrapper {
  border: 1px solid;
  display: 100%;
  padding: 1.5em;
  background-color: sandybrown;
  display: flex;
  flex-direction: row;
  justify-content: space-around;
}
</style>
