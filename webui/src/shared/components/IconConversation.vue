<script setup>

import { useConversationStore } from '../stores/conversation_store';

// const conversationStore = useConversationStore()

const props = defineProps({
    conversation: Object,
})

</script>

<template>
    <div class="avatar-container">
		<img
				  v-if="props.conversation.photo"
				  :src="`${props.conversation.displayPhoto}`"
				  class="avatar-image"
		/>
		<span v-else class="avatar-placeholder">
		  {{ props.conversation.name.charAt(0).toUpperCase() }}
		</span>
			</div>
</template>

<style>
.avatar-container {
    width: 50px;
    height: 50px;
    position: relative;
}

.avatar-image {
    width: 100%;
    height: 100%;
    border-radius: 50%;
    object-fit: cover;
}

.avatar-placeholder {
    width: 40px;
    height: 40px;
    background-color: #e5e7eb;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 500;
    color: #374151;
}</style>