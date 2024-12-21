<script setup lang="ts">
import type { Message } from '@/modules/message/models/message'
import { useUserStore } from '@/shared/stores/user'
import { computed } from 'vue'

const props = defineProps<{
  message: Message
}>()

const useStore = useUserStore()

const isFromUser = computed(() => {
  return props.message.sender.userId === useStore.user?.userId
})
</script>

<template>
  <div class="message-wrapper">
    <div class="message-box" :class="isFromUser ? 'own-message' : 'not-own-message'">
      <p>{{ message.content }}</p>
    </div>
  </div>
  <br />
</template>

<style scoped>
.message-wrapper {
  display: flex;
  width: 100%;
}

.message-box {
  background-color: yellow;
  height: 80px;
  width: 50%;
}

.own-message {
  margin-left: auto;
  text-align: right;
}

.not-own-message {
  margin-right: auto;
  text-align: left;
}
</style>
