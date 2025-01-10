<script setup>
import { useUserStore } from '@/shared/stores/user_store'
import { computed } from 'vue'
import { useConversationStore } from '../../../shared/stores/conversation_store'
import ForwardMessageModal from './ForwardMessageModal.vue'
import {ref} from 'vue'

const props = defineProps({
  message: Object,
})

const useStore = useUserStore()
const conversationStore = useConversationStore()

const showModal = ref(false)

const isFromUser = computed(() => {
  return props.message.sender.userId === useStore.user.userId
})

const replyMessage = () => {
  conversationStore.setReplyMessage(props.message)
}

</script>

<template>
  <div class="message-wrapper">
    <div class="message-box" :class="isFromUser ? 'own-message' : 'not-own-message'">
      <div class="replied-message" v-if="message.repliedToMessage">{{ message.repliedToMessage.content }}</div>
      <p v-if="message.isForwarded === true"><i>Forwarded</i></p>
      <div class="wrap">
        <button @click="replyMessage">reply</button>
        <button @click="showModal = true">forward</button>
        <ForwardMessageModal :message="message.content" :show="showModal" @close="showModal = false"/>
        <p> {{ message.content }}</p>
      </div>
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

.replied-message {
  border: 1px solid red;
}

.wrap {
  display: flex;
}
</style>
