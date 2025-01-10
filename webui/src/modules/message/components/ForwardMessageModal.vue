<script setup>
import { ref } from 'vue'
import { useConversationStore } from '../../../shared/stores/conversation_store'

const props = defineProps({
  show: Boolean,
  message: String
})

const emit = defineEmits(['close'])

const searchInput = ref('')
const selectedConversation = ref(null)
const conversationStore = useConversationStore()

const selectConversation = (conversation) => {
  selectedConversation.value = conversation
}

function closeModal() {
  searchInput.value = ''
  selectedConversation.value = null
  emit('close')
}

async function forwardMessage() {

  if (selectedConversation.value === null) {
    alert('Please select a conversation to forward the message to.')
  }

  const source_conversation_id = conversationStore.currentConversation.conversationId
  const destination_conversation_id = selectedConversation.value.conversationId
  const message = props.message

  console.log("entra aqui")
  await conversationStore.sendForwardedMessage(source_conversation_id, destination_conversation_id, message)

  closeModal()
}
</script>

<template>
  <BaseModal :show="show" title="Forward Message" @close="closeModal">
    <div class="search-section">
      <input
        type="text"
        placeholder="Search conversations..."
        v-model="searchInput"
      />
    </div>

    <div class="possible-conversations">
    <div 
      v-for="conversation in conversationStore.conversations"
      :key="conversation.conversationId"
      class="conversation-item"
      :class="{ 'selected': selectedConversation === conversation }"
      @click="selectConversation(conversation)"
    >
      <div class="conversation-content">
        <span v-if="selectedConversation === conversation" class="check-mark">✓</span>
        {{ conversation.name }}
      </div>
    </div>
  </div>

    <template #footer>
      <button @click="closeModal">Cancel</button>
      <button
        @click="forwardMessage"
        :disabled="selectedConversation === null"
      >
        Forward
      </button>
    </template>
  </BaseModal>
</template>

<style scoped>
.search-section {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.possible-conversations {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  width: 100%;
}

.conversation-item {
  width: 100%;
  padding: 8px;
  text-align: left;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.conversation-item:hover {
  background-color: #f5f5f5;
}

.selected {
  background-color: #e3e3e3;
}

.conversation-content {
  display: flex;
  align-items: center;
  gap: 8px;
}

.check-mark {
  color: #4CAF50;
  font-weight: bold;
}

</style>