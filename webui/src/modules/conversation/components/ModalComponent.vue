<script setup>
import { userAPI } from "@/modules/auth/api/user-api";
import { ref } from "vue";
import { conversationAPI } from "../api/conversation-api";
import { useConversationStore } from "@/shared/stores/conversation_store";
import { useUserStore } from "@/shared/stores/user_store";

defineProps({
	show: Boolean,
});

const conversationStore = useConversationStore();
const userStore = useUserStore();

const emit = defineEmits(["close", "submit"]);

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
	currentUsers.value = [];
	emit("close");
}

async function createConversation() {
	if (currentUsers.value.length > 2 && groupName.value == "") {
		alert("Please, insert a group name");
		return;
	} else if (currentUsers.value.length === 2) {
		groupName.value = currentUsers.value[1].username;
	}

	const conversation = await conversationAPI.createConversation(
		currentUsers.value.map((user) => user.userId),
		groupName.value,
	);
	conversationStore.addConversation(conversation);
	closeModal();
}
</script>

<template>
	<div v-if="show" class="modal-overlay">
		<div class="my-modal">
			<header>
				<input
					type="text"
					placeholder="Type a message..."
					v-model="searchInput"
				/>
				<button @click="addUser">ADD</button>
				<button @click="closeModal">X</button>
			</header>
			<div v-if="currentUsers.length > 2" class="group-name">
				<b>Group name:</b>
				<input
					type="text"
					placeholder="Group name"
					v-model="groupName"
				/>
			</div>
			<div class="current-users" v-if="currentUsers.length != 0">
				<div v-for="user in currentUsers" :key="user.userId">
					{{ user.username }}
				</div>
				<button @click="createConversation">CREATE</button>
			</div>
			<div class="test-nothing" v-else></div>
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
	z-index: 1000;
}

.my-modal {
	background-color: white;
}

.current-users {
	display: flex;
	flex-direction: column;
	gap: 5px;
}

.teste {
	background-color: white;
}

.test-nothing {
	background-color: red;
}
</style>
