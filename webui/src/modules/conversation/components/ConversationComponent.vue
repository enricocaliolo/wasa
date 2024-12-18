<script setup lang="ts">
import type { Conversation } from '../models/conversation'

const props = defineProps<{
  conversation: Conversation
}>()
</script>

<template>
  <div class="conversation-preview">
    <p>{{ conversation.name == '' ? conversation.participants[0].name : conversation.name }}</p>
  </div>
</template>

<style scoped>
.conversation-preview {
  height: 72px;
  width: 100%;
  background-color: green;
  padding: 1rem;
  margin-top: 1rem;
}
</style>
