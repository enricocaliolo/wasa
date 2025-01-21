<script setup>
import InputComponent from '@/modules/conversation/components/InputComponent.vue';
import MessageComponent from '../../message/components/MessageComponent.vue';
import { useConversationStore } from '../../../shared/stores/conversation_store';
import { useUserStore } from '../../../shared/stores/user_store';
import { watch, onMounted } from 'vue';
import GroupDetails from '../components/GroupDetails.vue';
import IconConversation from '../../../shared/components/IconConversation.vue';

const props = defineProps({
	conversation: Object,
});

const conversationStore = useConversationStore();
const userStore = useUserStore();
console.log(userStore.user);

function changeGroupName() {
	conversationStore.toggleGroupDetails(!conversationStore.showGroupDetails);
}
onMounted(async () => {
	if (conversationStore.currentConversation?.messages.length > 0) {
		const unseenMessages = props.conversation.messages
			.filter(
				(message) =>
					!message.isSeenBy(userStore.user.userId) &&
					message.sender.userId !== userStore.user.userId
			)
			.map((message) => message.messageId);

		if (unseenMessages.length > 0) {
			conversationStore.markMessagesSeen(unseenMessages);
		}
	}
});

watch(
	() => conversationStore.currentConversation?.messages,
	(newMessages, oldMessages) => {
		if (!newMessages || !oldMessages) return;

		const unseenMessages = newMessages
			.filter(
				(message) =>
					!message.isSeenBy(userStore.user.userId) &&
					message.sender.userId !== userStore.user.userId
			)
			.map((message) => message.messageId);

		if (unseenMessages.length > 0) {
			conversationStore.markMessagesSeen(unseenMessages);
		}
	},
	{ deep: true }
);
</script>

<template>
	<div class="conversation-wrapper">
		<header @click="changeGroupName">
			<IconConversation :conversation="conversationStore.currentConversation" />
			<h1>{{ conversation.name }}</h1>
		</header>
		<GroupDetails
			v-if="conversationStore.showGroupDetails"
			:conversation="conversation"
		/>
		<div v-else class="messages-box">
			<MessageComponent
				v-for="message in conversation.messages"
				:key="message.messageId"
				:message="message"
				:data-message-id="message.messageId"
			/>
		</div>
		<InputComponent v-if="!conversationStore.showGroupDetails" />
	</div>
</template>

<style scoped>
header {
	padding: 1em;
	background-color: grey;
	display: flex;
	align-items: center;
	gap: 1em;
	cursor: pointer;
}

.conversation-wrapper {
	display: grid;
	grid-template-rows: auto 1fr auto;
	min-height: 0;
}

.messages-box {
	background-color: lightblue;
	overflow-y: auto;
	padding: 1em;
	display: flex;
	flex-direction: column;
	justify-content: flex-end;
	min-height: 0;
	gap: 10px;
}
</style>
