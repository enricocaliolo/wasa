<script setup>

import { ref, onBeforeUnmount, computed } from 'vue'
import { useConversationStore } from '../../../shared/stores/conversation_store';
import IconConversation from '../../../shared/components/IconConversation.vue';
import AddUsersModal from './AddUsersModal.vue';

const props = defineProps({
    conversation: Object,
});

const conversationStore = useConversationStore()

const newGroupName = ref('')
const photo = ref(null)
const photoURL = ref(null)

const showAddUsersModal = ref(false);

const photoChanged = async (event) => {
	photo.value = event.target.files[0]
	if (photo) {
		if (photoURL.value) {
			URL.revokeObjectURL(photoURL.value);
		}
		photoURL.value = URL.createObjectURL(photo.value);
	}
	await conversationStore.updateGroupPhoto(photo.value)
	
}

const hasNewName = computed(() => {
    return newGroupName.value.trim().length > 0
})

const isGroup = computed(() => conversationStore.currentConversation?.isGroup);

const changeGroupName = async () => {
	await conversationStore.updateGroupName(newGroupName.value)
	newGroupName.value = ''
}

const closeDetails = () => {
	photo.value = null
	if(photoURL.value) {
		URL.revokeObjectURL(photoURL.value)
	}
	conversationStore.toggleGroupDetails(false)
}

const leaveGroup = async() => {
	try{
		await conversationStore.leaveGroup()
	} catch(e) {
		alert(e.message)
	}
}

onBeforeUnmount(() => {
	if (photoURL.value) {
		URL.revokeObjectURL(photoURL.value);
	}
	photo.value = null;
  	photoURL.value = null;
});

</script>

<template>
	<div class="group-details-container">
	  <h3 class="group-title">Edit Group Details</h3>
	  
	  <div class="group-header">
		<div class="avatar-container">
		  <IconConversation :conversation="conversationStore.currentConversation" />
		  <label class="avatar-upload-label">
			<input 
			  type="file" 
			  @change="photoChanged"
			  accept="image/*"
			  class="avatar-input"
			/>
			<span class="upload-text">Change Photo</span>
		  </label>
		</div>
		
		<div class="group-info">
		  <div class="input-container">
			<input 
				type="text"
				class="name-input"
				:placeholder="props.conversation.name"
				v-model="newGroupName"
			/>
			<button 
			v-if="hasNewName" 
			class="submit-button"
			aria-label="Update name"
			@click="changeGroupName"
			>
			<svg 
				class="submit-icon" 
				width="20" 
				height="20" 
				viewBox="0 0 24 24" 
				fill="none" 
				stroke="currentColor" 
				stroke-width="2" 
				stroke-linecap="round" 
				stroke-linejoin="round"
				>
				<polyline points="20 6 9 17 4 12"></polyline>
    		</svg></button>
		  </div>
		  <p class="member-count">{{ props.conversation.participants.length }} members</p>
		</div>
	  </div>
  
	  <div class="participants-list">
		<div v-for="(participant, index) in props.conversation.participants" 
			:key="participant" 
			class="participant-item">
			<p>{{participant.username}}</p>
			<button v-if="index === 0" @click="leaveGroup">Leave</button>
		</div>
	  </div>

	  <AddUsersModal 
			:show="showAddUsersModal"
			:conversation="conversation"
			@close="showAddUsersModal = false"
		/>
	  <div class="buttons-container">
		<div class="add-user-button">
			<button @click="showAddUsersModal = true">ADD USER</button>
		</div>
		<div>
		  <button @click="closeDetails">Close</button>
		</div>
	  </div>

	</div>
  </template>
  
  <style scoped>
  .group-details-container {
	padding: 24px;
	background-color: #ffffff;
	border-radius: 8px;
	box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }
  
  .group-title {
	font-size: 24px;
	font-weight: 600;
	color: #333333;
	margin-bottom: 24px;
  }
  
  .group-header {
	display: flex;
	align-items: flex-start;
	gap: 16px;
	margin-bottom: 24px;
  }
  
  .avatar-container {
	display: flex;
	flex-direction: column;
	align-items: center;
	gap: 8px;
  }
  
  .group-avatar {
	width: 80px;
	height: 80px;
	border-radius: 50%;
	background-color: #f0f0f0;
	display: flex;
	align-items: center;
	justify-content: center;
	overflow: hidden;
	position: relative;
  }
  
  .avatar-image {
	width: 100%;
	height: 100%;
	object-fit: cover;
  }
  
  .avatar-placeholder {
	font-size: 32px;
	color: #666666;
  }
  
  .avatar-upload-label {
	cursor: pointer;
	display: inline-block;
  }
  
  .avatar-input {
	display: none;
  }
  
  .upload-text {
	font-size: 14px;
	color: #2563eb;
	text-decoration: underline;
	cursor: pointer;
  }
  
  .upload-text:hover {
	color: #1d4ed8;
  }
  
  .group-info {
	flex-grow: 1;
	display: flex;
	flex-direction: column;
	gap: 8px;
  }
  
  .name-input {
	font-size: 18px;
	font-weight: 500;
	color: #333333;
	padding: 8px 12px;
	border: 1px solid #e5e5e5;
	border-radius: 4px;
	width: 100%;
	outline: none;
  }
  
  .name-input:focus {
	border-color: #2563eb;
	box-shadow: 0 0 0 2px rgba(37, 99, 235, 0.1);
  }
  
  .member-count {
	font-size: 14px;
	color: #666666;
	margin: 0;
  }
  
  .participants-list {
	display: flex;
	flex-direction: column;
	gap: 8px;
	margin-top: 16px;
  }
  
  .participant-item {
	padding: 8px 12px;
	background-color: #f5f5f5;
	border-radius: 4px;
	font-size: 14px;
	color: #444444;
	margin: 0;
	display: flex;
	justify-content: space-between;
  }
  
  .participant-item:hover {
	background-color: #eeeeee;
  }

  .input-container {
  position: relative;
  width: 100%;
  display: flex;
  align-items: center;
}

.submit-button {
  position: absolute;
  right: 8px;
  background: none;
  border: none;
  padding: 4px;
  cursor: pointer;
  color: #2563eb;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.submit-button:hover {
  background-color: rgba(37, 99, 235, 0.1);
}

.submit-icon {
  width: 20px;
  height: 20px;
}
  </style>