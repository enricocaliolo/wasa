<script setup>
import { computed } from 'vue';
import { useConversationStore } from '@/shared/stores/conversation_store';
import { conversationAPI } from '../api/conversation-api';
import IconConversation from '../../../shared/components/IconConversation.vue';

const props = defineProps({
	conversation: Object,
});

// const getLastMessage = computed(() => {
//   const lastMessage = props.conversation.messages[props.conversation.messages.length - 1]
//   return lastMessage ? lastMessage.content : ''
// })

// const getLastMessageSender = computed(() => {
//   const lastMessage = props.conversation.messages[props.conversation.messages.length - 1]
//   return lastMessage ? lastMessage.sender.username : ''
// })

const conversationStore = useConversationStore();

async function getConversation(_conversation) {
	try {
		const conversation = conversationStore.conversations.find(
			(c) => c.conversationId === _conversation.conversationId
		);
		if (!conversation) {
			const messages = await conversationAPI.getConversation(
				conversation.conversationId
			);
			conversation.messages = messages || [];
		}

		conversationStore.setCurrentConversation(conversation);
	} catch (e) {
		console.log(e);
	}
}

const getLastMessage = computed(() => {
	const conversation = conversationStore.conversations.find(
		(c) => c.conversationId === props.conversation.conversationId
	);

	if (conversation.messages.length === 0) {
		return '';
	}

	const lastMessage = conversation.messages[conversation.messages.length - 1];
	const messageContent = lastMessage.displayContent;

	return conversation.isGroup
		? `${lastMessage.sender.username}: ${messageContent} - ${formattedDate(lastMessage.sentTime)}`
		: messageContent + ' - ' + formattedDate(lastMessage.sentTime);
});

const formattedDate = (sentTime) => {
	return new Intl.DateTimeFormat('en-US', {
		hour: '2-digit',
		minute: '2-digit',
	}).format(sentTime);
};
</script>

<template>
	<div class="conversation-preview" @click="getConversation(conversation)">
		<IconConversation :conversation="props.conversation" />
		<div class="last-message-container">
			<span class="name">
				{{ conversation.name }}
			</span>

			<span>{{ getLastMessage }} </span>
		</div>
	</div>
</template>

<style scoped>
.conversation-preview {
	height: 72px;
	width: 100%;
	background-color: green;
	padding: 0.25rem;
	margin-top: 1rem;
	border: 1px solid gold;
	display: flex;
}

.name {
	font-size: 1.5rem;
	font-weight: bold;
}

.last-message-container {
	display: flex;
	flex-direction: column;
	justify-content: space-between;
	margin-left: 1rem;
}
</style>
