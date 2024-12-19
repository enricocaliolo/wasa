<script setup lang="ts">
import { computed } from 'vue'
import type { Conversation } from '../models/conversation'

const props = defineProps<{
  conversation: Conversation
}>()

const getLastMessage = computed(() => {
  const lastMessage = props.conversation.messages[props.conversation.messages.length - 1]
  return lastMessage ? lastMessage.content : ''
})

const getLastMessageSender = computed(() => {
  const lastMessage = props.conversation.messages[props.conversation.messages.length - 1]
  return lastMessage ? lastMessage.sender.username : ''
})
</script>

<template>
  <div class="conversation-preview">
    <span class="name">
      {{ conversation.name }}
    </span>
    <p>
      <span>{{ getLastMessageSender }}: </span>{{ getLastMessage }}
    </p>
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
