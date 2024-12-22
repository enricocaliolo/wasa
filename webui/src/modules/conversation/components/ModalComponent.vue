<script setup lang="ts">
import { userAPI } from '@/modules/auth/api/user-api'
import type { User } from '@/modules/auth/models/user'
import { ref } from 'vue'
import { conversationAPI } from '../api/conversation-api'
import { useConversationStore } from '@/shared/stores/conversation_store'
import { useUserStore } from '@/shared/stores/user'

defineProps({
  show: Boolean,
})

const conversationStore = useConversationStore()
const userStore = useUserStore()

const emit = defineEmits(['close', 'submit'])

const searchInput = ref('')
const currentUsers = ref<User[]>([userStore.user])

async function addUser() {
  const user = await userAPI.findUser(searchInput.value)
  currentUsers.value.push(user)
  searchInput.value = ''
}

function closeModal() {
  searchInput.value = ''
  currentUsers.value = []
  emit('close')
}

async function createConversation() {
  const conversation = await conversationAPI.createConversation(
    currentUsers.value.map((user) => user.userId),
  )
  conversationStore.addConversation(conversation)
  closeModal()
}
</script>

<template>
  <div v-if="show" class="modal-overlay">
    <div class="modal">
      <header>
        <input type="text" placeholder="Type a message..." v-model="searchInput" />
        <button @click="addUser">ADD</button>
        <button @click="closeModal">X</button>
      </header>
      <div class="current-users" v-if="currentUsers.length != 0">
        <div v-for="user in currentUsers" :key="user.userId">
          {{ user.username }}
        </div>
        <button @click="createConversation">CREATE</button>
      </div>
      <div class="test-nothing" v-else>teste</div>
    </div>
  </div>
</template>

<style scoped>
header {
  display: flex;
  padding: 0.5em;
  gap: 5px;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: grid;
  place-items: center;
}

.modal {
  background: white;
  padding: 2rem;
  border-radius: 8px;
}

.current-users {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.test-nothing {
  background-color: red;
}
</style>
