<script setup lang="ts">
import type { Message } from '@/modules/message/models/message'

const props = defineProps<{
  message: Message
}>()
</script>

<template>
  <div>
    <p>{{ message.content }}</p>
  </div>
</template>

<style scoped>
.message-component {
  height: 80px;
  background-color: purple;
}
</style>
