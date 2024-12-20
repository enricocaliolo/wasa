<script setup lang="ts">
import type { Conversation } from '../models/conversation'
import { InputComponent } from '../components/index'
import MessageComponent from '@/modules/message/components/MessageComponent.vue'

const props = defineProps<{
  conversation: Conversation
}>()
</script>

<template>
  <div class="conversation-wrapper">
    <header>
      <h1>{{ conversation.name }}</h1>
    </header>
    <div class="messages-box">
      <MessageComponent
        v-for="message in conversation.messages"
        :key="message.messageId"
        :message="message"
      />
    </div>
    <InputComponent />
  </div>
</template>

<style scoped>
header {
  padding: 1em;
  background-color: grey;
}

.conversation-wrapper {
  display: grid;
  grid-template-rows: auto 1fr auto;
  min-height: 0;
}

.messages-box {
  background-color: lightblue;
  overflow-y: auto;
  padding: 1em;
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  min-height: 0;
  gap: 10px;
}
</style>
