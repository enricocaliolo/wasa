<script setup>
import { useUserStore } from '@/shared/stores/user_store'
import { computed } from 'vue'
import { useConversationStore } from '../../../shared/stores/conversation_store'
import ForwardMessageModal from './ForwardMessageModal.vue'
import { ref, onMounted } from 'vue'

const props = defineProps({
  message: Object,
})

const userStore = useUserStore()
const conversationStore = useConversationStore()
const showModal = ref(false)

const showEmojiPicker = ref(false)
const emojis = ['👍', '❤️', '😊', '😂', '😮', '😢']
const userReaction = ref(null)

onMounted(() => {
  if (props.message.reactions) {
    const reaction = props.message.reactions.find(r => r.user.userId === userStore.user.userId)
    if (reaction) {
      userReaction.value = reaction
    }
  }
})

const isFromUser = computed(() => {
  return props.message.sender.userId === userStore.user.userId
})

const replyMessage = () => {
  conversationStore.setReplyMessage(props.message)
}

const addReaction = async (emoji) => {
  try{
    const reaction = await conversationStore.addReaction(props.message.messageId, emoji)
    userReaction.value = reaction
    
    showEmojiPicker.value = false
  } catch(e) {
    console.log(e)
  }
}

const handleEmojiButton =async() => {
  try{
    if (userReaction.value) {
      await conversationStore.deleteReaction(props.message.messageId, userReaction.value.reactionId)
      userReaction.value = null
      return
  }} catch(e) {
    console.log(e)
  }

  showEmojiPicker.value = !showEmojiPicker.value
}

const userHasReaction = computed(() => {
  if (userReaction.value) {
    return true
  } 
  return false
})


</script>

<template>
  <div class="message-wrapper">
    <div class="message-container" :class="isFromUser ? 'own-message' : 'not-own-message'">
      <div class="message-actions">
        <button @click="replyMessage" class="action-button" title="Reply">
          ↩️
        </button>
        <button @click="showModal = true" class="action-button" title="Forward">
          ↪️
        </button>
        <div class="emoji-picker-container">
          <button @click="handleEmojiButton" class="action-button" title="Reaction">
            {{ userHasReaction ? 'x' : '+' }}
          </button>
          </div>
      </div>
      
      <div class="message-box">
        <div v-if="message.isForwarded" class="forwarded-label">
          <i>Forwarded</i>
        </div>

        <div v-if="message.repliedToMessage" class="replied-message">
          {{ message.repliedToMessage.content }}
        </div>

        <div class="message-content">
          {{ message.content }}
        </div>

        <div v-if="message.reactions.length > 0" class="reactions-container">
          <div class="reactions-bubble" @click="showEmojis">
            <div class="reactions-list">
              <span v-for="reaction in message.reactions.slice(0, 2)" :key="reaction.id">
                {{ reaction.reaction }}
              </span>
              <span v-if="message.reactions.length > 2">...</span>
            </div>
            <span  class="reactions-count">{{ message.reactions.length }}</span>
          </div>
        </div>

        <div v-if="showEmojiPicker" class="emoji-picker">
              <button 
                v-for="emoji in emojis" 
                :key="emoji"
                @click="addReaction(emoji)"
                class="emoji-button"
              >
                {{ emoji }}
              </button>
            </div>
      </div>
      
    </div>
    <ForwardMessageModal 
      :message="message.content" 
      :show="showModal" 
      @close="showModal = false"
    />
  </div>
</template>

<style scoped>
.message-wrapper {
  display: flex;
  width: 100%;
  margin: 8px 0;
  padding: 0 16px;
}

.message-container {
  display: flex;
  align-items: center;
  max-width: 70%;
  gap: 8px;
}

.message-box {
  background-color: #f0f0f0;
  border-radius: 12px;
  padding: 12px;
  min-width: 80px;
}

.own-message {
  margin-left: auto;
  flex-direction: row-reverse;
}

.own-message .message-box {
  background-color: #0084ff;
  color: white;
}

.not-own-message {
  margin-right: auto;
}

.forwarded-label {
  font-size: 0.85rem;
  color: #666;
  margin-bottom: 4px;
}

.own-message .forwarded-label {
  color: rgba(255, 255, 255, 0.8);
}

.replied-message {
  border-left: 3px solid #ccc;
  padding-left: 8px;
  margin-bottom: 8px;
  font-size: 0.9rem;
  color: #666;
  background-color: rgba(0, 0, 0, 0.05);
  border-radius: 4px;
}

.own-message .replied-message {
  border-left-color: rgba(255, 255, 255, 0.4);
  color: rgba(255, 255, 255, 0.9);
  background-color: rgba(255, 255, 255, 0.1);
}

.message-content {
  word-break: break-word;
  line-height: 1.4;
}

.reactions-container {
  margin-top: 4px;
}

.reactions-bubble {
  display: inline-flex;
  align-items: center;
  background-color: rgba(255, 255, 255, 0.9);
  border-radius: 12px;
  padding: 2px 8px;
  font-size: 0.8rem;
  color: #333;
}

.message-actions {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.action-button {
  background-color: transparent;
  border: none;
  border-radius: 50%;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  padding: 0;
  font-size: 1.2rem;
  transition: transform 0.2s ease;
}

.action-button:hover {
  transform: scale(1.2);
}
</style>