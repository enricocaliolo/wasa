<script setup>
import { ref, onBeforeUnmount } from "vue";
import { useUserStore } from "../../../shared/stores/user_store";
import { useConversationStore } from "../../../shared/stores/conversation_store";

const props = defineProps({
	user: Object,
});

const emit = defineEmits(["close"]);

const userStore = useUserStore();
const conversationStore = useConversationStore();

const showError = ref(false)
const errorMessage = ref('')

const icon = ref(null);
const iconURL = ref(null);

const newUsername = ref('')

const onIconChange = (event) => {
	icon.value = event.target.files[0];

	if (icon) {
		if (iconURL.value) {
			URL.revokeObjectURL(iconURL.value);
		}
		iconURL.value = URL.createObjectURL(icon.value);
	}
};

onBeforeUnmount(() => {
	if (iconURL.value) {
		URL.revokeObjectURL(iconURL.value);
	}
	icon.value = null;
	iconURL.value = null;
});

const saveChanges = async () => {
	try {
		if (newUsername.value !== '') {
			const user = await userStore.updateUsername(newUsername.value);
			if (user instanceof Error) {
				showError.value = true;
				errorMessage.value = "Username already taken!";
				setTimeout(() => {
					showError.value = false;
				}, 5000);
				return;
			}
		} else if (icon.value) {
			await userStore.updateIcon(icon.value);
		}
		conversationStore.setCurrentConversation(null)
		icon.value = null;
		iconURL.value = null;
		emit("close");
	} catch (e) {
		console.log(e);
	}
};
console.log(userStore.user)
console.log(icon.value)
console.log(iconURL.value)
</script>

<template>
	<div class="user-settings">
		<h2 class="title">User Settings</h2>
		<p v-if="showError" class="error-message">{{ errorMessage }}</p>

		<div class="settings-form">
			<div class="form-group">
				<label for="username" class="label"> Username </label>
				<input
					type="text"
					id="username"
					class="input"
					:placeholder="props.user.username"
					v-model="newUsername"
				/>
			</div>

			<div class="form-group">
				<label class="label"> Profile Icon </label>
				<div class="icon-container">
					<img
						:src="iconURL || `${userStore.user.displayIcon}`"
						alt="Current icon"
						class="icon-preview"
					/>	
					<input
						type="file"
						@change="onIconChange"
						accept="image/*"
						ref="iconInput"
					/>
				</div>
			</div>

			<button class="save-button" @click="emit('close')">Cancel</button>
			<button class="save-button" @click="saveChanges">
				Save changes
			</button>
		</div>
	</div>
</template>

<style scoped>
.user-settings {
	max-width: 600px;
	margin: 0 auto;
	padding: 20px;
}

.title {
	font-size: 24px;
	font-weight: bold;
	margin-bottom: 24px;
}

.settings-form {
	width: 100%;
}

.form-group {
	margin-bottom: 16px;
}

.label {
	display: block;
	font-size: 14px;
	font-weight: 500;
	margin-bottom: 8px;
}

.input {
	width: 100%;
	padding: 8px 12px;
	border: 1px solid #e2e8f0;
	border-radius: 8px;
	font-size: 16px;
	transition:
		border-color 0.2s,
		box-shadow 0.2s;
}

.input:focus {
	outline: none;
	border-color: #3b82f6;
	box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2);
}

.input-error {
	border-color: #ef4444;
}

.error-message {
	margin-top: 4px;
	font-size: 14px;
	color: #ef4444;
}

.icon-container {
	display: flex;
	align-items: center;
	gap: 16px;
}

.icon-preview {
	width: 48px;
	height: 48px;
	border-radius: 50%;
	object-fit: cover;
}

.icon-button {
	padding: 8px 16px;
	background-color: #f3f4f6;
	border: none;
	border-radius: 8px;
	cursor: pointer;
	transition: background-color 0.2s;
}

.icon-button:hover {
	background-color: #e5e7eb;
}

.save-button {
	width: 100%;
	padding: 8px 16px;
	background-color: #3b82f6;
	color: white;
	border: none;
	border-radius: 8px;
	font-size: 16px;
	cursor: pointer;
	transition: background-color 0.2s;
}

.save-button:hover:not(:disabled) {
	background-color: #2563eb;
}

.button-disabled {
	opacity: 0.5;
	cursor: not-allowed;
}

.hidden {
	display: none;
}

.error-message {
    background-color: red;
    color: white;
    padding: 10px;
    border-radius: 4px;
    margin-bottom: 16px;
}
</style>
