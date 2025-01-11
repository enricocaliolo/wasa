<script setup>
import { useConversationStore } from '@/shared/stores/conversation_store'
import { ref, onBeforeUnmount } from 'vue'

const conversationStore = useConversationStore()
const messageInput = ref('')
const selectedImage = ref(null)
const file = ref(null)
const fileInput = ref(null)

const sendMessage = async () => {
  try{
    if(conversationStore.replyMessage) {
      // await conversationStore.sendRepliedMessage(messageInput.value)
      // await conversationStore.sendRepliedMessage(
      //   messageInput.value,
      //   content_type = file
      // )
      messageInput.value = ''
      conversationStore.setReplyMessage(null)
    return
    }

  if(selectedImage) {
    // await conversationStore.sendImage(file.value)
    selectedImage.value = null
    file.value = null
    messageInput.value = ''
    fileInput.value.value = ''
    return
  } 

  await conversationStore.sendMessage(messageInput.value)
  messageInput.value = ''
  } catch (e) {
    console.log(e)
  }
}

const onFileSelected = (event) => {
 file.value = event.target.files[0]
 selectedImage.value = URL.createObjectURL(file.value)
}

onBeforeUnmount(() => {
 if (selectedImage.value) {
   URL.revokeObjectURL(selectedImage.value)
 }
})

</script>

<template>
  <footer class="input-wrapper">
    <div v-if="conversationStore.replyMessage">
    {{ conversationStore.replyMessage.contentType === 'image' ? 'Image' : conversationStore.replyMessage.content }}
     <button @click="conversationStore.setReplyMessage(null)">RESET</button>
    </div>
    <div>
      <div class="container">
        <input type="file" ref="fileInput" @change="onFileSelected" accept="image/*">
        
        <!-- <div v-if="selectedImage" class="preview">
          <img :src="selectedImage" alt="Preview" />
        </div> -->
      </div>
      <input v-if="selectedImage === null" type="text" placeholder="Type a message..." v-model="messageInput" />
      <button @click="sendMessage">Send</button>
    </div>
  </footer>
</template>

<style scoped>
.input-wrapper {
  border: 1px solid;
  display: 100%;
  padding: 1.5em;
  background-color: sandybrown;
  display: flex;
  flex-direction: row;
  justify-content: space-around;
}

.container {
  margin: 20px;
}

.preview {
  margin-top: 20px;
}

.preview img {
  max-width: 300px;
  height: auto;
}
</style>
