<script setup>
import { useConversationStore } from '@/shared/stores/conversation_store';
import { ref, onBeforeUnmount } from 'vue';

const conversationStore = useConversationStore();
const messageInput = ref('');
const selectedImage = ref(null);
const file = ref(null);
const fileInput = ref(null);

const sendMessage = async () => {
	if (conversationStore.replyMessage) {
		await conversationStore.sendMessage({
			content: messageInput.value,
			replied_to_message: -1,
		});
		messageInput.value = '';
		conversationStore.setReplyMessage(null);
		return;
	}

	if (file.value) {
		await conversationStore.sendMessage({
			content: file.value,
			content_type: 'image',
		});
		selectedImage.value = null;
		file.value = null;
		messageInput.value = '';
		fileInput.value.value = '';
		return;
	}
	await conversationStore.sendMessage({ content: messageInput.value });
	messageInput.value = '';
};

const onFileSelected = (event) => {
	file.value = event.target.files[0];
	selectedImage.value = URL.createObjectURL(file.value);
};

onBeforeUnmount(() => {
	if (selectedImage.value) {
		URL.revokeObjectURL(selectedImage.value);
	}
});
</script>

<template>
	<footer class="input-wrapper">
		<div v-if="conversationStore.replyMessage" class="reply-container">
			<p>Replying to:</p>
			<div>
				{{
					conversationStore.replyMessage.contentType === 'image'
						? 'Image'
						: conversationStore.replyMessage.content
				}}
				<button @click="conversationStore.setReplyMessage(null)">RESET</button>
			</div>
		</div>
		<div class="row-input">
			<div class="container">
				<div v-if="selectedImage" class="preview">
					<img :src="selectedImage" alt="Preview" />
				</div>
				<div v-else>
					<input
						id="img"
						type="file"
						ref="fileInput"
						@change="onFileSelected"
						accept="image/*"
						class="file-input"
					/>
					<label class="file-label" for="img">Upload</label>
				</div>

				
			</div>
			<div class="msg-input-container">
				<input
					v-if="selectedImage === null"
					type="text"
					placeholder="Type a message..."
					v-model="messageInput"
					class="msg-input"
				/>
				<button @click="sendMessage" class="btn-msg-input">Send</button>
			</div>
		</div>
	</footer>
</template>

<style scoped>
.input-wrapper {
	border: 1px solid red;
	background-color: sandybrown;
	display: flex;
	flex-direction: row;
	justify-content: space-around;
	align-items: center;
}

.container {
	margin: 20px;
}

.row-input {
	display: flex;
	justify-content: center;
}

.msg-input-container {
	display: flex;
	flex-direction: row;
	justify-content: space-between;
}

.msg-input {
	height: 50%;
	align-self: center;
	margin-right: 10px;
}

.btn-msg-input {
	height: 50%;
	align-self: center;
}

.file-input {
	display: none;
}

.file-label {
	border: 1px solid black;
	padding: 10px;
}

.file-label:hover {
	cursor: pointer;
}

.preview {
  margin-top: 20px;
  display: flex;
  justify-content: center;
  align-items: center;
  border: 2px solid #ddd;
  border-radius: 10px;
  padding: 10px;
  background-color: #f9f9f9;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  width: 150px;
  height: 150px;
}

.preview img {
  max-width: 100%;
  height: auto;
  border-radius: 8px;
  object-fit: cover;
}

.reply-container {
	display: flex;
	flex-direction: column;
	align-items: center;
	padding: 10px;
	box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}
</style>
