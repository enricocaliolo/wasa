<script setup>
import { ref } from "vue";
// import BaseModal from './BaseModal.vue'
import { userAPI } from "@/modules/auth/api/user-api";
import { conversationAPI } from "../api/conversation-api";
import { useConversationStore } from "@/shared/stores/conversation_store";
import { useUserStore } from "@/shared/stores/user_store";

const props = defineProps({
	show: Boolean,
});

const conversationStore = useConversationStore();
const userStore = useUserStore();

const emit = defineEmits(["close"]);

const searchInput = ref("");
const groupName = ref("");
const currentUsers = ref([userStore.user]);

async function addUser() {
	const user = await userAPI.findUser(searchInput.value);
	currentUsers.value.push(user);
	searchInput.value = "";
}

function closeModal() {
	searchInput.value = "";
	groupName.value = "";
	currentUsers.value = [userStore.user];
	emit("close");
}

async function createConversation() {
	await conversationStore.createConversation({
		currentUsers: currentUsers.value,
		groupName: groupName.value,
	});

	closeModal();
}
</script>

<template>
	<BaseModal :show="show" title="Create Conversation" @close="closeModal">
		<div class="search-section">
			<input
				type="text"
				placeholder="Search users..."
				v-model="searchInput"
			/>
			<button @click="addUser">Add</button>
		</div>

		<div v-if="currentUsers.length > 2" class="group-name">
			<label>Group name:</label>
			<input
				type="text"
				placeholder="Enter group name"
				v-model="groupName"
			/>
		</div>

		<div class="current-users" v-if="currentUsers.length > 0">
			<div
				v-for="user in currentUsers"
				:key="user.userId"
				class="user-item"
			>
				{{ user.username }}
			</div>
		</div>

		<template #footer>
			<button @click="closeModal">Cancel</button>
			<button
				@click="createConversation"
				:disabled="currentUsers.length < 2"
			>
				Create
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

.group-name {
	margin: 1rem 0;
}

.current-users {
	display: flex;
	flex-direction: column;
	gap: 0.5rem;
}

.user-item {
	padding: 0.5rem;
	background-color: #f5f5f5;
	border-radius: 4px;
}
</style>
