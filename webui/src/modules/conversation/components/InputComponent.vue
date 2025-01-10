<script setup>
import { useConversationStore } from '@/shared/stores/conversation_store'
import { ref } from 'vue'

const currentConversationStore = useConversationStore()
const messageInput = ref('')

const sendMessage = async () => {
  if(currentConversationStore.replyMessage) {
    await currentConversationStore.sendRepliedMessage(messageInput.value).then(() => {
      messageInput.value = ''
      currentConversationStore.setReplyMessage(null)
    })
    return
  }

  await currentConversationStore.sendMessage(messageInput.value).then(() => {
    messageInput.value = ''
  })
}
</script>

<template>
  <footer class="input-wrapper">
    <div v-if="currentConversationStore.replyMessage">
    {{ currentConversationStore.replyMessage.content}}
     <button @click="currentConversationStore.setReplyMessage(null)">RESET</button>
    </div>
    <div>
      <input type="text" placeholder="Type a message..." v-model="messageInput" />
      <button @click="sendMessage">Send</button>
    </div>
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
