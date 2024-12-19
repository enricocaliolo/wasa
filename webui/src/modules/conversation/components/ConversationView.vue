<script setup lang="ts">
import type { Conversation } from '../models/conversation'

const props = defineProps<{
  conversation: Conversation
}>()
</script>

<template>
  <header>
    <h1>{{ conversation.name == '' ? conversation.participants[0].name : conversation.name }}</h1>
  </header>
  <div class="messages-box">
    <!-- <MessageComponent
      v-for="message in conversation.messages"
      :key="message.messageId"
      :message="message"
    /> -->
  </div>
</template>

<style scoped></style>
