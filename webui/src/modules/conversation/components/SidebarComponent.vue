<script setup>
import { computed, onMounted, ref } from 'vue';
import { useConversationStore } from '@/shared/stores/conversation_store';
import ConversationListItem from './ConversationListItem.vue';
import ConversationModal from './ConversationModal.vue';
import { useUserStore } from '../../../shared/stores/user_store';
import ChangeUserSetting from '../../auth/components/ChangeUserSetting.vue';

const conversationStore = useConversationStore();
const userStore = useUserStore();
const user = userStore.getUser();

const searchInput = ref('');
const showModal = ref(false);
const showUserConfig = ref(false);

onMounted(async () => {
	try {
		await conversationStore.getUserConversations();
	} catch (error) {
		console.error('Failed to fetch conversations:', error);
	}
});

const filteredConversations = computed(() => {
	return conversationStore.conversations.filter((conv) =>
		conv.name.toLowerCase().includes(searchInput.value.toLowerCase())
	);
});
</script>

<template>
	<ChangeUserSetting
		v-if="showUserConfig"
		:user="user"
		@close="showUserConfig = false"
	/>
	<div v-else class="conversations-box">
		<header>
			<input
				type="text"
				placeholder="Type a message..."
				v-model="searchInput"
			/>
			<button @click="showModal = true">+</button>
			<button @click="showUserConfig = true">USER</button>
			<ConversationModal :show="showModal" @close="showModal = false" />
		</header>
		<div
			v-if="
				conversationStore.conversations &&
				conversationStore.conversations.length > 0
			"
		>
			<ConversationListItem
				v-for="conversation in filteredConversations"
				:key="conversation.conversationId"
				:conversation="conversation"
			>
			</ConversationListItem>
		</div>
		<div v-else>NO CONVERSATIONS</div>
	</div>
</template>

<style scoped>
header {
	display: flex;
	padding: 0.5em;
	gap: 5px;
}

.conversations-box {
	background-color: blue;
}
</style>
