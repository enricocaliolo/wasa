<script setup lang="ts">
import { conversationAPI } from '@/modules/conversation/api/conversation-api'
import LoadingSpinner from '@/shared/components/LoadingSpinner.vue'
import { onMounted, ref } from 'vue'

const isLoading = ref(false)

onMounted(async () => {
  const messages = await conversationAPI.getUserConversation()
})
</script>

<template>
  <main v-if="isLoading">
    <div class="conversations-box"></div>
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
